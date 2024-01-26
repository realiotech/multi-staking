package types_test

import (
	"testing"

	"github.com/realio-tech/multi-staking-module/testing"
	multistakingtypes "github.com/realio-tech/multi-staking-module/x/multi-staking/types"
	"github.com/stretchr/testify/require"
)

func TestDelAddrAndValAddrFromLockID(t *testing.T) {
	val := testing.GenValAddress()
	del := testing.GenAddress()

	lockID := multistakingtypes.MultiStakingLockID(del.String(), val.String())
	lockBytes := lockID.ToBytes()
	rsDel, rsVal := multistakingtypes.DelAddrAndValAddrFromLockID(lockBytes)

	require.Equal(t, del, rsDel)
	require.Equal(t, val, rsVal)
}

func TestMultiStakingLockIterator(t *testing.T) {
	val := testing.GenValAddress()
	delA := testing.GenAddress()
	delB := testing.GenAddress()

	lockIDA := multistakingtypes.LockID{
		MultiStakerAddr: delA.String(),
		ValAddr:         val.String(),
	}

	lockIDB := multistakingtypes.LockID{
		MultiStakerAddr: delB.String(),
		ValAddr:         val.String(),
	}

	require.NotEqual(t, lockIDA, lockIDB)
	require.NotEqual(t, lockIDA.ToBytes(), lockIDB.ToBytes())
}
