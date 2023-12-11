package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/testutil"
	multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (suite *KeeperTestSuite) TestSetBondTokenWeight() {
	suite.SetupTest()

	gasDenom := "ario"
	govDenom := "arst"
	gasWeight := sdk.NewDec(1)
	govWeight := sdk.MustNewDecFromStr("0.5")

	suite.msKeeper.SetBondTokenWeight(suite.ctx, gasDenom, gasWeight)
	suite.msKeeper.SetBondTokenWeight(suite.ctx, govDenom, govWeight)

	btw, _ := suite.msKeeper.GetBondTokenWeight(suite.ctx, gasDenom)
	suite.Equal(gasWeight, btw)

	btw, _ = suite.msKeeper.GetBondTokenWeight(suite.ctx, govDenom)
	suite.Equal(govWeight, btw)
}

func (suite *KeeperTestSuite) TestSetValidatorAllowedToken() {
	valA := testutil.GenValAddress()
	valB := testutil.GenValAddress()
	gasDenom := "ario"
	govDenom := "arst"
	testCases := []struct {
		name     string
		malleate func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []string
		vals     []sdk.ValAddress
		expPanic bool
	}{
		{
			name: "1 val, 1 denom, success",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []string {
				msKeeper.SetValidatorAllowedToken(ctx, valA, gasDenom)
				return []string{gasDenom}
			},
			vals:     []sdk.ValAddress{valA},
			expPanic: false,
		},
		{
			name: "2 val, 2 denom, success",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []string {
				msKeeper.SetValidatorAllowedToken(ctx, valA, gasDenom)
				msKeeper.SetValidatorAllowedToken(ctx, valB, govDenom)
				return []string{gasDenom, govDenom}
			},
			vals:     []sdk.ValAddress{valA, valB},
			expPanic: false,
		},
		{
			name: "1 val, 2 denom, failed",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []string {
				msKeeper.SetValidatorAllowedToken(ctx, valA, gasDenom)
				msKeeper.SetValidatorAllowedToken(ctx, valA, govDenom)
				return []string{gasDenom, govDenom}
			},
			vals:     []sdk.ValAddress{valA, valB},
			expPanic: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			if tc.expPanic {
				suite.Require().PanicsWithValue("validator denom already set", func() {
					tc.malleate(suite.ctx, suite.msKeeper)
				})
			} else {
				inputs := tc.malleate(suite.ctx, suite.msKeeper)
				for idx, val := range tc.vals {
					actualDenom := suite.msKeeper.GetValidatorAllowedToken(suite.ctx, val)
					suite.Require().Equal(inputs[idx], actualDenom)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestSetMultiStakingLock() {
	suite.SetupTest()
	delAddr := testutil.GenAddress()
	valAddr := testutil.GenValAddress()

	lockAmt := sdk.NewInt(3001)
	weight := sdk.MustNewDecFromStr("0.3")
	lockRecord := types.NewMultiStakingLock(lockAmt, weight, delAddr, valAddr)
	suite.msKeeper.SetMultiStakingLock(suite.ctx, types.MultiStakingLockID(delAddr, valAddr), lockRecord)

	actualLockRecord, found := suite.msKeeper.GetMultiStakingLock(suite.ctx, types.MultiStakingLockID(delAddr, valAddr))
	suite.Require().True(found)
	suite.Equal(actualLockRecord.ConversionRatio, weight)
	suite.Equal(actualLockRecord.LockedAmount, lockAmt)
	suite.Equal(actualLockRecord.DelAddr, delAddr.String())
	suite.Equal(actualLockRecord.ValAddr, valAddr.String())
}

func (suite *KeeperTestSuite) TestRemoveMultiStakingLock() {
	suite.SetupTest()
	delAddr := testutil.GenAddress()
	valAddr := testutil.GenValAddress()

	lockAmt := sdk.NewInt(3001)
	weight := sdk.MustNewDecFromStr("0.3")
	lockRecord := types.NewMultiStakingLock(lockAmt, weight, delAddr, valAddr)
	suite.msKeeper.SetMultiStakingLock(suite.ctx, types.MultiStakingLockID(delAddr, valAddr), lockRecord)

	actualLockRecord, found := suite.msKeeper.GetMultiStakingLock(suite.ctx, types.MultiStakingLockID(delAddr, valAddr))
	suite.Require().True(found)
	suite.Equal(actualLockRecord.ConversionRatio, weight)
	suite.Equal(actualLockRecord.LockedAmount, lockAmt)
	suite.Equal(actualLockRecord.DelAddr, delAddr.String())
	suite.Equal(actualLockRecord.ValAddr, valAddr.String())

	suite.msKeeper.RemoveMultiStakingLock(suite.ctx, delAddr, valAddr)

	_, found = suite.msKeeper.GetMultiStakingLock(suite.ctx, types.MultiStakingLockID(delAddr, valAddr))
	suite.Require().False(found)
}

func (suite *KeeperTestSuite) TestSetUnbondMultiStaking() {
	suite.SetupTest()
	delAddr := testutil.GenAddress()
	valAddr := testutil.GenValAddress()

	unbondAmt := sdk.NewInt(3001)
	weight := sdk.MustNewDecFromStr("0.3")
	unbondRecord := types.NewMultiStakingUnlock(delAddr, valAddr, 1, weight, unbondAmt)
	suite.msKeeper.SetMultiStakingUnlock(suite.ctx, unbondRecord)

	actualUnbondRecord, found := suite.msKeeper.GetMultiStakingUnlock(suite.ctx, delAddr, valAddr)
	suite.Require().True(found)
	suite.Equal(actualUnbondRecord.Entries[0].ConversionRatio, weight)
	suite.Equal(actualUnbondRecord.Entries[0].Balance, unbondAmt)
	suite.Equal(actualUnbondRecord.DelegatorAddress, delAddr.String())
	suite.Equal(actualUnbondRecord.ValidatorAddress, valAddr.String())
}

func (suite *KeeperTestSuite) TestRemoveUnbondMultiStaking() {
	suite.SetupTest()
	delAddr := testutil.GenAddress()
	valAddr := testutil.GenValAddress()

	unbondAmt := sdk.NewInt(3001)
	weight := sdk.MustNewDecFromStr("0.3")
	unbondRecord := types.NewMultiStakingUnlock(delAddr, valAddr, 1, weight, unbondAmt)
	suite.msKeeper.SetMultiStakingUnlock(suite.ctx, unbondRecord)

	actualUnbondRecord, found := suite.msKeeper.GetMultiStakingUnlock(suite.ctx, delAddr, valAddr)
	suite.Require().True(found)
	suite.Equal(actualUnbondRecord.Entries[0].ConversionRatio, weight)
	suite.Equal(actualUnbondRecord.Entries[0].Balance, unbondAmt)
	suite.Equal(actualUnbondRecord.DelegatorAddress, delAddr.String())
	suite.Equal(actualUnbondRecord.ValidatorAddress, valAddr.String())

	suite.msKeeper.RemoveMultiStakingUnlock(suite.ctx, actualUnbondRecord)

	_, found = suite.msKeeper.GetMultiStakingUnlock(suite.ctx, delAddr, valAddr)
	suite.Require().False(found)
}

func (suite *KeeperTestSuite) TestSetUnbondMultiStakingEntry() {
	suite.SetupTest()
	delAddr := testutil.GenAddress()
	valAddr := testutil.GenValAddress()
	minTime := time.Now()
	unbondAmt := sdk.NewInt(3001)
	weight := sdk.MustNewDecFromStr("0.3")
	suite.msKeeper.SetMultiStakingUnlockEntry(suite.ctx, delAddr, valAddr, 1, weight, minTime, unbondAmt)

	actualUnbondRecord, found := suite.msKeeper.GetMultiStakingUnlock(suite.ctx, delAddr, valAddr)
	suite.Require().True(found)
	suite.Equal(actualUnbondRecord.Entries[0].ConversionRatio, weight)
	suite.Equal(actualUnbondRecord.Entries[0].Balance, unbondAmt)
	suite.Equal(actualUnbondRecord.DelegatorAddress, delAddr.String())
	suite.Equal(actualUnbondRecord.ValidatorAddress, valAddr.String())

	suite.msKeeper.SetMultiStakingUnlockEntry(suite.ctx, delAddr, valAddr, 2, weight, minTime, unbondAmt)
	actualUnbondRecordAfter, found := suite.msKeeper.GetMultiStakingUnlock(suite.ctx, delAddr, valAddr)
	suite.Require().True(found)
	suite.Equal(actualUnbondRecordAfter.Entries[1].ConversionRatio, weight)
	suite.Equal(actualUnbondRecordAfter.Entries[1].Balance, unbondAmt)

	suite.msKeeper.SetMultiStakingUnlockEntry(suite.ctx, delAddr, valAddr, 1, weight, minTime, unbondAmt)
	actualUnbondRecordAfter1, found := suite.msKeeper.GetMultiStakingUnlock(suite.ctx, delAddr, valAddr)
	suite.Require().True(found)
	suite.Equal(actualUnbondRecordAfter1.Entries[0].ConversionRatio, weight)
	suite.Equal(actualUnbondRecordAfter1.Entries[0].Balance, unbondAmt.Add(unbondAmt))
}
