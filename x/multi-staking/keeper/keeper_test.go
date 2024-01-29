package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/realio-tech/multi-staking-module/test"
	"github.com/realio-tech/multi-staking-module/test/simapp"
	multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app       *simapp.SimApp
	ctx       sdk.Context
	msKeeper  *multistakingkeeper.Keeper
	msgServer stakingtypes.MsgServer
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1})
	multiStakingMsgServer := multistakingkeeper.NewMsgServerImpl(app.MultiStakingKeeper)

	suite.app, suite.ctx, suite.msKeeper, suite.msgServer = app, ctx, &app.MultiStakingKeeper, multiStakingMsgServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) NextBlock(jumpTime time.Duration) {
	app := suite.app

	app.EndBlocker(suite.ctx, abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

	newBlockTime := suite.ctx.BlockTime().Add(jumpTime)
	nextHeight := suite.ctx.BlockHeight() + 1

	newCtx := suite.ctx.WithBlockTime(newBlockTime).WithBlockHeight(nextHeight)

	app.BeginBlocker(newCtx, abci.RequestBeginBlock{Header: newCtx.BlockHeader()})

	suite.ctx = app.NewContext(false, newCtx.BlockHeader())
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
