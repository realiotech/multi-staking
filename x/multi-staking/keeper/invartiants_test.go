package keeper_test

import (
	"github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
)

func (suite *KeeperTestSuite) TestModuleAccountInvariants() {
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			name:     "Success",
			malleate: func() {},
			expPass:  true,
		},
	}
	for _, tc := range testCases {
		suite.SetupTest() // reset
		tc.malleate()
		_, broken := keeper.ModuleAccountInvariants(*suite.msKeeper)(suite.ctx)

		if tc.expPass {
			suite.Require().False(broken)
		}
	}
}
