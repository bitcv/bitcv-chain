package app

import (


	"github.com/bitcv-chain/bitcv-chain/x/staking"

	abci "github.com/tendermint/tendermint/abci/types"
	sdk "github.com/bitcv-chain/bitcv-chain/types"

	"testing"
	"github.com/stretchr/testify/require"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/bitcv-chain/bitcv-chain/x/auth"
)
//have TestSlashRedelation when slash
func TestSlashRedelation(t *testing.T)  {
	var sh sdk.Handler
	app,ctx:= CreateTestInput([]auth.BaseAccount{})
	bankKeeper := app.bankKeeper
	slashKeeper := app.slashingKeeper
	stakingKeeper := app.stakingKeeper
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
	fmt.Println(accAddr0)
	staking.EndBlocker(ctx,stakingKeeper)
	height := int64(1)
	//40 blocks sign
	for ; height < slashKeeper.SignedBlocksWindow(ctx); height++ {
		ctx = ctx.WithBlockHeight(height)
		slashKeeper.HandleValidatorSignature(ctx, val0.Address(), power, true)
	}

	// 40 blocks missed
	for ; height < slashKeeper.SignedBlocksWindow(ctx)+(slashKeeper.SignedBlocksWindow(ctx)-slashKeeper.MinSignedPerWindow(ctx)); height++ {
		ctx = ctx.WithBlockHeight(height)
		slashKeeper.HandleValidatorSignature(ctx, val0.Address(), power, false)
	}


	info, found := slashKeeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(val0.Address()))
	require.True(t, found)
	require.Equal(t, int64(1), info.StartHeight)
	require.Equal(t, slashKeeper.SignedBlocksWindow(ctx)-slashKeeper.MinSignedPerWindow(ctx), info.MissedBlocksCounter)

	// validator should be bonded still
	validator, _ := stakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(val0))
	require.Equal(t, sdk.Bonded, validator.GetStatus())
	pool := stakingKeeper.GetPool(ctx)
	require.True(sdk.IntEq(t, amt.Add(amt), pool.BondedTokens))

	//validator redelate
	redelegateAmt := sdk.Coin{Denom:sdk.CHAIN_COIN_NAME_BCVSTAKE,Amount:sdk.NewInt(1000)}

	got = sh(ctx, staking.NewMsgBeginRedelegate(accAddr0,
		sdk.ValAddress(val0.Address()),sdk.ValAddress(val1.Address()), redelegateAmt))
	assert.EqualValues(t,got.Code,sdk.CodeNotEnoughEnergy)

	got = burnBcvToGetEnergy(ctx,app,accAddr0,sdk.NewInt(580000))
	got = sh(ctx, staking.NewMsgBeginRedelegate(accAddr0,
		sdk.ValAddress(val0.Address()),sdk.ValAddress(val1.Address()), redelegateAmt))
	assert.EqualValues(t,got.Code,sdk.CodeOK)
	fmt.Println(bankKeeper.GetCoins(ctx,accAddr0))


	ctx = ctx.WithBlockHeight(height)
	slashKeeper.HandleValidatorSignature(ctx, val0.Address(), power, false)
	info, found = slashKeeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(val0.Address()))


}



//have TestSlashRedelation when slash
func TestSlashUnbond(t *testing.T)  {
	var sh sdk.Handler
	app,ctx:= CreateTestInput([]auth.BaseAccount{})
	bankKeeper := app.bankKeeper
	slashKeeper := app.slashingKeeper
	stakingKeeper := app.stakingKeeper
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
	fmt.Println(accAddr0)
	staking.EndBlocker(ctx,stakingKeeper)
	height := int64(1)
	//40 blocks sign
	for ; height < slashKeeper.SignedBlocksWindow(ctx); height++ {
		ctx = ctx.WithBlockHeight(height)
		slashKeeper.HandleValidatorSignature(ctx, val0.Address(), power, true)
	}

	// 40 blocks missed
	for ; height < slashKeeper.SignedBlocksWindow(ctx)+(slashKeeper.SignedBlocksWindow(ctx)-slashKeeper.MinSignedPerWindow(ctx)); height++ {
		ctx = ctx.WithBlockHeight(height)
		slashKeeper.HandleValidatorSignature(ctx, val0.Address(), power, false)
	}


	info, found := slashKeeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(val0.Address()))
	require.True(t, found)
	require.Equal(t, slashKeeper.SignedBlocksWindow(ctx)-slashKeeper.MinSignedPerWindow(ctx), info.MissedBlocksCounter)

	// validator should be bonded still
	validator, _ := stakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(val0))
	require.Equal(t, sdk.Bonded, validator.GetStatus())
	pool := stakingKeeper.GetPool(ctx)
	require.True(sdk.IntEq(t, amt.Add(amt), pool.BondedTokens))

	//validator redelate
	unbondAmt := sdk.Coin{Denom:sdk.CHAIN_COIN_NAME_BCVSTAKE,Amount:sdk.NewInt(1000)}
	got = sh(ctx, staking.NewMsgUndelegate(accAddr0,
		sdk.ValAddress(val0.Address()), unbondAmt))
	assert.EqualValues(t,got.Code,sdk.CodeNotEnoughEnergy)
	got = burnBcvToGetEnergy(ctx,app,accAddr0,sdk.NewInt(580000))
	got = sh(ctx, staking.NewMsgUndelegate(accAddr0,
		sdk.ValAddress(val0.Address()), unbondAmt))
	assert.EqualValues(t,got.Code,sdk.CodeOK)

	staking.EndBlocker(ctx,stakingKeeper)

	ctx = ctx.WithBlockHeight(height)
	slashKeeper.HandleValidatorSignature(ctx, val0.Address(), power, false)
	info, found = slashKeeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(val0.Address()))


	fmt.Println(stakingKeeper.GetAllValidators(ctx))
	fmt.Println(bankKeeper.GetCoins(ctx,accAddr0))


	ctx = ctx.WithBlockTime(ctx.BlockHeader().Time.Add(stakingKeeper.UnbondingTime(ctx)))
	staking.EndBlocker(ctx,stakingKeeper)
	fmt.Println(bankKeeper.GetCoins(ctx,accAddr0))
	fmt.Println(app.slashingKeeper.GetValidatorSigningInfo(ctx,sdk.ConsAddress(val0.Address())))

}
