package types

import (
	"fmt"

	"sigs.k8s.io/yaml"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewUnlockEntry(creationHeight int64, weightedCoin MultiStakingCoin) UnlockEntry {
	return UnlockEntry{
		CreationHeight: creationHeight,
		UnlockingCoin:  weightedCoin,
	}
}

// String implements the stringer interface for a UnlockEntry.
func (e UnlockEntry) String() string {
	out, _ := yaml.Marshal(e)
	return string(out)
}

// NewMultiStakingUnlock - create a new MultiStaking unlock object
//
//nolint:interfacer
func NewMultiStakingUnlock(
	creationHeight int64, weightedCoin MultiStakingCoin,
) MultiStakingUnlock {
	return MultiStakingUnlock{
		Entries: []UnlockEntry{
			NewUnlockEntry(creationHeight, weightedCoin),
		},
	}
}

// AddEntry - append entry to the unbonding delegation
func (unlock *MultiStakingUnlock) AddEntry(creationHeight int64, weightedCoin MultiStakingCoin) {
	// Check the entries exists with creation_height and complete_time
	entryIndex := -1
	for index, ubdEntry := range unlock.Entries {
		if ubdEntry.CreationHeight == creationHeight {
			entryIndex = index
			break
		}
	}
	// entryIndex exists
	if entryIndex != -1 {
		ubdEntry := unlock.Entries[entryIndex]
		ubdEntry.UnlockingCoin = ubdEntry.UnlockingCoin.Add(weightedCoin)

		// update the entry
		unlock.Entries[entryIndex] = ubdEntry
	} else {
		// append the new unbond delegation entry
		entry := NewUnlockEntry(creationHeight, weightedCoin)
		unlock.Entries = append(unlock.Entries, entry)
	}
}

// RemoveEntry - remove entry at index i to the unbonding delegation
func (unlock *MultiStakingUnlock) RemoveEntryAtCreationHeight(creationHeight int64) {
	// Check the entries exists with creation_height and complete_time
	entryIndex := -1
	for index, ubdEntry := range unlock.Entries {
		if ubdEntry.CreationHeight == creationHeight {
			entryIndex = index
			break
		}
	}
	// entryIndex exists
	if entryIndex != -1 {
		unlock.Entries = append(unlock.Entries[:entryIndex], unlock.Entries[entryIndex+1:]...)
	}
}

func (u UnlockEntry) GetBondWeight() sdk.Dec {
	return u.UnlockingCoin.BondWeight
}

func (unlockEntry UnlockEntry) UnbondAmountToUnlockAmount(unbondAmount math.Int) math.Int {
	return sdk.NewDecFromInt(unbondAmount).Quo(unlockEntry.GetBondWeight()).RoundInt()
}

func (unlockEntry UnlockEntry) UnlockAmountToUnbondAmount(unlockAmount math.Int) math.Int {
	return unlockEntry.GetBondWeight().MulInt(unlockAmount).RoundInt()
}

// String returns a human readable string representation of an MultiStakingUnlock.
func (unlock MultiStakingUnlock) String() string {
	out := fmt.Sprintf(`Unlock ID: %s
	Entries:`, unlock.UnlockID)
	for i, entry := range unlock.Entries {
		out += fmt.Sprintf(`    Unbonding Delegation %d:
      Creation Height:           %v
     `, i, entry.CreationHeight,
		)
	}

	return out
}

// MultiStakingUnlocks is a collection of MultiStakingUnlock
// type MultiStakingUnlocks []UnlockEntry

// func (ubds MultiStakingUnlocks) String() (out string) {
// 	for _, u := range ubds {
// 		out += u.String() + "\n"
// 	}

// 	return strings.TrimSpace(out)
// }

// func NewUnbonedMultiStakingRecord( // ?
// 	delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, creationHeight int64,
// 	completionTime time.Time, rate sdk.Dec, balance math.Int,
// ) Unlock {
// 	return UnbonedMultiStakingRecord{
// 		CreationHeight:  creationHeight,
// 		CompletionTime:  completionTime,
// 		ConversionRatio: rate,
// 		InitialBalance:  balance,
// 		Balance:         balance,
// 	}
// }

// // String implements the stringer interface for a UnlockEntry.
// func (e UnbonedMultiStakingRecord) String() string {
// 	out, _ := yaml.Marshal(e)
// 	return string(out)
// }
