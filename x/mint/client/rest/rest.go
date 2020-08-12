package rest

import (
	"github.com/bitcv-chain/bitcv-chain/client/context"
	"github.com/bitcv-chain/bitcv-chain/codec"
	"github.com/gorilla/mux"
)

// RegisterRoutes registers minting module REST handlers on the provided router.
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	registerQueryRoutes(cliCtx, r, cdc)
}
