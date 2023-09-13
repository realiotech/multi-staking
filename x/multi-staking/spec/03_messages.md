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

1. Setting `ValidatorBondDenom`.

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
The `multi-staking` inturns mint a calculated amount of `sdkstaking.bondtoken` and
create an `IntermediaryAccount` to delegate on behalf of the delegator.

Logic flow:

* Get `IntermediaryAccount` for the delegator.

* Set `IntermediaryAccountDelegator` if it's not set yet.

* Send delegated coins from user to `IntermediaryAccount`.

* Caculate the `sdkbond token` to be minted using `BondTokenWeight`.

* Mint `sdkbond token` to `IntermediaryAccount`

* Update `DVPairSDKBondCoins`.

* Update `DelegatorTotalSDKBondToken`.

* Create `sdk delegation` with `IntermediaryAccount` using the minted `sdkbond token`

## MsgUndelegate

The `MsgUndelegate` message allows delegators to undelegate their tokens from
validator.

Logic flow:

* Calculate ammount of `sdkbond token` need to be `sdk undelegated`

* Call `stakingkeeper.Undelegate()` with the calculated amount of `sdkbond token`

* Update `DelegatorTotalSDKBondToken`.

The rest of the unbonding logic such as sending locked coins back to user will happens at `EndBlock()`

## MsgCancelUnbondingDelegation 

The `MsgCancelUnbondingDelegation` message allows delegators to cancel the `unbondingDelegation` entry and deleagate back to a previous validator.

Logic flow:

* Calculate amount of `sdkbond token` need to be `sdk cancel undelegation`

* Call `stakingkeeper.CancelUnbondingDelegation()` with the calculated amount of `sdkbond token`

* Update `DelegatorTotalSDKBondToken`.

## MsgBeginRedelegate

The `MsgBeginRedelegate` message allows delegators to instantly switch validators. Once
the unbonding period has passed, the redelegation is automatically completed in
the EndBlocker.

Logic flow:

* Calculate amount of `sdkbond token` need to be `sdk redelegate`

* Call `stakingkeeper.BeginRedelegate()` with the calculated amount of `sdkbond token`

* Update `DVPairSDKBondTokens`