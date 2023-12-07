package keeper_test

import (
	// sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (suite *KeeperTestSuite) TestHandleAddBondDenomProposal() {
	tests := []struct {
		name    string
		p       *types.AddBondDenomProposal
		wantErr bool
	}{
		{},
	}
	for _, test := range tests {
		suite.Run(test.name, func() {
			suite.SetupTest()
			err := keeper.HandlerAddBondDenomProposal(suite.ctx, suite.msKeeper, test.p)
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
		p       *types.UpdateBondTokenWeightProposals
		wantErr bool
	}{
		{},
	}
	for _, test := range tests {
		suite.Run(test.name, func() {
			suite.SetupTest()
			err := keeper.HandlerUpdateBondTokenWeightProposals(suite.ctx, suite.msKeeper, test.p)
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
		p       *types.RemoveBondTokenProposal
		wantErr bool
	}{
		{},
	}
	for _, test := range tests {
		suite.Run(test.name, func() {
			suite.SetupTest()
			keeper.HandlerRemoveBondTokenProposal(suite.ctx, suite.msKeeper, test.p)

			// get and check

		})
	}
}
