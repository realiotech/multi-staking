package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
	"github.com/stretchr/testify/require"
)

func TestMigrateStakingModule(t *testing.T) {
	file, err := os.Open("state.json")
	require.NoError(t, err)
	bs, err := io.ReadAll(file)
	require.NoError(t, err)

	var oldState AppMap
	err = json.Unmarshal(bs, &oldState)
	require.NoError(t, err)

	//Migrate staking module
	newState, err := migrateStaking(oldState)
	require.NoError(t, err)
	fmt.Println(newState[types.ModuleName])

}

func TestMigrateBankModule(t *testing.T) {
	file, err := os.Open("state.json")
	require.NoError(t, err)
	bs, err := io.ReadAll(file)
	require.NoError(t, err)

	var oldState AppMap
	err = json.Unmarshal(bs, &oldState)
	require.NoError(t, err)

	//Migrate bank module
	_, err = migrateBank(oldState)
	require.NoError(t, err)
}

func TestMigrateDistributionModule(t *testing.T) {
	file, err := os.Open("state.json")
	require.NoError(t, err)
	bs, err := io.ReadAll(file)
	require.NoError(t, err)

	var oldState AppMap
	err = json.Unmarshal(bs, &oldState)
	require.NoError(t, err)

	//Migrate staking module
	newState, err := migrateDistribution(oldState)

	// newStateData, _ := json.MarshalIndent(newState[distrtypes.ModuleName], "", "")
	// _ = os.WriteFile("new_state.json", newStateData, 0644)
	require.NoError(t, err)
	fmt.Println(newState[distrtypes.ModuleName])
}
