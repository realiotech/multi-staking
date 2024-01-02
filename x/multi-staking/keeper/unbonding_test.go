package keeper_test

// import (
// 	"time"

// 	"cosmossdk.io/math"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
// 	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

// 	"github.com/realio-tech/multi-staking-module/testutil"
// 	multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
// 	multistakingtypes "github.com/realio-tech/multi-staking-module/x/multi-staking/types"
// )

// func (suite *KeeperTestSuite) TestCompleteUnbonding() {
// 	suite.SetupTest()
// 	delAddr := testutil.GenAddress()
// 	valAddr := testutil.GenValAddress()
// 	imAddr := multistakingtypes.IntermediaryDelegator(delAddr)
// 	gasDenom := "ario"
// 	minTime := time.Now()
// 	testCases := []struct {
// 		name     string
// 		malleate func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error)
// 		expOut   math.Int
// 		expErr   bool
// 	}{
// 		{
// 			name: "lock 3001 coin, rate 0.3, expect unlock 3000",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				unbondAmt := sdk.NewInt(3001)
// 				weight := sdk.MustNewDecFromStr("0.3")
// 				msKeeper.SetMultiStakingUnlockEntry(ctx, delAddr, valAddr, 1, weight, minTime, unbondAmt)
// 				unlockAmt, err := msKeeper.CompleteUnbonding(ctx, imAddr, valAddr, 1, sdk.NewInt(900))

// 				return unlockAmt[0].Amount, err
// 			},
// 			expOut: sdk.NewInt(3000),
// 			expErr: false,
// 		},
// 		{
// 			name: "lock 25 coin, weight 0.5, expect unlock 24",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				unbondAmt := sdk.NewInt(25)
// 				weight := sdk.MustNewDecFromStr("0.5")
// 				msKeeper.SetMultiStakingUnlockEntry(ctx, delAddr, valAddr, 1, weight, minTime, unbondAmt)
// 				unlockAmt, err := msKeeper.CompleteUnbonding(ctx, imAddr, valAddr, 1, sdk.NewInt(12))

// 				return unlockAmt[0].Amount, err
// 			},
// 			expOut: sdk.NewInt(24),
// 			expErr: false,
// 		},
// 		{
// 			name: "lock 25 coin, weight 0.5, unlock 11 bond coin because of slashing, expect unlock 22 coin",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				unbondAmt := sdk.NewInt(25)
// 				weight := sdk.MustNewDecFromStr("0.5")
// 				msKeeper.SetMultiStakingUnlockEntry(ctx, delAddr, valAddr, 1, weight, minTime, unbondAmt)
// 				unlockAmt, err := msKeeper.CompleteUnbonding(ctx, imAddr, valAddr, 1, sdk.NewInt(11))

// 				return unlockAmt[0].Amount, err
// 			},
// 			expOut: sdk.NewInt(22),
// 			expErr: false,
// 		},
// 		{
// 			name: "unbond record not exist",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				_, err := msKeeper.CompleteUnbonding(ctx, imAddr, valAddr, 1, sdk.NewInt(11))

// 				return sdk.NewInt(0), err
// 			},
// 			expErr: true,
// 		},
// 		{
// 			name: "entry at height not exist",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				unbondAmt := sdk.NewInt(25)
// 				weight := sdk.MustNewDecFromStr("0.5")
// 				msKeeper.SetMultiStakingUnlockEntry(ctx, delAddr, valAddr, 1, weight, minTime, unbondAmt)
// 				_, err := msKeeper.CompleteUnbonding(ctx, imAddr, valAddr, 4, sdk.NewInt(11))

// 				return sdk.NewInt(0), err
// 			},
// 			expErr: true,
// 		},
// 		{
// 			name: "unlock too much",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				unbondAmt := sdk.NewInt(25)
// 				weight := sdk.MustNewDecFromStr("0.5")
// 				msKeeper.SetMultiStakingUnlockEntry(ctx, delAddr, valAddr, 1, weight, minTime, unbondAmt)
// 				_, err := msKeeper.CompleteUnbonding(ctx, imAddr, valAddr, 1, sdk.NewInt(13))

// 				return sdk.NewInt(0), err
// 			},
// 			expErr: true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc
// 		suite.Run(tc.name, func() {
// 			suite.SetupTest()
// 			suite.msKeeper.SetValidatorMultiStakingCoin(suite.ctx, valAddr, gasDenom)
// 			imAccBalance := sdk.NewCoins(sdk.NewCoin(stakingtypes.DefaultParams().BondDenom, sdk.NewInt(10000)), sdk.NewCoin(gasDenom, sdk.NewInt(10000)))
// 			suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, imAccBalance)
// 			suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, imAddr, imAccBalance)
// 			unlockAmt, err := tc.malleate(suite.ctx, suite.msKeeper)
// 			if tc.expErr {
// 				suite.Require().Error(err)
// 			} else {
// 				suite.Require().NoError(err)
// 				suite.Require().Equal(unlockAmt, tc.expOut)
// 			}
// 		})
// 	}
// }
