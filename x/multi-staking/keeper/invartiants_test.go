package keeper_test

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/realio-tech/multi-staking-module/testing/simapp"
	"github.com/realio-tech/multi-staking-module/testutil"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
)

var (
	MultiStakingDenomA = "ario"
	MultiStakingDenomB = "arst"
)

func (suite *KeeperTestSuite) TestModuleAccountInvariants() {
	delAddr := testutil.GenAddress()
	priv, valAddr := testutil.GenValAddressWithPrivKey()
	valPubKey := priv.PubKey()

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name:     "Success",
			malleate: func() {},
			expPass:  true,
		},
		{
			name: "Success Create Validator",
			malleate: func() {
				valCoins := sdk.NewCoins(sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(10000)), sdk.NewCoin(MultiStakingDenomB, sdk.NewInt(10000)))
				err := simapp.FundAccount(suite.app, suite.ctx, delAddr, valCoins)
				suite.Require().NoError(err)

				suite.msKeeper.SetBondWeight(suite.ctx, MultiStakingDenomA, sdk.MustNewDecFromStr("0.3"))
				bondAmount := sdk.NewCoin(MultiStakingDenomA, sdk.NewInt(3001))
				msg := stakingtypes.MsgCreateValidator{
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

				_, err = suite.msgServer.CreateValidator(suite.ctx, &msg)
				suite.Require().NoError(err)
			},
			expPass: true,
		},
		{
			name: "Success Edit Validator",
			malleate: func() {
				newRate := sdk.MustNewDecFromStr("0.03")
				newMinSelfDelegation := sdk.NewInt(300)
				editMsg := stakingtypes.NewMsgEditValidator(valAddr, stakingtypes.Description{
					Moniker:         "test 1",
					Identity:        "test 1",
					Website:         "test 1",
					SecurityContact: "test 1",
					Details:         "test 1",
				},
					&newRate,
					&newMinSelfDelegation,
				)
				_, err := suite.msgServer.EditValidator(suite.ctx, editMsg)
				suite.Require().NoError(err)
			},
			expPass: true,
		},
	}
	for _, tc := range testCases {
		suite.SetupTest() // reset
		tc.malleate()
		_, broken := keeper.ModuleAccountInvariants(*suite.msKeeper)(suite.ctx)

		if tc.expPass {
			suite.Require().False(broken)
		}
	}
}
