package types

import (
	"github.com/realio-tech/multi-staking-module/testutil"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestIntermediaryAccount(t *testing.T) {
	mapAddr := make(map[string]bool)
	isExist := false
	maxTest := rand.Intn(50000) + 50000
	for i := 0; i < maxTest; i++ {
		delAddr := testutil.GenAddress()
		interAddr := IntermediaryAccount(delAddr)
		if mapAddr[interAddr.String()] {
			isExist = true
			break
		}
		mapAddr[interAddr.String()] = true
	}
	assert.Equal(t, false, isExist)
}

func TestDelegatorAccount(t *testing.T) {
	maxTest := rand.Intn(10000) + 10000
	for i := 0; i < maxTest; i++ {
		delAddr := testutil.GenAddress()
		interAddr := IntermediaryAccount(delAddr)
		checkDelAddr := DelegatorAccount(interAddr)
		assert.Equal(t, delAddr, checkDelAddr)
	}
}
