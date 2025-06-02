package keeper_test

import (
	"time"

	multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"

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
		lockAmount  math.Int
		slashFactor math.LegacyDec
	}{
		{
			name:        "no slashing",
			lockAmount:  math.NewInt(3788),
			slashFactor: math.LegacyZeroDec(),
		},
		{
			name:        "slash half of lock coin",
			lockAmount:  math.NewInt(123),
			slashFactor: math.LegacyMustNewDecFromStr("0.5"),
		},
		{
			name:        "slash all of lock coin",
			lockAmount:  math.NewInt(19090),
			slashFactor: math.LegacyZeroDec(),
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			// height 1
			suite.SetupTest()

			vals, err := suite.app.StakingKeeper.GetAllValidators(suite.ctx)
			val := vals[0]
			operatorAddr, err := sdk.ValAddressFromBech32(val.OperatorAddress)

			msDenom := suite.msKeeper.GetValidatorMultiStakingCoin(suite.ctx, operatorAddr)

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
			suite.ctx = suite.ctx.WithBlockHeight(2).WithBlockTime(suite.ctx.BlockTime().Add(time.Second))

			if !tc.slashFactor.IsZero() {
				val, err := suite.app.StakingKeeper.GetValidator(suite.ctx, operatorAddr)
				require.NoError(suite.T(), err)

				slashedPow := suite.app.StakingKeeper.TokensToConsensusPower(suite.ctx, val.Tokens)

				valConsAddr, err := val.GetConsAddr()
				require.NoError(suite.T(), err)

				// height 3
				suite.ctx = suite.ctx.WithBlockHeight(3).WithBlockTime(suite.ctx.BlockTime().Add(time.Second))
				suite.app.SlashingKeeper.Slash(suite.ctx, valConsAddr, tc.slashFactor, slashedPow, 2)
			} else {
				// height 3
				suite.ctx = suite.ctx.WithBlockHeight(3).WithBlockTime(suite.ctx.BlockTime().Add(time.Second))
			}

			undelegateMsg := stakingtypes.MsgUndelegate{
				DelegatorAddress: msStaker.String(),
				ValidatorAddress: val.OperatorAddress,
				Amount:           msCoin,
			}

			msgServer := multistakingkeeper.NewMsgServerImpl(suite.app.MultiStakingKeeper)
			_, err = msgServer.Undelegate(suite.ctx, &undelegateMsg)
			suite.NoError(err)

			ubds, err := suite.app.StakingKeeper.GetAllUnbondingDelegations(suite.ctx, msStaker)
			suite.NoError(err)

			// pass unbonding period
			suite.ctx = suite.ctx.WithBlockTime(ubds[0].Entries[0].CompletionTime)
			_, err = suite.app.EndBlocker(suite.ctx)
			require.NoError(suite.T(), err)

			unlockAmount := suite.app.BankKeeper.GetBalance(suite.ctx, msStaker, msDenom).Amount
			expectedUnlockAmount := math.LegacyNewDecFromInt(tc.lockAmount).Mul(math.LegacyOneDec().Sub(tc.slashFactor)).TruncateInt()
			suite.True(SoftEqualInt(unlockAmount, expectedUnlockAmount) || DiffLTEThanOne(unlockAmount, expectedUnlockAmount))
		})
	}
}
