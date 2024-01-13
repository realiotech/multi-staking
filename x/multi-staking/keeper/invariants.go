package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

// RegisterInvariants registers all staking invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "module-accounts",
		ModuleAccountInvariants(k))
	ir.RegisterRoute(types.ModuleName, "validator-lock-denom",
		ValidatorLockDenomInvariants(k))
}

func ModuleAccountInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		totalLockCoinAmount := sdk.NewCoins()

		// calculate lock amount
		lockCoinAmount := sdk.NewCoins()
		k.MultiStakingLockIterator(ctx, func(stakingLock types.MultiStakingLock) bool {
			lockCoinAmount = lockCoinAmount.Add(sdk.NewCoin(stakingLock.LockedCoin.Denom, stakingLock.LockedCoin.Amount))
			return false
		})
		totalLockCoinAmount = totalLockCoinAmount.Add(lockCoinAmount...)

		// calculate unlocking amount
		unlockingCoinAmount := sdk.NewCoins()
		k.MultiStakingUnlockIterator(ctx, func(unlock types.MultiStakingUnlock) bool {
			for _, entry := range unlock.Entries {
				unlockingCoinAmount = unlockingCoinAmount.Add(sdk.NewCoin(entry.UnlockingCoin.Denom, entry.UnlockingCoin.Amount))
			}
			return false
		})
		totalLockCoinAmount = lockCoinAmount.Add(unlockingCoinAmount...)

		moduleAccount := authtypes.NewModuleAddress(types.ModuleName)
		escrowBalances := k.bankKeeper.GetAllBalances(ctx, moduleAccount)

		broken := !escrowBalances.IsEqual(totalLockCoinAmount)

		return sdk.FormatInvariant(
			types.ModuleName,
			"ModuleAccountInvariants",
			fmt.Sprintf(
				"\tescrow coins balances: %v\n"+
					"\ttotal lock coin amount: %v\n",
				escrowBalances, totalLockCoinAmount),
		), broken
	}
}

func ValidatorLockDenomInvariants(k Keeper) sdk.Invariant {

}
