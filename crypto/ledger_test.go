package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/encoding/amino"

	"github.com/bitcv-chain/bitcv-chain/crypto/keys/hd"
	"github.com/bitcv-chain/bitcv-chain/tests"
	sdk "github.com/bitcv-chain/bitcv-chain/types"
)

func TestLedgerErrorHandling(t *testing.T) {
	// first, try to generate a key, must return an error
	// (no panic)
	path := *hd.NewParams(44, 555, 0, false, 0)
	_, err := NewPrivKeyLedgerSecp256k1Unsafe(path)
	require.Error(t, err)
}

func TestPublicKeyUnsafe(t *testing.T) {
	path := *hd.NewFundraiserParams(0, 0)
	priv, err := NewPrivKeyLedgerSecp256k1Unsafe(path)
	require.Nil(t, err, "%s", err)
	require.NotNil(t, priv)

	require.Equal(t, "eb5ae9872102655896ea66c5ad0d63216365ee5c116aa89e710740db5c8751f3dd7092556fac",
		fmt.Sprintf("%x", priv.PubKey().Bytes()),
		"Is your device using test mnemonic: %s ?", tests.TestMnemonic)

	pubKeyAddr, err := sdk.Bech32ifyAccPub(priv.PubKey())
	require.NoError(t, err)
	require.Equal(t, "bacpub1addwnpepqfj439h2vmz66rtry93ktmjuz94238n3qaqdkhy828ea6uyj24h6c68n6qz",
		pubKeyAddr, "Is your device using test mnemonic: %s ?", tests.TestMnemonic)

	addr := sdk.AccAddress(priv.PubKey().Address()).String()
	require.Equal(t, "bac1sc86sla6y9gld4al3dhggca8l5q58hymz3vtye",
		addr, "Is your device using test mnemonic: %s ?", tests.TestMnemonic)
}

func TestPublicKeyUnsafeHDPath(t *testing.T) {
	expectedAnswers := []string{
		"bacpub1addwnpepqfj439h2vmz66rtry93ktmjuz94238n3qaqdkhy828ea6uyj24h6c68n6qz",
		"bacpub1addwnpepqfhx4xv439su8x6tkq93je56qaa52y8659swjx5x9ts8h4zdfc8wuhgl75z",
		"bacpub1addwnpepqvf5ha0w07ejkl5n6u5cz7kpxd2splzxynkx6dmayp3wugs6h6plz7vejaa",
		"bacpub1addwnpepq2lgk0f82mdhv9swawz0fzfvhsnak56jza7f2w5ru4ksjlyxm6uak5ny55w",
	}

	const numIters = 4

	privKeys := make([]tmcrypto.PrivKey, numIters)

	// Check with device
	for i := uint32(0); i < 4; i++ {
		path := *hd.NewFundraiserParams(0, i)
		fmt.Printf("Checking keys at %v\n", path)

		priv, err := NewPrivKeyLedgerSecp256k1Unsafe(path)
		require.Nil(t, err, "%s", err)
		require.NotNil(t, priv)

		// Check other methods
		require.NoError(t, priv.(PrivKeyLedgerSecp256k1).ValidateKey())
		tmp := priv.(PrivKeyLedgerSecp256k1)
		(&tmp).AssertIsPrivKeyInner()

		pubKeyAddr, err := sdk.Bech32ifyAccPub(priv.PubKey())
		require.NoError(t, err)
		require.Equal(t,
			expectedAnswers[i], pubKeyAddr,
			"Is your device using test mnemonic: %s ?", tests.TestMnemonic)

		// Store and restore
		serializedPk := priv.Bytes()
		require.NotNil(t, serializedPk)
		require.True(t, len(serializedPk) >= 50)

		privKeys[i] = priv
	}

	// Now check equality
	for i := 0; i < 4; i++ {
		for j := 0; j <4; j++ {
			require.Equal(t, i == j, privKeys[i].Equals(privKeys[j]))
			require.Equal(t, i == j, privKeys[j].Equals(privKeys[i]))
		}
	}
}

func TestPublicKeySafe(t *testing.T) {
	path := *hd.NewFundraiserParams(0, 0)
	priv, addr, err := NewPrivKeyLedgerSecp256k1(path, "bac")

	require.Nil(t, err, "%s", err)
	require.NotNil(t, priv)

	require.Equal(t, "eb5ae9872102655896ea66c5ad0d63216365ee5c116aa89e710740db5c8751f3dd7092556fac",
		fmt.Sprintf("%x", priv.PubKey().Bytes()),
		"Is your device using test mnemonic: %s ?", tests.TestMnemonic)

	pubKeyAddr, err := sdk.Bech32ifyAccPub(priv.PubKey())
	require.NoError(t, err)
	require.Equal(t, "bacpub1addwnpepqfj439h2vmz66rtry93ktmjuz94238n3qaqdkhy828ea6uyj24h6c68n6qz",
		pubKeyAddr, "Is your device using test mnemonic: %s ?", tests.TestMnemonic)

	require.Equal(t, "bac1sc86sla6y9gld4al3dhggca8l5q58hymz3vtye",
		addr, "Is your device using test mnemonic: %s ?", tests.TestMnemonic)

	addr2 := sdk.AccAddress(priv.PubKey().Address()).String()
	require.Equal(t, addr, addr2)
}

func TestPublicKeyHDPath(t *testing.T) {
	expectedPubKeys := []string{
		"bacpub1addwnpepqfj439h2vmz66rtry93ktmjuz94238n3qaqdkhy828ea6uyj24h6c68n6qz",
		"bacpub1addwnpepqfhx4xv439su8x6tkq93je56qaa52y8659swjx5x9ts8h4zdfc8wuhgl75z",
		"bacpub1addwnpepqvf5ha0w07ejkl5n6u5cz7kpxd2splzxynkx6dmayp3wugs6h6plz7vejaa",
		"bacpub1addwnpepq2lgk0f82mdhv9swawz0fzfvhsnak56jza7f2w5ru4ksjlyxm6uak5ny55w",
		"bacpub1addwnpepqfvv8pdsmpxpvm86mjgq4ymq668x5thxqz9mpk4gpy0d2jgc75u0snaefru",
	}

	expectedAddrs := []string{
		"bac1sc86sla6y9gld4al3dhggca8l5q58hymz3vtye",
		"bac1dj77hrdcrkgfkpvhpde3snsxac5shmtpl4tvm8",
		"bac17f8kvvx0f7rjtr4p7c2vf2trtn9r8jmnepxtd8",
		"bac18sgxq8j2h0me03pvmg5259wl0alfqpsac3cekg",
		"bac1t4rca7vulwcvqj9crc79fmkhra2kxwr0nwy7lr",
	}

	const numIters = 5

	privKeys := make([]tmcrypto.PrivKey, numIters)

	// Check with device
	for i := uint32(0); i < numIters; i++ {
		path := *hd.NewFundraiserParams(0, i)
		fmt.Printf("Checking keys at %v\n", path)

		priv, addr, err := NewPrivKeyLedgerSecp256k1(path, "bac")
		require.Nil(t, err, "%s", err)
		require.NotNil(t, addr)
		require.NotNil(t, priv)

		addr2 := sdk.AccAddress(priv.PubKey().Address()).String()
		require.Equal(t, addr2, addr)
		require.Equal(t,
			expectedAddrs[i], addr,
			"Is your device using test mnemonic: %s ?", tests.TestMnemonic)

		// Check other methods
		require.NoError(t, priv.(PrivKeyLedgerSecp256k1).ValidateKey())
		tmp := priv.(PrivKeyLedgerSecp256k1)
		(&tmp).AssertIsPrivKeyInner()

		pubKeyAddr, err := sdk.Bech32ifyAccPub(priv.PubKey())
		require.NoError(t, err)
		require.Equal(t,
			expectedPubKeys[i], pubKeyAddr,
			"Is your device using test mnemonic: %s ?", tests.TestMnemonic)

		// Store and restore
		serializedPk := priv.Bytes()
		require.NotNil(t, serializedPk)
		require.True(t, len(serializedPk) >= 50)

		privKeys[i] = priv
	}

	// Now check equality
	for i := 0; i < numIters; i++ {
		for j := 0; j < numIters; j++ {
			require.Equal(t, i == j, privKeys[i].Equals(privKeys[j]))
			require.Equal(t, i == j, privKeys[j].Equals(privKeys[i]))
		}
	}
}

func getFakeTx(accountNumber uint32) []byte {
	tmp := fmt.Sprintf(
		`{"account_number":"%d","chain_id":"1234","fee":{"amount":[{"amount":"150","denom":"atom"}],"gas":"5000"},"memo":"memo","msgs":[[""]],"sequence":"6"}`,
		accountNumber)

	return []byte(tmp)
}

func TestSignaturesHD(t *testing.T) {
	for account := uint32(0); account < 100; account += 30 {
		msg := getFakeTx(account)

		path := *hd.NewFundraiserParams(account, account/5)
		fmt.Printf("Checking signature at %v    ---   PLEASE REVIEW AND ACCEPT IN THE DEVICE\n", path)

		priv, err := NewPrivKeyLedgerSecp256k1Unsafe(path)
		require.Nil(t, err, "%s", err)

		pub := priv.PubKey()
		sig, err := priv.Sign(msg)
		require.Nil(t, err)

		valid := pub.VerifyBytes(msg, sig)
		require.True(t, valid, "Is your device using test mnemonic: %s ?", tests.TestMnemonic)
	}
}

func TestRealLedgerSecp256k1(t *testing.T) {
	msg := getFakeTx(50)
	path := *hd.NewFundraiserParams(0, 0)
	priv, err := NewPrivKeyLedgerSecp256k1Unsafe(path)
	require.Nil(t, err, "%s", err)

	pub := priv.PubKey()
	sig, err := priv.Sign(msg)
	require.Nil(t, err)

	valid := pub.VerifyBytes(msg, sig)
	require.True(t, valid)

	// now, let's serialize the public key and make sure it still works
	bs := priv.PubKey().Bytes()
	pub2, err := cryptoAmino.PubKeyFromBytes(bs)
	require.Nil(t, err, "%+v", err)

	// make sure we get the same pubkey when we load from disk
	require.Equal(t, pub, pub2)

	// signing with the loaded key should match the original pubkey
	sig, err = priv.Sign(msg)
	require.Nil(t, err)
	valid = pub.VerifyBytes(msg, sig)
	require.True(t, valid)

	// make sure pubkeys serialize properly as well
	bs = pub.Bytes()
	bpub, err := cryptoAmino.PubKeyFromBytes(bs)
	require.NoError(t, err)
	require.Equal(t, pub, bpub)
}
