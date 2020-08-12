package crisis

import (
	sdk "github.com/bitcv-chain/bitcv-chain/types"
)

// expected bank keeper
type DistrKeeper interface {
	DistributeFeePool(ctx sdk.Context, amount sdk.Coins, receiveAddr sdk.AccAddress) sdk.Error
}

// expected fee collection keeper
type FeeCollectionKeeper interface {
	AddCollectedFees(ctx sdk.Context, coins sdk.Coins) sdk.Coins
	BurnSupply(ctx sdk.Context, tokens sdk.Coins)
}

// expected bank keeper
type BankKeeper interface {
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error)
}
