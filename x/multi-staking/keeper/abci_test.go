package keeper_test

import (
	"github.com/realio-tech/multi-staking-module/testutil"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (suite *KeeperTestSuite) TestEndBlocker() {
	suite.SetupTest()

	mulStaker := testutil.GenAddress()
	valAddr := testutil.GenValAddress()
	const multiStakingDenom = "ario"

	suite.msKeeper.SetValidatorMultiStakingCoin(suite.ctx, valAddr, multiStakingDenom)

	mulBalance := sdk.NewCoins(sdk.NewCoin(stakingtypes.DefaultParams().BondDenom, sdk.NewInt(10000)), sdk.NewCoin(multiStakingDenom, sdk.NewInt(10000)))
	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, mulBalance)
	suite.NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, mulStaker, mulBalance)
	suite.NoError(err)

	unlockEntry := types.UnlockEntry{
		CreationHeight: suite.ctx.BlockHeight(),
		UnlockingCoin: types.MultiStakingCoin{
			Denom:      multiStakingDenom,
			Amount:     sdk.NewInt(10000),
			BondWeight: sdk.NewDec(1),
		},
	}
	newUbd := types.MultiStakingUnlock{
		UnlockID: types.UnlockID{
			MultiStakerAddr: mulStaker.String(),
			ValAddr:         valAddr.String(),
		},
		Entries: []types.UnlockEntry{unlockEntry},
	}

	suite.msKeeper.SetMultiStakingUnlock(suite.ctx, newUbd)
	matureUnbondingDelegations := suite.msKeeper.GetMatureUnbondingDelegations(suite.ctx)
	completionTime := suite.ctx.BlockHeader().Time

	unbondingDelegationEntry := stakingtypes.UnbondingDelegationEntry{
		CreationHeight: suite.ctx.BlockHeight(),
		CompletionTime: completionTime,
		InitialBalance: sdk.NewInt(1000),
		Balance:        sdk.NewInt(1000),
	}

	unbondingDelegation := stakingtypes.UnbondingDelegation{
		DelegatorAddress: mulStaker.String(),
		ValidatorAddress: valAddr.String(),
		Entries:          []stakingtypes.UnbondingDelegationEntry{unbondingDelegationEntry},
	}

	matureUnbondingDelegations = append(matureUnbondingDelegations, unbondingDelegation)

	suite.msKeeper.EndBlocker(suite.ctx, matureUnbondingDelegations)
}
