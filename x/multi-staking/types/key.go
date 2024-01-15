package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

	MultiStakingLockPrefix = []byte{0x02}

	MultiStakingUnlockPrefix = []byte{0x11} // key for an unbonding-delegation
)

func KeyPrefix(key string) []byte {
	return []byte(key)
}

// GetBondWeightKeyKey returns a key for an index containing the bond coin weight
func GetBondWeightKey(tokenDenom string) []byte {
	return append(BondWeightKey, []byte(tokenDenom)...)
}

// GetValidatorMultiStakingCoinKey returns a key for an index containing the bond denom of a validator
func GetValidatorMultiStakingCoinKey(valAddr sdk.ValAddress) []byte {
	return append(ValidatorMultiStakingCoinKey, []byte(valAddr)...)
}

func MultiStakingLockID(multiStakerAddr string, valAddr string) LockID {
	return LockID{MultiStakerAddr: multiStakerAddr, ValAddr: valAddr}
}

func MultiStakingUnlockID(multiStakerAddr string, valAddr string) UnlockID {
	return UnlockID{MultiStakerAddr: multiStakerAddr, ValAddr: valAddr}
}

func DelAddrAndValAddrFromLockID(lockIDByte []byte) (multiStakerAddr sdk.AccAddress, valAddr sdk.ValAddress, err error) {
	var newLockID LockID

	err = json.Unmarshal(lockIDByte, &newLockID)
	if err != nil {
		return nil, nil, err
	}

	return AccAddrAndValAddrFromStrings(newLockID.MultiStakerAddr, newLockID.ValAddr)
}

func DelAddrAndValAddrFromUnlockID(unlockIDByte []byte) (multiStakerAddr sdk.AccAddress, valAddr sdk.ValAddress, err error) {
	var newUnlockID UnlockID

	err = json.Unmarshal(unlockIDByte, &newUnlockID)
	if err != nil {
		return nil, nil, err
	}

	return AccAddrAndValAddrFromStrings(newUnlockID.MultiStakerAddr, newUnlockID.ValAddr)
}

// // GetUBDKey creates the key for an unbonding delegation by delegator and validator addr
// // VALUE: multi-staking/MultiStakingUnlock
// func GetUBDKey(multiStakerAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
// 	return append(GetUBDsKey(delAddr.Bytes()), address.MustLengthPrefix(valAddr)...)
// }

func (l LockID) ToByte() []byte {
	bz, err := json.Marshal(l)

	if err != nil {
		panic("can not Marshal")
	}
	return bz
}

func (l UnlockID) ToByte() []byte {
	bz, err := json.Marshal(l)

	if err != nil {
		panic("can not Marshal")
	}
	return bz
}
