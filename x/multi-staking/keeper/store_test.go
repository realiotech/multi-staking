package keeper_test

import (
	"github.com/realio-tech/multi-staking-module/testutil"
	multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	gasDenom = "ario"
	govDenom = "arst"
)

func (suite *KeeperTestSuite) TestSetBondWeight() {
	suite.SetupTest()

	gasWeight := sdk.OneDec()
	govWeight := sdk.NewDecWithPrec(2, 4)

	suite.msKeeper.SetBondWeight(suite.ctx, gasDenom, gasWeight)
	suite.msKeeper.SetBondWeight(suite.ctx, govDenom, govWeight)

	expectedGasWeight, _ := suite.msKeeper.GetBondWeight(suite.ctx, gasDenom)
	expectedGovWeight, _ := suite.msKeeper.GetBondWeight(suite.ctx, govDenom)

	suite.Equal(gasWeight, expectedGasWeight)
	suite.Equal(govWeight, expectedGovWeight)
}

func (suite *KeeperTestSuite) TestSetValidatorMultiStakingCoin() {
	valA := testutil.GenValAddress()
	valB := testutil.GenValAddress()

	testCases := []struct {
		name     string
		malleate func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []string
		vals     []sdk.ValAddress
		expPanic bool
	}{
		{
			name: "1 val, 1 denom, success",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []string {
				msKeeper.SetValidatorMultiStakingCoin(ctx, valA, gasDenom)
				return []string{gasDenom}
			},
			vals:     []sdk.ValAddress{valA},
			expPanic: false,
		},
		{
			name: "2 val, 2 denom, success",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []string {
				msKeeper.SetValidatorMultiStakingCoin(ctx, valA, gasDenom)
				msKeeper.SetValidatorMultiStakingCoin(ctx, valB, govDenom)
				return []string{gasDenom, govDenom}
			},
			vals:     []sdk.ValAddress{valA, valB},
			expPanic: false,
		},
		{
			name: "1 val, 2 denom, failed",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []string {
				msKeeper.SetValidatorMultiStakingCoin(ctx, valA, gasDenom)
				msKeeper.SetValidatorMultiStakingCoin(ctx, valA, govDenom)
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
				suite.Require().PanicsWithValue("validator multi staking coin already set", func() {
					tc.malleate(suite.ctx, suite.msKeeper)
				})
			} else {
				inputs := tc.malleate(suite.ctx, suite.msKeeper)
				for idx, val := range tc.vals {
					actualDenom := suite.msKeeper.GetValidatorMultiStakingCoin(suite.ctx, val)
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

	lock := types.MultiStakingLock{
		LockID: types.LockID{
			MultiStakerAddr: delAddr.String(),
			ValAddr:         valAddr.String(),
		},
		LockedCoin: types.MultiStakingCoin{
			Denom:      gasDenom,
			Amount:     sdk.NewIntFromUint64(1000000),
			BondWeight: sdk.NewDec(1),
		},
	}

	testCases := []struct {
		name     string
		malleate func()
		expError bool
	}{
		{
			"Success",
			func() {
				suite.msKeeper.SetMultiStakingLock(suite.ctx, lock)
			},
			false,
		},
	}
	for _, tc := range testCases {
		if !tc.expError {
			tc.malleate()
			msLock, found := suite.msKeeper.GetMultiStakingLock(suite.ctx, lock.LockID)
			suite.Require().True(found)
			suite.Require().Equal(lock, msLock)
		}
	}
}
