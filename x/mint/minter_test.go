package mint

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/bitcv-chain/bitcv-chain/types"
	"github.com/stretchr/testify/assert"
	"github.com/bitcv-chain/bitcv-chain/bacchain/v1_00"
	"fmt"
)

func TestNextInflation(t *testing.T) {
	minter := DefaultInitialMinter()
	params := DefaultParams()
	blocksPerYr := sdk.NewDec(int64(params.BlocksPerYear))

	// Governing Mechanism:
	//    inflationRateChangePerYear = (1- BondedRatio/ GoalBonded) * MaxInflationRateChange

	tests := []struct {
		bondedRatio, setInflation, expChange sdk.Dec
	}{
		// with 0% bonded atom supply the inflation should increase by InflationRateChange
		{sdk.ZeroDec(), sdk.NewDecWithPrec(7, 2), params.InflationRateChange.Quo(blocksPerYr)},

		// 100% bonded, starting at 20% inflation and being reduced
		// (1 - (1/0.67))*(0.13/8667)
		{sdk.OneDec(), sdk.NewDecWithPrec(20, 2),
			sdk.OneDec().Sub(sdk.OneDec().Quo(params.GoalBonded)).Mul(params.InflationRateChange).Quo(blocksPerYr)},

		// 50% bonded, starting at 10% inflation and being increased
		{sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(10, 2),
			sdk.OneDec().Sub(sdk.NewDecWithPrec(5, 1).Quo(params.GoalBonded)).Mul(params.InflationRateChange).Quo(blocksPerYr)},

		// test 7% minimum stop (testing with 100% bonded)
		{sdk.OneDec(), sdk.NewDecWithPrec(7, 2), sdk.ZeroDec()},
		{sdk.OneDec(), sdk.NewDecWithPrec(700000001, 10), sdk.NewDecWithPrec(-1, 10)},

		// test 20% maximum stop (testing with 0% bonded)
		{sdk.ZeroDec(), sdk.NewDecWithPrec(20, 2), sdk.ZeroDec()},
		{sdk.ZeroDec(), sdk.NewDecWithPrec(1999999999, 10), sdk.NewDecWithPrec(1, 10)},

		// perfect balance shouldn't change inflation
		{sdk.NewDecWithPrec(67, 2), sdk.NewDecWithPrec(15, 2), sdk.ZeroDec()},
	}
	for i, tc := range tests {
		minter.Inflation = tc.setInflation

		inflation := minter.NextInflationRate(params, tc.bondedRatio)
		diffInflation := inflation.Sub(tc.setInflation)

		require.True(t, diffInflation.Equal(tc.expChange),
			"Test Index: %v\nDiff:  %v\nExpected: %v\n", i, diffInflation, tc.expChange)
	}
}

//go test  -v -run="TestBlockProvision"
func TestBlockProvision(t *testing.T) {
	minter := InitialMinter(sdk.NewDec(1))
	sum  := sdk.NewInt(0)
	for i:=1;i<=int(FourthReduceHeight - v1_00.StartParamInitHeight );i++{
		a:= minter.NextReduceAnnualProvisions(Params{},int64(i))
		sum = sum.Add(a)
	}
	fmt.Println(sum)
	assert.Equal(t,sdk.NewDecFromInt(sum).Add(sdk.MustNewDecFromStr(v1_00.StartParamBacAlreadyProduce)),sdk.MustNewDecFromStr("43695126000000000"))
}

// Benchmarking :)
// previously using sdk.Int operations:
// BenchmarkBlockProvision-4 5000000 220 ns/op
//
// using sdk.Dec operations: (current implementation)
// BenchmarkBlockProvision-4 3000000 429 ns/op
func BenchmarkBlockProvision(b *testing.B) {
	minter := InitialMinter(sdk.NewDecWithPrec(1, 1))
	params := DefaultParams()

	s1 := rand.NewSource(100)
	r1 := rand.New(s1)
	minter.AnnualProvisions = sdk.NewDec(r1.Int63n(1000000))

	// run the BlockProvision function b.N times
	for n := 0; n < b.N; n++ {
		minter.BlockProvision(params)
	}
}

// Next inflation benchmarking
// BenchmarkNextInflation-4 1000000 1828 ns/op
func BenchmarkNextInflation(b *testing.B) {
	minter := InitialMinter(sdk.NewDecWithPrec(1, 1))
	params := DefaultParams()
	bondedRatio := sdk.NewDecWithPrec(1, 1)

	// run the NextInflationRate function b.N times
	for n := 0; n < b.N; n++ {
		minter.NextInflationRate(params, bondedRatio)
	}

}

// Next annual provisions benchmarking
// BenchmarkNextAnnualProvisions-4 5000000 251 ns/op
func BenchmarkNextAnnualProvisions(b *testing.B) {
	minter := InitialMinter(sdk.NewDecWithPrec(1, 1))
	params := DefaultParams()
	totalSupply := sdk.NewInt(100000000000000)

	// run the NextAnnualProvisions function b.N times
	for n := 0; n < b.N; n++ {
		minter.NextAnnualProvisions(params, totalSupply)
	}

}
