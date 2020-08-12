package gov

import (
	"github.com/bitcv-chain/bitcv-chain/codec"
)

var msgCdc = codec.New()

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSubmitProposal{}, "bacchain/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "bacchain/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgVote{}, "bacchain/MsgVote", nil)

	cdc.RegisterInterface((*ProposalContent)(nil), nil)
	cdc.RegisterConcrete(TextProposal{}, "gov/TextProposal", nil)
	cdc.RegisterConcrete(SoftwareUpgradeProposal{}, "gov/SoftwareUpgradeProposal", nil)
}

func init() {
	RegisterCodec(msgCdc)
}
