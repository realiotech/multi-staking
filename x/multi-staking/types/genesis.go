package types

import (
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func DefaultGenesis() *GenesisState {
	stakingGenesis := stakingtypes.DefaultGenesisState()

	return &GenesisState{
		StakingGenesisState: *stakingGenesis,
	}
}

// Validate performs basic genesis state validation, returning an error upon any failure.
// TODO: Add validate genesis
func (gs GenesisState) Validate() error {
	return nil
}
