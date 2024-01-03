package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/realio-tech/multi-staking-module/testutil"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func TestIntermediaryDelegatorAndMultiStakerAddress(t *testing.T) {
	multiStakerAddress := testutil.GenAddress()
	// IntermediaryDelegator from multiStakerAddress
	intermediaryAddress := types.IntermediaryDelegator(multiStakerAddress)

	// MultiStakerAddress from intermediaryAddress
	actualMultiStakerAddress := types.MultiStakerAddress(intermediaryAddress)
	require.Equal(t, actualMultiStakerAddress, multiStakerAddress)

	// IntermediaryDelegator from actualMultiStakerAddress
	actualIntermediaryAddress := types.IntermediaryDelegator(actualMultiStakerAddress)
	require.Equal(t, actualIntermediaryAddress, intermediaryAddress)
}
