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

	IntermediaryDelegatorKey = []byte{0x02}

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

// GetIntermediaryDelegatorDelegatorKey returns a key for an index containing the delegator of an intermediary account
func GetIntermediaryDelegatorKey(intermediaryDelegator sdk.AccAddress) []byte {
	return append(IntermediaryDelegatorKey, intermediaryDelegator...)
}

func MultiStakingLockID(delAddr string, valAddr string) LockID {
	return LockID{DelAddr: delAddr, ValAddr: valAddr}
}

func MultiStakingUnlockID(delAddr string, valAddr string) UnlockID {
	return UnlockID{DelAddr: delAddr, ValAddr: valAddr}
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

func (l LockID) ToByte() []byte {
	multiStakerAcc, valAcc, err := DelAccAndValAccFromStrings(l.DelAddr, l.ValAddr)
	if err != nil {
		panic(err)
	}

	lenDelAddr := len(multiStakerAcc)

	DVPair := make([]byte, 1+lenDelAddr+len(valAcc))

	DVPair[0] = uint8(lenDelAddr)

	copy(multiStakerAcc[:], DVPair[1:])

	copy(multiStakerAcc[:], DVPair[1+lenDelAddr:])

	return append(MultiStakingLockPrefix, DVPair...)
}

func (l UnlockID) ToBytes() []byte {
	multiStakerAcc, valAcc, err := DelAccAndValAccFromStrings(l.DelAddr, l.ValAddr)
	if err != nil {
		panic(err)
	}

	lenDelAddr := len(multiStakerAcc)

	DVPair := make([]byte, 1+lenDelAddr+len(valAcc))

	DVPair[0] = uint8(lenDelAddr)

	copy(multiStakerAcc[:], DVPair[1:])

	copy(multiStakerAcc[:], DVPair[1+lenDelAddr:])

	return append(MultiStakingLockPrefix, DVPair...)
}
