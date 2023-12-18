package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/realio-tech/multi-staking-module/testutil"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (suite *KeeperTestSuite) TestInitGenesis() {
	delAddr := testutil.GenAddress()
	valPubKey := testutil.GenPubKey()
	valAddr := sdk.ValAddress(valPubKey.Address())
	multiStakingLock := types.MultiStakingLock{
		ConversionRatio: sdk.NewDec(1),
		LockedAmount:    sdk.NewInt(1),
		DelAddr:         delAddr.String(),
		ValAddr:         valAddr.String(),
	}
	validatorAllowedToken := types.ValidatorAllowedToken{
		ValAddr:    valAddr.String(),
		TokenDenom: "stake",
	}

	var delegations []stakingtypes.Delegation
	genesisDelegations := suite.stakingKeeper.GetAllDelegations(suite.ctx)
	delegations = append(delegations, genesisDelegations...)

	validators := suite.stakingKeeper.GetAllValidators(suite.ctx)

	params := suite.stakingKeeper.GetParams(suite.ctx)

	stakingGenesisState := stakingtypes.NewGenesisState(params, validators, delegations)

	expectedGenesisState := types.GenesisState{
		MultiStakingLocks:     []types.MultiStakingLock{multiStakingLock},
		ValidatorAllowedToken: []types.ValidatorAllowedToken{validatorAllowedToken},
		StakingGenesisState:   stakingGenesisState,
	}

	suite.msKeeper.InitGenesis(suite.ctx, expectedGenesisState)

	actualGenesisState := suite.msKeeper.ExportGenesis(suite.ctx)
	suite.Require().NotNil(actualGenesisState)
	suite.Require().Equal(expectedGenesisState.MultiStakingLocks, actualGenesisState.MultiStakingLocks)
	suite.Require().Equal(expectedGenesisState.ValidatorAllowedToken, actualGenesisState.ValidatorAllowedToken)
	suite.Require().Equal(expectedGenesisState.StakingGenesisState.Delegations, actualGenesisState.StakingGenesisState.Delegations)
	suite.Require().Equal(expectedGenesisState.StakingGenesisState.Validators, actualGenesisState.StakingGenesisState.Validators)

}
