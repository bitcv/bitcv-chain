package bank

import (
	"testing"

	sdk "github.com/bitcv-chain/bitcv-chain/types"
	"github.com/bitcv-chain/bitcv-chain/x/auth"
	"github.com/bitcv-chain/bitcv-chain/x/mock"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type (
	expectedBalance struct {
		addr  sdk.AccAddress
		coins sdk.Coins
	}

	appTestCase struct {
		expSimPass       bool
		expPass          bool
		msgs             []sdk.Msg
		accNums          []uint64
		accSeqs          []uint64
		privKeys         []crypto.PrivKey
		expectedBalances []expectedBalance
	}
)

var (
	priv1 = secp256k1.GenPrivKey()
	addr1 = sdk.AccAddress(priv1.PubKey().Address())
	priv2 = secp256k1.GenPrivKey()
	addr2 = sdk.AccAddress(priv2.PubKey().Address())
	addr3 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	priv4 = secp256k1.GenPrivKey()
	addr4 = sdk.AccAddress(priv4.PubKey().Address())

	coins     = sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}
	halfCoins = sdk.Coins{sdk.NewInt64Coin("foocoin", 5)}
	manyCoins = sdk.Coins{sdk.NewInt64Coin("foocoin", 1), sdk.NewInt64Coin("barcoin", 1)}
	freeFee   = auth.NewStdFee(100000, sdk.Coins{sdk.NewInt64Coin("foocoin", 0)})

	sendMsg1 = NewMsgSend(addr1, addr2, coins)

	multiSendMsg1 = MsgMultiSend{
		Inputs:  []Input{NewInput(addr1, coins)},
		Outputs: []Output{NewOutput(addr2, coins)},
	}
	multiSendMsg2 = MsgMultiSend{
		Inputs: []Input{NewInput(addr1, coins)},
		Outputs: []Output{
			NewOutput(addr2, halfCoins),
			NewOutput(addr3, halfCoins),
		},
	}
	multiSendMsg3 = MsgMultiSend{
		Inputs: []Input{
			NewInput(addr1, coins),
			NewInput(addr4, coins),
		},
		Outputs: []Output{
			NewOutput(addr2, coins),
			NewOutput(addr3, coins),
		},
	}
	multiSendMsg4 = MsgMultiSend{
		Inputs: []Input{
			NewInput(addr2, coins),
		},
		Outputs: []Output{
			NewOutput(addr1, coins),
		},
	}
	multiSendMsg5 = MsgMultiSend{
		Inputs: []Input{
			NewInput(addr1, manyCoins),
		},
		Outputs: []Output{
			NewOutput(addr2, manyCoins),
		},
	}
)

// initialize the mock application for this module
func getMockApp(t *testing.T) *mock.App {
	mapp, err := getBenchmarkMockApp()
	require.NoError(t, err)
	return mapp
}

// overwrite the mock init chainer
func getInitChainer(mapp *mock.App, keeper BaseKeeper) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)
		bankGenesis := DefaultGenesisState()
		InitGenesis(ctx, keeper, bankGenesis)

		return abci.ResponseInitChain{}
	}
}

func TestSendNotEnoughBalance(t *testing.T) {
	mapp := getMockApp(t)
	acc := &auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{sdk.NewInt64Coin("foocoin", 67)},
	}

	mock.SetGenesis(mapp, []auth.Account{acc})

	ctxCheck := mapp.BaseApp.NewContext(true, abci.Header{})

	res1 := mapp.AccountKeeper.GetAccount(ctxCheck, addr1)
	require.NotNil(t, res1)
	require.Equal(t, acc, res1.(*auth.BaseAccount))

	origAccNum := res1.GetAccountNumber()
	origSeq := res1.GetSequence()

	sendMsg := NewMsgSend(addr1, addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 100)})
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mock.SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, header, []sdk.Msg{sendMsg}, []uint64{origAccNum}, []uint64{origSeq}, false, false, priv1)

	mock.CheckBalance(t, mapp, addr1, sdk.Coins{sdk.NewInt64Coin("foocoin", 67)})

	res2 := mapp.AccountKeeper.GetAccount(mapp.NewContext(true, abci.Header{}), addr1)
	require.NotNil(t, res2)

	require.True(t, res2.GetAccountNumber() == origAccNum)
	require.True(t, res2.GetSequence() == origSeq+1)
}

