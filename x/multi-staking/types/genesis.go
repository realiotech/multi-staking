package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
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

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (g GenesisState) UnpackInterfaces(c codectypes.AnyUnpacker) error {
	for i := range g.StakingGenesisState.Validators {
		if err := g.StakingGenesisState.Validators[i].UnpackInterfaces(c); err != nil {
			return err
		}
	}
	return nil
}
