package keeper_test

import (
	"fmt"
	"testing"
	"time"

	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/realio-tech/multi-staking-module/test"
	"github.com/realio-tech/multi-staking-module/test/simapp"
	multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"cosmossdk.io/math"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var (
	MultiStakingDenomA = "ario"
	MultiStakingDenomB = "arst"
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
	fmt.Println("Run SetupTest")
	configurator := evmtypes.NewEVMConfigurator()
	configurator.ResetTestConfig()
	app := simapp.Setup(suite.T(), false)
	ctx := app.BaseApp.NewContext(false)
	blockHeader := ctx.BlockHeader()
	blockHeader.Time = time.Now()
	blockHeader.Height = 1
	ctx = ctx.WithBlockHeader(blockHeader).WithBlockTime(blockHeader.Time).WithBlockHeight(1)
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
	valDelAddr := sdk.AccAddress(valPubKey.Address())

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
			suite.FundAccount(valDelAddr, sdk.NewCoins(userBalance))

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
				DelegatorAddress:  valDelAddr.String(),
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
	valPubKey := test.GenPubKey()
	valAddr := sdk.ValAddress(valPubKey.Address())
	delAddr := sdk.AccAddress(valPubKey.Address())

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
				DelegatorAddress:  delAddr.String(),
				ValidatorAddress:  valAddr.String(),
				Pubkey:            codectypes.UnsafePackAny(valPubKey),
				Value:             bondAmount,
			}
			_, err = suite.msgServer.CreateValidator(suite.ctx, &createMsg)
			suite.Require().NoError(err)
			_, err = tc.malleate(suite.ctx, suite.msgServer, *suite.msKeeper)
			suite.Require().NoError(err)

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
