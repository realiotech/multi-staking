## Government Proposal

### Add Bond Coin Proposals

We can make a token to be one of the `bond token` by submiting a `AddMultiStakingCoinProposal`. In this proposal we specified the token's denom and its `BondWeight`, if the proposal is passed the specified token will become a `bond token` with the specified `BondWeight`.

### Change Bond Coin Weight Proposals

We can change a bond token `BondWeight` by submiting a `Proposal`. In this proposal we specified the token's denom and its `BondWeight`, if the proposal is passed the specified token will have its `BondWeight` changed.

### Remove Bond Coin Proposals

We can remove a bond token by submiting a `RemoveMultiStakingCoinProposal`. In this proposal we specified the token's denom, if the proposal is passed the specified token will be remove from the list of bond token