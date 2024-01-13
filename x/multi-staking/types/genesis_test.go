package types_test

import (
	"testing"

	"github.com/realio-tech/multi-staking-module/testutil"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGenesisState_Validate(t *testing.T) {
	multiStakerAddress := testutil.GenAddress()
	valAddr := testutil.GenValAddress()
	denom := "ario"

	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				MultiStakingLocks: []types.MultiStakingLock{
					{
						LockID: &types.LockID{
							MultiStakerAddr: multiStakerAddress.String(),
							ValAddr:         valAddr.String(),
						},
						LockedCoin: types.MultiStakingCoin{
							Denom:      denom,
							Amount:     sdk.NewInt(1000),
							BondWeight: sdk.NewDec(1),
						},
					},
				},
				ValidatorMultiStakingCoins: []types.ValidatorMultiStakingCoin{
					{
						ValAddr:   valAddr.String(),
						CoinDenom: denom,
					},
				},
				StakingGenesisState: types.DefaultGenesis().StakingGenesisState,
			},
			valid: true,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
