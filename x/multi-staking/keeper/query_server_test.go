package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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
