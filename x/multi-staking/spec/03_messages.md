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

1. Setting `ValidatorBondDenom` in state.

2. Converting `MsgCreateValidator` to `stakingtypes.MsgCreateValidator` and
calling `stakingkeeper.CreateValidator()`.

This message is expected to fail if:

* `ValOperatorAddr` already exists in state.
* The call to `stakingkeeper.CreateValidator()` returns an error.

## MsgEditValidator

The `Description`, `CommissionRate` of a validator can be updated using the
`MsgEditValidator` message.

Logic flow:

1. Use `SdkCreateValidator()` to create `stakingtypes.MsgEditValidator`, calling `stakingkeeper.EditValidator()`

This message is expected to fail if:

* The call to `stakingkeeper.EditValidator()` returns an error.

## MsgDelegate

Within this message the delegator locked up coins in the `multi-staking` module account. 
The `multi-staking` inturns mint a calculated amount of `sdkstaking.bondtoken` and
create an `IntermediaryAccount` to delegate on behalf of the delegator.

Logic flow:

* Create/Get `IntermediaryAccount` for the delegation.

* Send delegated coins from user to `IntermediaryAccount`.

* Caculate the `sdktaking.bondtoken` to be minted using `BondTokenWeight`.
amountMinted = delegatedCoins * bondTokenWeight

* 













