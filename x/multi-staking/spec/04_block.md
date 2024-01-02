# Begin-Block

## Complete Unbonding Delegations

Check if there's any completed unbonding delegations. 
If so, for each of the unbonding delegation:

* Get the `delegator account` from `IntermediaryDelegatorDelegator` store.

* Update `CompletedDelegations`.

# End-Block

## Complete Unbonding Delegations

Check if there's any entries in `CompletedDelegations`.
If so, for each entry:

* Calculate the amount of `bond coin` to be unlocked.

* Send the calculated amount of `bond coin` from `IntermediaryDelegator` to `delegator`

* Update `DVPairSDKBondCoins`.

* Delete the entry in `CompletetedDelegations`.
