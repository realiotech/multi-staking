package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/testutil"
)

func (suite *KeeperTestSuite) TestSetBondTokenWeight() {
	suite.SetupTest()

	gasDenom := "ario"
	govDenom := "arst"
	gasWeight := math.LegacyNewDec(1)
	govWeight := math.LegacyNewDecWithPrec(2, 4)

	suite.msKeeper.SetBondTokenWeight(suite.ctx, gasDenom, gasWeight)
	suite.msKeeper.SetBondTokenWeight(suite.ctx, govDenom, govWeight)

	suite.Equal(gasWeight, suite.msKeeper.GetBondTokenWeight(suite.ctx, gasDenom))
	suite.Equal(govWeight, suite.msKeeper.GetBondTokenWeight(suite.ctx, govDenom))
}

func (suite *KeeperTestSuite) TestSetValidatorBondDenom() {
	suite.SetupTest()

	gasDenom := "ario"
	govDenom := "arst"

	valA := testutil.GenValAddress()
	valB := testutil.GenValAddress()

	suite.msKeeper.SetValidatorBondDenom(suite.ctx, valA, gasDenom)
	suite.msKeeper.SetValidatorBondDenom(suite.ctx, valB, govDenom)

	suite.Equal(gasDenom, suite.msKeeper.GetValidatorBondDenom(suite.ctx, valA))
	suite.Equal(govDenom, suite.msKeeper.GetValidatorBondDenom(suite.ctx, valB))

	suite.msKeeper.SetValidatorBondDenom(suite.ctx, valA, govDenom)
	suite.Equal(govDenom, suite.msKeeper.GetValidatorBondDenom(suite.ctx, valA))
}

func (suite *KeeperTestSuite) TestSetIntermediaryAccountDelegator() {
	suite.SetupTest()

	delA := testutil.GenAddress()
	delB := testutil.GenAddress()
	imAddrressA := testutil.GenAddress()
	imAddrressB := testutil.GenAddress()

	suite.msKeeper.SetIntermediaryAccountDelegator(suite.ctx, imAddrressA, delA)
	suite.msKeeper.SetIntermediaryAccountDelegator(suite.ctx, imAddrressB, delB)

	suite.Equal(delA, suite.msKeeper.GetIntermediaryAccountDelegator(suite.ctx, imAddrressA))
	suite.Equal(delB, suite.msKeeper.GetIntermediaryAccountDelegator(suite.ctx, imAddrressB))

	suite.msKeeper.SetIntermediaryAccountDelegator(suite.ctx, imAddrressB, delA)
	suite.Equal(delA, suite.msKeeper.GetIntermediaryAccountDelegator(suite.ctx, imAddrressB))
}

func (suite *KeeperTestSuite) TestSetDVPairSDKBondTokens() {
	suite.SetupTest()

	delA := testutil.GenAddress()
	delB := testutil.GenAddress()
	valA := testutil.GenValAddress()
	valB := testutil.GenValAddress()

	bondSDKAmountA := sdk.NewCoin("ario", sdk.NewInt(100))
	bondSDKAmountB := sdk.NewCoin("ario", sdk.NewInt(200))

	suite.msKeeper.SetDVPairSDKBondTokens(suite.ctx, delA, valA, bondSDKAmountA)
	suite.msKeeper.SetDVPairSDKBondTokens(suite.ctx, delB, valB, bondSDKAmountB)

	suite.Equal(bondSDKAmountA, suite.msKeeper.GetDVPairSDKBondTokens(suite.ctx, delA, valA))
	suite.Equal(bondSDKAmountB, suite.msKeeper.GetDVPairSDKBondTokens(suite.ctx, delB, valB))

	suite.msKeeper.SetDVPairSDKBondTokens(suite.ctx, delA, valB, bondSDKAmountB)
	suite.Equal(bondSDKAmountB, suite.msKeeper.GetDVPairSDKBondTokens(suite.ctx, delA, valB))
}

func (suite *KeeperTestSuite) TestSetDVPairBondTokens() {
	suite.SetupTest()

	delA := testutil.GenAddress()
	delB := testutil.GenAddress()
	valA := testutil.GenValAddress()
	valB := testutil.GenValAddress()

	bondAmountA := sdk.NewCoin("ario", sdk.NewInt(100))
	bondAmountB := sdk.NewCoin("arst", sdk.NewInt(200))

	suite.msKeeper.SetDVPairBondTokens(suite.ctx, delA, valA, bondAmountA)
	suite.msKeeper.SetDVPairBondTokens(suite.ctx, delB, valB, bondAmountB)

	suite.Equal(bondAmountA, suite.msKeeper.GetDVPairBondTokens(suite.ctx, delA, valA))
	suite.Equal(bondAmountB, suite.msKeeper.GetDVPairBondTokens(suite.ctx, delB, valB))

	suite.msKeeper.SetDVPairBondTokens(suite.ctx, delA, valB, bondAmountB)
	suite.Equal(bondAmountB, suite.msKeeper.GetDVPairBondTokens(suite.ctx, delA, valB))
}
