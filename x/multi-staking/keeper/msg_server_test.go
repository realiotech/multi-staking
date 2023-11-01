package keeper_test

import (
	"time"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/realio-tech/multi-staking-module/testutil"
	multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	multistakingtypes "github.com/realio-tech/multi-staking-module/x/multi-staking/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func (suite *KeeperTestSuite) TestCreateValidator() {
	delAddr := testutil.GenAddress()
	valPubKey := testutil.GenPubKey()
	valAddr := sdk.ValAddress(valPubKey.Address())
	gasDenom := "ario"
	govDenom := "arst"
	testCases := []struct {
		name     string
		malleate func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper, msgServer multistakingtypes.MsgServer) (sdk.Coin, error)
		expOut   sdk.Coin
		expErr   bool
	}{
		{
			name: "3001 token, weight 0.3, expect 900",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper, msgServer multistakingtypes.MsgServer) (sdk.Coin, error) {
				msKeeper.SetBondTokenWeight(ctx, gasDenom, math.LegacyMustNewDecFromStr("0.3"))
				bondAmount := sdk.NewCoin(gasDenom, sdk.NewInt(3001))

				msg := multistakingtypes.MsgCreateValidator{
					Description: stakingtypes.Description{
						Moniker:         "test",
						Identity:        "test",
						Website:         "test",
						SecurityContact: "test",
						Details:         "test",
					},
					Commission: stakingtypes.CommissionRates{
						Rate:          sdk.MustNewDecFromStr("0.05"),
						MaxRate:       sdk.MustNewDecFromStr("0.1"),
						MaxChangeRate: sdk.MustNewDecFromStr("0.05"),
					},
					MinSelfDelegation: sdk.NewInt(1),
					DelegatorAddress:  delAddr.String(),
					ValidatorAddress:  valAddr.String(),
					Pubkey:            codectypes.UnsafePackAny(valPubKey),
					Value:             bondAmount,
				}

				_, err := msgServer.CreateValidator(ctx, &msg)
				return bondAmount, err
			},
			expOut: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(900)),
			expErr: false,
		},
		{
			name: "25 token, weight 0.5, expect 12",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper, msgServer multistakingtypes.MsgServer) (sdk.Coin, error) {
				msKeeper.SetBondTokenWeight(ctx, govDenom, math.LegacyMustNewDecFromStr("0.5"))
				bondAmount := sdk.NewCoin(govDenom, sdk.NewInt(25))

				msg := multistakingtypes.MsgCreateValidator{
					Description: stakingtypes.Description{
						Moniker:         "test",
						Identity:        "test",
						Website:         "test",
						SecurityContact: "test",
						Details:         "test",
					},
					Commission: stakingtypes.CommissionRates{
						Rate:          sdk.MustNewDecFromStr("0.05"),
						MaxRate:       sdk.MustNewDecFromStr("0.1"),
						MaxChangeRate: sdk.MustNewDecFromStr("0.05"),
					},
					MinSelfDelegation: sdk.NewInt(1),
					DelegatorAddress:  delAddr.String(),
					ValidatorAddress:  valAddr.String(),
					Pubkey:            codectypes.UnsafePackAny(valPubKey),
					Value:             bondAmount,
				}

				_, err := msgServer.CreateValidator(ctx, &msg)
				return bondAmount, err
			},
			expOut: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(12)),
			expErr: false,
		},
		{
			name: "invalid bond token",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper, msgServer multistakingtypes.MsgServer) (sdk.Coin, error) {
				bondAmount := sdk.NewCoin(gasDenom, sdk.NewInt(25))
				msg := multistakingtypes.MsgCreateValidator{
					Description: stakingtypes.Description{
						Moniker:         "test",
						Identity:        "test",
						Website:         "test",
						SecurityContact: "test",
						Details:         "test",
					},
					Commission: stakingtypes.CommissionRates{
						Rate:          sdk.MustNewDecFromStr("0.05"),
						MaxRate:       sdk.MustNewDecFromStr("0.1"),
						MaxChangeRate: sdk.MustNewDecFromStr("0.05"),
					},
					MinSelfDelegation: sdk.NewInt(1),
					DelegatorAddress:  delAddr.String(),
					ValidatorAddress:  valAddr.String(),
					Pubkey:            codectypes.UnsafePackAny(valPubKey),
					Value:             bondAmount,
				}
				_, err := msgServer.CreateValidator(ctx, &msg)
				return bondAmount, err
			},
			expOut: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(12)),
			expErr: true,
		},
		{
			name: "invalid validator address",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper, msgServer multistakingtypes.MsgServer) (sdk.Coin, error) {
				msKeeper.SetBondTokenWeight(ctx, gasDenom, math.LegacyMustNewDecFromStr("0.3"))
				bondAmount := sdk.NewCoin(gasDenom, sdk.NewInt(3001))

				msg := multistakingtypes.MsgCreateValidator{
					Description: stakingtypes.Description{
						Moniker: "NewValidator",
					},
					Commission: stakingtypes.CommissionRates{
						Rate:          sdk.MustNewDecFromStr("0.05"),
						MaxRate:       sdk.MustNewDecFromStr("0.1"),
						MaxChangeRate: sdk.MustNewDecFromStr("0.05"),
					},
					MinSelfDelegation: sdk.NewInt(1),
					DelegatorAddress:  delAddr.String(),
					ValidatorAddress:  sdk.AccAddress([]byte("invalid")).String(),
					Pubkey:            codectypes.UnsafePackAny(valPubKey),
					Value:             bondAmount,
				}

				_, err := msgServer.CreateValidator(ctx, &msg)
				return bondAmount, err
			},
			expOut: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(12)),
			expErr: true,
		},
		{
			name: "nil delegation amount",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper, msgServer multistakingtypes.MsgServer) (sdk.Coin, error) {
				msKeeper.SetBondTokenWeight(ctx, gasDenom, math.LegacyMustNewDecFromStr("0.3"))

				msg := multistakingtypes.MsgCreateValidator{
					Description: stakingtypes.Description{
						Moniker: "NewValidator",
					},
					Commission: stakingtypes.CommissionRates{
						Rate:          sdk.MustNewDecFromStr("0.05"),
						MaxRate:       sdk.MustNewDecFromStr("0.1"),
						MaxChangeRate: sdk.MustNewDecFromStr("0.05"),
					},
					MinSelfDelegation: sdk.NewInt(1),
					DelegatorAddress:  delAddr.String(),
					ValidatorAddress:  valAddr.String(),
					Pubkey:            codectypes.UnsafePackAny(valPubKey),
					Value:             sdk.Coin{},
				}

				_, err := msgServer.CreateValidator(ctx, &msg)
				return sdk.Coin{}, err
			},
			expOut: sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0)),
			expErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			bondAmount, err := tc.malleate(suite.ctx, suite.msKeeper, multistakingkeeper.NewMsgServerImpl(*suite.msKeeper))

			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				actualBond := suite.msKeeper.GetDVPairBondAmount(suite.ctx, delAddr, valAddr)
				actualSDKBond := suite.msKeeper.GetDVPairSDKBondAmount(suite.ctx, delAddr, valAddr)
				suite.Require().Equal(bondAmount.Amount, actualBond)
				suite.Require().Equal(tc.expOut.Amount, actualSDKBond)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestEditValidator() {
	delAddr := testutil.GenAddress()
	valPubKey := testutil.GenPubKey()
	valAddr := sdk.ValAddress(valPubKey.Address())
	gasDenom := "ario"

	testCases := []struct {
		name     string
		malleate func(ctx sdk.Context, msgServer multistakingtypes.MsgServer) (multistakingtypes.MsgEditValidator, error)
		expErr   bool
	}{
		{
			name: "success",
			malleate: func(ctx sdk.Context, msgServer multistakingtypes.MsgServer) (multistakingtypes.MsgEditValidator, error) {
				newRate := sdk.MustNewDecFromStr("0.03")
				newMinSelfDelegation := sdk.NewInt(300)
				editMsg := multistakingtypes.NewMsgEditValidator(valAddr, stakingtypes.Description{
					Moniker:         "test 1",
					Identity:        "test 1",
					Website:         "test 1",
					SecurityContact: "test 1",
					Details:         "test 1",
				},
					&newRate,
					&newMinSelfDelegation,
				)
				_, err := msgServer.EditValidator(ctx, editMsg)
				return *editMsg, err
			},
			expErr: false,
		},
		{
			name: "not found validator",
			malleate: func(ctx sdk.Context, msgServer multistakingtypes.MsgServer) (multistakingtypes.MsgEditValidator, error) {
				newRate := sdk.MustNewDecFromStr("0.03")
				newMinSelfDelegation := sdk.NewInt(300)
				editMsg := multistakingtypes.NewMsgEditValidator(testutil.GenValAddress(), stakingtypes.Description{
					Moniker:         "test",
					Identity:        "test",
					Website:         "test",
					SecurityContact: "test",
					Details:         "test",
				},
					&newRate,
					&newMinSelfDelegation,
				)
				_, err := msgServer.EditValidator(ctx, editMsg)
				return *editMsg, err
			},
			expErr: true,
		},
		{
			name: "negative rate",
			malleate: func(ctx sdk.Context, msgServer multistakingtypes.MsgServer) (multistakingtypes.MsgEditValidator, error) {
				newRate := sdk.MustNewDecFromStr("-0.01")
				newMinSelfDelegation := sdk.NewInt(300)
				editMsg := multistakingtypes.NewMsgEditValidator(valAddr, stakingtypes.Description{
					Moniker:         "test 1",
					Identity:        "test 1",
					Website:         "test 1",
					SecurityContact: "test 1",
					Details:         "test 1",
				},
					&newRate,
					&newMinSelfDelegation,
				)
				_, err := msgServer.EditValidator(ctx, editMsg)
				return *editMsg, err
			},
			expErr: true,
		},
		{
			name: "less than minimum rate",
			malleate: func(ctx sdk.Context, msgServer multistakingtypes.MsgServer) (multistakingtypes.MsgEditValidator, error) {
				newRate := sdk.MustNewDecFromStr("0.01")
				newMinSelfDelegation := sdk.NewInt(300)
				editMsg := multistakingtypes.NewMsgEditValidator(valAddr, stakingtypes.Description{
					Moniker:         "test 1",
					Identity:        "test 1",
					Website:         "test 1",
					SecurityContact: "test 1",
					Details:         "test 1",
				},
					&newRate,
					&newMinSelfDelegation,
				)
				_, err := msgServer.EditValidator(ctx, editMsg)
				return *editMsg, err
			},
			expErr: true,
		},
		{
			name: "more than max rate",
			malleate: func(ctx sdk.Context, msgServer multistakingtypes.MsgServer) (multistakingtypes.MsgEditValidator, error) {
				newRate := sdk.MustNewDecFromStr("0.11")
				newMinSelfDelegation := sdk.NewInt(300)
				editMsg := multistakingtypes.NewMsgEditValidator(valAddr, stakingtypes.Description{
					Moniker:         "test 1",
					Identity:        "test 1",
					Website:         "test 1",
					SecurityContact: "test 1",
					Details:         "test 1",
				},
					&newRate,
					&newMinSelfDelegation,
				)
				_, err := msgServer.EditValidator(ctx, editMsg)
				return *editMsg, err
			},
			expErr: true,
		},
		{
			name: "min self delegation more than validator tokens",
			malleate: func(ctx sdk.Context, msgServer multistakingtypes.MsgServer) (multistakingtypes.MsgEditValidator, error) {
				newRate := sdk.MustNewDecFromStr("0.03")
				newMinSelfDelegation := sdk.NewInt(10000)
				editMsg := multistakingtypes.NewMsgEditValidator(valAddr, stakingtypes.Description{
					Moniker:         "test 1",
					Identity:        "test 1",
					Website:         "test 1",
					SecurityContact: "test 1",
					Details:         "test 1",
				},
					&newRate,
					&newMinSelfDelegation,
				)
				_, err := msgServer.EditValidator(ctx, editMsg)
				return *editMsg, err
			},
			expErr: true,
		},
		{
			name: "min self delegation more than old min delegation value",
			malleate: func(ctx sdk.Context, msgServer multistakingtypes.MsgServer) (multistakingtypes.MsgEditValidator, error) {
				newRate := sdk.MustNewDecFromStr("0.03")
				newMinSelfDelegation := sdk.NewInt(100)
				editMsg := multistakingtypes.NewMsgEditValidator(valAddr, stakingtypes.Description{
					Moniker:         "test 1",
					Identity:        "test 1",
					Website:         "test 1",
					SecurityContact: "test 1",
					Details:         "test 1",
				},
					&newRate,
					&newMinSelfDelegation,
				)
				_, err := msgServer.EditValidator(ctx, editMsg)
				return *editMsg, err
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			newParam := stakingtypes.DefaultParams()
			newParam.MinCommissionRate = sdk.MustNewDecFromStr("0.02")
			suite.stakingKeeper.SetParams(suite.ctx, newParam)
			msgServer := multistakingkeeper.NewMsgServerImpl(*suite.msKeeper)
			suite.msKeeper.SetBondTokenWeight(suite.ctx, gasDenom, math.LegacyOneDec())
			bondAmount := sdk.NewCoin(gasDenom, sdk.NewInt(1000))

			createMsg := multistakingtypes.MsgCreateValidator{
				Description: stakingtypes.Description{
					Moniker:         "test",
					Identity:        "test",
					Website:         "test",
					SecurityContact: "test",
					Details:         "test",
				},
				Commission: stakingtypes.CommissionRates{
					Rate:          sdk.MustNewDecFromStr("0.05"),
					MaxRate:       sdk.MustNewDecFromStr("0.1"),
					MaxChangeRate: sdk.MustNewDecFromStr("0.05"),
				},
				MinSelfDelegation: sdk.NewInt(200),
				DelegatorAddress:  delAddr.String(),
				ValidatorAddress:  valAddr.String(),
				Pubkey:            codectypes.UnsafePackAny(valPubKey),
				Value:             bondAmount,
			}
			msgServer.CreateValidator(suite.ctx, &createMsg)

			suite.ctx = suite.ctx.WithBlockHeader(tmproto.Header{Time: time.Now()})
			originMsg, err := tc.malleate(suite.ctx, msgServer)

			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				validatorInfo, found := suite.stakingKeeper.GetValidator(suite.ctx, sdk.ValAddress(originMsg.ValidatorAddress))
				if found {
					suite.Require().Equal(validatorInfo.Description, originMsg.Description)
					suite.Require().Equal(validatorInfo.MinSelfDelegation, &originMsg.MinSelfDelegation)
					suite.Require().Equal(validatorInfo.Commission.CommissionRates.Rate, &originMsg.CommissionRate)
				}
			}
		})
	}
}
