# End-Block

## Complete Unbonding Delegations

### Calculate total `UnbondedAmount`

* Retrieve `matureUnbondingDelegations` which is the array of all `UnbondingDelegations` that complete in this block

### Staking module EndBlock

* Call `Staking` module `EndBlock` to `CompleteUnbonding`

### MultiStaking module EndBlock

* Iterate over `matureUnbondingDelegations` which was retrieve above

* For each iteration, we will:

    * Calculate amount of `unlockedCoin` that will be return to user by multiply the amount of `unbonded coin` and `bonded weight`

    * Burn the `remainingCoin` that remain on the `Lock` after send `unlockedCoin` to user

    * Check if `unlockedCoin` is registered as erc20 token pair. If yes, convert it back to erc20 token in evm state. If no, that mean it is cosmos native token already in user account, we will skip the conversion step.

    * Delete `UnlockEntry`.

    
