package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
	
	"github.com/stretchr/testify/require"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
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