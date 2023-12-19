package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/multistaking module sentinel errors
var (
	ErrInvalidAddBondTokenProposal = sdkerrors.Register(ModuleName, 2,
		"invalid add bond token proposal",
	)
	ErrInvalidChangeBondTokenWeightProposal = sdkerrors.Register(ModuleName, 3,
		"invalid change bond token weight proposal",
	)
	ErrNotFoundMultiStaking = sdkerrors.Register(ModuleName, 4,
		"can't find multi staking",
	)
	ErrStakingNotExitsts = sdkerrors.Register(ModuleName, 5,
		"StakingLock not exists",
	)
	ErrCheckInsufficientAmount = sdkerrors.Register(ModuleName, 6,
		"unlock amount greater than lock amount",
	)
	ErrRecordNotExists = sdkerrors.Register(ModuleName, 7,
		"record not exists",
	)
	ErrLessThanZero = sdkerrors.Register(ModuleName, 8,
		"cannot be less than 0",
	)
	ErrNotAllowedToken = sdkerrors.Register(ModuleName, 9,
		"not allowed token",
	)
	ErrDelegationNotFound = sdkerrors.Register(ModuleName, 10,
		"delegation not found",
	)
	ErrValidatorNotFound = sdkerrors.Register(ModuleName, 10,
		"validator not found",
	)
	ErrUnrecognized = sdkerrors.Register(ModuleName, 11,
		"unrecognized brond denom proposal content type",
	)
	ErrBrondDenomAlreadyExists = sdkerrors.Register(ModuleName, 12,
		"brond denom already exists",
	)
	ErrBrondDenomDoesNotExist = sdkerrors.Register(ModuleName, 13,
		"brond denom does not exist",
	)
)
