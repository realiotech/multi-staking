package types

import (
	"fmt"
	"strings"
	"time"

	"cosmossdk.io/math"
	"sigs.k8s.io/yaml"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewUnlockEntry(creationHeight int64, rate sdk.Dec, balance math.Int) UnlockEntry {
	return UnlockEntry{
		CreationHeight:  creationHeight,
		ConversionRatio: rate,
		Balance:         balance,
	}
}

// String implements the stringer interface for a UnlockEntry.
func (e UnlockEntry) String() string {
	out, _ := yaml.Marshal(e)
	return string(out)
}

// NewMultiStakingUnlock - create a new unbonding delegation object
//
//nolint:interfacer
func NewMultiStakingUnlock(
	delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress,
	creationHeight int64, conversionRatio sdk.Dec, balance math.Int,
) MultiStakingUnlock {
	return MultiStakingUnlock{
		DelegatorAddress: delegatorAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		Entries: []UnlockEntry{
			NewUnlockEntry(creationHeight, conversionRatio, balance),
		},
	}
}

// AddEntry - append entry to the unbonding delegation
func (ubd *MultiStakingUnlock) AddEntry(creationHeight int64, rate sdk.Dec, balance math.Int) {
	// Check the entries exists with creation_height and complete_time
	entryIndex := -1
	for index, ubdEntry := range ubd.Entries {
		if ubdEntry.CreationHeight == creationHeight {
			entryIndex = index
			break
		}
	}
	// entryIndex exists
	if entryIndex != -1 {
		ubdEntry := ubd.Entries[entryIndex]
		ubdEntry.Balance = ubdEntry.Balance.Add(balance)

		// update the entry
		ubd.Entries[entryIndex] = ubdEntry
	} else {
		// append the new unbond delegation entry
		entry := NewUnlockEntry(creationHeight, rate, balance)
		ubd.Entries = append(ubd.Entries, entry)
	}
}

// RemoveEntry - remove entry at index i to the unbonding delegation
func (ubd *MultiStakingUnlock) RemoveEntry(i int64) {
	ubd.Entries = append(ubd.Entries[:i], ubd.Entries[i+1:]...)
}

// String returns a human readable string representation of an MultiStakingUnlock.
func (ubd MultiStakingUnlock) String() string {
	out := fmt.Sprintf(`Unbonding Delegations between:
  Delegator:                 %s
  Validator:                 %s
	Entries:`, ubd.DelegatorAddress, ubd.ValidatorAddress)
	for i, entry := range ubd.Entries {
		out += fmt.Sprintf(`    Unbonding Delegation %d:
      Creation Height:           %v
      Expected balance:          %s`, i, entry.CreationHeight,
			entry.Balance)
	}

	return out
}

// MultiStakingUnlocks is a collection of MultiStakingUnlock
type MultiStakingUnlocks []UnbonedMultiStakingRecord

func (ubds MultiStakingUnlocks) String() (out string) {
	for _, u := range ubds {
		out += u.String() + "\n"
	}

	return strings.TrimSpace(out)
}

func NewUnbonedMultiStakingRecord( // ?
	delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, creationHeight int64,
	completionTime time.Time, rate sdk.Dec, balance math.Int,
) UnbonedMultiStakingRecord {
	return UnbonedMultiStakingRecord{
		CreationHeight:  creationHeight,
		CompletionTime:  completionTime,
		ConversionRatio: rate,
		InitialBalance:  balance,
		Balance:         balance,
	}
}

// String implements the stringer interface for a UnlockEntry.
func (e UnbonedMultiStakingRecord) String() string {
	out, _ := yaml.Marshal(e)
	return string(out)
}
