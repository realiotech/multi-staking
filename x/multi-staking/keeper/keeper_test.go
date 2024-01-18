package keeper_test

import (
	"testing"

	"github.com/realio-tech/multi-staking-module/testing/simapp"
	multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app      *simapp.SimApp
	ctx      sdk.Context
	msKeeper *multistakingkeeper.Keeper
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	suite.app, suite.ctx, suite.msKeeper = app, ctx, &app.MultiStakingKeeper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
