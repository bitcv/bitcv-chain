package mint

import (
	sdk "github.com/bitcv-chain/bitcv-chain/types"
	"github.com/bitcv-chain/bitcv-chain/x/auth"
)

// expected staking keeper
type StakingKeeper interface {
	TotalTokens(ctx sdk.Context) sdk.Int
	BondedRatio(ctx sdk.Context) sdk.Dec
	InflateSupply(ctx sdk.Context, newTokens sdk.Int)
}

// expected fee collection keeper interface
type FeeCollectionKeeper interface {
	AddCollectedFees(sdk.Context, sdk.Coins) sdk.Coins
	NewCoinInflateSupply(ctx sdk.Context, newTokens sdk.Coins)
	GetNewTokenCirculation(ctx sdk.Context) sdk.Coins
	GetCollectedFees(ctx sdk.Context) sdk.Coins
	GetBacManagePool(ctx sdk.Context) auth.BacManagePool
	GetBacManagePoolForApi(ctx sdk.Context) auth.BacManagePoolForApi

}
