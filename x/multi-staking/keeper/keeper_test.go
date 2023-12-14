package keeper_test

import (
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/realio-tech/multi-staking-module/testing/simapp"
	multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx           sdk.Context
	app           simapp.SimApp
	msKeeper      *multistakingkeeper.Keeper
	stakingKeeper *stakingkeeper.Keeper
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false,
		tmproto.Header{
			Height:  2,
			ChainID: "multi-staking-1",
			Time:    time.Unix(1, 1).UTC(),
		},
	)

	suite.app, suite.ctx, suite.msKeeper, suite.stakingKeeper = *app, ctx, &app.MultiStakingKeeper, &app.StakingKeeper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// CommitAndBeginBlock commits the current state and begins the next block (without committing).
func (suite *KeeperTestSuite) CommitAndBeginBlock() {
	suite.CommitAndBeginBlocks(1)
}

// CommitAndBeginBlocks commits the current state and begins subsequent blocks.
func (suite *KeeperTestSuite) CommitAndBeginBlocks(numBlocks int64) {
	for i := int64(0); i < numBlocks; i++ {
		_ = suite.app.Commit()
		header := suite.ctx.BlockHeader()
		header.Height += 1
		header.Time = header.Time.Add(6 * time.Second)
		suite.app.BeginBlock(abci.RequestBeginBlock{
			Header: header,
		})

		// update ctx
		suite.ctx = suite.app.BaseApp.NewContext(false, header)
	}
}
