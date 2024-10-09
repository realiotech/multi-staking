package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	stkm stakingkeeper.Migrator
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper *stakingkeeper.Keeper, legacySubspace exported.Subspace) Migrator {
	return Migrator{
		stkm: stakingkeeper.NewMigrator(keeper, legacySubspace),
	}
}

// Migrate1to2 migrates multi-staking state from consensus version 1 to 2. (sdk46 to sdk47)
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return m.stkm.Migrate3to4(ctx)
}

// Migrate2to3 migrates multi-staking state from consensus version 2 to 3. (sdk47 to sdk50)
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return m.stkm.Migrate4to5(ctx)
}
