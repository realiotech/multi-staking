package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/testutil"
	multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	multistakingtypes "github.com/realio-tech/multi-staking-module/x/multi-staking/types"
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

func (suite *KeeperTestSuite) TestQueryMultiStakingLock() {
	multiStakingLock := &multistakingtypes.MultiStakingLock{
		ConversionRatio: sdk.MustNewDecFromStr("0.3"),
		LockedAmount:    sdk.NewInt(10000),
		DelAddr:         testutil.GenAddress().String(),
		ValAddr:         testutil.GenValAddress().String(),
	}

	testCases := []struct {
		name               string
		multiStakingLockID []byte
		setupFunc          func(ctx sdk.Context, k *multistakingkeeper.Keeper)
		expectedFound      bool
		expectedLock       *multistakingtypes.MultiStakingLock
	}{
		{
			name:               "existing multi staking lock",
			multiStakingLockID: []byte("lock_id_1"),
			setupFunc: func(ctx sdk.Context, k *multistakingkeeper.Keeper) {
				k.SetMultiStakingLock(ctx, []byte("lock_id_1"), *multiStakingLock)
			},
			expectedFound: true,
			expectedLock:  multiStakingLock,
		},
		{
			name:               "non-existing multi staking lock",
			multiStakingLockID: []byte("lock_id_2"),
			setupFunc:          nil, // no setup required
			expectedFound:      false,
			expectedLock: &multistakingtypes.MultiStakingLock{
				ConversionRatio: sdk.Dec{},
				LockedAmount:    sdk.Int{},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest() // setup the test environment

			if tc.setupFunc != nil {
				tc.setupFunc(suite.ctx, suite.msKeeper)
			}

			queryServer := multistakingkeeper.NewQueryServerImpl(*suite.msKeeper)
			response, err := queryServer.MultiStakingLock(sdk.WrapSDKContext(suite.ctx), &multistakingtypes.QueryMultiStakingLockRequest{
				MultiStakingLockId: tc.multiStakingLockID,
			})

			suite.Require().NoError(err)
			suite.Require().NotNil(response)
			suite.Require().Equal(tc.expectedFound, response.Found)
			suite.Require().Equal(tc.expectedLock, response.MultiStakingLock)
		})
	}
}

func (suite *KeeperTestSuite) TestQueryMultiStakingUnlock() {
	delegatorAddress := testutil.GenAddress()
	validatorAddress := testutil.GenValAddress()

	testCases := []struct {
		name           string
		setupFunc      func(ctx sdk.Context, k *multistakingkeeper.Keeper)
		expectedFound  bool
		expectedUnlock *multistakingtypes.MultiStakingUnlock
	}{
		{
			name: "existing multi staking unlock",
			setupFunc: func(ctx sdk.Context, k *multistakingkeeper.Keeper) {
				multiStakingUnlock := &multistakingtypes.MultiStakingUnlock{
					DelegatorAddress: delegatorAddress.String(),
					ValidatorAddress: validatorAddress.String(),
					Entries: []multistakingtypes.UnlockEntry{
						{
							CreationHeight:  1,
							ConversionRatio: sdk.MustNewDecFromStr("0.3"),
							Balance:         sdk.NewInt(10000),
						},
					},
				}

				k.SetMultiStakingUnlock(ctx, *multiStakingUnlock)
			},
			expectedFound: true,
			expectedUnlock: &multistakingtypes.MultiStakingUnlock{
				DelegatorAddress: delegatorAddress.String(),
				ValidatorAddress: validatorAddress.String(),
				Entries: []multistakingtypes.UnlockEntry{
					{
						CreationHeight:  1,
						ConversionRatio: sdk.MustNewDecFromStr("0.3"),
						Balance:         sdk.NewInt(10000),
					},
				},
			},
		},
		{
			name: "existing multi staking unlock",
			setupFunc: func(ctx sdk.Context, k *multistakingkeeper.Keeper) {
				multiStakingUnlock := &multistakingtypes.MultiStakingUnlock{
					DelegatorAddress: delegatorAddress.String(),
					ValidatorAddress: validatorAddress.String(),
					Entries: []multistakingtypes.UnlockEntry{
						{
							CreationHeight:  1,
							ConversionRatio: sdk.MustNewDecFromStr("0.3"),
							Balance:         sdk.NewInt(10000),
						},
					},
				}

				k.SetMultiStakingUnlock(ctx, *multiStakingUnlock)
			},
			expectedFound: true,
			expectedUnlock: &multistakingtypes.MultiStakingUnlock{
				DelegatorAddress: delegatorAddress.String(),
				ValidatorAddress: validatorAddress.String(),
				Entries: []multistakingtypes.UnlockEntry{
					{
						CreationHeight:  1,
						ConversionRatio: sdk.MustNewDecFromStr("0.3"),
						Balance:         sdk.NewInt(10000),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			if tc.setupFunc != nil {
				tc.setupFunc(suite.ctx, suite.msKeeper)
			}

			queryServer := multistakingkeeper.NewQueryServerImpl(*suite.msKeeper)
			response, err := queryServer.MultiStakingUnlock(sdk.WrapSDKContext(suite.ctx), &multistakingtypes.QueryMultiStakingUnlockRequest{
				DelegatorAddress: delegatorAddress.String(),
				ValidatorAddress: validatorAddress.String(),
			})

			fmt.Println("res: ", response)

			suite.Require().NoError(err)
			suite.Require().NotNil(response)
			suite.Require().Equal(tc.expectedFound, response.Found)
			if tc.expectedFound {
				suite.Require().Equal(tc.expectedUnlock, response.MultiStakingUnlock)
			} else {
				suite.Require().Nil(response.MultiStakingUnlock)
			}
		})
	}
}
