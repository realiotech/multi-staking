package keeper_test

import (
	dbm "github.com/cosmos/cosmos-db"
	"github.com/realio-tech/multi-staking-module/test/simapp"

	"cosmossdk.io/log"

	abci "github.com/cometbft/cometbft/abci/types"
)

func (suite *KeeperTestSuite) TestImportExportGenesis() {
	appState, err := suite.app.ExportAppStateAndValidators(false, []string{}, []string{})
	suite.NoError(err)

	encConfig := simapp.MakeTestEncodingConfig()

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
			ConsensusParams: &appState.ConsensusParams,
			AppStateBytes:   appState.AppState,
		},
	)
	suite.NoError(err)

	_, err = emptyApp.FinalizeBlock(&abci.RequestFinalizeBlock{Height: emptyApp.LastBlockHeight() + 1})
	suite.NoError(err)

	newAppState, err := emptyApp.ExportAppStateAndValidators(false, []string{}, []string{})
	suite.NoError(err)

	_, err = suite.app.FinalizeBlock(&abci.RequestFinalizeBlock{Height: suite.app.LastBlockHeight() + 1})
	suite.NoError(err)
	_, err = suite.app.Commit()
	suite.NoError(err)
	appState2, err := suite.app.ExportAppStateAndValidators(false, []string{}, []string{})
	suite.NoError(err)

	suite.Equal(appState2.AppState, newAppState.AppState)
}
