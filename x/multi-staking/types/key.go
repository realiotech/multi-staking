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

	// key prefix
	MultiStakingLockPrefix = KeyPrefix("multi-staking-lock")

	UnbondDelKey    = []byte{0x11} // key for an unbonding-delegation
	UnbondQueueKey         = []byte{0x16} // prefix for the timestamps in unbonding queue
)

func KeyPrefix(key string) []byte {
	return []byte(key)
}

// GetBondTokenWeightKeyKey returns a key for an index containing the bond token weight
func GetBondTokenWeightKey(tokenDenom string) []byte {
	return append(BondTokenWeightKey, []byte(tokenDenom)...)
}

// GetValidatorAllowedTokenKey returns a key for an index containing the bond denom of a validator
func GetValidatorAllowedTokenKey(valAddr string) []byte {
	return append(ValidatorAllowedTokenKey, []byte(valAddr)...)
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

// GetUBDsKey creates the prefix for all unbonding delegations from a delegator
func GetUBDsKey(delAddr sdk.AccAddress) []byte {
	return append(UnbondDelKey, address.MustLengthPrefix(delAddr)...)
}

// GetUBDKey creates the key for an unbonding delegation by delegator and validator addr
// VALUE: staking/UnbondingDelegation
func GetUBDKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetUBDsKey(delAddr.Bytes()), address.MustLengthPrefix(valAddr)...)
}