package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/stretchr/testify/suite"

	"cosmossdk.io/math"

	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

type ProposalTestSuite struct {
	suite.Suite
}

func TestProposalTestSuite(t *testing.T) {
	suite.Run(t, new(ProposalTestSuite))
}

func (suite *ProposalTestSuite) TestKeysTypes() {
	suite.Require().Equal("multistaking", (&types.AddMultiStakingCoinProposal{}).ProposalRoute())
	suite.Require().Equal("AddMultiStakingCoin", (&types.AddMultiStakingCoinProposal{}).ProposalType())
	suite.Require().Equal("multistaking", (&types.UpdateBondWeightProposal{}).ProposalRoute())
	suite.Require().Equal("UpdateBondWeight", (&types.UpdateBondWeightProposal{}).ProposalType())
}

func (suite *ProposalTestSuite) TestProposalString() {
	testTokenWeight := math.LegacyNewDec(1)
	testCases := []struct {
		msg           string
		proposal      govv1beta1.Content
		expectedValue string
	}{
		{msg: "Add Bond Token Proposal", proposal: &types.AddMultiStakingCoinProposal{BondToken: "token", TokenWeight: &testTokenWeight, Description: "Add token", Title: "Add #1"},
			expectedValue: "AddMultiStakingCoinProposal: Title: Add #1 Description: Add token BondToken: token TokenWeight: 1.000000000000000000"},

		{msg: "Change Bond Token Weight Proposal", proposal: &types.UpdateBondWeightProposal{BondToken: "token", TokenWeight: &testTokenWeight, Description: "Change Bond token weight", Title: "Change #2"},
			expectedValue: "UpdateBondWeightProposal: Title: Change #2 Description: Change Bond token weight BondToken: token TokenWeight: 1.000000000000000000"},
	}

	for _, tc := range testCases {
		str_result := tc.proposal.String()
		suite.Require().Equal(str_result, tc.expectedValue)
	}
}

func (suite *ProposalTestSuite) TestAddMultiStakingCoinProposal() {
	testCases := []struct {
		msg         string
		title       string
		description string
		bondToken   string
		tokenWeight sdk.Dec
		expectPass  bool
	}{
		// Valid tests
		{msg: "Add bond token", title: "test", description: "test desc", bondToken: "token", tokenWeight: math.LegacyNewDec(1), expectPass: true},

		// Invalid tests
		{msg: "Add bond token - invalid token", title: "test", description: "test desc", bondToken: "", tokenWeight: math.LegacyNewDec(1), expectPass: false},
		{msg: "Add bond token - negative weight", title: "test", description: "test desc", bondToken: "token", tokenWeight: math.LegacyNewDec(-1), expectPass: false},
		{msg: "Add bond token - zero weight", title: "test", description: "test desc", bondToken: "token", tokenWeight: math.LegacyNewDec(0), expectPass: false},
	}

	for i, tc := range testCases {
		tx := types.NewAddMultiStakingCoinProposal(tc.title, tc.description, tc.Denom, tc.BondWeight)
		err := tx.ValidateBasic()

		if tc.expectPass {
			suite.Require().NoError(err, "valid test %d failed: %s", i, tc.msg)
		} else {
			suite.Require().Error(err, "invalid test %d passed: %s", i, tc.msg)
		}
	}
}

func (suite *ProposalTestSuite) TestUpdateBondWeightProposal() {
	testCases := []struct {
		msg         string
		title       string
		description string
		bondToken   string
		tokenWeight sdk.Dec
		expectPass  bool
	}{
		// Valid tests
		{msg: "Change bond token weight", title: "test", description: "test desc", bondToken: "token", tokenWeight: math.LegacyNewDec(1), expectPass: true},

		// Invalid tests
		{msg: "Change bond token weight - invalid token", title: "test", description: "test desc", bondToken: "", tokenWeight: math.LegacyNewDec(1), expectPass: false},
		{msg: "Change bond token weight - negative weight", title: "test", description: "test desc", bondToken: "token", tokenWeight: math.LegacyNewDec(-1), expectPass: false},
		{msg: "Change bond token weight - zero weight", title: "test", description: "test desc", bondToken: "token", tokenWeight: math.LegacyNewDec(0), expectPass: false},
	}

	for i, tc := range testCases {
		tx := types.NewUpdateBondWeightProposal(tc.title, tc.description, tc.Denom, tc.BondWeight)
		err := tx.ValidateBasic()

		if tc.expectPass {
			suite.Require().NoError(err, "valid test %d failed: %s", i, tc.msg)
		} else {
			suite.Require().Error(err, "invalid test %d passed: %s", i, tc.msg)
		}
	}
}
