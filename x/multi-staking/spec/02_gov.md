## Government Proposal

### Add Bond Token Proposals

We can designate a token as a `bond token` by submiting an `AddBondDenomProposal`. In this proposal, we specify the token's denom and its `BondTokenWeight`, if the proposal passes, the specified token will become a `bond token` with the designated `BondTokenWeight`.

### Change Bond Token Weight Proposals

We can alter the `BondTokenWeight` of a `bond token` by submiting a `ChangeBondTokenWeightProposal`. This proposal requires specifying the `bond token` and the new `BondTokenWeight`, if the proposal is passed the specified `bond token` will have its `BondTokenWeight` changed.

### Remove Bond Token Proposals

We can remove a `bond token` by submiting a `RemoveBondTokenProposal`. This proposal requires specifying the `bond token`.  If the proposal passes, the specified token will be removed from the list of bond tokens.
We can remove a bond token by submitting a `RemoveBondTokenProposal`. In this proposal we specified the token's denom, if the proposal is passed the specified token will be remove from the list of bond token