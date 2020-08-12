package slashing

import (
	"github.com/bitcv-chain/bitcv-chain/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgUnjail{}, "bacchain/MsgUnjail", nil)
}

var cdcEmpty = codec.New()
