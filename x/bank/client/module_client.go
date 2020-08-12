package client

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"

	"github.com/bitcv-chain/bitcv-chain/client"
	"github.com/bitcv-chain/bitcv-chain/x/bank/types"
	"github.com/bitcv-chain/bitcv-chain/x/bank/client/cli"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	bankTxCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "bank transaction subcommands",
	}

	bankTxCmd.AddCommand(client.PostCommands(
		cli.GetSendTxCmd(mc.cdc),
		cli.GetIssueTokenTxCmd(mc.cdc),
		cli.GetRedeemTxCmd(mc.cdc),
		cli.GetAddMaginTxCmd(mc.cdc),
		cli.GetSaveTxCmd(mc.cdc),
	)...)

	return bankTxCmd
}


// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {

	return nil

}