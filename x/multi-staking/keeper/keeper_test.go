package keeper_test

import (
	"testing"

	"github.com/realio-tech/multi-staking-module/testing/simapp"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx       sdk.Context
	app       *simapp.SimApp
	msKeeper  *multistakingkeeper.Keeper
	msgServer stakingtypes.MsgServer
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	suite.app = app
	suite.ctx, suite.msKeeper = ctx, &app.MultiStakingKeeper
	suite.msgServer = keeper.NewMsgServerImpl(*suite.msKeeper)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
