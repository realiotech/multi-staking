package keeper_test

import (
	"testing"

	"github.com/realio-tech/multi-staking-module/testing/simapp"
	multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

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
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	multiStakingMsgServer := multistakingkeeper.NewMsgServerImpl(app.MultiStakingKeeper)

	suite.app, suite.ctx, suite.msKeeper, suite.msgServer = app, ctx, &app.MultiStakingKeeper, multiStakingMsgServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// Todo: add CheckBalance; AddAccountWithCoin; FundAccount
func (suite *KeeperTestSuite) FundAccount(addr sdk.AccAddress, amounts sdk.Coins) error {
	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, amounts)
	if err != nil {
		return err
	}
	return suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr, amounts)
}
