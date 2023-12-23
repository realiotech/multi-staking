package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DelAccAndValAccFromLockID(lockID []byte) (delAcc []byte, valAcc []byte) {
	return
}

func NewMultiStakingLock(lockID *LockID, lockedCoin WeightedCoin) MultiStakingLock {
	return MultiStakingLock{
		LockID:     lockID,
		LockedCoin: lockedCoin,
	}
}

func (lock MultiStakingLock) ToWeightedCoin(coin sdk.Coin) WeightedCoin {
	return lock.LockedCoin.WithAmount(coin.Amount)
}

func (lock MultiStakingLock) RemoveCoinFromMultiStakingLock(removedCoin sdk.Coin) (MultiStakingLock, error) {
	lockedCoinAfter, err := lock.LockedCoin.SafeSubCoin(removedCoin)
	lock.LockedCoin = lockedCoinAfter
	return lock, err
}

func (lock MultiStakingLock) IsEmpty() bool {
	return lock.LockedCoin.Amount.IsZero()
}

func (multiStakingLock MultiStakingLock) AddWeightedCoinToMultiStakingLock(addedCoin WeightedCoin) (MultiStakingLock, error) {
	lockedCoinAfter, err := multiStakingLock.LockedCoin.SafeAdd(addedCoin)
	multiStakingLock.LockedCoin = lockedCoinAfter
	return multiStakingLock, err
}

// func (m MultiStakingLock) To

func (multiStakingLock MultiStakingLock) LockedAmountToBondAmount(lockedAmount math.Int) sdk.Dec {
	bondWeight := multiStakingLock.LockedCoin.Weight

	return bondWeight.MulInt(lockedAmount)
}
