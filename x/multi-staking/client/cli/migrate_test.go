package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func TestMigrateStakingModule(t *testing.T) {
	file, err := os.Open("state.json")
	require.NoError(t, err)
	bs, err := io.ReadAll(file)
	require.NoError(t, err)

	var oldState AppMap
	err = json.Unmarshal(bs, &oldState)
	fmt.Println(oldState["staking"])
	require.NoError(t, err)

	//Migrate staking module
	newState, err := migrateStaking(oldState)
	require.NoError(t, err)
	fmt.Println(newState[types.ModuleName])
}