package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/testutil"
	multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
)

func (suite *KeeperTestSuite) TestSetBondWeight() {
	suite.SetupTest()

	gasDenom := "ario"
	govDenom := "arst"
	gasWeight := sdk.OneDec()
	govWeight := sdk.NewDecWithPrec(2, 4)

	suite.msKeeper.SetBondWeight(suite.ctx, gasDenom, gasWeight)
	suite.msKeeper.SetBondWeight(suite.ctx, govDenom, govWeight)

	suite.Equal(gasWeight, suite.msKeeper.GetBondWeight(suite.ctx, gasDenom))
	suite.Equal(govWeight, suite.msKeeper.GetBondWeight(suite.ctx, govDenom))
}

func (suite *KeeperTestSuite) TestSetValidatorMultiStakingCoin() {
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

func (suite *KeeperTestSuite) TestSetIntermediaryDelegator() {
	delA := testutil.GenAddress()
	delB := testutil.GenAddress()
	imAddrressA := testutil.GenAddress()
	imAddrressB := testutil.GenAddress()

	testCases := []struct {
		name     string
		malleate func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []sdk.AccAddress
		imAccs   []sdk.AccAddress
		expPanic bool
	}{
		{
			name: "1 delegator, 1 intermediary account, success",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []sdk.AccAddress {
				msKeeper.SetIntermediaryDelegator(ctx, imAddrressA, delA)
				return []sdk.AccAddress{delA}
			},
			imAccs:   []sdk.AccAddress{imAddrressA},
			expPanic: false,
		},
		{
			name: "2 delegator, 2 intermediary account, success",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []sdk.AccAddress {
				msKeeper.SetIntermediaryDelegator(ctx, imAddrressA, delA)
				msKeeper.SetIntermediaryDelegator(ctx, imAddrressB, delB)
				return []sdk.AccAddress{delA, delB}
			},
			imAccs:   []sdk.AccAddress{imAddrressA, imAddrressB},
			expPanic: false,
		},
		{
			name: "2 delegator, 2 intermediary account, failed",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []sdk.AccAddress {
				msKeeper.SetIntermediaryDelegator(ctx, imAddrressA, delA)
				msKeeper.SetIntermediaryDelegator(ctx, imAddrressA, delA)
				return []sdk.AccAddress{delA, delB}
			},
			imAccs:   []sdk.AccAddress{imAddrressA, imAddrressB},
			expPanic: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			if tc.expPanic {
				suite.Require().PanicsWithValue("intermediary delegator already set", func() {
					tc.malleate(suite.ctx, suite.msKeeper)
				})
			} else {
				inputs := tc.malleate(suite.ctx, suite.msKeeper)
				for idx, imAcc := range tc.imAccs {
					actualDel := suite.msKeeper.GetIntermediaryDelegator(suite.ctx, imAcc)
					suite.Require().Equal(inputs[idx], actualDel)
				}
			}
		})
	}
}
