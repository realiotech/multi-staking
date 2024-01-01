## Government Proposal

### Add Bond Token Proposals

We can designate a token as a `bond token` by submiting an `AddBondDenomProposal`. In this proposal, we specify the token's denom and its `BondWeight`, if the proposal passes, the specified token will become a `bond token` with the designated `BondWeight`.

### Change Bond Token Weight Proposals

We can alter the `BondWeight` of a `bond token` by submiting a `UpdateBondWeightProposal`. This proposal requires specifying the `bond token` and the new `BondWeight`, if the proposal is passed the specified `bond token` will have its `BondWeight` changed.

### Remove Bond Token Proposals

We can remove a `bond token` by submiting a `RemoveBondTokenProposal`. This proposal requires specifying the `bond token`.  If the proposal passes, the specified token will be removed from the list of bond tokens.