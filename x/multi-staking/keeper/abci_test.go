package keeper_test

import (
	"time"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (suite *KeeperTestSuite) TestMsUnlockEndBlocker() {
	// val A
	// delegate to val A with X ario
	// undelegate from val A
	// val A got slash
	// nextblock
	// check A balance has X ario/ zero stake

	testCases := []struct {
		name        string
		lockAmount  sdkmath.Int
		slashFactor sdkmath.LegacyDec
	}{
		{
			name:        "no slashing",
			lockAmount:  sdkmath.NewInt(3788),
			slashFactor: sdkmath.LegacyZeroDec(),
		},
		{
			name:        "slash half of lock coin",
			lockAmount:  sdkmath.NewInt(123),
			slashFactor: sdkmath.LegacyMustNewDecFromStr("0.5"),
		},
		{
			name:        "slash all of lock coin",
			lockAmount:  sdkmath.NewInt(19090),
			slashFactor: sdkmath.LegacyZeroDec(),
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			// height 1
			suite.SetupTest()

			vals, err := suite.app.StakingKeeper.GetAllValidators(suite.ctx)
			suite.NoError(err)
			val := vals[0]

			valAddr, err := sdk.ValAddressFromBech32(val.GetOperator())
			suite.NoError(err)
			msDenom := suite.msKeeper.GetValidatorMultiStakingCoin(suite.ctx, valAddr)

			msCoin := sdk.NewCoin(msDenom, tc.lockAmount)

			msStaker := suite.CreateAndFundAccount(sdk.NewCoins(msCoin))

			delegateMsg := &stakingtypes.MsgDelegate{
				DelegatorAddress: msStaker.String(),
				ValidatorAddress: val.OperatorAddress,
				Amount:           msCoin,
			}
			_, err = suite.msgServer.Delegate(suite.ctx, delegateMsg)
			suite.NoError(err)

			// height 2
			suite.NextBlock(time.Second)

			if !tc.slashFactor.IsZero() {
				valAddr, _ := sdk.ValAddressFromBech32(val.GetOperator())
				val, err := suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
				suite.NoError(err)
				require.NotNil(suite.T(), val)

				slashedPow := suite.app.StakingKeeper.TokensToConsensusPower(suite.ctx, val.Tokens)

				valConsAddr, err := val.GetConsAddr()
				require.NoError(suite.T(), err)

				// height 3
				suite.NextBlock(time.Second)

				err = suite.app.SlashingKeeper.Slash(suite.ctx, valConsAddr, tc.slashFactor, slashedPow, 2)
				require.NoError(suite.T(), err)
			} else {
				// height 3
				suite.NextBlock(time.Second)
			}

			undelegateMsg := stakingtypes.MsgUndelegate{
				DelegatorAddress: msStaker.String(),
				ValidatorAddress: val.OperatorAddress,
				Amount:           msCoin,
			}

			_, err = suite.msgServer.Undelegate(suite.ctx, &undelegateMsg)
			suite.NoError(err)

			// pass unbonding period
			suite.NextBlock(time.Duration(1000000000000000000))
			suite.NextBlock(time.Duration(1))

			unlockAmount := suite.app.BankKeeper.GetBalance(suite.ctx, msStaker, msDenom).Amount

			expectedUnlockAmount := sdkmath.LegacyNewDecFromInt(tc.lockAmount).Mul(sdkmath.LegacyOneDec().Sub(tc.slashFactor)).TruncateInt()

			suite.True(SoftEqualInt(unlockAmount, expectedUnlockAmount) || DiffLTEThanOne(unlockAmount, expectedUnlockAmount))
		})
	}
}
