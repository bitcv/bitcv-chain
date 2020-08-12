package types
import (
	 "github.com/tendermint/tendermint/crypto"
)
var (
	DefaultFeeDenom = "nbac"
	DefaultMinimunFee = "1"
	DefaultEnergyDenom = "energy"
	PRECISON_G  =   NewInt(1000000000) //10^9
	PRECISON_U = NewInt(1000000) // 1bcv = 1000000ubcv 10^6

	ISSUE_TOKEN_TEST_SUPPLY_NUM = NewInt(10000)//2
	ISSUE_TOKEN_TEST_MARGIN_NUM = NewInt(100000000)
	ISSUE_TOKEN_TEST_PRECISION = 2
	ISSUE_TOKEN_MIN_MARGIN_NUM = 1000 //发币抵押最小BCV


	ISSUE_TOKEN_PRO_MIN_MARGIN_NUM  = NewInt(100000000000)


	AccAddrBurnFromEdataSave = AccAddress(crypto.AddressHash([]byte("burnaddrfromedatasave")))


)


const (
	DEFAULT_FEE_COIN = "nbac"
	DEFAULT_MIN_FEE = "1"
	AccuracyG =    1000000000

	CHAIN_COIN_NAME_BCVSTAKE = "ubcvstake"
	CHAIN_COIN_NAME_ENERGY   = "energy"
	CHAIN_COIN_NAME_BAC  = "nbac"
	CHAIN_COIN_NAME_BCV  = "ubcv"

	CHAIN_PARAM_BCV_AMOUNT = "1200000000000000"
	CHAIN_PARAM_ENERGY_AMOUNT = "100000000000000000000000000000000"
	CHAIN_PARAM_BCVSTAKE_AMOUNT = "1560000000000000"
)


func GetChainParamBcvAmount() Int{
	amount,_ := NewIntFromString(CHAIN_PARAM_BCV_AMOUNT)
	return amount;
}

func GetChainParamEnergyAmount() Int{
	amount,_ := NewIntFromString(CHAIN_PARAM_ENERGY_AMOUNT)
	return amount;
}


func GetChainParamBcvStakeAmount() Int{
	amount,_ := NewIntFromString(CHAIN_PARAM_BCVSTAKE_AMOUNT)
	return amount;
}