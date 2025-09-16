## Government Proposal

### Add Bond Token Proposals (for cosmos base coin)

We can designate a token as a `multistaking coin` by submiting an `AddMultiStakingCoinProposal`. In this proposal, we specify the token's `denom` and its `BondWeight`, if the proposal passes, the specified token will become a `multistaking coin` with the designated `BondWeight`.

### Add EVM Bond Token Proposals (for erc20 coin)

We can designate an ERC20 token as a `multistaking coin` by submiting an `AddMultiStakingEVMCoinProposal`. In this proposal, we specify the ERC20 token's `contract_address` and its `BondWeight`. If the proposal passes, the specified ERC20 token will be registered in the ERC20 module and become a `multistaking coin` with the designated `BondWeight`.

The proposal performs the following actions:
1. Checks if the ERC20 contract address is already registered
2. Registers the ERC20 token using the ERC20 module's `RegisterERC20` functionality
3. Adds the token as a multi-staking coin with the specified bond weight

### Change Bond Token Weight Proposals

We can alter the `BondWeight` of a `multistaking coin` by submiting a `UpdateBondWeightProposal`. This proposal requires specifying `denom` of the `multistaking coin` and the new `BondWeight`, if the proposal is passed the specified `multistaking coin` have its `BondWeight` changed to new value that decleared by the proposal.

### Remove Bond Token Proposals (for cosmos base coin)

We can remove a `multistaking coin` boned token by submiting an `RemoveMultiStakingCoinProposal`. In this proposal, we specify the token's `denom`. If the proposal passes, we will force undelegate all the delegation of the removed bond token and remove the bond token from store.