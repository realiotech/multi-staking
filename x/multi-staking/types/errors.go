package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrUnrecognized = errorsmod.Register(ModuleName, 1,
		"unrecognized brond denom proposal content type",
	)
	ErrBrondDenomAlreadyExists = errorsmod.Register(ModuleName, 2,
		"brond denom already exists",
	)
	ErrBrondDenomDoesNotExist = errorsmod.Register(ModuleName, 3,
		"brond denom does not exist",
	)
	ErrNotFoundMultiStaking = errorsmod.Register(ModuleName, 4,
		"can't find multi staking",
	)
	ErrStakingNotExitsts = errorsmod.Register(ModuleName, 5,
		"StakingLock not exists",
	)
	ErrCheckInsufficientAmount = errorsmod.Register(ModuleName, 6,
		"unlock amount greater than lock amount",
	)
	ErrRecordNotExists = errorsmod.Register(ModuleName, 7,
		"record not exists",
	)
	ErrLessThanZero = errorsmod.Register(ModuleName, 8,
		"cannot be less than 0",
	)
	ErrNotAllowedToken = errorsmod.Register(ModuleName, 9,
		"not allowed token",
	)
	ErrDelegationNotFound = errorsmod.Register(ModuleName, 10,
		"delegation not found",
	)
	ErrValidatorNotFound = errorsmod.Register(ModuleName, 10,
		"validator not found",
	)
)
