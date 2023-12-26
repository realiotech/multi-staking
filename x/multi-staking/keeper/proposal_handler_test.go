package keeper_test

import (
	// sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (suite *KeeperTestSuite) TestHandleAddBondDenomProposal() {
	tests := []struct {
		name    string
		p       *types.AddMultiStakingCoinProposal
		wantErr bool
	}{
		{},
	}
	for _, test := range tests {
		suite.Run(test.name, func() {
			suite.SetupTest()
			err := keeper.HandlerAddMultiStakingCoinProposal(suite.ctx, suite.msKeeper, test.p)
			if test.wantErr {
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)
		})
	}
}

func (suite *KeeperTestSuite) TestHandleChangeBondDenomProposal() {
	tests := []struct {
		name    string
		p       *types.UpdateBondWeightProposal
		wantErr bool
	}{
		{},
	}
	for _, test := range tests {
		suite.Run(test.name, func() {
			suite.SetupTest()
			err := keeper.HandlerUpdateBondWeightProposals(suite.ctx, suite.msKeeper, test.p)
			if test.wantErr {
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err)
		})
	}
}

func (suite *KeeperTestSuite) TestHandleRemoveBondDenomProposal() {
	tests := []struct {
		name    string
		p       *types.RemoveMultiStakingCoinProposal
		wantErr bool
	}{
		{},
	}
	for _, test := range tests {
		suite.Run(test.name, func() {
			suite.SetupTest()
			keeper.HandlerRemoveMultiStakingCoinProposal(suite.ctx, suite.msKeeper, test.p)

			// get and check

		})
	}
}
