package mint

import (
	"fmt"

	"github.com/bitcv-chain/bitcv-chain/codec"
	sdk "github.com/bitcv-chain/bitcv-chain/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Query endpoints supported by the minting querier
const (
	QueryParameters       = "parameters"
	QueryInflation        = "inflation"
	QueryAnnualProvisions = "annual_provisions"
	QueryBacPool	      = "bac_pool"

)

// NewQuerier returns a minting Querier handler.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, _ abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryParameters:
			return queryParams(ctx, k)

		case QueryInflation:
			return queryInflation(ctx, k)

		case QueryAnnualProvisions:
			return queryAnnualProvisions(ctx, k)
		case QueryBacPool:
			return queryBacPool(ctx, k)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("unknown minting query endpoint: %s", path[0]))
		}
	}
}

func queryParams(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(k.cdc, params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal JSON", err.Error()))
	}

	return res, nil
}

func queryInflation(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	minter := k.GetMinter(ctx)

	res, err := codec.MarshalJSONIndent(k.cdc, minter.Inflation)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal JSON", err.Error()))
	}

	return res, nil
}

func queryAnnualProvisions(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	minter := k.GetMinter(ctx)

	res, err := codec.MarshalJSONIndent(k.cdc, minter.AnnualProvisions)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal JSON", err.Error()))
	}

	return res, nil
}

func queryBacPool(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	bacManagepool := k.fck.GetBacManagePoolForApi(ctx)
	res, err := codec.MarshalJSONIndent(k.cdc, bacManagepool)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal JSON", err.Error()))
	}

	return res, nil
}