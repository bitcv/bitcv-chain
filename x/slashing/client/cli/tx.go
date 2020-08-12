package cli

import (
	"github.com/bitcv-chain/bitcv-chain/client/context"
	"github.com/bitcv-chain/bitcv-chain/client/utils"
	"github.com/bitcv-chain/bitcv-chain/codec"
	sdk "github.com/bitcv-chain/bitcv-chain/types"
	authtxb "github.com/bitcv-chain/bitcv-chain/x/auth/client/txbuilder"
	"github.com/bitcv-chain/bitcv-chain/x/slashing"

	"github.com/spf13/cobra"
)

// GetCmdUnjail implements the create unjail validator command.
func GetCmdUnjail(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "unjail",
		Args:  cobra.NoArgs,
		Short: "unjail validator previously jailed for downtime",
		Long: `unjail a jailed validator:

$ baccli tx slashing unjail --from mykey
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			valAddr := cliCtx.GetFromAddress()

			msg := slashing.NewMsgUnjail(sdk.ValAddress(valAddr))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}
}
