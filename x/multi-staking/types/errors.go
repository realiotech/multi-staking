package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/multistaking module sentinel errors
var (
	ErrInvalidAddBondTokenProposal          = sdkerrors.Register(ModuleName, 2, "invalid add bond token proposal")
	ErrInvalidChangeBondTokenWeightProposal = sdkerrors.Register(ModuleName, 3, "invalid change bond token weight proposal")
	ErrInvalidRemoveBondTokenProposal       = sdkerrors.Register(ModuleName, 4, "invalid remove bond token weight proposal")
)
