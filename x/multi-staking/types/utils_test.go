package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/realio-tech/multi-staking-module/testutil"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func TestAccAddrAndValAddrFromStrings(t *testing.T) {
	accountAddress := testutil.GenAddress()
	valAddress := testutil.GenValAddress()

	actualAccAddr, actualValAddr, err := types.AccAddrAndValAddrFromStrings(accountAddress.String(), valAddress.String())
	require.NoError(t, err)
	require.Equal(t, accountAddress, actualAccAddr)
	require.Equal(t, valAddress, actualValAddr)
}
