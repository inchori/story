package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name.
	ModuleName = "evmengine"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key.
	MemStoreKey = "mem_evmengine"
)

// KVStore key prefixes.
var (
	ParamsKey = collections.NewPrefix(0)
)
