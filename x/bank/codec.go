package bank

import (
	"github.com/bitcv-chain/bitcv-chain/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "bacchain/MsgSend", nil)
	cdc.RegisterConcrete(MsgMultiSend{}, "bacchain/MsgMultiSend", nil)
	cdc.RegisterConcrete(MsgIssueToken{}, "bacchain/MsgIssueToken", nil)
	cdc.RegisterConcrete(MsgRedeem{}, "bacchain/MsgRedeem", nil)
	cdc.RegisterConcrete(MsgAddMargin{}, "bacchain/MsgAddMargin", nil)
	cdc.RegisterConcrete(MsgEdata{}, "bacchain/MsgEdata", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
