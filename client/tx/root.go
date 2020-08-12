package tx

import (
	"github.com/gorilla/mux"

	"github.com/bitcv-chain/bitcv-chain/client/context"
	"github.com/bitcv-chain/bitcv-chain/codec"
)

// register REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/txs/{hash}", QueryTxRequestHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/txs", QueryTxsByTagsRequestHandlerFn(cliCtx, cdc)).Methods("GET")
	r.HandleFunc("/txs", BroadcastTxRequest(cliCtx, cdc)).Methods("POST")
	r.HandleFunc("/txs/encode", EncodeTxRequestHandlerFn(cdc, cliCtx)).Methods("POST")
}
