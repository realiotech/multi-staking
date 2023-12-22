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

func (lock MultiStakingLock) RemoveCoinFromMultiStakingLock(removedCoin WeightedCoin) (MultiStakingLock, error) {
	lockedCoinAfter, err := lock.LockedCoin.SafeSub(removedCoin)
	lock.LockedCoin = lockedCoinAfter
	return lock, err
}

func (lock MultiStakingLock) IsEmpty() bool {
	return lock.LockedCoin.Amount.IsZero()
}

func (multiStakingLock MultiStakingLock) AddCoinToMultiStakingLock(addedCoin WeightedCoin) (MultiStakingLock, error) {
	lockedCoinAfter, err := multiStakingLock.LockedCoin.SafeAdd(addedCoin)
	multiStakingLock.LockedCoin = lockedCoinAfter
	return multiStakingLock, err
}

func (multiStakingLock MultiStakingLock) LockedAmountToBondAmount(lockedAmount math.Int) sdk.Dec {
	conversionRatio := multiStakingLock.LockedCoin.Weight

	return conversionRatio.MulInt(lockedAmount)
}
