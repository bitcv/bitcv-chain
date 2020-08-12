package slashing

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/bitcv-chain/bitcv-chain/types"
)

func TestMsgUnjailGetSignBytes(t *testing.T) {
	addr := sdk.AccAddress("abcd")
	msg := NewMsgUnjail(sdk.ValAddress(addr))
	bytes := msg.GetSignBytes()
	require.Equal(t, string(bytes), `{"address":"bacvaloper1v93xxeqhg9nn6"}`)
}
