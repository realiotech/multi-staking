package keeper_test

import (
	"time"

	"github.com/realio-tech/multi-staking-module/test"
	multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	multistakingtypes "github.com/realio-tech/multi-staking-module/x/multi-staking/types"

	"cosmossdk.io/math"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
)

var (
	MultiStakingDenomA = "ario"
	MultiStakingDenomB = "arst"
)

func (suite *KeeperTestSuite) TestCreateValidator() {
	valPubKey := test.GenPubKey()
	valAddr := sdk.ValAddress(valPubKey.Address())

	testCases := []struct {
		name     string
		malleate func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper, msgServer stakingtypes.MsgServer) (sdk.Coin, error)
		expOut   sdk.Coin
		expErr   bool
	}{
		{
			name: "3001 token, weight 0.3, expect 900",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper, msgServer stakingtypes.MsgServer) (sdk.Coin, error) {
				msKeeper.SetBondWeight(ctx, MultiStakingDenomA, math.LegacyMustNewDecFromStr("0.3"))
				bondAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(3001))
				msg := stakingtypes.MsgCreateValidator{
					Description: stakingtypes.Description{
						Moniker:         "test",
						Identity:        "test",
						Website:         "test",
						SecurityContact: "test",
						Details:         "test",
					},
					Commission: stakingtypes.CommissionRates{
						Rate:          math.LegacyMustNewDecFromStr("0.05"),
						MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
						MaxChangeRate: math.LegacyMustNewDecFromStr("0.05"),
					},
					MinSelfDelegation: math.NewInt(1),
					DelegatorAddress:  sdk.AccAddress(valAddr).String(),
					ValidatorAddress:  valAddr.String(),
					Pubkey:            codectypes.UnsafePackAny(valPubKey),
					Value:             bondAmount,
				}

				_, err := msgServer.CreateValidator(ctx, &msg)
				return bondAmount, err
			},
			expOut: sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(900)),
			expErr: false,
		},
		{
			name: "25 token, weight 0.5, expect 12",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper, msgServer stakingtypes.MsgServer) (sdk.Coin, error) {
				msKeeper.SetBondWeight(ctx, MultiStakingDenomB, math.LegacyMustNewDecFromStr("0.5"))
				bondAmount := sdk.NewCoin(MultiStakingDenomB, math.NewInt(25))

				msg := stakingtypes.MsgCreateValidator{
					Description: stakingtypes.Description{
						Moniker:         "test",
						Identity:        "test",
						Website:         "test",
						SecurityContact: "test",
						Details:         "test",
					},
					Commission: stakingtypes.CommissionRates{
						Rate:          math.LegacyMustNewDecFromStr("0.05"),
						MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
						MaxChangeRate: math.LegacyMustNewDecFromStr("0.05"),
					},
					MinSelfDelegation: math.NewInt(1),
					DelegatorAddress:  sdk.AccAddress(valAddr).String(),
					ValidatorAddress:  valAddr.String(),
					Pubkey:            codectypes.UnsafePackAny(valPubKey),
					Value:             bondAmount,
				}

				_, err := msgServer.CreateValidator(ctx, &msg)
				return bondAmount, err
			},
			expOut: sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(12)),
			expErr: false,
		},
		{
			name: "invalid bond token",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper, msgServer stakingtypes.MsgServer) (sdk.Coin, error) {
				msKeeper.RemoveBondWeight(ctx, MultiStakingDenomA)
				bondAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(25))

				msg := stakingtypes.MsgCreateValidator{
					Description: stakingtypes.Description{
						Moniker:         "test",
						Identity:        "test",
						Website:         "test",
						SecurityContact: "test",
						Details:         "test",
					},
					Commission: stakingtypes.CommissionRates{
						Rate:          math.LegacyMustNewDecFromStr("0.05"),
						MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
						MaxChangeRate: math.LegacyMustNewDecFromStr("0.05"),
					},
					MinSelfDelegation: math.NewInt(1),
					DelegatorAddress:  sdk.AccAddress(valAddr).String(),
					ValidatorAddress:  valAddr.String(),
					Pubkey:            codectypes.UnsafePackAny(valPubKey),
					Value:             bondAmount,
				}
				_, err := msgServer.CreateValidator(ctx, &msg)
				return bondAmount, err
			},
			expOut: sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(12)),
			expErr: true,
		},
		{
			name: "invalid validator address",
			malleate: func(ctx sdk.Context, msKeeper *multistakingkeeper.Keeper, msgServer stakingtypes.MsgServer) (sdk.Coin, error) {
				msKeeper.SetBondWeight(ctx, MultiStakingDenomA, math.LegacyMustNewDecFromStr("0.3"))
				bondAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(3001))

				msg := stakingtypes.MsgCreateValidator{
					Description: stakingtypes.Description{
						Moniker: "NewValidator",
					},
					Commission: stakingtypes.CommissionRates{
						Rate:          math.LegacyMustNewDecFromStr("0.05"),
						MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
						MaxChangeRate: math.LegacyMustNewDecFromStr("0.05"),
					},
					MinSelfDelegation: math.NewInt(1),
					DelegatorAddress:  sdk.AccAddress(valAddr).String(),
					ValidatorAddress:  sdk.AccAddress([]byte("invalid")).String(),
					Pubkey:            codectypes.UnsafePackAny(valPubKey),
					Value:             bondAmount,
				}

				_, err := msgServer.CreateValidator(ctx, &msg)
				return bondAmount, err
			},
			expOut: sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(12)),
			expErr: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			valCoins := sdk.NewCoins(sdk.NewCoin(MultiStakingDenomA, math.NewInt(10000)), sdk.NewCoin(MultiStakingDenomB, math.NewInt(10000)))
			suite.FundAccount(sdk.AccAddress(valAddr), valCoins)

			bondAmount, err := tc.malleate(suite.ctx, suite.msKeeper, suite.msgServer)
			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				lockId := multistakingtypes.MultiStakingLockID(sdk.AccAddress(valAddr).String(), valAddr.String())
				lockRecord, found := suite.msKeeper.GetMultiStakingLock(suite.ctx, lockId)
				suite.Require().True(found)
				actualBond, err := suite.app.StakingKeeper.GetDelegation(suite.ctx, sdk.AccAddress(valAddr), valAddr)
				suite.Require().NoError(err)
				suite.Require().Equal(bondAmount.Amount, lockRecord.LockedCoin.Amount)
				suite.Require().Equal(tc.expOut.Amount, actualBond.Shares.TruncateInt())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestEditValidator() {
	valPubKey := test.GenPubKey()
	valAddr := sdk.ValAddress(valPubKey.Address())

	testCases := []struct {
		name     string
		malleate func(ctx sdk.Context, msgServer stakingtypes.MsgServer) (stakingtypes.MsgEditValidator, error)
		expErr   bool
	}{
		{
			name: "success",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer) (stakingtypes.MsgEditValidator, error) {
				newRate := math.LegacyMustNewDecFromStr("0.03")
				newMinSelfDelegation := math.NewInt(300)
				editMsg := stakingtypes.NewMsgEditValidator(valAddr.String(), stakingtypes.Description{
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
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer) (stakingtypes.MsgEditValidator, error) {
				newRate := math.LegacyMustNewDecFromStr("0.03")
				newMinSelfDelegation := math.NewInt(300)
				editMsg := stakingtypes.NewMsgEditValidator(test.GenValAddress().String(), stakingtypes.Description{
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
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer) (stakingtypes.MsgEditValidator, error) {
				newRate := math.LegacyMustNewDecFromStr("-0.01")
				newMinSelfDelegation := math.NewInt(300)
				editMsg := stakingtypes.NewMsgEditValidator(valAddr.String(), stakingtypes.Description{
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
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer) (stakingtypes.MsgEditValidator, error) {
				newRate := math.LegacyMustNewDecFromStr("0.01")
				newMinSelfDelegation := math.NewInt(300)
				editMsg := stakingtypes.NewMsgEditValidator(valAddr.String(), stakingtypes.Description{
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
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer) (stakingtypes.MsgEditValidator, error) {
				newRate := math.LegacyMustNewDecFromStr("0.11")
				newMinSelfDelegation := math.NewInt(300)
				editMsg := stakingtypes.NewMsgEditValidator(valAddr.String(), stakingtypes.Description{
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
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer) (stakingtypes.MsgEditValidator, error) {
				newRate := math.LegacyMustNewDecFromStr("0.03")
				newMinSelfDelegation := math.NewInt(10000)
				editMsg := stakingtypes.NewMsgEditValidator(valAddr.String(), stakingtypes.Description{
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
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer) (stakingtypes.MsgEditValidator, error) {
				newRate := math.LegacyMustNewDecFromStr("0.03")
				newMinSelfDelegation := math.NewInt(100)
				editMsg := stakingtypes.NewMsgEditValidator(valAddr.String(), stakingtypes.Description{
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
			newParam.MinCommissionRate = math.LegacyMustNewDecFromStr("0.02")
			err := suite.app.StakingKeeper.SetParams(suite.ctx, newParam)
			suite.Require().NoError(err)
			suite.msKeeper.SetBondWeight(suite.ctx, MultiStakingDenomA, math.LegacyOneDec())
			bondAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(1000))
			suite.FundAccount(sdk.AccAddress(valAddr), sdk.NewCoins(bondAmount))

			createMsg := stakingtypes.MsgCreateValidator{
				Description: stakingtypes.Description{
					Moniker:         "test",
					Identity:        "test",
					Website:         "test",
					SecurityContact: "test",
					Details:         "test",
				},
				Commission: stakingtypes.CommissionRates{
					Rate:          math.LegacyMustNewDecFromStr("0.05"),
					MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
					MaxChangeRate: math.LegacyMustNewDecFromStr("0.05"),
				},
				MinSelfDelegation: math.NewInt(200),
				DelegatorAddress:  sdk.AccAddress(valAddr).String(),
				ValidatorAddress:  valAddr.String(),
				Pubkey:            codectypes.UnsafePackAny(valPubKey),
				Value:             bondAmount,
			}
			_, err = suite.msgServer.CreateValidator(suite.ctx, &createMsg)
			suite.Require().NoError(err)

			suite.ctx = suite.ctx.WithBlockHeader(tmproto.Header{Time: time.Now()})
			originMsg, err := tc.malleate(suite.ctx, suite.msgServer)

			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				valCodec := authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix())
				msgValAddr, err := valCodec.StringToBytes(originMsg.ValidatorAddress)
				suite.Require().NoError(err)
				validatorInfo, err := suite.app.StakingKeeper.GetValidator(suite.ctx, msgValAddr)
				if err != nil {
					suite.Require().Equal(validatorInfo.MinSelfDelegation, &originMsg.MinSelfDelegation)
					suite.Require().Equal(validatorInfo.Commission.CommissionRates.Rate, &originMsg.CommissionRate)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestDelegate() {
	delAddr := test.GenAddress()
	valPubKey := test.GenPubKey()
	valAddr := sdk.ValAddress(valPubKey.Address())

	testCases := []struct {
		name     string
		malleate func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (sdk.Coin, error)
		expRate  math.LegacyDec
		expErr   bool
	}{
		{
			name: "success and not change rate",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (sdk.Coin, error) {
				multiStakingAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(1000))
				delMsg := stakingtypes.NewMsgDelegate(delAddr.String(), valAddr.String(), multiStakingAmount)
				_, err := msgServer.Delegate(ctx, delMsg)
				return multiStakingAmount, err
			},
			expRate: math.LegacyOneDec(),
			expErr:  false,
		},
		{
			name: "rate change from 1 to 0.75 (1000 * 1 + 3000 * 0.5 = 4000 * 0.625)",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (sdk.Coin, error) {
				multiStakingAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(1000))
				delMsg := stakingtypes.NewMsgDelegate(delAddr.String(), valAddr.String(), multiStakingAmount)
				_, err := msgServer.Delegate(ctx, delMsg)
				if err != nil {
					return multiStakingAmount, err
				}
				msKeeper.SetBondWeight(ctx, MultiStakingDenomA, math.LegacyMustNewDecFromStr("0.5"))
				multiStakingAmount1 := sdk.NewCoin(MultiStakingDenomA, math.NewInt(3000))
				delMsg1 := stakingtypes.NewMsgDelegate(delAddr.String(), valAddr.String(), multiStakingAmount1)
				_, err = msgServer.Delegate(ctx, delMsg1)
				return multiStakingAmount.Add(multiStakingAmount1), err
			},
			expRate: math.LegacyMustNewDecFromStr("0.625"),
			expErr:  false,
		},
		{
			name: "not found validator",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (sdk.Coin, error) {
				multiStakingAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(1000))

				delMsg := stakingtypes.NewMsgDelegate(delAddr.String(), test.GenValAddress().String(), multiStakingAmount)
				_, err := msgServer.Delegate(ctx, delMsg)
				return multiStakingAmount, err
			},
			expErr: true,
		},
		{
			name: "not allow token",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) (sdk.Coin, error) {
				multiStakingAmount := sdk.NewCoin(MultiStakingDenomB, math.NewInt(1000))

				delMsg := stakingtypes.NewMsgDelegate(delAddr.String(), valAddr.String(), multiStakingAmount)
				_, err := msgServer.Delegate(ctx, delMsg)
				return multiStakingAmount, err
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			newParam := stakingtypes.DefaultParams()
			newParam.MinCommissionRate = math.LegacyMustNewDecFromStr("0.02")
			err := suite.app.StakingKeeper.SetParams(suite.ctx, newParam)
			suite.Require().NoError(err)
			suite.msKeeper.SetBondWeight(suite.ctx, MultiStakingDenomA, math.LegacyOneDec())
			bondAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(1000))
			userBalance := sdk.NewCoin(MultiStakingDenomA, math.NewInt(10000))
			suite.FundAccount(delAddr, sdk.NewCoins(userBalance))
			suite.FundAccount(sdk.AccAddress(valAddr), sdk.NewCoins(userBalance))

			createMsg := stakingtypes.MsgCreateValidator{
				Description: stakingtypes.Description{
					Moniker:         "test",
					Identity:        "test",
					Website:         "test",
					SecurityContact: "test",
					Details:         "test",
				},
				Commission: stakingtypes.CommissionRates{
					Rate:          math.LegacyMustNewDecFromStr("0.05"),
					MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
					MaxChangeRate: math.LegacyMustNewDecFromStr("0.05"),
				},
				MinSelfDelegation: math.NewInt(200),
				DelegatorAddress:  sdk.AccAddress(valAddr).String(),
				ValidatorAddress:  valAddr.String(),
				Pubkey:            codectypes.UnsafePackAny(valPubKey),
				Value:             bondAmount,
			}
			_, err = suite.msgServer.CreateValidator(suite.ctx, &createMsg)
			suite.Require().NoError(err)

			multiStakingAmount, err := tc.malleate(suite.ctx, suite.msgServer, *suite.msKeeper)

			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				lockId := multistakingtypes.MultiStakingLockID(delAddr.String(), valAddr.String())
				lockRecord, found := suite.msKeeper.GetMultiStakingLock(suite.ctx, lockId)
				suite.Require().True(found)
				suite.Require().Equal(tc.expRate, lockRecord.GetBondWeight())

				delegation, err := suite.app.StakingKeeper.GetDelegation(suite.ctx, delAddr, valAddr)
				suite.Require().NoError(err)
				validator, err := suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
				suite.Require().NoError(err)

				multiStakingCoin := multistakingtypes.NewMultiStakingCoin(multiStakingAmount.Denom, multiStakingAmount.Amount, tc.expRate)
				expShares, err := validator.SharesFromTokens(multiStakingCoin.BondValue())
				suite.Require().NoError(err)
				suite.Require().Equal(expShares, delegation.GetShares())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestBeginRedelegate() {
	delAddr := test.GenAddress()
	valPubKey1 := test.GenPubKey()
	valPubKey2 := test.GenPubKey()

	valAddr1 := sdk.ValAddress(valPubKey1.Address())
	valAddr2 := sdk.ValAddress(valPubKey2.Address())

	testCases := []struct {
		name     string
		malleate func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) ([]sdk.Coin, error)
		expRate  []math.LegacyDec
		expLock  []math.Int
		expErr   bool
	}{
		{
			name: "redelegate from val1 to val2",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) ([]sdk.Coin, error) {
				multiStakingAmount1 := sdk.NewCoin(MultiStakingDenomA, math.NewInt(1000))
				delMsg := stakingtypes.NewMsgDelegate(delAddr.String(), valAddr1.String(), multiStakingAmount1)
				_, err := msgServer.Delegate(ctx, delMsg)
				suite.Require().NoError(err)

				multiStakingAmount2 := sdk.NewCoin(MultiStakingDenomA, math.NewInt(500))
				redelegateMsg := stakingtypes.NewMsgBeginRedelegate(delAddr.String(), valAddr1.String(), valAddr2.String(), multiStakingAmount2)
				_, err = msgServer.BeginRedelegate(ctx, redelegateMsg)
				return []sdk.Coin{multiStakingAmount1.Sub(multiStakingAmount2), multiStakingAmount2}, err
			},
			expRate: []math.LegacyDec{math.LegacyOneDec(), math.LegacyOneDec()},
			expLock: []math.Int{math.NewInt(500), math.NewInt(500)},
			expErr:  false,
		},
		{
			name: "delegate 2000 more to val1 then change rate and redelegate 600 to val2",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) ([]sdk.Coin, error) {
				multiStakingAmount1 := sdk.NewCoin(MultiStakingDenomA, math.NewInt(1000))
				delMsg1 := stakingtypes.NewMsgDelegate(delAddr.String(), valAddr1.String(), multiStakingAmount1)
				_, err := msgServer.Delegate(ctx, delMsg1)
				suite.Require().NoError(err)

				multiStakingAmount2 := sdk.NewCoin(MultiStakingDenomA, math.NewInt(1000))
				delMsg3 := stakingtypes.NewMsgDelegate(delAddr.String(), valAddr2.String(), multiStakingAmount2)
				_, err = msgServer.Delegate(ctx, delMsg3)
				suite.Require().NoError(err)

				msKeeper.SetBondWeight(ctx, MultiStakingDenomA, math.LegacyMustNewDecFromStr("0.25"))
				multiStakingAmount3 := sdk.NewCoin(MultiStakingDenomA, math.NewInt(2000))
				delMsg2 := stakingtypes.NewMsgDelegate(delAddr.String(), valAddr1.String(), multiStakingAmount3)
				_, err = msgServer.Delegate(ctx, delMsg2)
				suite.Require().NoError(err)

				multiStakingAmount4 := sdk.NewCoin(MultiStakingDenomA, math.NewInt(600))
				redelMsg := stakingtypes.NewMsgBeginRedelegate(delAddr.String(), valAddr1.String(), valAddr2.String(), multiStakingAmount4)
				if err != nil {
					return []sdk.Coin{}, err
				}

				_, err = msgServer.BeginRedelegate(ctx, redelMsg)
				return []sdk.Coin{multiStakingAmount1.Add(multiStakingAmount3).Sub(multiStakingAmount4), multiStakingAmount2.Add(multiStakingAmount4)}, err
			},
			expRate: []math.LegacyDec{math.LegacyMustNewDecFromStr("0.5"), math.LegacyMustNewDecFromStr("0.8125")},
			expLock: []math.Int{math.NewInt(2400), math.NewInt(1600)},
			expErr:  false,
		},
		{
			name: "not found validator",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) ([]sdk.Coin, error) {
				bondAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(500))
				multiStakingMsg := stakingtypes.NewMsgBeginRedelegate(delAddr.String(), valAddr1.String(), test.GenValAddress().String(), bondAmount)
				_, err := msgServer.BeginRedelegate(ctx, multiStakingMsg)
				return []sdk.Coin{}, err
			},
			expErr: true,
		},
		{
			name: "not allow token",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) ([]sdk.Coin, error) {
				bondAmount := sdk.NewCoin(MultiStakingDenomB, math.NewInt(1000))

				multiStakingMsg := stakingtypes.NewMsgBeginRedelegate(delAddr.String(), valAddr1.String(), valAddr2.String(), bondAmount)
				_, err := msgServer.BeginRedelegate(ctx, multiStakingMsg)
				return []sdk.Coin{}, err
			},
			expErr: true,
		},
		{
			name: "setup val3 with bond denom is arst then redelgate from val1 to val3",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) ([]sdk.Coin, error) {
				valPubKey3 := test.GenPubKey()
				bondAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(500))
				valAddr3 := sdk.ValAddress(valPubKey3.Address())

				userBalance := sdk.NewCoin(MultiStakingDenomA, math.NewInt(10000))
				suite.FundAccount(sdk.AccAddress(valAddr3), sdk.NewCoins(userBalance, sdk.NewCoin(MultiStakingDenomB, math.NewInt(10000))))
				createMsg := stakingtypes.MsgCreateValidator{Description: stakingtypes.Description{Moniker: "test", Identity: "test", Website: "test", SecurityContact: "test", Details: "test"}, Commission: stakingtypes.CommissionRates{Rate: math.LegacyMustNewDecFromStr("0.05"), MaxRate: math.LegacyMustNewDecFromStr("0.1"), MaxChangeRate: math.LegacyMustNewDecFromStr("0.1")}, MinSelfDelegation: math.NewInt(200), DelegatorAddress: sdk.AccAddress(valAddr3).String(), ValidatorAddress: valAddr3.String(), Pubkey: codectypes.UnsafePackAny(valPubKey3), Value: sdk.NewCoin(MultiStakingDenomB, math.NewInt(1000))}
				_, err := msgServer.CreateValidator(suite.ctx, &createMsg)
				suite.Require().NoError(err)

				multiStakingMsg := stakingtypes.NewMsgBeginRedelegate(delAddr.String(), valAddr1.String(), valAddr3.String(), bondAmount)
				_, err = msgServer.BeginRedelegate(ctx, multiStakingMsg)
				return []sdk.Coin{}, err
			},
			expRate: []math.LegacyDec{},
			expLock: []math.Int{},
			expErr:  true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			newParam := stakingtypes.DefaultParams()
			newParam.MinCommissionRate = math.LegacyMustNewDecFromStr("0.02")
			err := suite.app.StakingKeeper.SetParams(suite.ctx, newParam)
			suite.Require().NoError(err)
			suite.msKeeper.SetBondWeight(suite.ctx, MultiStakingDenomA, math.LegacyOneDec())
			suite.msKeeper.SetBondWeight(suite.ctx, MultiStakingDenomB, math.LegacyOneDec())

			bondAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(1000))
			userBalance := sdk.NewCoin(MultiStakingDenomA, math.NewInt(10000))
			suite.FundAccount(delAddr, sdk.NewCoins(userBalance, sdk.NewCoin(MultiStakingDenomB, math.NewInt(10000))))
			suite.FundAccount(sdk.AccAddress(valAddr1), sdk.NewCoins(userBalance, sdk.NewCoin(MultiStakingDenomB, math.NewInt(10000))))
			suite.FundAccount(sdk.AccAddress(valAddr2), sdk.NewCoins(userBalance, sdk.NewCoin(MultiStakingDenomB, math.NewInt(10000))))

			createMsg := stakingtypes.MsgCreateValidator{
				Description: stakingtypes.Description{
					Moniker:         "test",
					Identity:        "test",
					Website:         "test",
					SecurityContact: "test",
					Details:         "test",
				},
				Commission: stakingtypes.CommissionRates{
					Rate:          math.LegacyMustNewDecFromStr("0.05"),
					MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
					MaxChangeRate: math.LegacyMustNewDecFromStr("0.05"),
				},
				MinSelfDelegation: math.NewInt(200),
				DelegatorAddress:  sdk.AccAddress(valAddr1).String(),
				ValidatorAddress:  valAddr1.String(),
				Pubkey:            codectypes.UnsafePackAny(valPubKey1),
				Value:             bondAmount,
			}
			createMsg2 := stakingtypes.MsgCreateValidator{
				Description: stakingtypes.Description{
					Moniker:         "test",
					Identity:        "test",
					Website:         "test",
					SecurityContact: "test",
					Details:         "test",
				},
				Commission: stakingtypes.CommissionRates{
					Rate:          math.LegacyMustNewDecFromStr("0.05"),
					MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
					MaxChangeRate: math.LegacyMustNewDecFromStr("0.1"),
				},
				MinSelfDelegation: math.NewInt(200),
				DelegatorAddress:  sdk.AccAddress(valAddr2).String(),
				ValidatorAddress:  valAddr2.String(),
				Pubkey:            codectypes.UnsafePackAny(valPubKey2),
				Value:             bondAmount,
			}
			_, err = suite.msgServer.CreateValidator(suite.ctx, &createMsg)
			suite.Require().NoError(err)

			_, err = suite.msgServer.CreateValidator(suite.ctx, &createMsg2)
			suite.Require().NoError(err)

			suite.ctx = suite.ctx.WithBlockHeader(tmproto.Header{Time: time.Now()})

			multiStakingAmounts, err := tc.malleate(suite.ctx, suite.msgServer, *suite.msKeeper)

			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				lockId1 := multistakingtypes.MultiStakingLockID(delAddr.String(), valAddr1.String())
				lockRecord1, found := suite.msKeeper.GetMultiStakingLock(suite.ctx, lockId1)
				suite.Require().True(found)
				suite.Require().Equal(tc.expRate[0], lockRecord1.GetBondWeight())
				suite.Require().Equal(tc.expLock[0], lockRecord1.LockedCoin.Amount)

				delegation1, err := suite.app.StakingKeeper.GetDelegation(suite.ctx, delAddr, valAddr1)
				suite.Require().NoError(err)
				validator1, err := suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr1)
				suite.Require().NoError(err)

				multiStakingCoin1 := multistakingtypes.NewMultiStakingCoin(multiStakingAmounts[0].Denom, multiStakingAmounts[0].Amount, tc.expRate[0])
				expShares1, err := validator1.SharesFromTokens(multiStakingCoin1.BondValue())
				suite.Require().NoError(err)
				suite.Require().Equal(expShares1, delegation1.GetShares())

				lockId2 := multistakingtypes.MultiStakingLockID(delAddr.String(), valAddr2.String())
				lockRecord2, found := suite.msKeeper.GetMultiStakingLock(suite.ctx, lockId2)
				suite.Require().True(found)
				suite.Require().Equal(tc.expRate[1], lockRecord2.GetBondWeight())
				suite.Require().Equal(tc.expLock[1], lockRecord2.LockedCoin.Amount)

				delegation2, err := suite.app.StakingKeeper.GetDelegation(suite.ctx, delAddr, valAddr2)
				suite.Require().NoError(err)
				validator2, err := suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr2)
				suite.Require().NoError(err)

				multiStakingCoin2 := multistakingtypes.NewMultiStakingCoin(multiStakingAmounts[1].Denom, multiStakingAmounts[1].Amount, tc.expRate[1])
				expShares2, err := validator2.SharesFromTokens(multiStakingCoin2.BondValue())
				suite.Require().NoError(err)
				suite.Require().Equal(expShares2, delegation2.GetShares())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestUndelegate() {
	delAddr := test.GenAddress()
	valPubKey := test.GenPubKey()
	valAddr := sdk.ValAddress(valPubKey.Address())

	testCases := []struct {
		name      string
		malleate  func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) error
		expUnlock math.Int
		expLock   math.Int
		expErr    bool
	}{
		{
			name: "undelegate success",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) error {
				undelegateAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(500))
				multiStakingMsg := stakingtypes.NewMsgUndelegate(delAddr.String(), valAddr.String(), undelegateAmount)
				_, err := msgServer.Undelegate(ctx, multiStakingMsg)
				return err
			},
			expUnlock: math.NewInt(500),
			expLock:   math.NewInt(500),
			expErr:    false,
		},
		{
			name: "undelegate 250 then undelegate 500",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) error {
				undelegateAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(250))
				multiStakingMsg := stakingtypes.NewMsgUndelegate(delAddr.String(), valAddr.String(), undelegateAmount)
				_, err := msgServer.Undelegate(ctx, multiStakingMsg)
				if err != nil {
					return err
				}
				undelegateAmount1 := sdk.NewCoin(MultiStakingDenomA, math.NewInt(500))
				multiStakingMsg1 := stakingtypes.NewMsgUndelegate(delAddr.String(), valAddr.String(), undelegateAmount1)
				_, err = msgServer.Undelegate(ctx, multiStakingMsg1)
				return err
			},
			expUnlock: math.NewInt(750),
			expLock:   math.NewInt(250),
			expErr:    false,
		},
		{
			name: "not found validator",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) error {
				undelegateAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(500))
				multiStakingMsg := stakingtypes.NewMsgUndelegate(delAddr.String(), test.GenValAddress().String(), undelegateAmount)
				_, err := msgServer.Undelegate(ctx, multiStakingMsg)
				return err
			},
			expErr: true,
		},
		{
			name: "not allow token",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) error {
				undelegateAmount := sdk.NewCoin(MultiStakingDenomB, math.NewInt(1000))

				multiStakingMsg := stakingtypes.NewMsgUndelegate(delAddr.String(), test.GenValAddress().String(), undelegateAmount)
				_, err := msgServer.Undelegate(ctx, multiStakingMsg)
				return err
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			newParam := stakingtypes.DefaultParams()
			newParam.MinCommissionRate = math.LegacyMustNewDecFromStr("0.02")
			err := suite.app.StakingKeeper.SetParams(suite.ctx, newParam)
			suite.Require().NoError(err)

			initialWeight := math.LegacyMustNewDecFromStr("0.5")
			suite.msKeeper.SetBondWeight(suite.ctx, MultiStakingDenomA, initialWeight)
			bondAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(1000))
			userBalance := sdk.NewCoin(MultiStakingDenomA, math.NewInt(10000))
			suite.FundAccount(delAddr, sdk.NewCoins(userBalance, sdk.NewCoin(MultiStakingDenomB, math.NewInt(10000))))
			suite.FundAccount(sdk.AccAddress(valAddr), sdk.NewCoins(userBalance, sdk.NewCoin(MultiStakingDenomB, math.NewInt(10000))))

			createMsg := stakingtypes.MsgCreateValidator{
				Description: stakingtypes.Description{
					Moniker:         "test",
					Identity:        "test",
					Website:         "test",
					SecurityContact: "test",
					Details:         "test",
				},
				Commission: stakingtypes.CommissionRates{
					Rate:          math.LegacyMustNewDecFromStr("0.05"),
					MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
					MaxChangeRate: math.LegacyMustNewDecFromStr("0.05"),
				},
				MinSelfDelegation: math.NewInt(200),
				DelegatorAddress:  sdk.AccAddress(valAddr).String(),
				ValidatorAddress:  valAddr.String(),
				Pubkey:            codectypes.UnsafePackAny(valPubKey),
				Value:             bondAmount,
			}

			_, err = suite.msgServer.CreateValidator(suite.ctx, &createMsg)
			suite.Require().NoError(err)

			delMsg := stakingtypes.NewMsgDelegate(delAddr.String(), valAddr.String(), bondAmount)
			_, err = suite.msgServer.Delegate(suite.ctx, delMsg)
			suite.Require().NoError(err)

			suite.ctx = suite.ctx.WithBlockHeader(tmproto.Header{Time: time.Now()})
			curHeight := suite.ctx.BlockHeight()
			err = tc.malleate(suite.ctx, suite.msgServer, *suite.msKeeper)

			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				lockId := multistakingtypes.MultiStakingLockID(delAddr.String(), valAddr.String())
				lockRecord, found := suite.msKeeper.GetMultiStakingLock(suite.ctx, lockId)
				suite.Require().True(found)
				suite.Require().Equal(tc.expLock, lockRecord.LockedCoin.Amount)

				delegation, err := suite.app.StakingKeeper.GetDelegation(suite.ctx, delAddr, valAddr)
				suite.Require().NoError(err)
				validator, err := suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
				suite.Require().NoError(err)

				multiStakingCoin := multistakingtypes.NewMultiStakingCoin(MultiStakingDenomA, tc.expLock, initialWeight)
				expShares, err := validator.SharesFromTokens(multiStakingCoin.BondValue())
				suite.Require().NoError(err)
				suite.Require().Equal(expShares, delegation.GetShares())

				unlockID := multistakingtypes.MultiStakingUnlockID(delAddr.String(), valAddr.String())
				unbondRecord, found := suite.msKeeper.GetMultiStakingUnlock(suite.ctx, unlockID)
				suite.Require().True(found)
				suite.Require().Equal(tc.expUnlock, unbondRecord.Entries[0].UnlockingCoin.Amount)

				ubd, err := suite.app.StakingKeeper.GetUnbondingDelegation(suite.ctx, delAddr, valAddr)
				suite.Require().NoError(err)
				unlockStakingCoin := multistakingtypes.NewMultiStakingCoin(MultiStakingDenomA, tc.expUnlock, initialWeight)
				totalUBDAmount := math.ZeroInt()

				for _, ubdEntry := range ubd.Entries {
					if ubdEntry.CreationHeight == curHeight {
						totalUBDAmount = totalUBDAmount.Add(ubdEntry.Balance)
					}
				}
				suite.Require().Equal(unlockStakingCoin.BondValue(), totalUBDAmount)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestCancelUnbondingDelegation() {
	delAddr := test.GenAddress()
	valPubKey := test.GenPubKey()
	valAddr := sdk.ValAddress(valPubKey.Address())

	testCases := []struct {
		name      string
		malleate  func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) error
		expUnlock math.Int
		expLock   math.Int
		expErr    bool
	}{
		{
			name: "cancel unbonding success",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) error {
				cancelAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(500))
				multiStakingMsg := stakingtypes.NewMsgCancelUnbondingDelegation(delAddr.String(), valAddr.String(), ctx.BlockHeight(), cancelAmount)
				_, err := msgServer.CancelUnbondingDelegation(ctx, multiStakingMsg)
				return err
			},
			expUnlock: math.NewInt(500),
			expLock:   math.NewInt(1500),
			expErr:    false,
		},
		{
			name: "cancel unbonding 250 then cancel unbonding 500",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) error {
				cancelAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(250))
				multiStakingMsg := stakingtypes.NewMsgCancelUnbondingDelegation(delAddr.String(), valAddr.String(), ctx.BlockHeight(), cancelAmount)
				_, err := msgServer.CancelUnbondingDelegation(ctx, multiStakingMsg)
				if err != nil {
					return err
				}
				cancelAmount1 := sdk.NewCoin(MultiStakingDenomA, math.NewInt(500))
				multiStakingMsg1 := stakingtypes.NewMsgCancelUnbondingDelegation(delAddr.String(), valAddr.String(), ctx.BlockHeight(), cancelAmount1)
				_, err = msgServer.CancelUnbondingDelegation(ctx, multiStakingMsg1)
				return err
			},
			expUnlock: math.NewInt(250),
			expLock:   math.NewInt(1750),
			expErr:    false,
		},
		{
			name: "not found validator",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) error {
				cancelAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(500))
				multiStakingMsg := stakingtypes.NewMsgCancelUnbondingDelegation(delAddr.String(), test.GenValAddress().String(), ctx.BlockHeight(), cancelAmount)
				_, err := msgServer.CancelUnbondingDelegation(ctx, multiStakingMsg)
				return err
			},
			expErr: true,
		},
		{
			name: "not allow token",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) error {
				cancelAmount := sdk.NewCoin(MultiStakingDenomB, math.NewInt(1000))

				multiStakingMsg := stakingtypes.NewMsgCancelUnbondingDelegation(delAddr.String(), valAddr.String(), ctx.BlockHeight(), cancelAmount)
				_, err := msgServer.CancelUnbondingDelegation(ctx, multiStakingMsg)
				return err
			},
			expErr: true,
		},
		{
			name: "not found entry at height 20",
			malleate: func(ctx sdk.Context, msgServer stakingtypes.MsgServer, msKeeper multistakingkeeper.Keeper) error {
				cancelAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(500))
				multiStakingMsg := stakingtypes.NewMsgCancelUnbondingDelegation(delAddr.String(), valAddr.String(), 20, cancelAmount)
				_, err := msgServer.CancelUnbondingDelegation(ctx, multiStakingMsg)
				return err
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			newParam := stakingtypes.DefaultParams()
			newParam.MinCommissionRate = math.LegacyMustNewDecFromStr("0.02")
			err := suite.app.StakingKeeper.SetParams(suite.ctx, newParam)
			suite.Require().NoError(err)

			initialWeight := math.LegacyMustNewDecFromStr("0.5")
			suite.msKeeper.SetBondWeight(suite.ctx, MultiStakingDenomA, initialWeight)
			bondAmount := sdk.NewCoin(MultiStakingDenomA, math.NewInt(2000))
			userBalance := sdk.NewCoin(MultiStakingDenomA, math.NewInt(10000))
			suite.FundAccount(delAddr, sdk.NewCoins(userBalance, sdk.NewCoin(MultiStakingDenomB, math.NewInt(10000))))
			suite.FundAccount(sdk.AccAddress(valAddr), sdk.NewCoins(userBalance, sdk.NewCoin(MultiStakingDenomB, math.NewInt(10000))))

			createMsg := stakingtypes.MsgCreateValidator{
				Description: stakingtypes.Description{
					Moniker:         "test",
					Identity:        "test",
					Website:         "test",
					SecurityContact: "test",
					Details:         "test",
				},
				Commission: stakingtypes.CommissionRates{
					Rate:          math.LegacyMustNewDecFromStr("0.05"),
					MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
					MaxChangeRate: math.LegacyMustNewDecFromStr("0.05"),
				},
				MinSelfDelegation: math.NewInt(200),
				DelegatorAddress:  sdk.AccAddress(valAddr).String(),
				ValidatorAddress:  valAddr.String(),
				Pubkey:            codectypes.UnsafePackAny(valPubKey),
				Value:             bondAmount,
			}

			_, err = suite.msgServer.CreateValidator(suite.ctx, &createMsg)
			suite.Require().NoError(err)

			delMsg := stakingtypes.NewMsgDelegate(delAddr.String(), valAddr.String(), bondAmount)
			_, err = suite.msgServer.Delegate(suite.ctx, delMsg)
			suite.Require().NoError(err)

			suite.ctx = suite.ctx.WithBlockHeader(tmproto.Header{Time: time.Now()}).WithBlockHeight(1)
			curHeight := suite.ctx.BlockHeight()

			unbondMsg := stakingtypes.NewMsgUndelegate(delAddr.String(), valAddr.String(), sdk.NewCoin(MultiStakingDenomA, math.NewInt(1000)))
			_, err = suite.msgServer.Undelegate(suite.ctx, unbondMsg)
			suite.Require().NoError(err)

			err = tc.malleate(suite.ctx, suite.msgServer, *suite.msKeeper)

			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				lockId := multistakingtypes.MultiStakingLockID(delAddr.String(), valAddr.String())
				lockRecord, found := suite.msKeeper.GetMultiStakingLock(suite.ctx, lockId)
				suite.Require().True(found)
				suite.Require().Equal(tc.expLock, lockRecord.LockedCoin.Amount)

				delegation, err := suite.app.StakingKeeper.GetDelegation(suite.ctx, delAddr, valAddr)
				suite.Require().NoError(err)
				validator, err := suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
				suite.Require().NoError(err)

				multiStakingCoin := multistakingtypes.NewMultiStakingCoin(MultiStakingDenomA, tc.expLock, initialWeight)
				expShares, err := validator.SharesFromTokens(multiStakingCoin.BondValue())
				suite.Require().NoError(err)
				suite.Require().Equal(expShares, delegation.GetShares())

				unlockID := multistakingtypes.MultiStakingUnlockID(delAddr.String(), valAddr.String())
				unbondRecord, found := suite.msKeeper.GetMultiStakingUnlock(suite.ctx, unlockID)
				suite.Require().True(found)
				suite.Require().Equal(tc.expUnlock, unbondRecord.Entries[0].UnlockingCoin.Amount)

				ubd, err := suite.app.StakingKeeper.GetUnbondingDelegation(suite.ctx, delAddr, valAddr)
				suite.Require().NoError(err)
				unlockStakingCoin := multistakingtypes.NewMultiStakingCoin(MultiStakingDenomA, tc.expUnlock, initialWeight)
				totalUBDAmount := math.ZeroInt()

				for _, ubdEntry := range ubd.Entries {
					if ubdEntry.CreationHeight == curHeight {
						totalUBDAmount = totalUBDAmount.Add(ubdEntry.Balance)
					}
				}
				suite.Require().Equal(unlockStakingCoin.BondValue(), totalUBDAmount)
			}
		})
	}
}
