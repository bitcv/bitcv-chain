package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/bitcv-chain/bitcv-chain/client/context"
	"github.com/bitcv-chain/bitcv-chain/codec"
	sdk "github.com/bitcv-chain/bitcv-chain/types"
	"github.com/bitcv-chain/bitcv-chain/types/rest"

	"github.com/bitcv-chain/bitcv-chain/x/auth"
)

// register REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {
	r.HandleFunc(
		"/auth/accounts/{address}",
		QueryAccountRequestHandlerFn(storeName, cdc, context.GetAccountDecoder(cdc), cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/bank/balances/{address}",
		QueryBalancesRequestHandlerFn(storeName, cdc, context.GetAccountDecoder(cdc), cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/bank/supply/bcv",
		QuerySupplyBcvRequestHandlerFn(storeName, cdc, context.GetAccountDecoder(cdc), cliCtx),
	).Methods("GET")
}

// query accountREST Handler
func QueryAccountRequestHandlerFn(
	storeName string, cdc *codec.Codec,
	decoder auth.AccountDecoder, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32addr := vars["address"]

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, err := cliCtx.QueryStore(auth.AddressStoreKey(addr), storeName)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// the query will return empty account if there is no data
		if len(res) == 0 {
			rest.PostProcessResponse(w, cdc, auth.BaseAccount{}, cliCtx.Indent)
			return
		}

		// decode the value
		account, err := decoder(res)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, account, cliCtx.Indent)
	}
}

// query accountREST Handler
func QueryBalancesRequestHandlerFn(
	storeName string, cdc *codec.Codec,
	decoder auth.AccountDecoder, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		bech32addr := vars["address"]

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, err := cliCtx.QueryStore(auth.AddressStoreKey(addr), storeName)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// the query will return empty if there is no data for this account
		if len(res) == 0 {
			rest.PostProcessResponse(w, cdc, sdk.Coins{}, cliCtx.Indent)
			return
		}

		// decode the value
		account, err := decoder(res)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, account.GetCoins(), cliCtx.Indent)
	}
}



// query accountREST Handler
func QuerySupplyBcvRequestHandlerFn(
	storeName string, cdc *codec.Codec,
	decoder auth.AccountDecoder, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")


		//换票销毁
		var burnBcvAmountFromBuyStake sdk.Int =  sdk.ZeroInt()
		res, err := cliCtx.QueryStore(auth.AddressStoreKey(sdk.AccAddrBcvBurnFromBuyBcvStake), storeName)
		if err == nil  && len(res) > 0 {
			account, err := decoder(res)
			if err == nil{
				burnBcvAmountFromBuyStake = account.GetCoins().AmountOf(sdk.DefaultBCVDemon)
			}
		}


		//购电销毁
		var burnBcvAmountFromBuyEnergy sdk.Int = sdk.ZeroInt()
		res, err = cliCtx.QueryStore(auth.AddressStoreKey(sdk.AccAddrEnergyPool), storeName)
		if err == nil  && len(res) > 0 {
			account, err := decoder(res)
			if err == nil{
				burnBcvAmountFromBuyEnergy = account.GetCoins().AmountOf(sdk.DefaultBCVDemon)
			}
		}

		//直接销毁
		var burnBcvAmountFromManual sdk.Int  = sdk.ZeroInt()
		res, err = cliCtx.QueryStore(auth.AddressStoreKey(sdk.AccAddrBurn), storeName)
		if err == nil  && len(res) > 0 {
			account, err := decoder(res)
			if err == nil{
				burnBcvAmountFromManual = account.GetCoins().AmountOf(sdk.DefaultBCVDemon)
			}
		}
		//总发行
		allSupply,_:= sdk.NewIntFromString(sdk.CHAIN_PARAM_BCV_AMOUNT)
		allBurn:= burnBcvAmountFromBuyStake.Add(burnBcvAmountFromBuyEnergy).Add(burnBcvAmountFromManual)
		burnBcvForApi := BurnBcvForApi{
			FromBuyStake:burnBcvAmountFromBuyStake,
			FromBuyEnergy:burnBcvAmountFromBuyEnergy,
			FromManual:burnBcvAmountFromManual,
			AllBurn:allBurn,
			AllSupply:allSupply,
		}




		rest.PostProcessResponse(w, cdc, burnBcvForApi, cliCtx.Indent)
	}
}

type BurnBcvForApi struct{
	FromBuyStake sdk.Int `json:"from_buy_stake"`
	FromBuyEnergy sdk.Int `json:"from_buy_energy"`
	FromManual 	sdk.Int `json:"from_manual"`
	AllBurn     sdk.Int	`json:"all_burn"`
	AllSupply   sdk.Int	`json:"all_supply"`
}
