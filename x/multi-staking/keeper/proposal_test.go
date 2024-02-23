package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (suite *KeeperTestSuite) TestAddHostZoneProposal() {
	bondWeight := sdk.NewDec(1)

	for _, tc := range []struct {
		desc      string
		malleate  func(p *types.AddMultiStakingCoinProposal)
		proposal  *types.AddMultiStakingCoinProposal
		shouldErr bool
	}{
		{
			desc: "Success",
			malleate: func(p *types.AddMultiStakingCoinProposal) {
				_, found := suite.msKeeper.GetBondWeight(suite.ctx, p.Denom)
				suite.Require().False(found)
			},
			proposal: &types.AddMultiStakingCoinProposal{
				Title:       "Add multistaking coin",
				Description: "Add new multistaking coin",
				Denom:       "stake",
				BondWeight:  &bondWeight,
			},
			shouldErr: false,
		},
		{
			desc: "Error multistaking coin already exists",
			malleate: func(p *types.AddMultiStakingCoinProposal) {
				suite.msKeeper.SetBondWeight(suite.ctx, p.Denom, *p.BondWeight)
			},
			proposal: &types.AddMultiStakingCoinProposal{
				Title:       "Add multistaking coin",
				Description: "Add new multistaking coin",
				Denom:       "stake",
				BondWeight:  &bondWeight,
			},
			shouldErr: false,
		},
	} {
		tc := tc
		suite.Run(tc.desc, func() {
			suite.SetupTest()
			tc.malleate(tc.proposal)

			legacyProposal, err := govv1types.NewLegacyContent(tc.proposal, authtypes.NewModuleAddress(govtypes.ModuleName).String())
			suite.Require().NoError(err)

			if !tc.shouldErr {
				// store proposal
				_, err = suite.govKeeper.SubmitProposal(suite.ctx, []sdk.Msg{legacyProposal}, "")
				suite.Require().NoError(err)

				// execute proposal
				handler := suite.govKeeper.LegacyRouter().GetRoute(tc.proposal.ProposalRoute())
				err = handler(suite.ctx, tc.proposal)
				suite.Require().NoError(err)

				_, found := suite.msKeeper.GetBondWeight(suite.ctx, tc.proposal.Denom)
				suite.Require().True(found)
			} else {
				// store proposal
				_, err = suite.govKeeper.SubmitProposal(suite.ctx, []sdk.Msg{legacyProposal}, "")
				suite.Require().Error(err)
			}
		})
	}
}
