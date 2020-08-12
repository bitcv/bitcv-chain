package bank

import (
	sdk "github.com/bitcv-chain/bitcv-chain/types"
)

// expected crisis keeper
type CrisisKeeper interface {
	RegisterRoute(moduleName, route string, invar sdk.Invariant)
}
