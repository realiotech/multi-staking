package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	BondCoinWeightKey       = []byte{0x00}
	ValidatorAllowedCoinKey = []byte{0x01}

	IntermediaryAccountKey = []byte{0x02}

	MultiStakingLockPrefix = []byte{0x03}

	MultiStakingUnlockPrefix = []byte{0x11} // key for an unbonding-delegation
)

func KeyPrefix(key string) []byte {
	return []byte(key)
}

// GetBondCoinWeightKeyKey returns a key for an index containing the bond coin weight
func GetBondCoinWeightKey(tokenDenom string) []byte {
	return append(BondCoinWeightKey, []byte(tokenDenom)...)
}

// GetValidatorAllowedCoinKey returns a key for an index containing the bond denom of a validator
func GetValidatorAllowedCoinKey(valAddr string) []byte {
	return append(ValidatorAllowedCoinKey, []byte(valAddr)...)
}

// GetIntermediaryAccountDelegatorKey returns a key for an index containing the delegator of an intermediary account
func GetIntermediaryAccountKey(intermediaryAccount sdk.AccAddress) []byte {
	return append(IntermediaryAccountKey, intermediaryAccount...)
}

func MultiStakingLockID(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	lenDelAddr := len(delAddr)

	DVPair := make([]byte, 1+lenDelAddr+len(valAddr))

	DVPair[0] = uint8(lenDelAddr)

	copy(delAddr[:], DVPair[1:])

	copy(valAddr[:], DVPair[1+lenDelAddr:])

	return append(MultiStakingLockPrefix, DVPair...)
}

func MultiStakingUnlockID(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	lenDelAddr := len(delAddr)

	DVPair := make([]byte, 1+lenDelAddr+len(valAddr))

	DVPair[0] = uint8(lenDelAddr)

	copy(delAddr[:], DVPair[1:])

	copy(valAddr[:], DVPair[1+lenDelAddr:])
	return append(MultiStakingUnlockPrefix, DVPair...)
}

func DelAddrAndValAddrFromLockID(lockID []byte) (delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	lenDelAddr := lockID[0]

	delAddr = lockID[1 : lenDelAddr+1]

	valAddr = lockID[1+lenDelAddr:]

	return delAddr, valAddr
}

func DelAddrAndValAddrFromUnlockID(unlockID []byte) (delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	lenDelAddr := unlockID[0]

	delAddr = unlockID[1 : lenDelAddr+1]

	valAddr = unlockID[1+lenDelAddr:]

	return delAddr, valAddr
}

// // GetUBDKey creates the key for an unbonding delegation by delegator and validator addr
// // VALUE: multi-staking/MultiStakingUnlock
// func GetUBDKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
// 	return append(GetUBDsKey(delAddr.Bytes()), address.MustLengthPrefix(valAddr)...)
// }
