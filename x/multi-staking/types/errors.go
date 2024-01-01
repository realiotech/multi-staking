package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/multistaking module sentinel errors
var (
	ErrInvalidAddMultiStakingCoinProposal = sdkerrors.Register(ModuleName, 2, "invalid add multi staking coin proposal")
	ErrInvalidUpdateBondWeightProposal    = sdkerrors.Register(ModuleName, 3, "invalid update bond weight proposal")
)
