# multi-staking-module

The multi-staking-module is a module that allows the cosmos-sdk staking system to support many types of token 

## Multi staking design

Given the fact that several core sdk modules such as distribution or slashing is dependent on the sdk staking module, we design the multi staking module as a wrapper around the sdk staking module so that there's no need to replace the sdk staking module and its related modules.

The multi staking module has the following features:
- Staking with many diffrent type of tokens
- Bond denom selection via Gov proposal
- A validator's delegations can only be in one denom

The mechanism of this module is that it still uses the sdk staking module for all the logic related to staking. But since the sdk staking module doesn't allow multiple bond token/denom, in order to support such feature, the multi-staking module will convert (lock and mint) those different bond token/denom into the one token/denom that is used by the sdk staking module and then stake with the converted token/denom. 

## Concepts and Terms

### Sdk bond token 

Since there're many bond denom/token stake-able via the multi-staking module but only one denom/token used by the underlying sdk staking module, let's refer to the former as `bond token/denom` and the latter as `sdkbond token/denom`.

### Delegation

Each delegation from a `delegator A` to a `validator B` is actually reprensented in the form of a `sdk delegation` which refers to the delegation happened at the sdk staking module layer. In other words, there's little to no logic related to delegation happens at the `multi-staking module` layer as well as delegation data being stored at `multi-staking module` store.

### Intermediary Account

For each delegation from a `delegator A` to a `validator B`, the underlying `sdk delegation` will be managed DIRECTLY by an unique `intermediary account C` instead of the `delegator A`, meaning that the `sdk delegator` has the `intermediary account` as its delegator. The `delegator A` though, can still d

Each delegation made by `delegator A` to a validator (who hasn't been delegated by `delegator A` before) will trigger 
For delegating, users won't `sdk delegate` directly with their account but through accounts called Intermediary Accounts which is managed by the multi-staking module. Each delegation made by `delegator A` to `validator B` will be represented in the form of a  created and managed by an unique `Intermediary Account C`. Everytime `delegator A` delegate to `validator B`, their delegation tokens will be locked by sending to the `Intermediary Account C`, then the multi-staking module will mint a calculated amount of `sdkbond token/denom` to the `Intermediary Account C` so that it can `sdk delegate` on `delegator A` behalf. The multi-staking module interior logic allows `delegator A` to dictate actions (such as unbonding or withdrawing reward) involving delegations made by `Intermediary Account C` so that `delegator A` still get full control of the delegation.

### Bond Token Weight

Each `bond token` is associated with a `bond token weight`. This `bond token weight` is specified via the gov proposal in which the `bond token` is accepted.
We mentioned above that for each delegation the multi-staking will lock the `bond token` and mint a calculated ammount of `sdkbond token`. The calculation is a multiplication : minted sdkbond token ammount = bond token amount * bond token weight.

