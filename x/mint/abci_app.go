package mint

import (
	sdk "github.com/bitcv-chain/bitcv-chain/types"
)

// Inflate every block, update inflation parameters once per hour
func BeginBlocker(ctx sdk.Context, k Keeper) {

	// fetch stored minter & params
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	totalSupply := k.sk.TotalTokens(ctx)
	k.SetMinter(ctx, minter)

	mintedCoin := minter.NextReduceAnnualProvisions(params, ctx.BlockHeight())
	mintedCoins := sdk.Coins{sdk.NewCoin(params.MintDenom,mintedCoin)}

	// Collect fee by fee pool
	k.fck.AddCollectedFees(ctx, mintedCoins)
	// Record total new coins
	k.fck.NewCoinInflateSupply(ctx, mintedCoins)
	// Get remain total new coin(bac)
	newCoinBalance := k.fck.GetNewTokenCirculation(ctx)

	ctx.Logger().Info("module","total_BCVSTAKE_OK",totalSupply, "BAC_GEN_EVERY_BLOCK",mintedCoin, "total_BAC",newCoinBalance,"collect_FEE",k.fck.GetCollectedFees(ctx))
}