# multi-staking-module

The multi-staking-module is a module that allows the cosmos-sdk staking system to support many types of token 

## Features

- Staking with many diffrent type of tokens
- Bond denom selection via Gov proposal
- A validator can only be delegated using its bonded denom
- All user usecases of the sdk staking module

## Multi staking design

Given the fact that several core sdk modules such as distribution or slashing is dependent on the sdk staking module, we design the multi staking module as a wrapper around the sdk staking module so that there's no need to replace the sdk staking module and its related modules.

The mechanism of this module is that it still uses the sdk staking module for all the logic related to staking. But since the sdk staking module doesn't allow multiple bond token/denom, in order to support such feature, the multi-staking module will convert (lock and mint) those different bond token/denom into the one token/denom that is used by the sdk staking module and then stake with the converted token/denom. 

## Concepts and Terms

### Sdk bond token 

Since there're many bond denom/token stake-able via the multi-staking module but only one denom/token used by the underlying sdk staking module, let's refer to the former as `bond token/denom` and the latter as `sdkbond token/denom`.

### Delegation

Each delegation from a `delegator A` is actually reprensented in the form of a `sdk delegation` which refers to the delegation happened at the sdk staking module layer. In other words, there's little to no logic related to the actual delegation system (validator power distr, slashing, distributing rewards...) happens at the `multi-staking module` layer as well as delegation data being stored at `multi-staking module` store.

### Intermediary Account

For each delegation from a `delegator A`, the underlying `sdk delegation` will be created and managed DIRECTLY by an unique `intermediary account C` instead of the `delegator A`, meaning that the `sdk delegation` will have the `intermediary account` as its delegator. The `delegator A` though, can still dictate what `intermediary account C` on what to do with the `sdk delegation` so that `delegator A` still have full controll over the delegation.

The `intermediary account` is also where the `bond token` from `delegator` is locked and the `sdkbond token` is minted to, the minted `sdkbond token` will then be used to create the `sdk delegation`.

### Bond Token Weight

Each `bond token` is associated with a `bond token weight`. This `bond token weight` is specified via the gov proposal in which the `bond token` is accepted.

We mentioned above that for each delegation the multi-staking will lock the `bond token` and mint a calculated ammount of `sdkbond token`. The calculation here is a multiplication : minted sdkbond token ammount = bond token amount * bond token weight.

