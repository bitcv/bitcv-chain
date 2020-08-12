package bank
import (

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/bitcv-chain/bitcv-chain/types"
	"github.com/bitcv-chain/bitcv-chain/codec"
	"fmt"
)

// query endpoints supported by the auth Querier
const (
	QueryToken = "token"
	QueryTokens = "tokens"
	QueryEdata	= "edata"
)


// creates a querier for auth REST endpoints
func NewQuerier(k Keeper,cdc *codec.Codec) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryToken:
			return queryToken(ctx, cdc,req, k)
		case QueryTokens:
			return queryTokens(ctx, cdc, k)
		case QueryEdata:
			return  queryEdata(ctx,cdc,req,k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown auth query endpoint")
		}
	}
}



func queryToken(ctx sdk.Context, cdc *codec.Codec,req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params QueryTokenParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", errRes.Error()))
	}

	isExist,token := k.GetTokenByInnerName(ctx,params.InnerName)
	if !isExist{
		return nil, sdk.ErrInternal("token is not exist")
	}
	bz, err := codec.MarshalJSONIndent(cdc, token)

	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryEdata(ctx sdk.Context, cdc *codec.Codec,req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params QueryEdataParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", errRes.Error()))
	}

	found,edatas := k.GetEdataByAccount(ctx,params.AccAddr)
	if !found{
		return nil, sdk.ErrInternal("edata is not exist")
	}
	bz, err := codec.MarshalJSONIndent(cdc, edatas)

	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryTokens(ctx sdk.Context, cdc *codec.Codec, k Keeper) ([]byte, sdk.Error) {
	tokens := k.GetTokens(ctx)
	bz, err := codec.MarshalJSONIndent(cdc, tokens)

	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}


type QueryTokenParams struct {
	InnerName string
}

type QueryEdataParams struct {
	AccAddr sdk.AccAddress
}

