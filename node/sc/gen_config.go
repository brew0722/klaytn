// Code generated by github.com/fjl/gencodec. DO NOT EDIT.

package sc

import (
	"time"
)

// MarshalTOML marshals as TOML.
func (s SCConfig) MarshalTOML() (interface{}, error) {
	type SCConfig struct {
		Name                    string `toml:"-"`
		EnabledMainBridge       bool
		EnabledSubBridge        bool
		DataDir                 string
		NetworkId               uint64
		SkipBcVersionCheck      bool `toml:"-"`
		DatabaseHandles         int  `toml:"-"`
		LevelDBCacheSize        int
		TrieCacheSize           int
		TrieTimeout             time.Duration
		TrieBlockInterval       uint
		ChildChainIndexing      bool
		MainBridgePort          string
		SubBridgePort           string
		MaxPeer                 int
		ServiceChainConsensus   string
		AnchoringPeriod         uint64
		SentChainTxsLimit       uint64
		ParentChainID           uint64
		VTRecovery              bool
		VTRecoveryInterval      uint64
		Anchoring               bool
    DefaultGasLimit         uint64
		KASAnchor               bool
		KASAnchorUrl            string
		KASAnchorPeriod         uint64
		KASAnchorOperator       string
		KASAccessKey            string
		KASSecretKey            string
		KASXChainId             string
		KASAnchorRequestTimeout time.Duration
	}
	var enc SCConfig
	enc.Name = s.Name
	enc.EnabledMainBridge = s.EnabledMainBridge
	enc.EnabledSubBridge = s.EnabledSubBridge
	enc.DataDir = s.DataDir
	enc.NetworkId = s.NetworkId
	enc.SkipBcVersionCheck = s.SkipBcVersionCheck
	enc.DatabaseHandles = s.DatabaseHandles
	enc.LevelDBCacheSize = s.LevelDBCacheSize
	enc.TrieCacheSize = s.TrieCacheSize
	enc.TrieTimeout = s.TrieTimeout
	enc.TrieBlockInterval = s.TrieBlockInterval
	enc.ChildChainIndexing = s.ChildChainIndexing
	enc.MainBridgePort = s.MainBridgePort
	enc.SubBridgePort = s.SubBridgePort
	enc.MaxPeer = s.MaxPeer
	enc.ServiceChainConsensus = s.ServiceChainConsensus
	enc.AnchoringPeriod = s.AnchoringPeriod
	enc.SentChainTxsLimit = s.SentChainTxsLimit
	enc.ParentChainID = s.ParentChainID
	enc.VTRecovery = s.VTRecovery
	enc.VTRecoveryInterval = s.VTRecoveryInterval
	enc.Anchoring = s.Anchoring
	enc.DefaultGasLimit = s.DefaultGasLimit
	enc.KASAnchor = s.KASAnchor
	enc.KASAnchorUrl = s.KASAnchorUrl
	enc.KASAnchorPeriod = s.KASAnchorPeriod
	enc.KASAnchorOperator = s.KASAnchorOperator
	enc.KASAccessKey = s.KASAccessKey
	enc.KASSecretKey = s.KASSecretKey
	enc.KASXChainId = s.KASXChainId
	enc.KASAnchorRequestTimeout = s.KASAnchorRequestTimeout
	return &enc, nil
}

// UnmarshalTOML unmarshals from TOML.
func (s *SCConfig) UnmarshalTOML(unmarshal func(interface{}) error) error {
	type SCConfig struct {
		Name                    *string `toml:"-"`
		EnabledMainBridge       *bool
		EnabledSubBridge        *bool
		DataDir                 *string
		NetworkId               *uint64
		SkipBcVersionCheck      *bool `toml:"-"`
		DatabaseHandles         *int  `toml:"-"`
		LevelDBCacheSize        *int
		TrieCacheSize           *int
		TrieTimeout             *time.Duration
		TrieBlockInterval       *uint
		ChildChainIndexing      *bool
		MainBridgePort          *string
		SubBridgePort           *string
		MaxPeer                 *int
		ServiceChainConsensus   *string
		AnchoringPeriod         *uint64
		SentChainTxsLimit       *uint64
		ParentChainID           *uint64
		VTRecovery              *bool
		VTRecoveryInterval      *uint64
		Anchoring               *bool
    DefaultGasLimit         *uint64
		KASAnchor               *bool
		KASAnchorUrl            *string
		KASAnchorPeriod         *uint64
		KASAnchorOperator       *string
		KASAccessKey            *string
		KASSecretKey            *string
		KASXChainId             *string
		KASAnchorRequestTimeout *time.Duration
	}
	var dec SCConfig
	if err := unmarshal(&dec); err != nil {
		return err
	}
	if dec.Name != nil {
		s.Name = *dec.Name
	}
	if dec.EnabledMainBridge != nil {
		s.EnabledMainBridge = *dec.EnabledMainBridge
	}
	if dec.EnabledSubBridge != nil {
		s.EnabledSubBridge = *dec.EnabledSubBridge
	}
	if dec.DataDir != nil {
		s.DataDir = *dec.DataDir
	}
	if dec.NetworkId != nil {
		s.NetworkId = *dec.NetworkId
	}
	if dec.SkipBcVersionCheck != nil {
		s.SkipBcVersionCheck = *dec.SkipBcVersionCheck
	}
	if dec.DatabaseHandles != nil {
		s.DatabaseHandles = *dec.DatabaseHandles
	}
	if dec.LevelDBCacheSize != nil {
		s.LevelDBCacheSize = *dec.LevelDBCacheSize
	}
	if dec.TrieCacheSize != nil {
		s.TrieCacheSize = *dec.TrieCacheSize
	}
	if dec.TrieTimeout != nil {
		s.TrieTimeout = *dec.TrieTimeout
	}
	if dec.TrieBlockInterval != nil {
		s.TrieBlockInterval = *dec.TrieBlockInterval
	}
	if dec.ChildChainIndexing != nil {
		s.ChildChainIndexing = *dec.ChildChainIndexing
	}
	if dec.MainBridgePort != nil {
		s.MainBridgePort = *dec.MainBridgePort
	}
	if dec.SubBridgePort != nil {
		s.SubBridgePort = *dec.SubBridgePort
	}
	if dec.MaxPeer != nil {
		s.MaxPeer = *dec.MaxPeer
	}
	if dec.ServiceChainConsensus != nil {
		s.ServiceChainConsensus = *dec.ServiceChainConsensus
	}
	if dec.AnchoringPeriod != nil {
		s.AnchoringPeriod = *dec.AnchoringPeriod
	}
	if dec.SentChainTxsLimit != nil {
		s.SentChainTxsLimit = *dec.SentChainTxsLimit
	}
	if dec.ParentChainID != nil {
		s.ParentChainID = *dec.ParentChainID
	}
	if dec.VTRecovery != nil {
		s.VTRecovery = *dec.VTRecovery
	}
	if dec.VTRecoveryInterval != nil {
		s.VTRecoveryInterval = *dec.VTRecoveryInterval
	}
	if dec.Anchoring != nil {
		s.Anchoring = *dec.Anchoring
	}
	if dec.DefaultGasLimit != nil {
		s.DefaultGasLimit = *dec.DefaultGasLimit
	}
	if dec.KASAnchor != nil {
		s.KASAnchor = *dec.KASAnchor
	}
	if dec.KASAnchorUrl != nil {
		s.KASAnchorUrl = *dec.KASAnchorUrl
	}
	if dec.KASAnchorPeriod != nil {
		s.KASAnchorPeriod = *dec.KASAnchorPeriod
	}
	if dec.KASAnchorOperator != nil {
		s.KASAnchorOperator = *dec.KASAnchorOperator
	}
	if dec.KASAccessKey != nil {
		s.KASAccessKey = *dec.KASAccessKey
	}
	if dec.KASSecretKey != nil {
		s.KASSecretKey = *dec.KASSecretKey
	}
	if dec.KASXChainId != nil {
		s.KASXChainId = *dec.KASXChainId
	}
	if dec.KASAnchorRequestTimeout != nil {
		s.KASAnchorRequestTimeout = *dec.KASAnchorRequestTimeout
	}
	return nil
}
