package types

import (
	"cosmossdk.io/math"
	"fmt"
	"sigs.k8s.io/yaml"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewUnbonedMultiStakingEntry(creationHeight int64, completionTime time.Time, rate sdk.Dec, balance math.Int) UnbonedMultiStakingEntry {
	return UnbonedMultiStakingEntry{
		CreationHeight:  creationHeight,
		CompletionTime:  completionTime,
		ConversionRatio: rate,
		InitialBalance:  balance,
		Balance:         balance,
	}
}

// String implements the stringer interface for a UnbonedMultiStakingEntry.
func (e UnbonedMultiStakingEntry) String() string {
	out, _ := yaml.Marshal(e)
	return string(out)
}

// IsMature - is the current entry mature
func (e UnbonedMultiStakingEntry) IsMature(currentTime time.Time) bool {
	return !e.CompletionTime.After(currentTime)
}

// NewUnbondedMultiStaking - create a new unbonding delegation object
//
//nolint:interfacer
func NewUnbondedMultiStaking(
	delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress, creationHeight int64,
	conversionRatio sdk.Dec, minTime time.Time, balance math.Int,
) UnbondedMultiStaking {
	return UnbondedMultiStaking{
		DelegatorAddress: delegatorAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		Entries: []UnbonedMultiStakingEntry{
			NewUnbonedMultiStakingEntry(creationHeight, minTime, conversionRatio, balance),
		},
	}
}

// AddEntry - append entry to the unbonding delegation
func (ubd *UnbondedMultiStaking) AddEntry(creationHeight int64, minTime time.Time, rate sdk.Dec, balance math.Int) {
	// Check the entries exists with creation_height and complete_time
	entryIndex := -1
	for index, ubdEntry := range ubd.Entries {
		if ubdEntry.CreationHeight == creationHeight && ubdEntry.CompletionTime.Equal(minTime) {
			entryIndex = index
			break
		}
	}
	// entryIndex exists
	if entryIndex != -1 {
		ubdEntry := ubd.Entries[entryIndex]
		ubdEntry.Balance = ubdEntry.Balance.Add(balance)
		ubdEntry.InitialBalance = ubdEntry.InitialBalance.Add(balance)

		// update the entry
		ubd.Entries[entryIndex] = ubdEntry
	} else {
		// append the new unbond delegation entry
		entry := NewUnbonedMultiStakingEntry(creationHeight, minTime, rate, balance)
		ubd.Entries = append(ubd.Entries, entry)
	}
}

// RemoveEntry - remove entry at index i to the unbonding delegation
func (ubd *UnbondedMultiStaking) RemoveEntry(i int64) {
	ubd.Entries = append(ubd.Entries[:i], ubd.Entries[i+1:]...)
}

// return the unbonding delegation
func MustMarshalUBD(cdc codec.BinaryCodec, ubd UnbondedMultiStaking) []byte {
	return cdc.MustMarshal(&ubd)
}

// unmarshal a unbonding delegation from a store value
func MustUnmarshalUBD(cdc codec.BinaryCodec, value []byte) UnbondedMultiStaking {
	ubd, err := UnmarshalUBD(cdc, value)
	if err != nil {
		panic(err)
	}

	return ubd
}

// unmarshal a unbonding delegation from a store value
func UnmarshalUBD(cdc codec.BinaryCodec, value []byte) (ubd UnbondedMultiStaking, err error) {
	err = cdc.Unmarshal(value, &ubd)
	return ubd, err
}

// String returns a human readable string representation of an UnbondedMultiStaking.
func (ubd UnbondedMultiStaking) String() string {
	out := fmt.Sprintf(`Unbonding Delegations between:
  Delegator:                 %s
  Validator:                 %s
	Entries:`, ubd.DelegatorAddress, ubd.ValidatorAddress)
	for i, entry := range ubd.Entries {
		out += fmt.Sprintf(`    Unbonding Delegation %d:
      Creation Height:           %v
      Min time to unbond (unix): %v
      Expected balance:          %s`, i, entry.CreationHeight,
			entry.CompletionTime, entry.Balance)
	}

	return out
}

// UnbondedMultiStakings is a collection of UnbondedMultiStaking
type UnbondedMultiStakings []UnbondedMultiStaking

func (ubds UnbondedMultiStakings) String() (out string) {
	for _, u := range ubds {
		out += u.String() + "\n"
	}

	return strings.TrimSpace(out)
}
