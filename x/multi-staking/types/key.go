package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "multistaking"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "memory:capability"
)

// KVStore keys
var (
	BondWeightKey = []byte{0x00}

	ValidatorMultiStakingCoinKey = []byte{0x01}

	IntermediaryDelegatorKey = []byte{0x02}

	MultiStakingLockPrefix = []byte{0x03}

	// DVPairSDKBondTokens = []byte{0x03}

	// DVPairBondTokens = []byte{0x04}

	// // mem store key
	// CompletedDelegationsPrefix = []byte{0x05}
)

// GetBondWeightKey returns a key for an index containing the bond token weight
func GetBondWeightKey(tokenDenom string) []byte {
	return append(BondWeightKey, []byte(tokenDenom)...)
}

// GetValidatorMultiStakingCoinKey returns a key for an index containing the bond denom of a validator
func GetValidatorMultiStakingCoinKey(operatorAddr sdk.ValAddress) []byte {
	return append(ValidatorMultiStakingCoinKey, address.MustLengthPrefix(operatorAddr)...)
}

// GetIntermediaryDelegatorKey returns a key for an index containing the delegator of an intermediary account
func GetIntermediaryDelegatorKey(intermediaryAccount sdk.AccAddress) []byte {
	return append(IntermediaryDelegatorKey, intermediaryAccount...)
}
