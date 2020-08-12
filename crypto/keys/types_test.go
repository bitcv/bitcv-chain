package keys

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/bitcv-chain/bitcv-chain/crypto/keys/hd"
	"github.com/bitcv-chain/bitcv-chain/types"
)

//mnemonic jeans antenna lucky way advice inherit sunset wild shock motion primary transfer exit excite design hope stage critic flush sister spell broom coach zoo
// pri 0ec4e0365b42f135ef63b48c02695ff930003e6568c43f27cce91745d4226b1c
// accAddr bac10tyhju9pfpfkt7hrd2zqr0vjn4k5sfrrhenf7v
func Test_writeReadLedgerInfo(t *testing.T) {
	var tmpKey secp256k1.PubKeySecp256k1
	bz, _ := hex.DecodeString("0ec4e0365b42f135ef63b48c02695ff930003e6568c43f27cce91745d4226b1c")
	copy(tmpKey[:], bz)

	lInfo := ledgerInfo{
		"some_name",
		tmpKey,
		*hd.NewFundraiserParams(5, 1)}
	assert.Equal(t, TypeLedger, lInfo.GetType())

	path, err := lInfo.GetPath()

	assert.NoError(t, err)
	assert.Equal(t, "44'/572'/5'/0/1", path.String())

	assert.Equal(t,
		"bacpub1addwnpeppmzwqdjmgtcntmmrkjxqy62llycqq0n9drzr7f7vayt5t4pzdvwqqglm5m8",
		types.MustBech32ifyAccPub(lInfo.GetPubKey()))
	// Serialize and restore
	serialized := writeInfo(lInfo)
	restoredInfo, err := readInfo(serialized)
	assert.NoError(t, err)
	assert.NotNil(t, restoredInfo)

	// Check both keys match
	assert.Equal(t, lInfo.GetName(), restoredInfo.GetName())
	assert.Equal(t, lInfo.GetType(), restoredInfo.GetType())
	assert.Equal(t, lInfo.GetPubKey(), restoredInfo.GetPubKey())

	restoredPath, err := restoredInfo.GetPath()
	assert.NoError(t, err)

	assert.Equal(t, path, restoredPath)
}
