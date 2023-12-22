package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (coin WeightedCoin) Validate() error {
	if !coin.Amount.IsPositive() {
		return fmt.Errorf("amount zero or negative")
	}

	if !coin.Amount.IsPositive() {
		return fmt.Errorf("weight zero or negative")
	}
	return nil
}

func NewWeightedCoin(denom string, amount sdk.Int, weight sdk.Dec) WeightedCoin {
	return WeightedCoin{Denom: denom, Amount: amount, Weight: weight}
}

func (coin WeightedCoin) SafeSub(coinB WeightedCoin) (WeightedCoin, error) {
	if coin.Denom != coinB.Denom {
		return WeightedCoin{}, fmt.Errorf("denom mismatch")
	}

	coin.Amount = coin.Amount.Sub(coinB.Amount)
	if coin.Amount.IsNegative() {
		return WeightedCoin{}, fmt.Errorf("insufficient amount")
	}

	return coin, nil
}

func (coinA WeightedCoin) SafeAdd(coinB WeightedCoin) (WeightedCoin, error) {
	if coinA.Denom != coinB.Denom {
		return WeightedCoin{}, fmt.Errorf("denom mismatch")
	}

	amountA := coinA.Amount
	weightA := coinA.Weight

	amountB := coinB.Amount
	weightB := coinB.Weight

	// amountAB = amountA + amountB
	amountAB := amountA.Add(amountB)
	// weightAB = (weightA * amountA + weightB * amountB) / (amountA + amountB)
	weightAB := ((weightA.MulInt(amountA)).Add(weightB.MulInt(amountB))).QuoInt(amountAB)

	return NewWeightedCoin(coinA.Denom, amountAB, weightAB), nil
}

func (coinA WeightedCoin) Add(coinB WeightedCoin) WeightedCoin {
	coinAB, err := coinA.SafeAdd(coinB)
	if err != nil {
		panic(err)
	}

	return coinAB
}
