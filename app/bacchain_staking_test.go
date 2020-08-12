package app

import (


	"github.com/bitcv-chain/bitcv-chain/x/staking"

	abci "github.com/tendermint/tendermint/abci/types"
	sdk "github.com/bitcv-chain/bitcv-chain/types"

	"testing"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
	"fmt"
	"github.com/bitcv-chain/bitcv-chain/x/auth"
	"github.com/bitcv-chain/bitcv-chain/x/mock"
)

var (
	commissionMsg = staking.NewCommissionMsg(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
)

func TestUnbond(t *testing.T)  {
	var sh sdk.Handler
	app,ctx:= CreateTestInput([]auth.BaseAccount{})
	bankKeeper := app.bankKeeper
	header := abci.Header{Height: 1}
	app.BeginBlocker(ctx,abci.RequestBeginBlock{Header: header})
	power := int64(100000)
	amt := sdk.TokensFromTendermintPower(power)

	addr0, val0 := addrs[0], pks[0]
	accAddr0 := sdk.AccAddress(addr0)
	sh = staking.NewHandler(app.stakingKeeper)
	got := sh(ctx, staking.NewTestMsgCreateValidator(addr0, val0, amt))
	require.True(t,got.IsOK())

	addr1, val1 := addrs[1], pks[1]
	sh = staking.NewHandler(app.stakingKeeper)
	got = sh(ctx, staking.NewTestMsgCreateValidator(addr1, val1, amt))
	require.True(t,got.IsOK())

	//need more energy
	got = sh(ctx,staking.NewMsgUndelegate(sdk.AccAddress(addr0),addr0,
		sdk.Coin{Denom:defaultBondDenom,Amount:sdk.NewInt(100)}))
	assert.EqualValues(t,got.Code,sdk.CodeNotEnoughEnergy)

	//exchange energy
	got = sh(ctx,staking.NewMsgBurnBcvToEnergy(accAddr0,sdk.Coin{Denom:"ubcv",Amount:sdk.NewInt(10000)}))
	energyCoins := sdk.Coins{sdk.Coin{Denom:"ubcv",Amount:sdk.NewInt(10000)}}
	require.True(t,bankKeeper.HasCoins(ctx,sdk.AccAddress(addr0),energyCoins))

	//unbond succ
	unbondAmt := sdk.NewInt(100)
	got = sh(ctx,staking.NewMsgUndelegate(sdk.AccAddress(addr0),addr0,
		sdk.Coin{Denom:defaultBondDenom,Amount:unbondAmt}))

	require.True(t,got.IsOK())
	validator0,_:= app.stakingKeeper.GetValidator(ctx,sdk.ValAddress(addr0.Bytes()))
	require.Equal(t,validator0.Tokens,amt.Sub(unbondAmt))

	a1 := bankKeeper.GetCoins(ctx,accAddr0).AmountOf(sdk.CHAIN_COIN_NAME_BCV)
	b1 := bankKeeper.GetCoins(ctx,sdk.AccAddrBcvstakePool).AmountOf(sdk.CHAIN_COIN_NAME_BCV)
	staking.EndBlocker(ctx,app.stakingKeeper)
	ctx = ctx.WithBlockTime(ctx.BlockHeader().Time.Add(app.stakingKeeper.UnbondingTime(ctx)))
	staking.EndBlocker(ctx,app.stakingKeeper)
	a2 := bankKeeper.GetCoins(ctx,accAddr0).AmountOf(sdk.CHAIN_COIN_NAME_BCV)
	b2 := bankKeeper.GetCoins(ctx,sdk.AccAddrBcvstakePool).AmountOf(sdk.CHAIN_COIN_NAME_BCV)
	require.Equal(t,unbondAmt,a2.Sub(a1))
	require.Equal(t,unbondAmt,b1.Sub(b2))
}


func TestHasMaxUnbond(t *testing.T)  {
	var sh sdk.Handler
	app,ctx:= CreateTestInput([]auth.BaseAccount{})
	bankKeeper := app.bankKeeper
	header := abci.Header{Height: 1}
	app.BeginBlocker(ctx,abci.RequestBeginBlock{Header: header})
	power := int64(100000)
	amt := sdk.TokensFromTendermintPower(power)

	addr0, val0 := addrs[0], pks[0]
	accAddr0 := sdk.AccAddress(addr0)
	sh = staking.NewHandler(app.stakingKeeper)
	got := sh(ctx, staking.NewTestMsgCreateValidator(addr0, val0, amt))
	require.True(t,got.IsOK())

	addr1, val1 := addrs[1], pks[1]
	sh = staking.NewHandler(app.stakingKeeper)
	got = sh(ctx, staking.NewTestMsgCreateValidator(addr1, val1, amt))
	require.True(t,got.IsOK())

	//need more energy
	got = sh(ctx,staking.NewMsgUndelegate(sdk.AccAddress(addr0),addr0,
		sdk.Coin{Denom:defaultBondDenom,Amount:sdk.NewInt(100)}))
	assert.EqualValues(t,got.Code,sdk.CodeNotEnoughEnergy)

	//exchange energy
	got = sh(ctx,staking.NewMsgBurnBcvToEnergy(accAddr0,sdk.Coin{Denom:"ubcv",Amount:sdk.NewInt(10000)}))
	energyCoins := sdk.Coins{sdk.Coin{Denom:"ubcv",Amount:sdk.NewInt(10000)}}
	require.True(t,bankKeeper.HasCoins(ctx,sdk.AccAddress(addr0),energyCoins))


	for i:=0; i< int(app.stakingKeeper.MaxEntries(ctx)) -1;i++{
		got = sh(ctx,staking.NewMsgUndelegate(sdk.AccAddress(addr0),addr0,
			sdk.Coin{Denom:defaultBondDenom,Amount:sdk.NewInt(100)}))
		require.True(t,got.IsOK())
	}
	got = sh(ctx,staking.NewMsgUndelegate(sdk.AccAddress(addr0),addr0,
		sdk.Coin{Denom:defaultBondDenom,Amount:sdk.NewInt(100)}))
	require.True(t,got.IsOK())

	got = sh(ctx,staking.NewMsgUndelegate(sdk.AccAddress(addr0),addr0,
		sdk.Coin{Denom:defaultBondDenom,Amount:sdk.NewInt(100)}))
	require.Equal(t,staking.CodeInvalidDelegation,got.Code)

}

func TestShareFromToken(t *testing.T)  {
	var sh sdk.Handler
	app,ctx:= CreateTestInput([]auth.BaseAccount{})
	header := abci.Header{Height: 1}
	app.BeginBlocker(ctx,abci.RequestBeginBlock{Header: header})
	power := int64(100001)
	amt := sdk.TokensFromTendermintPower(power)
	stakingKeeper := app.stakingKeeper
	addr0, val0 := addrs[0], pks[0]
	//accAddr0 := sdk.AccAddress(addr0)
	sh = staking.NewHandler(app.stakingKeeper)
	got := sh(ctx, staking.NewTestMsgCreateValidator(addr0, val0, amt))
	require.True(t,got.IsOK())
	validator0,_:=stakingKeeper.GetValidator(ctx,sdk.ValAddress(sdk.AccAddress(addr0).Bytes()))
	validator0.DelegatorShares= sdk.NewDecFromIntWithPrec(sdk.NewInt(1000020000009999),4)
	fmt.Println(int(validator0.GetTokens().Int64()))
	for i:=1;i<int(validator0.GetTokens().Int64());i=i+99999 {
		testAmt := sdk.NewInt(int64(i))
		sharesAmt1,_  := validator0.SharesFromTokens(testAmt)
		sharesAmt2,_  := validator0.SharesFromTokensTruncated(testAmt)
		require.True(t,sharesAmt1.Equal(sharesAmt2))
	}

	
}

func TestRedelagete(t *testing.T)  {
	var sh sdk.Handler
	app,ctx:= CreateTestInput([]auth.BaseAccount{})
	header := abci.Header{Height: 1}
	app.BeginBlocker(ctx,abci.RequestBeginBlock{Header: header})
	power := int64(100000)
	amt := sdk.TokensFromTendermintPower(power)
	bankKeeper := app.bankKeeper
	stakeKeeper := app.stakingKeeper

	addr0, val0 := addrs[0], pks[0]
	addr1, val1 := addrs[1], pks[1]
	accAddr0 := sdk.AccAddress(addr0)
	accAddr1 := sdk.AccAddress(addr1)

	sh = staking.NewHandler(app.stakingKeeper)
	got := sh(ctx, staking.NewTestMsgCreateValidator(addr0, val0, amt))
	require.True(t,got.IsOK())

	sh = staking.NewHandler(app.stakingKeeper)
	got = sh(ctx, staking.NewTestMsgCreateValidator(addr1, val1, amt))
	require.True(t,got.IsOK())

	amtExchangeEnergy := sdk.NewInt(10000)
	amtBeginRedelate := sdk.NewInt(100)
	fmt.Println(bankKeeper.GetCoins(ctx,sdk.AccAddress(addr0)))

	got = sh(ctx,staking.NewMsgBurnBcvToEnergy(accAddr0,sdk.Coin{Denom:sdk.CHAIN_COIN_NAME_BCV,Amount:sdk.NewInt(10000)}))
	require.Equal(t,stakeKeeper.GetEnergyFromBcv(amtExchangeEnergy),
		bankKeeper.GetCoins(ctx,sdk.AccAddress(addr0)).AmountOf(sdk.CHAIN_COIN_NAME_ENERGY))


	got = sh(ctx,staking.NewMsgBeginRedelegate(accAddr0,
								sdk.ValAddress(accAddr0.Bytes()),
								sdk.ValAddress(accAddr1.Bytes()),
								sdk.Coin{Denom:defaultBondDenom,Amount:amtBeginRedelate}))



	}


func TestStakingMsgs(t *testing.T) {

	genTokens := sdk.TokensFromTendermintPower( 420000)
	bondTokens := sdk.TokensFromTendermintPower(100000)
	genCoin := sdk.NewCoin(sdk.DefaultBCVDemon, genTokens)
	bondCoin := sdk.NewCoin(sdk.DefaultBCVDemon, bondTokens)
	//feeCoins := sdk.NewCoin(sdk.DefaultMintDenom, sdk.NewInt(10000000000))

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

	// create validator
	description := staking.NewDescription("foo_moniker", "123", "123", "123")
	createValidatorMsg := staking.NewMsgCreateValidator(
		sdk.ValAddress(addr1), priv1.PubKey(), bondCoin, description, commissionMsg, sdk.OneInt(),
	)

	header := abci.Header{Height: app.LastBlockHeight() + 1}
	SignCheckDeliver(t, app.cdc, app.BaseApp, header, []sdk.Msg{createValidatorMsg}, []uint64{GetAccountNum(ctx,app,addr1)}, []uint64{GetAccountSequence(ctx,app,addr1)}, true, true, priv1)
	CheckBalance(t, app, addr1, sdk.Coins{genCoin.Sub(bondCoin)})



	header = abci.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	validator := checkValidator(t, app, app.stakingKeeper, addr1, true)
	require.Equal(t, sdk.ValAddress(addr1), validator.OperatorAddress)
	require.Equal(t, sdk.Bonded, validator.Status)
	require.True(sdk.IntEq(t, bondTokens, validator.BondedTokens()))


	header = abci.Header{Height: app.LastBlockHeight() + 1}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	// edit the validator
	description = staking.NewDescription("bar_moniker", "", "", "")
	editValidatorMsg := staking.NewMsgEditValidator(sdk.ValAddress(addr1), description, nil, nil)

	header = abci.Header{Height: app.LastBlockHeight() + 1}
	mock.SignCheckDeliver(t, app.cdc, app.BaseApp, header, []sdk.Msg{editValidatorMsg}, []uint64{GetAccountNum(ctx,app,addr1)}, []uint64{GetAccountSequence(ctx,app,addr1)}, true, true, priv1)

	validator = checkValidator(t, app, app.stakingKeeper, addr1, true)
	require.Equal(t, description, validator.Description)

	// delegate
	CheckBalance(t, app, addr2, sdk.Coins{ubcvStakeCoin})
	delegateMsg := staking.NewMsgDelegate(addr2, sdk.ValAddress(addr1), ubcvStakeCoin)
	header = abci.Header{Height: app.LastBlockHeight() + 1}
	SignCheckDeliver(t, app.cdc, app.BaseApp, header, []sdk.Msg{delegateMsg}, []uint64{GetAccountNum(ctx,app,addr2)}, []uint64{GetAccountSequence(ctx,app,addr2)}, true, true, priv2)
	CheckBalance(t, app, addr2, nil)
	checkDelegation(t, app, app.stakingKeeper, addr2, sdk.ValAddress(addr1), true, ubcvStakeCoin.Amount.ToDec())
	//
	// begin unbonding
	beginUnbondingMsg := staking.NewMsgUndelegate(addr2, sdk.ValAddress(addr1), ubcvStakeCoin)
	header = abci.Header{Height: app.LastBlockHeight() + 1}
	SignCheckDeliver(t, app.cdc, app.BaseApp, header, []sdk.Msg{beginUnbondingMsg}, []uint64{GetAccountNum(ctx,app,addr2)}, []uint64{GetAccountSequence(ctx,app,addr2)}, false, false, priv2)

	CheckBalance(t, app, addr2, nil)

}



func TestCoins(t *testing.T) {

	genCoin := sdk.NewCoin(sdk.DefaultBCVDemon, sdk.NewInt(1))
	feeCoins := sdk.NewCoin(sdk.DefaultMintDenom, sdk.NewInt(2))
	ubcvStakeCoin := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(3))
	coin1 := sdk.NewCoin("zabcd", sdk.NewInt(3))
	coin2 := sdk.NewCoin("abc", sdk.NewInt(3))

	acc1 := auth.BaseAccount{
		Address: addr1,
		Coins:   sdk.Coins{genCoin},
	}
	accs := []auth.BaseAccount{acc1}
	app,ctx := CreateTestInput(accs)
	app.bankKeeper.AddCoins(ctx,addr1,sdk.Coins{feeCoins})
	app.bankKeeper.AddCoins(ctx,addr1,sdk.Coins{ubcvStakeCoin})
	app.bankKeeper.AddCoins(ctx,addr1,sdk.Coins{coin1})
	app.bankKeeper.AddCoins(ctx,addr1,sdk.Coins{coin2})

	fmt.Println(app.accountKeeper.GetAccount(ctx,addr1).GetCoins().AmountOf("abc"))
	require.Equal(t,app.accountKeeper.GetAccount(ctx,addr1).GetCoins().AmountOf("abc"),sdk.NewInt(3))
}

func checkValidator(t *testing.T,app *BacApp, keeper staking.Keeper,
	addr sdk.AccAddress, expFound bool) staking.Validator {
	ctxCheck := app.BaseApp.NewContext(true, abci.Header{})
	validator, found := keeper.GetValidator(ctxCheck, sdk.ValAddress(addr1))
	require.Equal(t, expFound, found)
	return validator
}


func checkDelegation(
	t *testing.T, app *BacApp, keeper staking.Keeper, delegatorAddr sdk.AccAddress,
	validatorAddr sdk.ValAddress, expFound bool, expShares sdk.Dec,
) {

	ctxCheck := app.BaseApp.NewContext(true, abci.Header{})
	delegation, found := keeper.GetDelegation(ctxCheck, delegatorAddr, validatorAddr)
	if expFound {
		require.True(t, found)
		require.True(sdk.DecEq(t, expShares, delegation.Shares))

		return
	}

	require.False(t, found)
}



func TestValidatorByPowerIndex(t *testing.T) {
	genTokens := sdk.TokensFromTendermintPower( 420000)
	bondTokens := sdk.TokensFromTendermintPower(100000)
	genCoin := sdk.NewCoin(sdk.DefaultBCVDemon, genTokens)
	bondCoin := sdk.NewCoin(sdk.DefaultBCVDemon, bondTokens)
	//feeCoins := sdk.NewCoin(sdk.DefaultMintDenom, sdk.NewInt(10000000000))

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

	// create validator
	description := staking.NewDescription("foo_moniker", "123", "123", "123")
	createValidatorMsg := staking.NewMsgCreateValidator(
		sdk.ValAddress(addr1), priv1.PubKey(), bondCoin, description, commissionMsg, sdk.OneInt(),
	)

	header := abci.Header{Height: app.LastBlockHeight() + 1}
	SignCheckDeliver(t, app.cdc, app.BaseApp, header, []sdk.Msg{createValidatorMsg}, []uint64{GetAccountNum(ctx,app,addr1)}, []uint64{GetAccountSequence(ctx,app,addr1)}, true, true, priv1)
	CheckBalance(t, app, addr1, sdk.Coins{genCoin.Sub(bondCoin)})


	keeper := app.stakingKeeper
	// must end-block
	require.Equal(t, 1, len(keeper.GetAllValidators(ctx)))
	// verify the self-delegation exists
	bond, found := keeper.GetDelegation(ctx, addr1, sdk.ValAddress(addr1))
	require.True(t, found)
	gotBond := bond.Shares.RoundInt()
	require.Equal(t, bondCoin.Amount, gotBond)




	// verify that the by power index exists
	validator, found := keeper.GetValidator(ctx, sdk.ValAddress(addr1))
	require.True(t, found)
	power := staking.GetValidatorsByPowerIndexKey(validator)
	require.True(t, staking.ValidatorByPowerIndexExists(ctx, keeper, power))
}



func TestDuplicatesMsgCreateValidator(t *testing.T) {

	genTokens := sdk.TokensFromTendermintPower( 420000)
	bondTokens := sdk.TokensFromTendermintPower(100000)
	genCoin := sdk.NewCoin(sdk.DefaultBCVDemon, genTokens)
	bondCoin := sdk.NewCoin(sdk.DefaultBCVDemon, bondTokens)
	//feeCoins := sdk.NewCoin(sdk.DefaultMintDenom, sdk.NewInt(10000000000))

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

	// create validator
	description := staking.NewDescription("foo_moniker", "123", "123", "123")
	createValidatorMsg := staking.NewMsgCreateValidator(
		sdk.ValAddress(addr1), priv1.PubKey(), bondCoin, description, commissionMsg, sdk.OneInt(),
	)

	header := abci.Header{Height: app.LastBlockHeight() + 1}
	SignCheckDeliver(t, app.cdc, app.BaseApp, header, []sdk.Msg{createValidatorMsg}, []uint64{GetAccountNum(ctx,app,addr1)}, []uint64{GetAccountSequence(ctx,app,addr1)}, true, true, priv1)
	CheckBalance(t, app, addr1, sdk.Coins{genCoin.Sub(bondCoin)})



	validator, found := app.stakingKeeper.GetValidator(ctx, sdk.ValAddress(addr1))

	require.True(t, found)
	assert.Equal(t, sdk.Bonded, validator.Status)
}
