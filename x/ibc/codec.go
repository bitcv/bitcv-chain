package ibc

import (
	"github.com/bitcv-chain/bitcv-chain/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgIBCTransfer{}, "bacchain/MsgIBCTransfer", nil)
	cdc.RegisterConcrete(MsgIBCReceive{}, "bacchain/MsgIBCReceive", nil)
}
