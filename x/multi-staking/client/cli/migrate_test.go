package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/realio-tech/multi-staking-module/testing/simapp"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
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

func TestMigrateFull(t *testing.T) {
	genDoc, err := tmtypes.GenesisDocFromFile("oldState.json")
	require.NoError(t, err)
	app := simapp.NewSimApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, map[int64]bool{}, "", 1, simapp.MakeEncodingConfig(), simapp.EmptyAppOptions{})
	res := app.InitChain(
		abci.RequestInitChain{
			Validators: []abci.ValidatorUpdate{},
			ConsensusParams: &abci.ConsensusParams{
				Block: &abci.BlockParams{
					MaxBytes: genDoc.ConsensusParams.Block.MaxBytes,
					MaxGas:   genDoc.ConsensusParams.Block.MaxGas,
				},
				Evidence:  &genDoc.ConsensusParams.Evidence,
				Validator: &genDoc.ConsensusParams.Validator,
				Version:   &genDoc.ConsensusParams.Version,
			},
			AppStateBytes: genDoc.AppState,
		},
	)
	fmt.Println("res", res)
}
