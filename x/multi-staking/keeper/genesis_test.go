package keeper_test

import (
	dbm "github.com/cosmos/cosmos-db"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/realio-tech/multi-staking-module/test/simapp"

	"cosmossdk.io/log"
	abci "github.com/cometbft/cometbft/abci/types"
)

func (suite *KeeperTestSuite) TestImportExportGenesis() {
	appState, err := suite.app.ExportAppStateAndValidators(false, []string{})
	suite.NoError(err)

	encConfig := simapp.MakeEncodingConfig()

	configurator := evmtypes.NewEVMConfigurator()
	configurator.ResetTestConfig()

	emptyApp := simapp.NewSimApp(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		map[int64]bool{},
		"temp",
		simapp.FlagPeriodValue,
		encConfig,
		simapp.EmptyAppOptions{},
	)

	_, err = emptyApp.InitChain(
		&abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: simapp.DefaultConsensusParams,
			AppStateBytes:   appState.AppState,
		},
	)
	suite.NoError(err)
	_, err = emptyApp.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height: emptyApp.LastBlockHeight() + 1,
		Hash:   emptyApp.LastCommitID().Hash,
	})
	suite.NoError(err)

	_, err = emptyApp.ExportAppStateAndValidators(false, []string{})
	suite.NoError(err)
}
