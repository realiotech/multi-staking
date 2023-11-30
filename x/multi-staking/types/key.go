package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "multi-staking"

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
	BondTokenWeightKey       = []byte{0x00}
	ValidatorAllowedTokenKey = []byte{0x01}

	IntermediaryAccountKey = []byte{0x02}

	DVPairSDKBondAmount = []byte{0x03}

	DVPairBondAmount = []byte{0x04}

	// mem store key
	CompletedDelegationsPrefix = []byte{0x05}
)

// GetBondTokenWeightKeyKey returns a key for an index containing the bond token weight
func GetBondTokenWeightKey(tokenDenom string) []byte {
	return append(BondTokenWeightKey, []byte(tokenDenom)...)
}

// GetValidatorAllowedTokenKey returns a key for an index containing the bond denom of a validator
func GetValidatorAllowedTokenKey(operatorAddr sdk.ValAddress) []byte {
	return append(ValidatorAllowedTokenKey, address.MustLengthPrefix(operatorAddr)...)
}

// GetIntermediaryAccountDelegatorKey returns a key for an index containing the delegator of an intermediary account
func GetIntermediaryAccountKey(intermediaryAccount sdk.AccAddress) []byte {
	return append(IntermediaryAccountKey, intermediaryAccount...)
}

func GetDVPairSDKBondAmountKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	DVPair := append(delAddr, address.MustLengthPrefix(valAddr)...)
	return append(DVPairSDKBondAmount, DVPair...)
}

func GetDVPairBondAmountKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	DVPair := append(delAddr, address.MustLengthPrefix(valAddr)...)
	return append(DVPairBondAmount, DVPair...)
}

func MultiStakingLockID(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	DVPair := append(delAddr, address.MustLengthPrefix(valAddr)...)
	return append(DVPairBondAmount, DVPair...)
}
