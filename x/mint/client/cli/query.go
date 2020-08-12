package cli

import (
	"fmt"

	"github.com/bitcv-chain/bitcv-chain/client/context"
	"github.com/bitcv-chain/bitcv-chain/codec"
	sdk "github.com/bitcv-chain/bitcv-chain/types"
	"github.com/bitcv-chain/bitcv-chain/x/mint"
	"github.com/spf13/cobra"
	"github.com/bitcv-chain/bitcv-chain/x/auth"
)

// GetCmdQueryParams implements a command to return the current minting
// parameters.
func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the current minting parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", mint.QuerierRoute, mint.QueryParameters)
			res, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params mint.Params
			if err := cdc.UnmarshalJSON(res, &params); err != nil {
				return err
			}

			return cliCtx.PrintOutput(params)
		},
	}
}

// GetCmdQueryInflation implements a command to return the current minting
// inflation value.
func GetCmdQueryInflation(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "inflation",
		Short: "Query the current minting inflation value",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", mint.QuerierRoute, mint.QueryInflation)
			res, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var inflation sdk.Dec
			if err := cdc.UnmarshalJSON(res, &inflation); err != nil {
				return err
			}

			return cliCtx.PrintOutput(inflation)
		},
	}
}

// GetCmdQueryAnnualProvisions implements a command to return the current minting
// annual provisions value.
func GetCmdQueryAnnualProvisions(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "annual-provisions",
		Short: "Query the current minting annual provisions value",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", mint.QuerierRoute, mint.QueryAnnualProvisions)
			res, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var inflation sdk.Dec
			if err := cdc.UnmarshalJSON(res, &inflation); err != nil {
				return err
			}

			return cliCtx.PrintOutput(inflation)
		},
	}
}



//query bac pool
func GetCmdBacPool(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "bac_pool",
		Short: "Query the current supply and  burn bac ",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", mint.QuerierRoute, mint.QueryBacPool)
			res, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var bacpool auth.BacManagePool

			if err := cdc.UnmarshalJSON(res, &bacpool); err != nil {
				return err
			}

			return cliCtx.PrintOutput(bacpool)
		},
	}
}