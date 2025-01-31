// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
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
// This file is derived from core/types/derive_sha.go (2018/06/04).
// Modified and improved for the klaytn development.

package derivesha

import (
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
)

type IDeriveSha interface {
	DeriveSha(list types.DerivableList) common.Hash
}

var (
	config     *params.ChainConfig
	instances  map[int]IDeriveSha
	emptyRoots map[int]common.Hash

	logger = log.NewModuleLogger(log.Blockchain)
)

func init() {
	instances = map[int]IDeriveSha{
		types.ImplDeriveShaOriginal: DeriveShaOrig{},
		types.ImplDeriveShaSimple:   DeriveShaSimple{},
		types.ImplDeriveShaConcat:   DeriveShaConcat{},
	}

	emptyRoots = make(map[int]common.Hash)
	for implType, instance := range instances {
		emptyRoots[implType] = instance.DeriveSha(types.Transactions{})
	}
}

func InitDeriveSha(chainConfig *params.ChainConfig) {
	config = chainConfig
	types.DeriveSha = DeriveShaMux
	types.EmptyRootHash = EmptyRootHashMux
}

func DeriveShaMux(list types.DerivableList, num *big.Int) common.Hash {
	return instances[getType()].DeriveSha(list)
}

func EmptyRootHashMux(num *big.Int) common.Hash {
	return emptyRoots[getType()]
}

// TODO: Choose appropriate DeriveShaImpl from governance based on block number
func getType() int {
	implType := config.DeriveShaImpl
	if _, ok := instances[implType]; ok {
		return implType
	} else {
		logger.Error("Unrecognized deriveShaImpl, falling back to Orig", "impl", implType)
		return types.ImplDeriveShaOriginal
	}
}
