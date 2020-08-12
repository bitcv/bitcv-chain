package app

import (
	sdk "github.com/bitcv-chain/bitcv-chain/types"
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestGetConsumeEnergy(t *testing.T)  {
	cases := []struct {
		amt      sdk.Dec//share
		diff   int64 //height
		ret    sdk.Int
	}{
		{
			amt:  sdk.NewDecFromIntWithPrec(sdk.NewInt(11111111),int64(1)),
			diff: int64(2),
			ret: sdk.NewInt(3),
		},
		{
			amt:  sdk.NewDecFromIntWithPrec(sdk.NewInt(1),int64(2)),
			diff: int64(2),
			ret: sdk.NewInt(1),
		},
	}
	for _, tc := range cases {
		acc := sdk.GetConsumeEnergy(tc.amt,tc.diff)
		assert.EqualValues(t,acc,tc.ret)
	}
}

func TestTruncateInt(t *testing.T)  {
	a1 := sdk.MustNewDecFromStr("1.999999999999999999")
	fmt.Println(a1.Mul(a1))
	fmt.Println(a1.Mul(a1).TruncateInt())
	fmt.Println(a1.Mul(a1).TruncateDec())


}