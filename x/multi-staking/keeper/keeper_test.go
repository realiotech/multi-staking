package keeper_test

import (
	"testing"
	"time"

	"github.com/realio-tech/multi-staking-module/test"
	"github.com/realio-tech/multi-staking-module/test/simapp"
	multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	coreheader "cosmossdk.io/core/header"
	"cosmossdk.io/math"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app       *simapp.SimApp
	ctx       sdk.Context
	msKeeper  *multistakingkeeper.Keeper
	govKeeper govkeeper.Keeper
	msgServer stakingtypes.MsgServer
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup()
	ctx := app.NewContextLegacy(false, tmproto.Header{Height: app.LastBlockHeight() + 1})

	_, err := app.CrisisKeeper.ConstantFee.Get(ctx)
	suite.Require().NoError(err)

	multiStakingMsgServer := multistakingkeeper.NewMsgServerImpl(app.MultiStakingKeeper)

	suite.app, suite.ctx, suite.msKeeper, suite.govKeeper, suite.msgServer = app, ctx, &app.MultiStakingKeeper, app.GovKeeper, multiStakingMsgServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestAdjustUnbondAmount() {
	delAddr := test.GenAddress()
	valPubKey := test.GenPubKey()
	valAddr := sdk.ValAddress(valPubKey.Address())

	testCases := []struct {
		name         string
		malleate     func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (sdk.Coin, error)
		adjustAmount math.Int
		expAmount    math.Int
		expErr       bool
	}{
		{
			name: "success and not change adjust amount",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (sdk.Coin, error) {
				multiStakingAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(1000))
				delMsg := stakingtypes.NewMsgDelegate(delAddr.String(), valAddr.String(), multiStakingAmount)
				_, err := msgServer.Delegate(ctx, delMsg)
				return multiStakingAmount, err
			},
			adjustAmount: math.NewInt(800),
			expAmount:    math.NewInt(800),
			expErr:       false,
		},
		{
			name: "success and reduce adjust amount",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (sdk.Coin, error) {
				multiStakingAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(1000))
				delMsg := stakingtypes.NewMsgDelegate(delAddr.String(), valAddr.String(), multiStakingAmount)
				_, err := msgServer.Delegate(ctx, delMsg)

				return multiStakingAmount, err
			},
			adjustAmount: math.NewInt(2000),
			expAmount:    math.NewInt(1000),
			expErr:       false,
		},
		{
			name: "not found delegation",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (sdk.Coin, error) {
				multiStakingAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(1000))
				return multiStakingAmount, nil
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			newParam := stakingtypes.DefaultParams()
			newParam.MinCommissionRate = math.LegacyMustNewDecFromStr("0.02")
			err := suite.app.StakingKeeper.SetParams(suite.ctx, newParam)
			suite.Require().NoError(err)
			suite.msKeeper.SetBondWeight(suite.ctx, MultiStakingDenomA, math.LegacyOneDec())
			bondAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(1000))
			userBalance := sdk.NewCoin(MultiStakingDenomA, math.NewInt(10000))
			suite.FundAccount(delAddr, sdk.NewCoins(userBalance))
			suite.FundAccount(sdk.AccAddress(valAddr), sdk.NewCoins(userBalance))

			createMsg := stakingtypes.MsgCreateValidator{
				Description: stakingtypes.Description{
					Moniker:         "test",
					Identity:        "test",
					Website:         "test",
					SecurityContact: "test",
					Details:         "test",
				},
				Commission: stakingtypes.CommissionRates{
					Rate:          math.LegacyMustNewDecFromStr("0.05"),
					MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
					MaxChangeRate: math.LegacyMustNewDecFromStr("0.05"),
				},
				MinSelfDelegation: math.NewInt(200),
				DelegatorAddress:  sdk.AccAddress(valAddr).String(),
				ValidatorAddress:  valAddr.String(),
				Pubkey:            codectypes.UnsafePackAny(valPubKey),
				Value:             bondAmount,
			}
			_, err = suite.msgServer.CreateValidator(suite.ctx, &createMsg)
			suite.Require().NoError(err)
			_, err = tc.malleate(suite.ctx, suite.msgServer, *suite.msKeeper)
			suite.Require().NoError(err)

			actualAmt, err := suite.msKeeper.AdjustUnbondAmount(suite.ctx, delAddr, valAddr, tc.adjustAmount)

			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(actualAmt, tc.expAmount)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestAdjustCancelUnbondAmount() {
	delAddr := test.GenAddress()
	valPubKey := test.GenPubKey()
	valAddr := sdk.ValAddress(valPubKey.Address())

	testCases := []struct {
		name         string
		malleate     func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (sdk.Coin, error)
		adjustAmount math.Int
		expAmount    math.Int
		expErr       bool
	}{
		{
			name: "success and not change adjust amount",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (sdk.Coin, error) {
				multiStakingAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(1000))
				undelMsg := stakingtypes.NewMsgUndelegate(delAddr.String(), valAddr.String(), multiStakingAmount)
				_, err := msgServer.Undelegate(ctx, undelMsg)
				return multiStakingAmount, err
			},
			adjustAmount: math.NewInt(800),
			expAmount:    math.NewInt(800),
			expErr:       false,
		},
		{
			name: "success with many unbonding delegations",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (sdk.Coin, error) {
				multiStakingAmount1 := sdk.NewCoin(MultiStakingDenomA, math.NewInt(400))
				undelMsg1 := stakingtypes.NewMsgUndelegate(delAddr.String(), valAddr.String(), multiStakingAmount1)
				_, err := msgServer.Undelegate(ctx, undelMsg1)
				suite.Require().NoError(err)

				multiStakingAmount2 := sdk.NewCoin(MultiStakingDenomA, math.NewInt(500))
				undelMsg2 := stakingtypes.NewMsgUndelegate(delAddr.String(), valAddr.String(), multiStakingAmount2)
				_, err = msgServer.Undelegate(ctx, undelMsg2)
				suite.Require().NoError(err)

				multiStakingAmount3 := sdk.NewCoin(MultiStakingDenomA, math.NewInt(600))
				undelMsg3 := stakingtypes.NewMsgUndelegate(delAddr.String(), valAddr.String(), multiStakingAmount3)
				_, err = msgServer.Undelegate(ctx, undelMsg3)
				suite.Require().NoError(err)

				return multiStakingAmount1, nil
			},
			adjustAmount: math.NewInt(1500),
			expAmount:    math.NewInt(1500),
			expErr:       false,
		},
		{
			name: "success and reduce adjust amount",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (sdk.Coin, error) {
				multiStakingAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(1000))
				undelMsg := stakingtypes.NewMsgUndelegate(delAddr.String(), valAddr.String(), multiStakingAmount)
				_, err := msgServer.Undelegate(ctx, undelMsg)

				return multiStakingAmount, err
			},
			adjustAmount: math.NewInt(2000),
			expAmount:    math.NewInt(1000),
			expErr:       false,
		},
		{
			name: "not found delegation",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (sdk.Coin, error) {
				multiStakingAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(1000))
				return multiStakingAmount, nil
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			newParam := stakingtypes.DefaultParams()
			newParam.MinCommissionRate = math.LegacyMustNewDecFromStr("0.02")
			err := suite.app.StakingKeeper.SetParams(suite.ctx, newParam)
			suite.Require().NoError(err)
			suite.msKeeper.SetBondWeight(suite.ctx, MultiStakingDenomA, math.LegacyOneDec())
			bondAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(5000))
			userBalance := sdk.NewCoin(MultiStakingDenomA, math.NewInt(10000))
			suite.FundAccount(delAddr, sdk.NewCoins(userBalance))
			suite.FundAccount(sdk.AccAddress(valAddr), sdk.NewCoins(userBalance))

			createMsg := stakingtypes.MsgCreateValidator{
				Description: stakingtypes.Description{
					Moniker:         "test",
					Identity:        "test",
					Website:         "test",
					SecurityContact: "test",
					Details:         "test",
				},
				Commission: stakingtypes.CommissionRates{
					Rate:          math.LegacyMustNewDecFromStr("0.05"),
					MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
					MaxChangeRate: math.LegacyMustNewDecFromStr("0.05"),
				},
				MinSelfDelegation: math.NewInt(200),
				DelegatorAddress:  sdk.AccAddress(valAddr).String(),
				ValidatorAddress:  valAddr.String(),
				Pubkey:            codectypes.UnsafePackAny(valPubKey),
				Value:             bondAmount,
			}
			_, err = suite.msgServer.CreateValidator(suite.ctx, &createMsg)
			suite.Require().NoError(err)

			delMsg := stakingtypes.NewMsgDelegate(delAddr.String(), valAddr.String(), bondAmount)
			_, err = suite.msgServer.Delegate(suite.ctx, delMsg)
			suite.Require().NoError(err)

			_, err = tc.malleate(suite.ctx, suite.msgServer, *suite.msKeeper)
			suite.Require().NoError(err)

			suite.ctx = suite.ctx.WithBlockHeader(tmproto.Header{Time: time.Now()}).WithBlockHeight(1)

			actualAmt, err := suite.msKeeper.AdjustCancelUnbondingAmount(suite.ctx, delAddr, valAddr, suite.ctx.BlockHeight(), tc.adjustAmount)

			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(actualAmt, tc.expAmount)
			}
		})
	}
}

// Todo: add CheckBalance; AddAccountWithCoin; FundAccount
func (suite *KeeperTestSuite) NextBlock(jumpTime time.Duration) {
	app := suite.app
	ctx := suite.ctx
	_, err := app.FinalizeBlock(&abci.RequestFinalizeBlock{Height: ctx.BlockHeight(), Time: ctx.BlockTime()})
	suite.Require().NoError(err)
	_, err = app.Commit()
	suite.Require().NoError(err)
	newBlockTime := ctx.BlockTime().Add(jumpTime)

	header := ctx.BlockHeader()
	header.Time = newBlockTime
	header.Height++

	newCtx := app.BaseApp.NewUncachedContext(false, header).WithHeaderInfo(coreheader.Info{
		Height: header.Height,
		Time:   header.Time,
	})

	suite.ctx = newCtx
}

// Todo: add CheckBalance; AddAccountWithCoin; FundAccount
func (suite *KeeperTestSuite) FundAccount(addr sdk.AccAddress, amounts sdk.Coins) {
	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, amounts)
	require.NoError(suite.T(), err)

	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr, amounts)
	require.NoError(suite.T(), err)
}

func (suite *KeeperTestSuite) CreateAndFundAccount(amounts sdk.Coins) sdk.AccAddress {
	addr := test.GenAddress()
	suite.FundAccount(addr, amounts)
	return addr
}

func (suite *KeeperTestSuite) CheckBalance(addr sdk.AccAddress, coins sdk.Coins) {
	accBalance := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr)

	require.Equal(suite.T(), accBalance, coins)
}

func SoftEqualInt(a math.Int, b math.Int) bool {
	biggerNum := math.MaxInt(a, b)
	smallerNum := math.MinInt(a, b)

	biggerNumDec := math.LegacyNewDecFromInt(biggerNum)
	smallerNumDec := math.LegacyNewDecFromInt(smallerNum)

	return smallerNumDec.Quo(biggerNumDec).GTE(math.LegacyMustNewDecFromStr("0.99"))
}

func DiffLTEThanOne(a, b math.Int) bool {
	return a.Sub(b).Abs().LTE(math.OneInt())
}
