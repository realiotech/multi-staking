package types

import (
	time "time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type StakingKeeper interface {
	DequeueAllMatureUBDQueue(ctx sdk.Context, currTime time.Time) (matureUnbonds []stakingtypes.DVPair)
	BondDenom(ctx sdk.Context) (res string)
	GetUnbondingDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (ubd stakingtypes.UnbondingDelegation, found bool)
}
