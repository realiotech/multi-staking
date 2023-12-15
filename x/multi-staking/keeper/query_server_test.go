package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/testutil"
	multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	multistakingtypes "github.com/realio-tech/multi-staking-module/x/multi-staking/types"
	// other necessary imports
)

func (suite *KeeperTestSuite) TestBondTokenWeightQuery() {
	testCases := []struct {
		name           string
		tokenDenom     string
		setupFunc      func(ctx sdk.Context, k *multistakingkeeper.Keeper)
		expectedWeight sdk.Dec
		isSet          bool
	}{
		{
			name:       "existing token weight",
			tokenDenom: "ario",
			setupFunc: func(ctx sdk.Context, k *multistakingkeeper.Keeper) {
				k.SetBondTokenWeight(ctx, "ario", sdk.MustNewDecFromStr("0.3"))
			},
			expectedWeight: sdk.MustNewDecFromStr("0.3"),
			isSet:          true,
		},
		{
			name:           "non-existing token weight",
			tokenDenom:     "nexit",
			setupFunc:      nil, // no setup required
			expectedWeight: sdk.Dec{},
			isSet:          false,
		},
		// other test cases
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		suite.Run(tc.name, func() {
			suite.SetupTest() // setup the test environment

			ctx := sdk.WrapSDKContext(suite.ctx)

			if tc.setupFunc != nil {
				tc.setupFunc(suite.ctx, suite.msKeeper)
			}

			queryServer := multistakingkeeper.NewQueryServerImpl(*suite.msKeeper)

			response, err := queryServer.BondTokenWeight(ctx, &multistakingtypes.QueryBondTokenWeightRequest{
				TokenDenom: tc.tokenDenom,
			})

			suite.Require().NoError(err)
			suite.Require().NotNil(response)
			suite.Require().Equal(tc.expectedWeight, response.Weight)
			suite.Require().Equal(tc.isSet, response.IsSet)
		})
	}
}

func (suite *KeeperTestSuite) TestValidatorAllowedTokenQuery() {
	testCases := []struct {
		name            string
		operatorAddress sdk.ValAddress
		setupFunc       func(ctx sdk.Context, k *multistakingkeeper.Keeper, oa sdk.ValAddress)
		expectedToken   string
	}{
		{
			name:            "existing validator allowed token",
			operatorAddress: testutil.GenValAddress(),
			setupFunc: func(ctx sdk.Context, k *multistakingkeeper.Keeper, oa sdk.ValAddress) {
				k.SetValidatorAllowedToken(ctx, oa, "ario")
			},
			expectedToken: "ario",
		},
		{
			name:            "not allowed token",
			operatorAddress: testutil.GenValAddress(),
			setupFunc:       nil, // no setup required
			expectedToken:   "",
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		suite.Run(tc.name, func() {
			suite.SetupTest() // setup the test environment

			ctx := sdk.WrapSDKContext(suite.ctx)

			if tc.setupFunc != nil {
				tc.setupFunc(suite.ctx, suite.msKeeper, tc.operatorAddress)
			}

			queryServer := multistakingkeeper.NewQueryServerImpl(*suite.msKeeper)

			response, err := queryServer.ValidatorAllowedToken(ctx, &multistakingtypes.QueryValidatorAllowedTokenRequest{
				OperatorAddress: tc.operatorAddress.String(), // bench32 format
			})

			suite.Require().NoError(err)
			suite.Require().NotNil(response)
			suite.Require().Equal(tc.expectedToken, response.Denom)
		})
	}
}
