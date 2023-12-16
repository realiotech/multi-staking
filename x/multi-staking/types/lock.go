package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DelAccAndValAccFromLockID(lockID []byte) (delAcc []byte, valAcc []byte) {
	return
}

func NewMultiStakingLock(lockedAmount math.Int, conversionRatio sdk.Dec, delAddr sdk.AccAddress, valAddr sdk.ValAddress) MultiStakingLock {
	return MultiStakingLock{
		LockedAmount:    lockedAmount,
		ConversionRatio: conversionRatio,
		DelAddr:         delAddr.String(),
		ValAddr:         valAddr.String(),
	}
}

func (lock MultiStakingLock) RemoveTokenFromMultiStakingLock(removedAmount math.Int) (MultiStakingLock, error) {
	if removedAmount.GT(lock.LockedAmount) {
		return MultiStakingLock{}, fmt.Errorf("removed amount greater than existing amount in lock")
	}

	lock.LockedAmount = lock.LockedAmount.Sub(removedAmount)

	return lock, nil
}

func (lock MultiStakingLock) IsEmpty() bool {
	return lock.LockedAmount.IsZero()
}

func (multiStakingLock MultiStakingLock) AddTokenToMultiStakingLock(addedAmount math.Int, currentConversionRatio sdk.Dec) MultiStakingLock {
	lockedAmountBefore := multiStakingLock.LockedAmount
	conversionRatioBefore := multiStakingLock.ConversionRatio

	lockedAmountAfter := lockedAmountBefore.Add(addedAmount)
	// conversionRatioAfter = ( (conversionRatioBefore * lockedAmountBefore) + (currentConversionRatio * addedAmount) ) / lockedAmountAfter
	conversionRatioAfter := ((conversionRatioBefore.MulInt(lockedAmountBefore)).Add(currentConversionRatio.MulInt(addedAmount))).QuoInt(lockedAmountAfter)

	multiStakingLock.LockedAmount = lockedAmountAfter
	multiStakingLock.ConversionRatio = conversionRatioAfter
	return multiStakingLock
}

func (multiStakingLock MultiStakingLock) LockedAmountToBondAmount(lockedAmount math.Int) sdk.Dec {
	return multiStakingLock.ConversionRatio.MulInt(lockedAmount)
}
