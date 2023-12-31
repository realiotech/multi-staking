package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/multistaking module sentinel errors
var (
	ErrInvalidAddMultiStakingCoinProposal = sdkerrors.Register(ModuleName, 2, "invalid add bond token proposal")
	ErrInvalidUpdateBondWeightProposal    = sdkerrors.Register(ModuleName, 3, "invalid change bond token weight proposal")
)
