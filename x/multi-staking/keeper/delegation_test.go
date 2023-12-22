package keeper_test

// import (
// "cosmossdk.io/math"
// sdk "github.com/cosmos/cosmos-sdk/types"
// "github.com/realio-tech/multi-staking-module/testutil"
// multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
// )

func (suite *KeeperTestSuite) TestUpdateDVPairBondAmount() {
	// 	delA := testutil.GenAddress()
	// 	delB := testutil.GenAddress()
	// 	valA := testutil.GenValAddress()
	// 	valB := testutil.GenValAddress()

	// 	bondAmountA := sdk.NewInt(100)
	// 	bondAmountB := sdk.NewInt(200)

	// 	testCases := []struct {
	// 		name     string
	// 		malleate func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []math.Int
	// 		dels     []sdk.AccAddress
	// 		vals     []sdk.ValAddress
	// 	}{
	// 		{
	// 			name: "1 delegator, 1 validator, new",
	// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []math.Int {
	// 				msKeeper.UpdateDVPairBondAmount(ctx, delA, valA, bondAmountA)

	// 				return []math.Int{bondAmountA}
	// 			},
	// 			dels: []sdk.AccAddress{delA},
	// 			vals: []sdk.ValAddress{valA},
	// 		},
	// 		{
	// 			name: "2 delegator, 2 validator, success",
	// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []math.Int {
	// 				msKeeper.UpdateDVPairBondAmount(ctx, delA, valA, bondAmountA)
	// 				msKeeper.UpdateDVPairBondAmount(ctx, delA, valA, bondAmountA)
	// 				msKeeper.UpdateDVPairBondAmount(ctx, delB, valB, bondAmountB)
	// 				return []math.Int{bondAmountA.Add(bondAmountA), bondAmountB}
	// 			},
	// 			dels: []sdk.AccAddress{delA, delB},
	// 			vals: []sdk.ValAddress{valA, valB},
	// 		},
	// 		{
	// 			name: "1 delegator, 2 validator, success",
	// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []math.Int {
	// 				msKeeper.UpdateDVPairBondAmount(ctx, delA, valA, bondAmountA)
	// 				msKeeper.UpdateDVPairBondAmount(ctx, delA, valB, bondAmountB)
	// 				msKeeper.UpdateDVPairBondAmount(ctx, delA, valB, bondAmountA)
	// 				return []math.Int{bondAmountA, bondAmountB.Add(bondAmountA)}
	// 			},
	// 			dels: []sdk.AccAddress{delA, delA},
	// 			vals: []sdk.ValAddress{valA, valB},
	// 		},
	// 		{
	// 			name: "1 delegator, 1 validator, 2 bond amounts",
	// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []math.Int {
	// 				msKeeper.UpdateDVPairBondAmount(ctx, delA, valA, bondAmountA)
	// 				msKeeper.UpdateDVPairBondAmount(ctx, delA, valA, bondAmountB)
	// 				return []math.Int{bondAmountB.Add(bondAmountA)}
	// 			},
	// 			dels: []sdk.AccAddress{delA},
	// 			vals: []sdk.ValAddress{valA},
	// 		},
	// 	}

	// 	for _, tc := range testCases {
	// 		tc := tc
	// 		suite.Run(tc.name, func() {
	// 			suite.SetupTest()

	// 			inputs := tc.malleate(suite.ctx, suite.msKeeper)
	// 			for idx, expOut := range inputs {
	// 				actualCoin := suite.msKeeper.GetDVPairBondAmount(suite.ctx, tc.dels[idx], tc.vals[idx])
	// 				suite.Require().Equal(expOut, actualCoin)
	// 			}
	// 		})
	// 	}
	// }

	// func (suite *KeeperTestSuite) TestUpdateDVPairSDKBondAmount() {
	// 	delA := testutil.GenAddress()
	// 	delB := testutil.GenAddress()
	// 	valA := testutil.GenValAddress()
	// 	valB := testutil.GenValAddress()

	// 	bondSDKAmountA := sdk.NewInt(100)
	// 	bondSDKAmountB := sdk.NewInt(200)

	// 	testCases := []struct {
	// 		name     string
	// 		malleate func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []math.Int
	// 		dels     []sdk.AccAddress
	// 		vals     []sdk.ValAddress
	// 	}{
	// 		{
	// 			name: "1 delegator, 1 validator, new",
	// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []math.Int {
	// 				msKeeper.UpdateDVPairSDKBondAmount(ctx, delA, valA, bondSDKAmountA)

	// 				return []math.Int{bondSDKAmountA}
	// 			},
	// 			dels: []sdk.AccAddress{delA},
	// 			vals: []sdk.ValAddress{valA},
	// 		},
	// 		{
	// 			name: "2 delegator, 2 validator, success",
	// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []math.Int {
	// 				msKeeper.UpdateDVPairSDKBondAmount(ctx, delA, valA, bondSDKAmountA)
	// 				msKeeper.UpdateDVPairSDKBondAmount(ctx, delA, valA, bondSDKAmountA)
	// 				msKeeper.UpdateDVPairSDKBondAmount(ctx, delB, valB, bondSDKAmountB)
	// 				return []math.Int{bondSDKAmountA.Add(bondSDKAmountA), bondSDKAmountB}
	// 			},
	// 			dels: []sdk.AccAddress{delA, delB},
	// 			vals: []sdk.ValAddress{valA, valB},
	// 		},
	// 		{
	// 			name: "1 delegator, 2 validator, success",
	// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []math.Int {
	// 				msKeeper.UpdateDVPairSDKBondAmount(ctx, delA, valA, bondSDKAmountA)
	// 				msKeeper.UpdateDVPairSDKBondAmount(ctx, delA, valB, bondSDKAmountB)
	// 				msKeeper.UpdateDVPairSDKBondAmount(ctx, delA, valB, bondSDKAmountA)
	// 				return []math.Int{bondSDKAmountA, bondSDKAmountB.Add(bondSDKAmountA)}
	// 			},
	// 			dels: []sdk.AccAddress{delA, delA},
	// 			vals: []sdk.ValAddress{valA, valB},
	// 		},
	// 		{
	// 			name: "1 delegator, 1 validator, 2 bond amounts",
	// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) []math.Int {
	// 				msKeeper.UpdateDVPairSDKBondAmount(ctx, delA, valA, bondSDKAmountA)
	// 				msKeeper.UpdateDVPairSDKBondAmount(ctx, delA, valA, bondSDKAmountB)
	// 				return []math.Int{bondSDKAmountB.Add(bondSDKAmountA)}
	// 			},
	// 			dels: []sdk.AccAddress{delA},
	// 			vals: []sdk.ValAddress{valA},
	// 		},
	// 	}

	// 	for _, tc := range testCases {
	// 		tc := tc
	// 		suite.Run(tc.name, func() {
	// 			suite.SetupTest()

	//			inputs := tc.malleate(suite.ctx, suite.msKeeper)
	//			for idx, expOut := range inputs {
	//				actualCoin := suite.msKeeper.GetDVPairSDKBondAmount(suite.ctx, tc.dels[idx], tc.vals[idx])
	//				suite.Require().Equal(expOut, actualCoin)
	//			}
	//		})
	//	}
}

func (suite *KeeperTestSuite) TestCalSDKBondCoin() {
	// 	gasDenom := "ario"
	// 	govDenom := "arst"
	// 	testCases := []struct {
	// 		name     string
	// 		malleate func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (sdk.Coin, error)
	// 		expOut   sdk.Coin
	// 		expErr   bool
	// 	}{
	// 		{
	// 			name: "3001 coin, weight 0.3, expect 900",
	// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (sdk.Coin, error) {
	// 				msKeeper.SetBondCoinWeight(ctx, gasDenom, math.LegacyMustNewDecFromStr("0.3"))
	// 				return msKeeper.CalSDKBondCoin(ctx, sdk.NewCoin(gasDenom, sdk.NewInt(3001)))
	// 			},
	// 			expOut: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(900)),
	// 			expErr: false,
	// 		},
	// 		{
	// 			name: "25 coin, weight 0.5, expect 12",
	// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (sdk.Coin, error) {
	// 				msKeeper.SetBondCoinWeight(ctx, govDenom, math.LegacyMustNewDecFromStr("0.5"))
	// 				return msKeeper.CalSDKBondCoin(ctx, sdk.NewCoin(govDenom, sdk.NewInt(25)))
	// 			},
	// 			expOut: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(12)),
	// 			expErr: false,
	// 		},
	// 		{
	// 			name: "invalid bond coin",
	// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (sdk.Coin, error) {
	// 				return msKeeper.CalSDKBondCoin(ctx, sdk.NewCoin(govDenom, sdk.NewInt(25)))
	// 			},
	// 			expOut: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(12)),
	// 			expErr: true,
	// 		},
	// 	}

	// 	for _, tc := range testCases {
	// 		tc := tc
	// 		suite.Run(tc.name, func() {
	// 			suite.SetupTest()
	// 			actualOut, err := tc.malleate(suite.ctx, suite.msKeeper)

	// 			if tc.expErr {
	// 				suite.Require().Error(err)
	// 			} else {
	// 				suite.Require().NoError(err)
	// 				suite.Require().Equal(tc.expOut, actualOut)
	// 			}
	// 		})
	// 	}
	// }

	// func (suite *KeeperTestSuite) TestPreDelegate() {
	// 	delAddr := testutil.GenAddress()
	// 	valAddr := testutil.GenValAddress()
	// 	gasDenom := "ario"
	// 	govDenom := "arst"
	// 	testCases := []struct {
	// 		name     string
	// 		malleate func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (sdk.Coin, error)
	// 		expOut   sdk.Coin
	// 		expErr   bool
	// 	}{
	// 		{
	// 			name: "3001 coin, weight 0.3, expect 900",
	// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (sdk.Coin, error) {
	// 				msKeeper.SetBondCoinWeight(ctx, gasDenom, math.LegacyMustNewDecFromStr("0.3"))
	// 				msKeeper.SetValidatorAllowedCoin(ctx, valAddr, gasDenom)
	// 				bondAmount := sdk.NewCoin(gasDenom, sdk.NewInt(3001))
	// 				err := msKeeper.PreDelegate(ctx, delAddr, valAddr, bondAmount)
	// 				return bondAmount, err
	// 			},
	// 			expOut: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(900)),
	// 			expErr: false,
	// 		},
	// 		{
	// 			name: "25 coin, weight 0.5, expect 12",
	// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (sdk.Coin, error) {
	// 				msKeeper.SetBondCoinWeight(ctx, govDenom, math.LegacyMustNewDecFromStr("0.5"))
	// 				msKeeper.SetValidatorAllowedCoin(ctx, valAddr, govDenom)
	// 				bondAmount := sdk.NewCoin(govDenom, sdk.NewInt(25))
	// 				err := msKeeper.PreDelegate(ctx, delAddr, valAddr, bondAmount)
	// 				return bondAmount, err
	// 			},
	// 			expOut: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(12)),
	// 			expErr: false,
	// 		},
	// 		{
	// 			name: "invalid bond coin",
	// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (sdk.Coin, error) {
	// 				msKeeper.SetValidatorAllowedCoin(ctx, valAddr, gasDenom)
	// 				bondAmount := sdk.NewCoin(govDenom, sdk.NewInt(25))
	// 				err := msKeeper.PreDelegate(ctx, delAddr, valAddr, bondAmount)
	// 				return bondAmount, err
	// 			},
	// 			expOut: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(12)),
	// 			expErr: true,
	// 		},
	// 		{
	// 			name: "invalid val address",
	// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (sdk.Coin, error) {
	// 				bondAmount := sdk.NewCoin(govDenom, sdk.NewInt(25))
	// 				err := msKeeper.PreDelegate(ctx, delAddr, []byte{}, bondAmount)
	// 				return bondAmount, err
	// 			},
	// 			expOut: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(12)),
	// 			expErr: true,
	// 		},
	// 		{
	// 			name: "invalid val address coin",
	// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (sdk.Coin, error) {
	// 				msKeeper.SetValidatorAllowedCoin(ctx, valAddr, gasDenom)
	// 				bondAmount := sdk.NewCoin(govDenom, sdk.NewInt(25))
	// 				err := msKeeper.PreDelegate(ctx, delAddr, valAddr, bondAmount)
	// 				return bondAmount, err
	// 			},
	// 			expOut: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(12)),
	// 			expErr: true,
	// 		},
	// 	}

	// 	for _, tc := range testCases {
	// 		tc := tc
	// 		suite.Run(tc.name, func() {
	// 			suite.SetupTest()
	// 			bondAmount, err := tc.malleate(suite.ctx, suite.msKeeper)

	//			if tc.expErr {
	//				suite.Require().Error(err)
	//			} else {
	//				suite.Require().NoError(err)
	//				actualBond := suite.msKeeper.GetDVPairBondAmount(suite.ctx, delAddr, valAddr)
	//				actualSDKBond := suite.msKeeper.GetDVPairSDKBondAmount(suite.ctx, delAddr, valAddr)
	//				suite.Require().Equal(bondAmount.Amount, actualBond)
	//				suite.Require().Equal(tc.expOut.Amount, actualSDKBond)
	//			}
	//		})
	//	}
}
