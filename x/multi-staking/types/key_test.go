package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/realio-tech/multi-staking-module/testutil"
	mulStakingtypes "github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func TestDelAddrAndValAddrFromLockID(t *testing.T) {
	val := testutil.GenValAddress()
	del := testutil.GenAddress()

	lockID := mulStakingtypes.MultiStakingLockID(del.String(), val.String())

	toByte := lockID.ToByte()

	rsDel, rsVal := mulStakingtypes.DelAddrAndValAddrFromLockID(toByte)
	require.NotEqual(t, del, rsDel)
	require.NotEqual(t, val, rsVal)
}

func TestMultiStakingLockIterator(t *testing.T) {
	val := testutil.GenValAddress()
	del1 := testutil.GenAddress()
	del2 := testutil.GenAddress()

	lockID1 := mulStakingtypes.LockID{
		MultiStakerAddr: del1.String(),
		ValAddr:         val.String(),
	}

	lockID2 := mulStakingtypes.LockID{
		MultiStakerAddr: del2.String(),
		ValAddr:         val.String(),
	}

	require.NotEqual(t, lockID1, lockID2)
	require.Equal(t, lockID1.ToByte(), lockID2.ToByte())
}
