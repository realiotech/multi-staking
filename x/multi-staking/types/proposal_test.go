package types_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

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
	suite.Require().Equal("multistaking", (&types.AddBondTokenProposal{}).ProposalRoute())
	suite.Require().Equal("AddBondToken", (&types.AddBondTokenProposal{}).ProposalType())
	suite.Require().Equal("multistaking", (&types.ChangeBondTokenWeightProposal{}).ProposalRoute())
	suite.Require().Equal("ChangeBondTokenWeight", (&types.ChangeBondTokenWeightProposal{}).ProposalType())
}

func (suite *ProposalTestSuite) TestProposalString() {
	testTokenWeight := math.LegacyNewDec(1)
    testCases := []struct {
		msg string
		proposal govv1beta1.Content
		expectedValue string
	} {
		{msg: "Add Bond Token Proposal", proposal: &types.AddBondTokenProposal{BondToken: "token", TokenWeight: &testTokenWeight, Description: "Add token", Title: "Add #1"}, 
		expectedValue: "AddBondTokenProposal: Title: Add #1 Description: Add token BondToken: token TokenWeight: 1.000000000000000000" },
	
		{msg: "Change Bond Token Weight Proposal", proposal: &types.ChangeBondTokenWeightProposal{BondToken: "token", TokenWeight: &testTokenWeight, Description: "Change Bond token weight", Title: "Change #2"}, 
		expectedValue: "ChangeBondTokenWeightProposal: Title: Change #2 Description: Change Bond token weight BondToken: token TokenWeight: 1.000000000000000000" },
	}

	for _, tc := range testCases {
		str_result := tc.proposal.String()
		suite.Require().Equal(str_result, tc.expectedValue)
	}
}


func (suite *ProposalTestSuite) TestAddBondTokenProposal() {
	testCases := []struct {
		msg         string
		title       string
		description string
		bondToken   string
		tokenWeight math.LegacyDec
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
		tx := types.NewAddBondTokenProposal(tc.title, tc.description, tc.bondToken, tc.tokenWeight)
		err := tx.ValidateBasic()

		if tc.expectPass {
			suite.Require().NoError(err, "valid test %d failed: %s", i, tc.msg)
		} else {
			suite.Require().Error(err, "invalid test %d passed: %s", i, tc.msg)
		}
	}
}

func (suite *ProposalTestSuite) TestChangeBondTokenWeightProposal() {
	testCases := []struct {
		msg         string
		title       string
		description string
		bondToken   string
		tokenWeight math.LegacyDec
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
		tx := types.NewChangeBondTokenWeightProposal(tc.title, tc.description, tc.bondToken, tc.tokenWeight)
		err := tx.ValidateBasic()

		if tc.expectPass {
			suite.Require().NoError(err, "valid test %d failed: %s", i, tc.msg)
		} else {
			suite.Require().Error(err, "invalid test %d passed: %s", i, tc.msg)
		}
	}
}