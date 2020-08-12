package app

import (
	sdk "github.com/bitcv-chain/bitcv-chain/types"
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/bitcv-chain/bitcv-chain/x/auth"
	"github.com/bitcv-chain/bitcv-chain/x/bank"
	abci "github.com/tendermint/tendermint/abci/types"
	"fmt"
)

func TestBankKeeper(t *testing.T) {

	genTokens := sdk.TokensFromTendermintPower( 420000)
	genCoin := sdk.NewCoin(sdk.DefaultBCVDemon, genTokens)

	ubcvStakeCoin := sdk.NewCoin(sdk.DefaultBondDenom, genTokens)

	acc1 := auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{genCoin},
	}
	acc2 := auth.BaseAccount{
		Address: addr2,
		Coins:   sdk.Coins{ubcvStakeCoin},
	}

	accs := []auth.BaseAccount{acc1, acc2}
	app,ctx := CreateTestInput(accs)

	CheckBalance(t, app, addr1, sdk.Coins{genCoin})
	CheckBalance(t, app, addr2, sdk.Coins{ubcvStakeCoin})




	//新币
	issueToken := bank.NewIssueToken(addr1,"ubcv", sdk.NewInt(100000000),sdk.Coin{Denom:"ubcv",Amount:sdk.NewInt(1000000000)},
		uint8(2),"http://www,bitcv.com","http://www,bitcv.com")
	_,err,issueToken := app.bankKeeper.IssueToken(ctx,issueToken)
	require.Nil(t,err)


	//增加抵押金
	var addBcvCoin  sdk.Coin

	addBcvCoin =sdk.Coin{Denom:"ubcv",Amount:sdk.NewInt(1000000000)}
	_,_,newSaveToken := app.bankKeeper.AddMarginByInnerName(ctx,addr1,issueToken.InnerName,addBcvCoin)
	require.Equal(t,newSaveToken.ExchangeRate,sdk.MustNewDecFromStr("0.2"))
	require.Equal(t,sdk.NewInt(418000000000),app.accountKeeper.GetAccount(ctx,addr1).GetCoins().AmountOf("ubcv"))
}


func TestBankMsg(t *testing.T) {

	genTokens := sdk.TokensFromTendermintPower( 420000)
	genCoin := sdk.NewCoin(sdk.DefaultBCVDemon, genTokens)

	ubcvStakeCoin := sdk.NewCoin(sdk.DefaultBondDenom, genTokens)

	acc1 := auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{genCoin},
	}
	acc2 := auth.BaseAccount{
		Address: addr2,
		Coins:   sdk.Coins{ubcvStakeCoin},
	}

	accs := []auth.BaseAccount{acc1, acc2}
	app,ctx := CreateTestInput(accs)

	CheckBalance(t, app, addr1, sdk.Coins{genCoin})
	CheckBalance(t, app, addr2, sdk.Coins{ubcvStakeCoin})


	var header abci.Header
	header = abci.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	//新币
	issueTokenMsg := bank.NewMsgIssueToken(addr1,"ubcv", sdk.NewInt(1000),sdk.Coin{Denom:"ubcv",Amount:sdk.NewInt(1000000000)},
		uint8(2),"http://www,bitcv.com","http://www,bitcv.com")

	SignCheckDeliver(t, app.cdc, app.BaseApp, header, []sdk.Msg{issueTokenMsg}, []uint64{GetAccountNum(ctx,app,addr1)}, []uint64{GetAccountSequence(ctx,app,addr1)}, true, true, priv1)

	var newSaveToken  bank.IssueToken
	newSaveToken = app.bankKeeper.GetTokens(ctx)[0]


	//增加抵押金
	var addBcvCoin  sdk.Coin
	addBcvCoin =sdk.Coin{Denom:"ubcv",Amount:sdk.NewInt(1000000000)}
	addMarginMsg := bank.NewMsgAddMargin(addr1,newSaveToken.InnerName,addBcvCoin)
	header = abci.Header{Height: app.LastBlockHeight() + 1}

	SignCheckDeliver(t, app.cdc, app.BaseApp, header, []sdk.Msg{addMarginMsg}, []uint64{GetAccountNum(ctx,app,addr1)}, []uint64{GetAccountSequence(ctx,app,addr1)}, true, true, priv1)
	ctx = app.BaseApp.NewContext(true, abci.Header{})
	newSaveToken = app.bankKeeper.GetTokens(ctx)[0]

	require.Equal(t,newSaveToken.ExchangeRate,sdk.MustNewDecFromStr("20000"))
	require.Equal(t,sdk.NewInt(418000000000),app.accountKeeper.GetAccount(ctx,addr1).GetCoins().AmountOf("ubcv"))


	//增加抵押金,币名错误
	ctx = app.BaseApp.NewContext(true, abci.Header{})
	addMarginMsg = bank.NewMsgAddMargin(addr1,newSaveToken.InnerName,sdk.Coin{Denom:newSaveToken.InnerName,Amount:sdk.NewInt(1000000000)})
	header = abci.Header{Height: app.LastBlockHeight() + 1}
	res := SignCheckDeliver(t, app.cdc, app.BaseApp, header, []sdk.Msg{addMarginMsg}, []uint64{GetAccountNum(ctx,app,addr1)}, []uint64{GetAccountSequence(ctx,app,addr1)}, false, false, priv1)
	ctx = app.BaseApp.NewContext(true, abci.Header{})
	require.Equal(t,res.Code,bank.CodeIssueTokenCoinErr)


	//增加抵押金,资产不足
	ctx = app.BaseApp.NewContext(true, abci.Header{})
	addMarginMsg = bank.NewMsgAddMargin(addr1,newSaveToken.InnerName,sdk.Coin{Denom:"ubcv",Amount:sdk.MustNewIntFromString("41800000000000000000000000000000000000000000")})
	header = abci.Header{Height: app.LastBlockHeight() + 1}
	res = SignCheckDeliver(t, app.cdc, app.BaseApp, header, []sdk.Msg{addMarginMsg}, []uint64{GetAccountNum(ctx,app,addr1)}, []uint64{GetAccountSequence(ctx,app,addr1)}, false, false, priv1)
	ctx = app.BaseApp.NewContext(true, abci.Header{})
	require.Equal(t,res.Code,sdk.CodeInsufficientCoins)


	//赎回
	ctx = app.BaseApp.NewContext(true, abci.Header{})
	redeemMsg := bank.NewMsgRedeem(addr1,sdk.Coin{Amount:sdk.NewInt(30000),Denom:newSaveToken.InnerName})
	header = abci.Header{Height: app.LastBlockHeight() + 1}
	res = SignCheckDeliver(t, app.cdc, app.BaseApp, header, []sdk.Msg{redeemMsg}, []uint64{GetAccountNum(ctx,app,addr1)}, []uint64{GetAccountSequence(ctx,app,addr1)}, true, true, priv1)
	ctx = app.BaseApp.NewContext(true, abci.Header{})
	require.Equal(t,sdk.NewInt(418600000000),app.accountKeeper.GetAccount(ctx, addr1).GetCoins().AmountOf("ubcv"))
	require.Equal(t,sdk.NewInt(70000),app.accountKeeper.GetAccount(ctx, addr1).GetCoins().AmountOf(newSaveToken.InnerName))
	require.Equal(t,sdk.NewInt(30000),app.accountKeeper.GetAccount(ctx, newSaveToken.ExchangeAddress).GetCoins().AmountOf(newSaveToken.InnerName))


	addBcvCoin =sdk.Coin{Denom:"ubcv",Amount:sdk.NewInt(1000000000)}
	addMarginMsg = bank.NewMsgAddMargin(addr1,newSaveToken.InnerName,addBcvCoin)
	header = abci.Header{Height: app.LastBlockHeight() + 1}
	res = SignCheckDeliver(t, app.cdc, app.BaseApp, header, []sdk.Msg{addMarginMsg}, []uint64{GetAccountNum(ctx,app,addr1)}, []uint64{GetAccountSequence(ctx,app,addr1)}, true, true, priv1)
	ctx = app.BaseApp.NewContext(true, abci.Header{})
	newSaveToken = app.bankKeeper.GetTokens(ctx)[0]
	require.Equal(t,newSaveToken.ExchangeRate,sdk.MustNewDecFromStr("30000"))
	require.Equal(t,sdk.NewInt(417600000000),app.accountKeeper.GetAccount(ctx, addr1).GetCoins().AmountOf("ubcv"))

	ctx = app.BaseApp.NewContext(true, abci.Header{})
	redeemMsg = bank.NewMsgRedeem(addr1,sdk.Coin{Amount:sdk.NewInt(70000),Denom:newSaveToken.InnerName})
	header = abci.Header{Height: app.LastBlockHeight() + 1}
	res = SignCheckDeliver(t, app.cdc, app.BaseApp, header, []sdk.Msg{redeemMsg}, []uint64{GetAccountNum(ctx,app,addr1)}, []uint64{GetAccountSequence(ctx,app,addr1)}, true, true, priv1)
	ctx = app.BaseApp.NewContext(true, abci.Header{})
	fmt.Println(app.accountKeeper.GetAccount(ctx, addr1).GetCoins())
	require.Equal(t,sdk.NewInt(419700000000),app.accountKeeper.GetAccount(ctx, addr1).GetCoins().AmountOf("ubcv"))
	require.Equal(t,sdk.NewInt(0),app.accountKeeper.GetAccount(ctx, addr1).GetCoins().AmountOf(newSaveToken.InnerName))
	require.Equal(t,sdk.NewInt(100000),app.accountKeeper.GetAccount(ctx, newSaveToken.ExchangeAddress).GetCoins().AmountOf(newSaveToken.InnerName))
	
}



func TestCtx(t *testing.T) {

	genTokens := sdk.TokensFromTendermintPower( 420000)
	genCoin := sdk.NewCoin(sdk.DefaultBCVDemon, genTokens)

	ubcvStakeCoin := sdk.NewCoin(sdk.DefaultBondDenom, genTokens)

	acc1 := auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{genCoin},
	}
	acc2 := auth.BaseAccount{
		Address: addr2,
		Coins:   sdk.Coins{ubcvStakeCoin},
	}

	accs := []auth.BaseAccount{acc1, acc2}
	app,ctx := CreateTestInput(accs)

	//CheckBalance(t, app, addr1, sdk.Coins{genCoin})
	//CheckBalance(t, app, addr2, sdk.Coins{ubcvStakeCoin})


	var header abci.Header
	header = abci.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	issueTokenMsg := bank.NewMsgIssueToken(addr1,"ubcv", sdk.NewInt(1000),sdk.Coin{Denom:"ubcv",Amount:sdk.NewInt(1000000000)},
		uint8(2),"http://www,bitcv.com","http://www,bitcv.com")

	SignCheckDeliver(t, app.cdc, app.BaseApp, header, []sdk.Msg{issueTokenMsg}, []uint64{GetAccountNum(ctx,app,addr1)}, []uint64{GetAccountSequence(ctx,app,addr1)}, true, true, priv1)

	var newSaveToken  bank.IssueToken
	newSaveToken = app.bankKeeper.GetTokens(ctx)[0]


	var addBcvCoin  sdk.Coin
	addBcvCoin =sdk.Coin{Denom:"ubcv",Amount:sdk.NewInt(1000000000)}
	addMarginMsg := bank.NewMsgAddMargin(addr1,newSaveToken.InnerName,addBcvCoin)
	header = abci.Header{Height: app.LastBlockHeight() + 1}

	SignCheckDeliver(t, app.cdc, app.BaseApp, header, []sdk.Msg{addMarginMsg}, []uint64{GetAccountNum(ctx,app,addr1)}, []uint64{GetAccountSequence(ctx,app,addr1)}, true, true, priv1)
	//ctx = app.BaseApp.NewContext(true, abci.Header{})
	newSaveToken = app.bankKeeper.GetTokens(ctx)[0]


	fmt.Println(app.accountKeeper.GetAccount(ctx,addr1).GetCoins())
	require.Equal(t,newSaveToken.ExchangeRate,sdk.MustNewDecFromStr("20000"))
	require.Equal(t,sdk.NewInt(418000000000),app.accountKeeper.GetAccount(ctx,addr1).GetCoins().AmountOf("ubcv"))


	ctx = app.BaseApp.NewContext(true, abci.Header{})
	addMarginMsg = bank.NewMsgAddMargin(addr1,newSaveToken.InnerName,sdk.Coin{Denom:newSaveToken.InnerName,Amount:sdk.NewInt(1000000000)})
	header = abci.Header{Height: app.LastBlockHeight() + 1}
	res := SignCheckDeliver(t, app.cdc, app.BaseApp, header, []sdk.Msg{addMarginMsg}, []uint64{GetAccountNum(ctx,app,addr1)}, []uint64{GetAccountSequence(ctx,app,addr1)}, false, false, priv1)
	ctx = app.BaseApp.NewContext(true, abci.Header{})
	require.Equal(t,res.Code,bank.CodeIssueTokenCoinErr)


	ctx = app.BaseApp.NewContext(true, abci.Header{})
	addMarginMsg = bank.NewMsgAddMargin(addr1,newSaveToken.InnerName,sdk.Coin{Denom:"ubcv",Amount:sdk.MustNewIntFromString("41800000000000000000000000000000000000000000")})
	header = abci.Header{Height: app.LastBlockHeight() + 1}
	res = SignCheckDeliver(t, app.cdc, app.BaseApp, header, []sdk.Msg{addMarginMsg}, []uint64{GetAccountNum(ctx,app,addr1)}, []uint64{GetAccountSequence(ctx,app,addr1)}, false, false, priv1)
	ctx = app.BaseApp.NewContext(true, abci.Header{})
	require.Equal(t,res.Code,sdk.CodeInsufficientCoins)


	ctx = app.BaseApp.NewContext(true, abci.Header{})
	redeemMsg := bank.NewMsgRedeem(addr1,sdk.Coin{Amount:sdk.NewInt(30000),Denom:newSaveToken.InnerName})
	header = abci.Header{Height: app.LastBlockHeight() + 1}
	SignCheckDeliver(t, app.cdc, app.BaseApp, header, []sdk.Msg{redeemMsg}, []uint64{GetAccountNum(ctx,app,addr1)}, []uint64{GetAccountSequence(ctx,app,addr1)}, true, true, priv1)


	header = abci.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	app.EndBlock(abci.RequestEndBlock{})
	app.Commit()

	fmt.Println(app.accountKeeper.GetAccount(ctx, addr1).GetCoins())
	fmt.Println(ctx.KVStore(app.keyAccount))

	ctx1 := app.BaseApp.NewContext(true, abci.Header{})
	fmt.Println(app.accountKeeper.GetAccount(ctx1, addr1).GetCoins())
	fmt.Println(ctx1.KVStore(app.keyAccount))

}