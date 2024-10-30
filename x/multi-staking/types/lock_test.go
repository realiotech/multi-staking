package types_test

import (
	"testing"

	"github.com/realio-tech/multi-staking-module/test"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"
)

var (
	MultiStakingDenomA = "ario"
	MultiStakingDenomB = "arst"
)

func TestAddCoinToMultiStakingLock(t *testing.T) {
	valAddr := test.GenValAddress()
	delAddr := test.GenAddress()

	testCases := []struct {
		name         string
		originMSCoin types.MultiStakingCoin
		addingMSCoin types.MultiStakingCoin
		expMSCoin    types.MultiStakingCoin
		expErr       bool
	}{
		{
			name:         "success",
			originMSCoin: types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(100000), math.LegacyMustNewDecFromStr("0.3")),
			addingMSCoin: types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(23456), math.LegacyMustNewDecFromStr("0.3")),
			expMSCoin:    types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(123456), math.LegacyMustNewDecFromStr("0.3")),
			expErr:       false,
		},
		{
			name:         "success and change bond weight",
			originMSCoin: types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(100000), math.LegacyOneDec()),
			addingMSCoin: types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(200000), math.LegacyMustNewDecFromStr("0.25")),
			expMSCoin:    types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(300000), math.LegacyMustNewDecFromStr("0.5")),
			expErr:       false,
		},
		{
			name:         "success from zero coin",
			originMSCoin: types.NewMultiStakingCoin(MultiStakingDenomA, math.ZeroInt(), math.LegacyMustNewDecFromStr("0.3")),
			addingMSCoin: types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(23456), math.LegacyMustNewDecFromStr("0.3")),
			expMSCoin:    types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(23456), math.LegacyMustNewDecFromStr("0.3")),
			expErr:       false,
		},
		{
			name:         "denom mismatch",
			originMSCoin: types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(100000), math.LegacyMustNewDecFromStr("0.3")),
			addingMSCoin: types.NewMultiStakingCoin(MultiStakingDenomB, math.NewInt(23456), math.LegacyMustNewDecFromStr("0.3")),
			expErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lockID := types.MultiStakingLockID(delAddr.String(), valAddr.String())
			lockRecord := types.NewMultiStakingLock(lockID, tc.originMSCoin)

			err := lockRecord.AddCoinToMultiStakingLock(tc.addingMSCoin)

			if tc.expErr {
				require.Error(t, err, tc.name)
			} else {
				require.NoError(t, err)
				require.Equal(t, lockRecord.LockedCoin.Amount, tc.expMSCoin.Amount)
				require.Equal(t, lockRecord.LockedCoin.Denom, tc.expMSCoin.Denom)
				require.Equal(t, lockRecord.LockedCoin.BondWeight, tc.expMSCoin.BondWeight)
			}
		})
	}
}

func TestRemoveCoinFromMultiStakingLock(t *testing.T) {
	valAddr := test.GenValAddress()
	delAddr := test.GenAddress()

	testCases := []struct {
		name         string
		originMSCoin types.MultiStakingCoin
		removeMSCoin types.MultiStakingCoin
		expMSCoin    types.MultiStakingCoin
		expErr       bool
	}{
		{
			name:         "success",
			originMSCoin: types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(123456), math.LegacyMustNewDecFromStr("0.3")),
			removeMSCoin: types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(23456), math.LegacyMustNewDecFromStr("0.3")),
			expMSCoin:    types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(100000), math.LegacyMustNewDecFromStr("0.3")),
			expErr:       false,
		},
		{
			name:         "denom mismatch",
			originMSCoin: types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(100000), math.LegacyMustNewDecFromStr("0.3")),
			removeMSCoin: types.NewMultiStakingCoin(MultiStakingDenomB, math.NewInt(23456), math.LegacyMustNewDecFromStr("0.3")),
			expErr:       true,
		},
		{
			name:         "insufficient amount",
			originMSCoin: types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(100000), math.LegacyMustNewDecFromStr("0.3")),
			removeMSCoin: types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(234567), math.LegacyMustNewDecFromStr("0.3")),
			expErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lockID := types.MultiStakingLockID(delAddr.String(), valAddr.String())
			lockRecord := types.NewMultiStakingLock(lockID, tc.originMSCoin)

			err := lockRecord.RemoveCoinFromMultiStakingLock(tc.removeMSCoin)

			if tc.expErr {
				require.Error(t, err, tc.name)
			} else {
				require.NoError(t, err)
				require.Equal(t, lockRecord.LockedCoin.Amount, tc.expMSCoin.Amount)
				require.Equal(t, lockRecord.LockedCoin.Denom, tc.expMSCoin.Denom)
				require.Equal(t, lockRecord.LockedCoin.BondWeight, tc.expMSCoin.BondWeight)
			}
		})
	}
}

func TestMoveCoinToLock(t *testing.T) {
	valAddrA := test.GenValAddress()
	valAddrB := test.GenValAddress()

	delAddr := test.GenAddress()

	testCases := []struct {
		name          string
		fromMSCoin    types.MultiStakingCoin
		toMSCoin      types.MultiStakingCoin
		moveMSCoin    types.MultiStakingCoin
		expFromMSCoin types.MultiStakingCoin
		expToMSCoin   types.MultiStakingCoin
		expErr        bool
	}{
		{
			name:          "success",
			fromMSCoin:    types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(123456), math.LegacyMustNewDecFromStr("0.3")),
			toMSCoin:      types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(100000), math.LegacyMustNewDecFromStr("0.3")),
			moveMSCoin:    types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(23456), math.LegacyMustNewDecFromStr("0.3")),
			expFromMSCoin: types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(100000), math.LegacyMustNewDecFromStr("0.3")),
			expToMSCoin:   types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(123456), math.LegacyMustNewDecFromStr("0.3")),
			expErr:        false,
		},
		{
			name:          "success and change rate",
			fromMSCoin:    types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(323456), math.LegacyMustNewDecFromStr("0.5")),
			toMSCoin:      types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(100000), math.LegacyOneDec()),
			moveMSCoin:    types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(300000), math.LegacyMustNewDecFromStr("0.5")),
			expFromMSCoin: types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(23456), math.LegacyMustNewDecFromStr("0.5")),
			expToMSCoin:   types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(400000), math.LegacyMustNewDecFromStr("0.625")),
			expErr:        false,
		},
		{
			name:       "denom mismatch at fromLock",
			fromMSCoin: types.NewMultiStakingCoin(MultiStakingDenomB, math.NewInt(123456), math.LegacyMustNewDecFromStr("0.3")),
			toMSCoin:   types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(100000), math.LegacyMustNewDecFromStr("0.3")),
			moveMSCoin: types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(23456), math.LegacyMustNewDecFromStr("0.3")),
			expErr:     true,
		},
		{
			name:       "denom mismatch at toLock",
			fromMSCoin: types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(123456), math.LegacyMustNewDecFromStr("0.3")),
			toMSCoin:   types.NewMultiStakingCoin(MultiStakingDenomB, math.NewInt(100000), math.LegacyMustNewDecFromStr("0.3")),
			moveMSCoin: types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(23456), math.LegacyMustNewDecFromStr("0.3")),
			expErr:     true,
		},
		{
			name:       "denom mismatch at move coin",
			fromMSCoin: types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(123456), math.LegacyMustNewDecFromStr("0.3")),
			toMSCoin:   types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(100000), math.LegacyMustNewDecFromStr("0.3")),
			moveMSCoin: types.NewMultiStakingCoin(MultiStakingDenomB, math.NewInt(23456), math.LegacyMustNewDecFromStr("0.3")),
			expErr:     true,
		},
		{
			name:       "insufficient amount",
			fromMSCoin: types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(123456), math.LegacyMustNewDecFromStr("0.3")),
			toMSCoin:   types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(200000), math.LegacyMustNewDecFromStr("0.3")),
			moveMSCoin: types.NewMultiStakingCoin(MultiStakingDenomA, math.NewInt(234567), math.LegacyMustNewDecFromStr("0.3")),
			expErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lockID1 := types.MultiStakingLockID(delAddr.String(), valAddrA.String())
			lockRecord1 := types.NewMultiStakingLock(lockID1, tc.fromMSCoin)

			lockID2 := types.MultiStakingLockID(delAddr.String(), valAddrB.String())
			lockRecord2 := types.NewMultiStakingLock(lockID2, tc.toMSCoin)

			err := lockRecord1.MoveCoinToLock(&lockRecord2, tc.moveMSCoin)

			if tc.expErr {
				require.Error(t, err, tc.name)
			} else {
				require.NoError(t, err)
				require.Equal(t, lockRecord1.LockedCoin.Amount, tc.expFromMSCoin.Amount)
				require.Equal(t, lockRecord1.LockedCoin.Denom, tc.expFromMSCoin.Denom)
				require.Equal(t, lockRecord1.LockedCoin.BondWeight, tc.expFromMSCoin.BondWeight)

				require.Equal(t, lockRecord2.LockedCoin.Amount, tc.expToMSCoin.Amount)
				require.Equal(t, lockRecord2.LockedCoin.Denom, tc.expToMSCoin.Denom)
				require.Equal(t, lockRecord2.LockedCoin.BondWeight, tc.expToMSCoin.BondWeight)
			}
		})
	}
}
