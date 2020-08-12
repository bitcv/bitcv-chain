package mint

import (
	"fmt"

	sdk "github.com/bitcv-chain/bitcv-chain/types"
	bacv1 "github.com/bitcv-chain/bitcv-chain/bacchain/v1_00"
)

var (
	secondRunHeight = int64(1363344) //gen 5bac
	ThirdRunHeight = int64(2620800)
	FourthRunHeight = int64(2620806)

	FirstReduceHeight = bacv1.StartParamInitHeight + bacv1.StartParamBeginGenBac
	SecondReduceHeight = FirstReduceHeight + secondRunHeight
	ThirdReduceHeight = SecondReduceHeight + ThirdRunHeight
	FourthReduceHeight = ThirdReduceHeight + FourthRunHeight
)

// Minter represents the minting state.
type Minter struct {
	Inflation        sdk.Dec `json:"inflation"`         // current annual inflation rate
	AnnualProvisions sdk.Dec `json:"annual_provisions"` // current annual expected provisions
}

// NewMinter returns a new Minter object with the given inflation and annual
// provisions values.
func NewMinter(inflation, annualProvisions sdk.Dec) Minter {
	return Minter{
		Inflation:        inflation,
		AnnualProvisions: annualProvisions,
	}
}

// InitialMinter returns an initial Minter object with a given inflation value.
func InitialMinter(inflation sdk.Dec) Minter {
	return NewMinter(
		inflation,
		sdk.NewDec(0),
	)
}

// DefaultInitialMinter returns a default initial Minter object for a new chain
// which uses an inflation rate of 13%.
func DefaultInitialMinter() Minter {
	return InitialMinter(
		sdk.NewDecWithPrec(13, 2),
	)
}

func validateMinter(minter Minter) error {
	if minter.Inflation.LT(sdk.ZeroDec()) {
		return fmt.Errorf("mint parameter Inflation should be positive, is %s",
			minter.Inflation.String())
	}
	return nil
}

// NextInflationRate returns the new inflation rate for the next hour.
func (m Minter) NextInflationRate(params Params, bondedRatio sdk.Dec) sdk.Dec {
	// The target annual inflation rate is recalculated for each previsions cycle. The
	// inflation is also subject to a rate change (positive or negative) depending on
	// the distance from the desired ratio (67%). The maximum rate change possible is
	// defined to be 13% per year, however the annual inflation is capped as between
	// 7% and 20%.

	// (1 - bondedRatio/GoalBonded) * InflationRateChange
	inflationRateChangePerYear := sdk.OneDec().
		Sub(bondedRatio.Quo(params.GoalBonded)).
		Mul(params.InflationRateChange)
	inflationRateChange := inflationRateChangePerYear.Quo(sdk.NewDec(int64(params.BlocksPerYear)))

	// adjust the new annual inflation for this next cycle
	inflation := m.Inflation.Add(inflationRateChange) // note inflationRateChange may be negative
	if inflation.GT(params.InflationMax) {
		inflation = params.InflationMax
	}
	if inflation.LT(params.InflationMin) {
		inflation = params.InflationMin
	}

	return inflation
}

// NextAnnualProvisions returns the annual provisions based on current total
// supply and inflation rate.
func (m Minter) NextAnnualProvisions(_ Params, totalSupply sdk.Int) sdk.Dec {
	return m.Inflation.MulInt(totalSupply)
}

// BlockProvision returns the provisions for a block based on the annual
// provisions rate.
func (m Minter) BlockProvision(params Params) sdk.Coin {
	provisionAmt := m.AnnualProvisions.QuoInt(sdk.NewInt(int64(params.BlocksPerYear)))
	return sdk.NewCoin(params.MintDenom, provisionAmt.TruncateInt())
}

//nbac  pow(10,6) 1000000
func (m Minter) NextReduceAnnualProvisions(_ Params, height int64) sdk.Int {
	//不产生BAC
	if height <= bacv1.StartParamBeginGenBac{
		return sdk.NewInt(0)
	}

	height += bacv1.StartParamInitHeight

	// According the block height, get the reward
	if height <= FirstReduceHeight {
		return sdk.PRECISON_G.Mul(sdk.NewInt(int64(10)))
	} else if height <= SecondReduceHeight {
		return sdk.PRECISON_G.Mul(sdk.NewInt(int64(5)))
	} else if height <= ThirdReduceHeight {
		return sdk.PRECISON_G.Mul(sdk.NewInt(int64(2)))
	} else if height <= FourthReduceHeight {
		return sdk.PRECISON_G.Mul(sdk.NewInt(int64(1)))
	} else {  // After two years, stop generate reward from block
		return sdk.ZeroInt()
	}
}


