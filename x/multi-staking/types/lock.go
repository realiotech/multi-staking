package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DelAccAndValAccFromLockID(lockID []byte) (delAcc []byte, valAcc []byte) {
	return
}

func NewMultiStakingLock(lockID *LockID, lockedCoin MultiStakingCoin) MultiStakingLock {
	return MultiStakingLock{
		LockID:     lockID,
		LockedCoin: lockedCoin,
	}
}

func (lock MultiStakingLock) ToMultiStakingCoin(coin sdk.Coin) MultiStakingCoin {
	return lock.LockedCoin.WithAmount(coin.Amount)
}

func (lock MultiStakingLock) RemoveCoinFromMultiStakingLock(removedCoin MultiStakingCoin) (MultiStakingLock, error) {
	lockedCoinAfter, err := lock.LockedCoin.SafeSub(removedCoin)
	lock.LockedCoin = lockedCoinAfter
	return lock, err
}

func (lock MultiStakingLock) IsEmpty() bool {
	return lock.LockedCoin.Amount.IsZero()
}

func (multiStakingLock MultiStakingLock) AddCoinToMultiStakingLock(addedCoin MultiStakingCoin) (MultiStakingLock, error) {
	lockedCoinAfter, err := multiStakingLock.LockedCoin.SafeAdd(addedCoin)
	multiStakingLock.LockedCoin = lockedCoinAfter
	return multiStakingLock, err
}

func (m MultiStakingLock) GetBondWeight() sdk.Dec {
	return m.LockedCoin.BondWeight
}

func (multiStakingLock MultiStakingLock) LockedAmountToBondAmount(lockedAmount math.Int) sdk.Int {
	return multiStakingLock.GetBondWeight().MulInt(lockedAmount).RoundInt()
}
