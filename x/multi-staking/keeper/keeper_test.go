package keeper_test

import (
	"github.com/stretchr/testify/suite"
	"testing"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/realio-tech/multi-staking-module/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/testing/simapp"
	multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	multistakingtypes "github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx      sdk.Context
	msKeeper *multistakingkeeper.Keeper
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(false)

	storeKey := storetypes.NewKVStoreKey(multistakingtypes.StoreKey)
	memKey := storetypes.NewMemoryStoreKey(multistakingtypes.MemStoreKey)
	encCfg := simapp.MakeTestEncodingConfig()
	testCtx := testutil.DefaultContext(storeKey, memKey, storetypes.NewTransientStoreKey("transient_test"))

	msKeeper := multistakingkeeper.NewKeeper(
		encCfg.Marshaler,
		app.StakingKeeper,
		storeKey,
		memKey,
	)

	suite.ctx, suite.msKeeper = testCtx, msKeeper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
