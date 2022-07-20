// Modifications Copyright 2022 The Klaytn Authors
// Copyright 2021 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from eth/gasprice/feehistory.go (2021/11/09).
// Modified and improved for the klaytn development.

package gasprice

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"sync/atomic"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/misc"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/params"
)

var (
	logger               = log.NewModuleLogger(log.NodeCnGasPrice)
	errInvalidPercentile = errors.New("invalid reward percentile")
	errRequestBeyondHead = errors.New("request beyond head block")
)

const (
	// maxBlockFetchers is the max number of goroutines to spin up to pull blocks
	// for the fee history calculation.
	maxBlockFetchers = 1
)

// blockFees represents a single block for processing
type blockFees struct {
	// set by the caller
	blockNumber uint64
	header      *types.Header
	block       *types.Block // only set if reward percentiles are requested
	receipts    types.Receipts
	// filled by processBlock
	results processedFees
	err     error
}

// processedFees contains the results of a processed block and is also used for caching
type processedFees struct {
	reward               []*big.Int
	baseFee, nextBaseFee *big.Int
	gasUsedRatio         float64
}

// txGasAndReward is sorted in ascending order based on reward
type (
	txGasAndReward struct {
		gasUsed uint64
		reward  *big.Int
	}
	sortGasAndReward []txGasAndReward
)

func (s sortGasAndReward) Len() int { return len(s) }
func (s sortGasAndReward) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortGasAndReward) Less(i, j int) bool {
	return s[i].reward.Cmp(s[j].reward) < 0
}

// processBlock takes a blockFees structure with the blockNumber, the header and optionally
// the block field filled in, retrieves the block from the backend if not present yet and
// fills in the rest of the fields.
func (oracle *Oracle) processBlock(bf *blockFees, percentiles []float64) {
	chainconfig := oracle.backend.ChainConfig()
	// TODO-Klaytn: If we implement baseFee feature like Ethereum does, we should set it from header, not constant.
	if bf.results.baseFee = bf.header.BaseFee; bf.results.baseFee == nil {
		bf.results.baseFee = new(big.Int).SetUint64(params.ZeroBaseFee)
	}
	// TODO-Klaytn: If we implement baseFee feature like Ethereum does, we should calculate nextBaseFee from parent block header.
	if chainconfig.IsMagmaForkEnabled(big.NewInt(int64(bf.blockNumber + 1))) {
		bf.results.nextBaseFee = misc.NextMagmaBlockBaseFee(bf.header, chainconfig.Governance.KIP71)
	} else {
		bf.results.nextBaseFee = new(big.Int).SetUint64(params.ZeroBaseFee)
	}

	// There is no GasLimit in Klaytn, so it is enough to use pre-defined constant in api package as now.
	bf.results.gasUsedRatio = float64(bf.header.GasUsed) / float64(params.UpperGasLimit)
	if len(percentiles) == 0 {
		// rewards were not requested, return null
		return
	}
	if bf.block == nil || (bf.receipts == nil && len(bf.block.Transactions()) != 0) {
		logger.Error("Block or receipts are missing while reward percentiles are requested")
		return
	}

	bf.results.reward = make([]*big.Int, len(percentiles))
	if len(bf.block.Transactions()) == 0 {
		// return an all zero row if there are no transactions to gather data from
		for i := range bf.results.reward {
			bf.results.reward[i] = new(big.Int)
		}
		return
	}

	sorter := make(sortGasAndReward, len(bf.block.Transactions()))
	for i := range bf.block.Transactions() {
		// TODO-Klaytn: If we change the fixed unit price policy and add baseFee feature, we should re-calculate reward.
		reward := bf.block.Header().BaseFee
		if reward == nil {
			reward = new(big.Int).SetUint64(chainconfig.UnitPrice)
		}
		sorter[i] = txGasAndReward{gasUsed: bf.receipts[i].GasUsed, reward: reward}
	}
	sort.Sort(sorter)

	var txIndex int
	sumGasUsed := sorter[0].gasUsed

	for i, p := range percentiles {
		thresholdGasUsed := uint64(float64(bf.block.GasUsed()) * p / 100)
		for sumGasUsed < thresholdGasUsed && txIndex < len(bf.block.Transactions())-1 {
			txIndex++
			sumGasUsed += sorter[txIndex].gasUsed
		}
		bf.results.reward[i] = sorter[txIndex].reward
	}
}

// resolveBlockRange resolves the specified block range to absolute block numbers while also
// enforcing backend specific limitations.
// Pending block does not exist in Klaytn, so there is no logic to look up pending blocks.
// This part has a different implementation with Ethereum.
// Note: an error is only returned if retrieving the head header has failed. If there are no
// retrievable blocks in the specified range then zero block count is returned with no error.
func (oracle *Oracle) resolveBlockRange(ctx context.Context, lastBlock rpc.BlockNumber, blocks int) (uint64, int, error) {
	var headBlock rpc.BlockNumber
	// query either pending block or head header and set headBlock
	if lastBlock == rpc.PendingBlockNumber {
		// pending block not supported by backend, process until latest block
		lastBlock = rpc.LatestBlockNumber
		blocks--
		if blocks == 0 {
			return 0, 0, nil
		}
	}
	// if pending block is not fetched then we retrieve the head header to get the head block number
	if latestHeader, err := oracle.backend.HeaderByNumber(ctx, rpc.LatestBlockNumber); err == nil {
		headBlock = rpc.BlockNumber(latestHeader.Number.Uint64())
	} else {
		return 0, 0, err
	}
	if lastBlock == rpc.LatestBlockNumber {
		lastBlock = headBlock
	} else if lastBlock > headBlock {
		return 0, 0, fmt.Errorf("%w: requested %d, head %d", errRequestBeyondHead, lastBlock, headBlock)
	}
	// ensure not trying to retrieve before genesis
	if rpc.BlockNumber(blocks) > lastBlock+1 {
		blocks = int(lastBlock + 1)
	}
	return uint64(lastBlock), blocks, nil
}

// FeeHistory returns data relevant for fee estimation based on the specified range of blocks.
// The range can be specified either with absolute block numbers or ending with the latest
// or pending block. Backends may or may not support gathering data from the pending block
// or blocks older than a certain age (specified in maxHistory). The first block of the
// actually processed range is returned to avoid ambiguity when parts of the requested range
// are not available or when the head has changed during processing this request.
// Three arrays are returned based on the processed blocks:
// - reward: the requested percentiles of effective priority fees per gas of transactions in each
//   block, sorted in ascending order and weighted by gas used.
// - baseFee: base fee per gas in the given block
// - gasUsedRatio: gasUsed/gasLimit in the given block
// Note: baseFee includes the next block after the newest of the returned range, because this
// value can be derived from the newest block.
func (oracle *Oracle) FeeHistory(
	ctx context.Context, blocks int,
	unresolvedLastBlock rpc.BlockNumber,
	rewardPercentiles []float64,
) (*big.Int, [][]*big.Int, []*big.Int, []float64, error) {
	if blocks < 1 {
		return common.Big0, nil, nil, nil, nil // returning with no data and no error means there are no retrievable blocks
	}
	maxFeeHistory := oracle.maxHeaderHistory
	if len(rewardPercentiles) != 0 {
		maxFeeHistory = oracle.maxBlockHistory
	}
	if blocks > maxFeeHistory {
		logger.Warn("Sanitizing fee history length", "requested", blocks, "truncated", maxFeeHistory)
		blocks = maxFeeHistory
	}
	for i, p := range rewardPercentiles {
		if p < 0 || p > 100 {
			return common.Big0, nil, nil, nil, fmt.Errorf("%w: %f", errInvalidPercentile, p)
		}
		if i > 0 && p < rewardPercentiles[i-1] {
			return common.Big0, nil, nil, nil, fmt.Errorf("%w: #%d:%f > #%d:%f", errInvalidPercentile, i-1, rewardPercentiles[i-1], i, p)
		}
	}
	var err error
	lastBlock, blocks, err := oracle.resolveBlockRange(ctx, unresolvedLastBlock, blocks)
	if err != nil || blocks == 0 {
		return common.Big0, nil, nil, nil, err
	}
	oldestBlock := lastBlock + 1 - uint64(blocks)

	var (
		next    = oldestBlock
		results = make(chan *blockFees, blocks)
	)
	for i := 0; i < maxBlockFetchers && i < blocks; i++ {
		go func() {
			for {
				// Retrieve the next block number to fetch with this goroutine
				blockNumber := atomic.AddUint64(&next, 1) - 1
				if blockNumber > lastBlock {
					return
				}

				fees := &blockFees{blockNumber: blockNumber}
				if len(rewardPercentiles) != 0 {
					fees.block, fees.err = oracle.backend.BlockByNumber(ctx, rpc.BlockNumber(blockNumber))
					if fees.block != nil && fees.err == nil {
						fees.receipts = oracle.backend.GetBlockReceipts(ctx, fees.block.Hash())
						fees.header = fees.block.Header()
					}
				} else {
					fees.header, fees.err = oracle.backend.HeaderByNumber(ctx, rpc.BlockNumber(blockNumber))
				}
				if fees.header != nil && fees.err == nil {
					oracle.processBlock(fees, rewardPercentiles)
				}
				// send to results even if empty to guarantee that blocks items are sent in total
				results <- fees
			}
		}()
	}
	var (
		reward       = make([][]*big.Int, blocks)
		baseFee      = make([]*big.Int, blocks+1)
		gasUsedRatio = make([]float64, blocks)
		blockCount   = blocks
	)
	for ; blocks > 0; blocks-- {
		fees := <-results
		if fees.err != nil {
			return common.Big0, nil, nil, nil, fees.err
		}
		i := int(fees.blockNumber - oldestBlock)
		reward[i], baseFee[i], baseFee[i+1], gasUsedRatio[i] = fees.results.reward, fees.results.baseFee, fees.results.nextBaseFee, fees.results.gasUsedRatio
	}
	if len(rewardPercentiles) != 0 {
		reward = reward[:blockCount]
	} else {
		reward = nil
	}
	baseFee, gasUsedRatio = baseFee[:blockCount+1], gasUsedRatio[:blockCount]
	return new(big.Int).SetUint64(oldestBlock), reward, baseFee, gasUsedRatio, nil
}
