package keeper_test

// import (
// 	"cosmossdk.io/math"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
// 	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

// 	"github.com/realio-tech/multi-staking-module/testutil"
// 	multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
// 	multistakingtypes "github.com/realio-tech/multi-staking-module/x/multi-staking/types"
// )

// func (suite *KeeperTestSuite) TestLockedAmountToBondAmount() {
// 	delAddr := testutil.GenAddress()
// 	valAddr := testutil.GenValAddress()

// 	testCases := []struct {
// 		name     string
// 		malleate func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error)
// 		expOut   math.Int
// 		expErr   bool
// 	}{
// 		{
// 			name: "3001 token, weight 0.3, expect 900",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				lockAmt := sdk.NewInt(3001)
// 				weight := sdk.MustNewDecFromStr("0.3")
// 				lockRecord := multistakingtypes.NewMultiStakingLock(lockAmt, weight, delAddr, valAddr)
// 				msKeeper.SetMultiStakingLock(ctx, multistakingtypes.MultiStakingLockID(delAddr, valAddr), lockRecord)
// 				return msKeeper.LockedAmountToBondAmount(ctx, delAddr, valAddr, lockAmt)
// 			},
// 			expOut: sdk.NewInt(900),
// 			expErr: false,
// 		},
// 		{
// 			name: "25 token, weight 0.5, expect 12",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				lockAmt := sdk.NewInt(25)
// 				weight := sdk.MustNewDecFromStr("0.5")
// 				lockRecord := multistakingtypes.NewMultiStakingLock(lockAmt, weight, delAddr, valAddr)
// 				msKeeper.SetMultiStakingLock(ctx, multistakingtypes.MultiStakingLockID(delAddr, valAddr), lockRecord)
// 				return msKeeper.LockedAmountToBondAmount(ctx, delAddr, valAddr, lockAmt)
// 			},
// 			expOut: sdk.NewInt(12),
// 			expErr: false,
// 		},
// 		{
// 			name: "lock not exist",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				lockAmt := sdk.NewInt(3000)
// 				return msKeeper.LockedAmountToBondAmount(ctx, delAddr, valAddr, lockAmt)
// 			},
// 			expErr: true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc
// 		suite.Run(tc.name, func() {
// 			suite.SetupTest()

// 			bondAmount, err := tc.malleate(suite.ctx, suite.msKeeper)
// 			if tc.expErr {
// 				suite.Require().Error(err)
// 			} else {
// 				suite.Require().NoError(err)
// 				lockId := multistakingtypes.MultiStakingLockID(delAddr, valAddr)
// 				_, found := suite.msKeeper.GetMultiStakingLock(suite.ctx, lockId)
// 				suite.Require().True(found)
// 				suite.Require().Equal(bondAmount, tc.expOut)
// 			}
// 		})
// 	}
// }

// func (suite *KeeperTestSuite) TestAddTokenToLock() {
// 	delAddr := testutil.GenAddress()
// 	valAddr := testutil.GenValAddress()

// 	testCases := []struct {
// 		name     string
// 		malleate func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) math.Int
// 		expAmt   math.Int
// 		expRate  sdk.Dec
// 		expErr   bool
// 	}{
// 		{
// 			name: "3000 token, rate 0.3",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) math.Int {
// 				lockAmt := sdk.NewInt(3000)
// 				weight := sdk.MustNewDecFromStr("0.3")
// 				lockRecord := msKeeper.AddTokenToLock(ctx, delAddr, valAddr, lockAmt, weight)
// 				return lockRecord.LockedAmount
// 			},
// 			expAmt:  sdk.NewInt(3000),
// 			expRate: sdk.MustNewDecFromStr("0.3"),
// 			expErr:  false,
// 		},
// 		{
// 			name: "1000 * 0.5  + 500 * 0.8 = 1500 * 0.6",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) math.Int {
// 				lockAmt := sdk.NewInt(1000)
// 				weight := sdk.MustNewDecFromStr("0.5")
// 				lockRecord := msKeeper.AddTokenToLock(ctx, delAddr, valAddr, lockAmt, weight)
// 				lockAmt1 := sdk.NewInt(500)
// 				weight1 := sdk.MustNewDecFromStr("0.8")
// 				lockRecord = msKeeper.AddTokenToLock(ctx, delAddr, valAddr, lockAmt1, weight1)
// 				return lockRecord.LockedAmount
// 			},
// 			expAmt:  sdk.NewInt(1500),
// 			expRate: sdk.MustNewDecFromStr("0.6"),
// 			expErr:  false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc
// 		suite.Run(tc.name, func() {
// 			suite.SetupTest()

// 			bondAmount := tc.malleate(suite.ctx, suite.msKeeper)
// 			lockId := multistakingtypes.MultiStakingLockID(delAddr, valAddr)
// 			lockRecord, found := suite.msKeeper.GetMultiStakingLock(suite.ctx, lockId)
// 			suite.Require().True(found)
// 			suite.Require().Equal(lockRecord.LockedAmount, tc.expAmt)
// 			suite.Require().Equal(bondAmount, tc.expAmt)
// 			suite.Require().Equal(lockRecord.ConversionRatio, tc.expRate)
// 		})
// 	}
// }

// func (suite *KeeperTestSuite) TestRemoveTokenFromLock() {
// 	delAddr := testutil.GenAddress()
// 	valAddr := testutil.GenValAddress()

// 	testCases := []struct {
// 		name     string
// 		malleate func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error)
// 		expOut   math.Int
// 		expErr   bool
// 	}{
// 		{
// 			name: "3001 token, weight 0.3, remove 2000",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				lockAmt := sdk.NewInt(3001)
// 				weight := sdk.MustNewDecFromStr("0.3")
// 				lockRecord := multistakingtypes.NewMultiStakingLock(lockAmt, weight, delAddr, valAddr)
// 				msKeeper.SetMultiStakingLock(ctx, multistakingtypes.MultiStakingLockID(delAddr, valAddr), lockRecord)
// 				lockRemoveAmt := sdk.NewInt(2000)
// 				lockRecord, err := msKeeper.RemoveTokenFromLock(ctx, delAddr, valAddr, lockRemoveAmt)
// 				return lockRecord.LockedAmount, err
// 			},
// 			expOut: sdk.NewInt(1001),
// 			expErr: false,
// 		},
// 		{
// 			name: "3001 token, weight 0.3, remove 4000",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				lockAmt := sdk.NewInt(3001)
// 				weight := sdk.MustNewDecFromStr("0.3")
// 				lockRecord := multistakingtypes.NewMultiStakingLock(lockAmt, weight, delAddr, valAddr)
// 				msKeeper.SetMultiStakingLock(ctx, multistakingtypes.MultiStakingLockID(delAddr, valAddr), lockRecord)
// 				lockRemoveAmt := sdk.NewInt(4000)
// 				lockRecord, err := msKeeper.RemoveTokenFromLock(ctx, delAddr, valAddr, lockRemoveAmt)
// 				return lockRecord.LockedAmount, err
// 			},
// 			expOut: sdk.NewInt(1001),
// 			expErr: true,
// 		},
// 		{
// 			name: "lock not exist",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				lockAmt := sdk.NewInt(3000)
// 				lockRecord, err := msKeeper.RemoveTokenFromLock(ctx, delAddr, valAddr, lockAmt)
// 				return lockRecord.LockedAmount, err
// 			},
// 			expErr: true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc
// 		suite.Run(tc.name, func() {
// 			suite.SetupTest()

// 			bondAmount, err := tc.malleate(suite.ctx, suite.msKeeper)
// 			if tc.expErr {
// 				suite.Require().Error(err)
// 			} else {
// 				suite.Require().NoError(err)
// 				lockId := multistakingtypes.MultiStakingLockID(delAddr, valAddr)
// 				_, found := suite.msKeeper.GetMultiStakingLock(suite.ctx, lockId)
// 				suite.Require().True(found)
// 				suite.Require().Equal(bondAmount, tc.expOut)
// 			}
// 		})
// 	}
// }

// func (suite *KeeperTestSuite) TestMoveLockedMultistakingToken() {
// 	delAddr := testutil.GenAddress()
// 	valAddr := testutil.GenValAddress()
// 	valAddr2 := testutil.GenValAddress()

// 	testCases := []struct {
// 		name     string
// 		malleate func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error)
// 		expOut   math.Int
// 		expErr   bool
// 	}{
// 		{
// 			name: "3001 token, weight 0.3, move 2000",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				lockAmt := sdk.NewInt(3001)
// 				weight := sdk.MustNewDecFromStr("0.3")
// 				lockRecord := multistakingtypes.NewMultiStakingLock(lockAmt, weight, delAddr, valAddr)
// 				msKeeper.SetMultiStakingLock(ctx, multistakingtypes.MultiStakingLockID(delAddr, valAddr), lockRecord)

// 				moveAmt := sdk.NewInt(2000)
// 				err := msKeeper.MoveLockedMultistakingToken(ctx, delAddr, valAddr, valAddr2, sdk.NewCoin("ario", moveAmt))

// 				return moveAmt, err
// 			},
// 			expOut: sdk.NewInt(1001),
// 			expErr: false,
// 		},
// 		{
// 			name: "move to existed delegation",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				lockAmt := sdk.NewInt(3001)
// 				weight := sdk.MustNewDecFromStr("0.3")
// 				lockRecord := multistakingtypes.NewMultiStakingLock(lockAmt, weight, delAddr, valAddr)
// 				msKeeper.AddTokenToLock(ctx, delAddr, valAddr2, lockAmt, weight)
// 				msKeeper.SetMultiStakingLock(ctx, multistakingtypes.MultiStakingLockID(delAddr, valAddr), lockRecord)

// 				moveAmt := sdk.NewInt(2000)
// 				err := msKeeper.MoveLockedMultistakingToken(ctx, delAddr, valAddr, valAddr2, sdk.NewCoin("ario", moveAmt))

// 				return moveAmt.Add(lockAmt), err
// 			},
// 			expOut: sdk.NewInt(1001),
// 			expErr: false,
// 		},
// 		{
// 			name: "lock not exist",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				moveAmt := sdk.NewInt(2000)
// 				err := msKeeper.MoveLockedMultistakingToken(ctx, delAddr, valAddr, valAddr2, sdk.NewCoin("ario", moveAmt))

// 				return moveAmt, err
// 			},
// 			expErr: true,
// 		},
// 		{
// 			name: "3001 token, weight 0.3, move 4000 failed",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				lockAmt := sdk.NewInt(3001)
// 				weight := sdk.MustNewDecFromStr("0.3")
// 				lockRecord := multistakingtypes.NewMultiStakingLock(lockAmt, weight, delAddr, valAddr)
// 				msKeeper.SetMultiStakingLock(ctx, multistakingtypes.MultiStakingLockID(delAddr, valAddr), lockRecord)

// 				moveAmt := sdk.NewInt(4000)
// 				err := msKeeper.MoveLockedMultistakingToken(ctx, delAddr, valAddr, valAddr2, sdk.NewCoin("ario", moveAmt))

// 				return moveAmt, err
// 			},
// 			expOut: sdk.NewInt(1001),
// 			expErr: true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc
// 		suite.Run(tc.name, func() {
// 			suite.SetupTest()

// 			val2BondAmt, err := tc.malleate(suite.ctx, suite.msKeeper)
// 			if tc.expErr {
// 				suite.Require().Error(err)
// 			} else {
// 				suite.Require().NoError(err)
// 				lockId1 := multistakingtypes.MultiStakingLockID(delAddr, valAddr)
// 				lockRecord1, found := suite.msKeeper.GetMultiStakingLock(suite.ctx, lockId1)
// 				suite.Require().True(found)
// 				lockId2 := multistakingtypes.MultiStakingLockID(delAddr, valAddr2)
// 				lockRecord2, found := suite.msKeeper.GetMultiStakingLock(suite.ctx, lockId2)
// 				suite.Require().True(found)
// 				suite.Require().Equal(val2BondAmt, lockRecord2.LockedAmount)
// 				suite.Require().Equal(tc.expOut, lockRecord1.LockedAmount)
// 			}
// 		})
// 	}
// }

// func (suite *KeeperTestSuite) TestLockMultiStakingTokenAndMintBondToken() {
// 	delAddr := testutil.GenAddress()
// 	valAddr := testutil.GenValAddress()
// 	gasDenom := "ario"
// 	govDenom := "arst"
// 	testCases := []struct {
// 		name     string
// 		malleate func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error)
// 		expOut   math.Int
// 		expErr   bool
// 	}{
// 		{
// 			name: "3001 token, weight 0.3, expect 900",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				lockAmt := sdk.NewInt(3001)
// 				weight := sdk.MustNewDecFromStr("0.3")
// 				msKeeper.SetBondTokenWeight(ctx, gasDenom, weight)
// 				mintAmt, err := msKeeper.LockMultiStakingTokenAndMintBondToken(ctx, delAddr, valAddr, sdk.NewCoin(gasDenom, lockAmt))
// 				return mintAmt.Amount, err
// 			},
// 			expOut: sdk.NewInt(900),
// 			expErr: false,
// 		},
// 		{
// 			name: "25 token, weight 0.5, expect 12",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				lockAmt := sdk.NewInt(25)
// 				weight := sdk.MustNewDecFromStr("0.5")
// 				msKeeper.SetBondTokenWeight(ctx, gasDenom, weight)
// 				mintAmt, err := msKeeper.LockMultiStakingTokenAndMintBondToken(ctx, delAddr, valAddr, sdk.NewCoin(gasDenom, lockAmt))
// 				return mintAmt.Amount, err
// 			},
// 			expOut: sdk.NewInt(12),
// 			expErr: false,
// 		},
// 		{
// 			name: "invalid coin",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				lockAmt := sdk.NewInt(3001)
// 				mintAmt, err := msKeeper.LockMultiStakingTokenAndMintBondToken(ctx, delAddr, valAddr, sdk.NewCoin(gasDenom, lockAmt))
// 				return mintAmt.Amount, err
// 			},
// 			expErr: true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc
// 		suite.Run(tc.name, func() {
// 			suite.SetupTest()
// 			valCoins := sdk.NewCoins(sdk.NewCoin(gasDenom, sdk.NewInt(10000)), sdk.NewCoin(govDenom, sdk.NewInt(10000)))
// 			suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, valCoins)
// 			suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, delAddr, valCoins)
// 			bondAmount, err := tc.malleate(suite.ctx, suite.msKeeper)
// 			if tc.expErr {
// 				suite.Require().Error(err)
// 			} else {
// 				suite.Require().NoError(err)
// 				suite.Require().Equal(bondAmount, tc.expOut)
// 			}
// 		})
// 	}
// }

// func (suite *KeeperTestSuite) TestBurnBondTokenAndUnlockMultiStakingToken() {
// 	delAddr := testutil.GenAddress()
// 	valAddr := testutil.GenValAddress()
// 	imAddr := multistakingtypes.IntermediaryAccount(delAddr)
// 	gasDenom := "ario"

// 	testCases := []struct {
// 		name      string
// 		malleate  func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error)
// 		expOut    math.Int
// 		expRemain math.Int
// 		expErr    bool
// 	}{
// 		{
// 			name: "lock 3001 token, rate 0.3, expect unlock 3000",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				lockAmt := sdk.NewInt(3001)
// 				weight := sdk.MustNewDecFromStr("0.3")
// 				lockRecord := multistakingtypes.NewMultiStakingLock(lockAmt, weight, delAddr, valAddr)
// 				msKeeper.SetMultiStakingLock(ctx, multistakingtypes.MultiStakingLockID(delAddr, valAddr), lockRecord)

// 				unlockSDKAmt := sdk.NewCoin(stakingtypes.DefaultParams().BondDenom, sdk.NewInt(900))
// 				unlockAmt, err := msKeeper.BurnBondTokenAndUnlockMultiStakingToken(ctx, imAddr, valAddr, unlockSDKAmt)

// 				return unlockAmt[0].Amount, err
// 			},
// 			expOut: sdk.NewInt(3000),
// 			expRemain: sdk.NewInt(3001),
// 			expErr: false,
// 		},
// 		{
// 			name: "lock 25 token, weight 0.5, expect unlock 24",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				lockAmt := sdk.NewInt(25)
// 				weight := sdk.MustNewDecFromStr("0.5")
// 				lockRecord := multistakingtypes.NewMultiStakingLock(lockAmt, weight, delAddr, valAddr)
// 				msKeeper.SetMultiStakingLock(ctx, multistakingtypes.MultiStakingLockID(delAddr, valAddr), lockRecord)

// 				unlockSDKAmt := sdk.NewCoin(stakingtypes.DefaultParams().BondDenom, sdk.NewInt(12))
// 				unlockAmt, err := msKeeper.BurnBondTokenAndUnlockMultiStakingToken(ctx, imAddr, valAddr, unlockSDKAmt)

// 				return unlockAmt[0].Amount, err
// 			},
// 			expOut: sdk.NewInt(24),
// 			expRemain: sdk.NewInt(25),
// 			expErr: false,
// 		},
// 		{
// 			name: "lock 25 token, weight 0.5, unlock 11 bond token, expect 22 token unlock, remain 3 token",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				lockAmt := sdk.NewInt(25)
// 				weight := sdk.MustNewDecFromStr("0.5")
// 				lockRecord := multistakingtypes.NewMultiStakingLock(lockAmt, weight, delAddr, valAddr)
// 				msKeeper.SetMultiStakingLock(ctx, multistakingtypes.MultiStakingLockID(delAddr, valAddr), lockRecord)

// 				unlockSDKAmt := sdk.NewCoin(stakingtypes.DefaultParams().BondDenom, sdk.NewInt(11))
// 				unlockAmt, err := msKeeper.BurnBondTokenAndUnlockMultiStakingToken(ctx, imAddr, valAddr, unlockSDKAmt)

// 				return unlockAmt[0].Amount, err
// 			},
// 			expOut: sdk.NewInt(22),
// 			expRemain: sdk.NewInt(25),
// 			expErr: false,
// 		},
// 		{
// 			name: "lock not exist",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				unlockSDKAmt := sdk.NewCoin(stakingtypes.DefaultParams().BondDenom, sdk.NewInt(11))
// 				_, err := msKeeper.BurnBondTokenAndUnlockMultiStakingToken(ctx, imAddr, valAddr, unlockSDKAmt)

// 				return sdk.NewInt(0), err
// 			},
// 			expErr: true,
// 		},
// 		{
// 			name: "unlock too much",
// 			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper) (math.Int, error) {
// 				lockAmt := sdk.NewInt(25)
// 				weight := sdk.MustNewDecFromStr("0.5")
// 				lockRecord := multistakingtypes.NewMultiStakingLock(lockAmt, weight, delAddr, valAddr)
// 				msKeeper.SetMultiStakingLock(ctx, multistakingtypes.MultiStakingLockID(delAddr, valAddr), lockRecord)

// 				unlockSDKAmt := sdk.NewCoin(stakingtypes.DefaultParams().BondDenom, sdk.NewInt(13))
// 				_, err := msKeeper.BurnBondTokenAndUnlockMultiStakingToken(ctx, imAddr, valAddr, unlockSDKAmt)

// 				return math.Int{}, err
// 			},
// 			expErr: true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc
// 		suite.Run(tc.name, func() {
// 			suite.SetupTest()
// 			suite.msKeeper.SetValidatorAllowedToken(suite.ctx, valAddr, gasDenom)
// 			imAccBalance := sdk.NewCoins(sdk.NewCoin(stakingtypes.DefaultParams().BondDenom, sdk.NewInt(10000)), sdk.NewCoin(gasDenom, sdk.NewInt(10000)))
// 			suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, imAccBalance)
// 			suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, imAddr, imAccBalance)
// 			unlockAmt, err := tc.malleate(suite.ctx, suite.msKeeper)
// 			if tc.expErr {
// 				suite.Require().Error(err)
// 			} else {
// 				suite.Require().NoError(err)
// 				lockId := multistakingtypes.MultiStakingLockID(delAddr, valAddr)
// 				lockRecord, found := suite.msKeeper.GetMultiStakingLock(suite.ctx, lockId)
// 				suite.Require().True(found)
// 				suite.Require().Equal(unlockAmt, tc.expOut)
// 				suite.Require().Equal(lockRecord.LockedAmount, tc.expRemain)
// 			}
// 		})
// 	}
// }
