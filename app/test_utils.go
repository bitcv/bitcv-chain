package app

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/bitcv-chain/bitcv-chain/codec"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/bitcv-chain/bitcv-chain/baseapp"
	sdk "github.com/bitcv-chain/bitcv-chain/types"
	"github.com/bitcv-chain/bitcv-chain/x/auth"
)
const chainID = ""

// BigInterval is a representation of the interval [lo, hi), where
// lo and hi are both of type sdk.Int
type BigInterval struct {
	lo sdk.Int
	hi sdk.Int
}

// RandFromBigInterval chooses an interval uniformly from the provided list of
// BigIntervals, and then chooses an element from an interval uniformly at random.
func RandFromBigInterval(r *rand.Rand, intervals []BigInterval) sdk.Int {
	if len(intervals) == 0 {
		return sdk.ZeroInt()
	}

	interval := intervals[r.Intn(len(intervals))]

	lo := interval.lo
	hi := interval.hi

	diff := hi.Sub(lo)
	result := sdk.NewIntFromBigInt(new(big.Int).Rand(r, diff.BigInt()))
	result = result.Add(lo)

	return result
}

// CheckBalance checks the balance of an account.
func CheckBalance(t *testing.T, app *BacApp, addr sdk.AccAddress, exp sdk.Coins) {
	ctxCheck := app.BaseApp.NewContext(true, abci.Header{})
	res := app.accountKeeper.GetAccount(ctxCheck, addr)

	require.Equal(t, exp, res.GetCoins())
}

// CheckGenTx checks a generated signed transaction. The result of the check is
// compared against the parameter 'expPass'. A test assertion is made using the
// parameter 'expPass' against the result. A corresponding result is returned.
func CheckGenTx(
	t *testing.T, app *baseapp.BaseApp, msgs []sdk.Msg, accNums []uint64,
	seq []uint64, expPass bool, priv ...crypto.PrivKey,
) sdk.Result {
	tx := GenTx(msgs, accNums, seq, priv...)
	res := app.Check(tx)

	if expPass {
		require.Equal(t, sdk.CodeOK, res.Code, res.Log)
	} else {
		require.NotEqual(t, sdk.CodeOK, res.Code, res.Log)
	}

	return res
}

// SignCheckDeliver checks a generated signed transaction and simulates a
// block commitment with the given transaction. A test assertion is made using
// the parameter 'expPass' against the result. A corresponding result is
// returned.
func SignCheckDeliver(
	t *testing.T, cdc *codec.Codec, app *baseapp.BaseApp, header abci.Header, msgs []sdk.Msg,
	accNums, seq []uint64, expSimPass, expPass bool, priv ...crypto.PrivKey,
) sdk.Result {

	tx := GenTx(msgs, accNums, seq, priv...)
	txBytes, err := cdc.MarshalBinaryLengthPrefixed(tx)
	require.Nil(t, err)

	// Must simulate now as CheckTx doesn't run Msgs anymore
	res := app.Simulate(txBytes, tx)

	if expSimPass {
		require.Equal(t, sdk.CodeOK, res.Code, res.Log)
	} else {
		require.NotEqual(t, sdk.CodeOK, res.Code, res.Log)
	}

	// Simulate a sending a transaction and committing a block
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	res = app.Deliver(tx)

	if expPass {
		require.Equal(t, sdk.CodeOK, res.Code, res.Log)
	} else {
		require.NotEqual(t, sdk.CodeOK, res.Code, res.Log)
	}

	app.EndBlock(abci.RequestEndBlock{})
	app.Commit()

	return res
}


// GenTx generates a signed mock transaction.
func GenTx(msgs [] sdk.Msg, accnums []uint64, seq []uint64, priv ...crypto.PrivKey) auth.StdTx {
	// Make the transaction free
	fee := auth.StdFee{
		Amount: sdk.Coins{sdk.Coin{Denom:"nbac", Amount:sdk.ZeroInt()}},
		Gas:    100000,
	}

	sigs := make([]auth.StdSignature, len(priv))
	memo := ""

	for i, p := range priv {
		sig, err := p.Sign(auth.StdSignBytes(chainID, accnums[i], seq[i], fee, msgs, memo))
		if err != nil {
			panic(err)
		}

		sigs[i] = auth.StdSignature{
			PubKey:    p.PubKey(),
			Signature: sig,
		}
	}

	return auth.NewStdTx(msgs, fee, sigs, memo)
}

func GetAccountNum(ctx sdk.Context,app *BacApp,accAddr sdk.AccAddress) uint64  {
	ctxCheck := app.BaseApp.NewContext(true, abci.Header{})
	return  app.accountKeeper.GetAccount(ctxCheck,accAddr).GetAccountNumber()
}

func GetAccountSequence(ctx sdk.Context,app *BacApp,accAddr sdk.AccAddress) uint64  {
	ctxCheck := app.BaseApp.NewContext(true, abci.Header{})
	return  app.accountKeeper.GetAccount(ctxCheck,accAddr).GetSequence()
}