package types

import (
	"github.com/bitcv-chain/bitcv-chain/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateValidator{}, "bacchain/MsgCreateValidator", nil)
	cdc.RegisterConcrete(MsgEditValidator{}, "bacchain/MsgEditValidator", nil)
	cdc.RegisterConcrete(MsgDelegate{}, "bacchain/MsgDelegate", nil)
	cdc.RegisterConcrete(MsgUndelegate{}, "bacchain/MsgUndelegate", nil)
	cdc.RegisterConcrete(MsgBeginRedelegate{}, "bacchain/MsgBeginRedelegate", nil)
	cdc.RegisterConcrete(MsgExchangeBcvstakeToBcv{},"bacchain/MsgExchangeBcvstakeToBcv", nil)
	cdc.RegisterConcrete(MsgBurnBcvToEnergy{},"bacchain/MsgBurnBcvToEnergy", nil)
	cdc.RegisterConcrete(MsgBurnBcvToStake{},"bacchain/MsgBurnBcvToStake", nil)

}

// generic sealed codec to be used throughout sdk
var MsgCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	MsgCdc = cdc.Seal()
}
