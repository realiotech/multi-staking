<!--
order: 3
-->

# Messages

In this section we describe the processing of the multi-staking messages and the corresponding updates to the state. 
All created/modified state objects specified by each message are defined within the [state](./02_state.md) section.

## MsgCreateValidator

A validator is created using the `MsgCreateValidator` message.
The validator must be created with an initial delegation from the operator. 
The Initial delegation token must match the `bond denom` specified in `MsgCreateValidator`.

Logic flow:

1. Setting `ValidatorMultiStakingCoin`.

2. Converting `MsgCreateValidator` to `stakingtypes.MsgCreateValidator` and
calling `stakingkeeper.CreateValidator()`.

This message is expected to fail if:

* `ValOperatorAddr` already exists in state.
* The call to `stakingkeeper.CreateValidator()` returns an error.

## MsgEditValidator

The `Description`, `CommissionRate` of a validator can be updated using the
`MsgEditValidator` message.

Logic flow:

1. Converting `MsgEditValidator` to `stakingtypes.MsgEditValidator` and
calling `stakingkeeper.EditValidator()`.

This message is expected to fail if:

* The call to `stakingkeeper.EditValidator()` returns an error.

## MsgDelegate

Within this message the delegator locked up coins in the `multi-staking` module account. 
The `multi-staking` inturns mint a calculated amount of `bondtoken` and delegate.

Logic flow:

* Lock `multi staking` coin in the `multi-staking` module account.

* Caculate the `bond token` to be minted using `BondWeight`.

* Mint `bond token` to `delegator`

* Update `multi staking lock`.

* `delegate` using the minted `sdkbond token`

## MsgUndelegate

The `MsgUndelegate` message allows delegators to undelegate their `multi-staking` tokens from
validator, after the unbonding period the module will unlock the `multi-staking` tokens to return to the delegator

Logic flow:

* Calculate ammount of `bond token` need to be `undelegate`

* Update `multi staking lock`

* Update `multi staking unlock`

* Call `stakingkeeper.Undelegate()` with the calculated amount of `bond token`

The rest of the unbonding logic such as sending locked coins back to user will happens at `EndBlock()`

## MsgCancelUnbonding 

The `MsgCancelUnbonding` message allows delegators to cancel the `unbondingDelegation` entry and deleagate back to a previous validator.

Logic flow:

* Calculate amount of `bond token` need to be `cancel undelegation`

* Update `multi staking lock`

* Update `multi staking unlock`

* Call `stakingkeeper.CancelUnbondingDelegation()` with the calculated amount of `bond token`

## MsgBeginRedelegate

The `MsgBeginRedelegate` message allows delegators to instantly switch validators. Once
the unbonding period has passed, the redelegation is automatically completed in
the EndBlocker.

Logic flow:

* Calculate amount of `bond token` need to be `redelegate`

* Update the src `multi-staking lock` and the dst `multi-staking lock`

* Call `stakingkeeper.BeginRedelegate()` with the calculated amount of `bond token`

## MsgCreateEVMValidator

A validator is created using ERC20 tokens with the `MsgCreateEVMValidator` message.
The validator must be created with an initial delegation from the operator using an ERC20 token.
The ERC20 token must be registered as a multi-staking coin via `AddMultiStakingEVMCoinProposal`.

Logic flow:

1. Convert ERC20 tokens to cosmos coins using the ERC20 keeper's `ConvertERC20` functionality.

2. Get the cosmos denom for the ERC20 token using `GetTokenDenom`.

3. Create a standard `stakingtypes.MsgCreateValidator` with the converted cosmos coin.

4. Call the existing `CreateValidator` method with the converted message.

This message is expected to fail if:

* `ValidatorAddress` already exists in state.
* The ERC20 contract is not registered as a multi-staking coin.
* The ERC20 conversion fails.
* The call to `CreateValidator()` returns an error.

## MsgDelegateEVM

Within this message the delegator delegates ERC20 tokens to a validator.
The ERC20 tokens are first converted to cosmos coins, then processed through the standard delegation flow.

Logic flow:

1. Convert ERC20 tokens to cosmos coins using `ConvertERC20`.

2. Get the cosmos denom for the ERC20 token.

3. Create a standard `stakingtypes.MsgDelegate` with the converted cosmos coin.

4. Call the existing `Delegate` method which handles multi-staking logic.

This message is expected to fail if:

* The ERC20 contract is not registered as a multi-staking coin.
* The ERC20 conversion fails.
* The call to `Delegate()` returns an error.

## MsgBeginRedelegateEVM

The `MsgBeginRedelegateEVM` message allows delegators to instantly switch validators using ERC20 tokens.
The ERC20 contract address is used to determine the cosmos denom for processing.

Logic flow:

1. Get the cosmos denom for the ERC20 token using `GetTokenDenom`.

2. Create a standard `stakingtypes.MsgBeginRedelegate` with the cosmos denom.

3. Call the existing `BeginRedelegate` method.

This message is expected to fail if:

* The ERC20 contract is not registered as a multi-staking coin.
* The call to `BeginRedelegate()` returns an error.

## MsgUndelegateEVM

The `MsgUndelegateEVM` message allows delegators to undelegate ERC20 tokens from a validator.
After the unbonding period, the tokens will be converted back to ERC20 format.

Logic flow:

1. Get the cosmos denom for the ERC20 token using `GetTokenDenom`.

2. Create a standard `stakingtypes.MsgUndelegate` with the cosmos denom.

3. Call the existing `Undelegate` method which handles multi-staking unlock logic.

This message is expected to fail if:

* The ERC20 contract is not registered as a multi-staking coin.
* The call to `Undelegate()` returns an error.

## MsgCancelUnbondingEVMDelegation

The `MsgCancelUnbondingEVMDelegation` message allows delegators to cancel unbonding delegation entries for ERC20 tokens and delegate back to the validator.

Logic flow:

1. Get the cosmos denom for the ERC20 token using `GetTokenDenom`.

2. Create a standard `stakingtypes.MsgCancelUnbondingDelegation` with the cosmos denom.

3. Call the existing `CancelUnbondingDelegation` method.

This message is expected to fail if:

* The ERC20 contract is not registered as a multi-staking coin.
* The call to `CancelUnbondingDelegation()` returns an error.