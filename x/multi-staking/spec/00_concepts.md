# multi-staking-module

The multi-staking-module is a module that allows the cosmos-sdk staking system to support many types of token 

## Multi staking design

Given the fact that several core sdk modules such as distribution or slashing is dependent on the sdk staking module, we design the multi staking module as a wrapper around the sdk staking module so that there's no need to replace the sdk staking module and its related modules.

The multi staking module has the following features:
- Staking with many diffrent type of tokens
- Bond denom selection via Gov proposal
- A validator's delegations can only be in one denom

This module

Since we're working around the sdk staking module, this module frequently calls the staking api. 
Most of the logic () will be handled by the sdk staking module. This module 

## Terms

### StakingBondToken 
