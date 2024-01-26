package keeper_test

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	"github.com/realio-tech/multi-staking-module/testing"
	"github.com/realio-tech/multi-staking-module/testing/simapp"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestMsUnlockEnblocker() {
	// val A
	// delegate to val A with X ario
	// undelegate from val A
	// val A got slash
	// nextblock
	// check A balance has X ario/ zero stake

	testCases := []struct {
		name        string
		lockAmount  math.Int
		slashFactor sdk.Dec
	}{
		{
			name:        "no slashing",
			lockAmount:  math.NewInt(3788),
			slashFactor: sdk.ZeroDec(),
		},
		// {
		// 	name:        "slash half of lock coin",
		// 	lockAmount:  math.NewInt(123),
		// 	slashFactor: sdk.MustNewDecFromStr("0.5"),
		// },
		// {
		// 	name:        "slash all of lock coin",
		// 	lockAmount:  math.NewInt(19090),
		// 	slashFactor: sdk.ZeroDec(),
		// },
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			// height 0
			msStaker := testing.GenAddress()

			vals := suite.app.StakingKeeper.GetAllValidators(suite.ctx)
			val := vals[0]

			msDenom := suite.msKeeper.GetValidatorMultiStakingCoin(suite.ctx, val.GetOperator())
			fmt.Println(val.GetOperator().String())
			fmt.Println(msDenom)
			msCoin := sdk.NewCoin(msDenom, tc.lockAmount)

			simapp.FundAccount(suite.app, suite.ctx, msStaker, sdk.NewCoins(msCoin))

			suite.NextBlock(time.Second)

			if !tc.slashFactor.IsZero() {
				val, found := suite.app.StakingKeeper.GetValidator(suite.ctx, val.GetOperator())
				require.True(suite.T(), found)

				slashedPow := suite.app.StakingKeeper.TokensToConsensusPower(suite.ctx, val.Tokens)

				valConsAddr, err := val.GetConsAddr()
				require.NoError(suite.T(), err)

				// height 1
				suite.NextBlock(time.Second)

				suite.app.SlashingKeeper.Slash(suite.ctx, valConsAddr, tc.slashFactor, slashedPow, 0)
			} else {

				// height 1
				suite.NextBlock(time.Second)
			}

			// pass unbonding period
			suite.NextBlock(time.Duration(1000000000000000000))
			suite.NextBlock(time.Duration(1))

			unlockAmount := suite.app.BankKeeper.GetBalance(suite.ctx, msStaker, msDenom).Amount

			expectedUnlockAmount := sdk.NewDecFromInt(tc.lockAmount).Mul(sdk.OneDec().Sub(tc.slashFactor)).TruncateInt()

			require.Equal(suite.T(), unlockAmount, expectedUnlockAmount)
		})
	}
}
