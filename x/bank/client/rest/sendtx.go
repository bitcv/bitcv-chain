package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/bitcv-chain/bitcv-chain/client/context"
	clientrest "github.com/bitcv-chain/bitcv-chain/client/rest"
	"github.com/bitcv-chain/bitcv-chain/codec"
	"github.com/bitcv-chain/bitcv-chain/crypto/keys"
	sdk "github.com/bitcv-chain/bitcv-chain/types"
	"github.com/bitcv-chain/bitcv-chain/types/rest"

	"github.com/bitcv-chain/bitcv-chain/x/bank"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, kb keys.Keybase) {
	r.HandleFunc("/bank/accounts/{address}/transfers", SendRequestHandlerFn(cdc, kb, cliCtx)).Methods("POST")
	r.HandleFunc("/bank/tokens", QueryTokensRequestHandlerFn(cliCtx,cdc)).Methods("GET")
	r.HandleFunc("/bank/token/{inner_name}", QueryTokenByInnerNameRequestHandlerFn(cliCtx,cdc)).Methods("GET")
	r.HandleFunc("/bank/edata/{address}", QueryEdataRequestHandlerFn(cliCtx,cdc)).Methods("GET")
}

// SendReq defines the properties of a send request's body.
type SendReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Amount  sdk.Coins    `json:"amount"`
}

var msgCdc = codec.New()

func init() {
	bank.RegisterCodec(msgCdc)
}

// SendRequestHandlerFn - http request handler to send coins to a address.
func SendRequestHandlerFn(cdc *codec.Codec, kb keys.Keybase, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32Addr := vars["address"]

		toAddr, err := sdk.AccAddressFromBech32(bech32Addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req SendReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := bank.NewMsgSend(fromAddr, toAddr, req.Amount)
		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func QueryTokensRequestHandlerFn(cliCtx context.CLIContext,cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryWithData("custom/bank/tokens", nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func QueryTokenByInnerNameRequestHandlerFn(cliCtx context.CLIContext,cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		innerName := vars["inner_name"]
		params := bank.QueryTokenParams{InnerName:innerName}
		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		res, err := cliCtx.QueryWithData("custom/bank/token", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}


func QueryEdataRequestHandlerFn(cliCtx context.CLIContext,cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		toAddr, err := sdk.AccAddressFromBech32(vars["address"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := bank.QueryEdataParams{AccAddr:toAddr}
		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		res, err := cliCtx.QueryWithData("custom/bank/edata", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
