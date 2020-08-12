package cli

import (
	"fmt"

	"github.com/bitcv-chain/bitcv-chain/client"
	"github.com/bitcv-chain/bitcv-chain/client/context"
	"github.com/bitcv-chain/bitcv-chain/client/utils"
	"github.com/bitcv-chain/bitcv-chain/codec"
	sdk "github.com/bitcv-chain/bitcv-chain/types"
	authtxb "github.com/bitcv-chain/bitcv-chain/x/auth/client/txbuilder"
	"github.com/bitcv-chain/bitcv-chain/x/bank"

	"github.com/spf13/cobra"
	"strconv"
)

const (
	flagTo     = "to"
	flagAmount = "amount"
)

// SendTxCmd will create a send tx and sign it with the given key.
func GetSendTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [to_address] [amount]",
		Short: "Create and sign a send tx",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			to, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			// parse coins trying to be sent
			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			account, err := cliCtx.GetAccount(from)
			if err != nil {
				return err
			}

			// ensure account has enough coins
			if !account.GetCoins().IsAllGTE(coins) {
				return fmt.Errorf("address %s doesn't have enough coins to pay for this transaction", from)
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := bank.NewMsgSend(from, to, coins)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.MarkFlagRequired(client.FlagFrom)

	return cmd
}



// baccli tx bank issue  lhy  100000000 1000000000ubcv   6   www.bitcv.com  test --from bac1hmkvg4ylpp2vk8aj0tzuagh7yrvjhpgh5durxl --fees 10000000nbac
func GetIssueTokenTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue [outer_name] [supply_num] [margin] [precision] [website] [description] --from  mykey",
		Short: "issue a new token",
		Long:  "baccli tx bank issue [outer_name] [supply_num] [margin] [precision] [website] [description] --from  mykey",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			//发行者
			ownerAddr := cliCtx.GetFromAddress()

			//币外部名字3-10个字符;大小写敏感
			outerName := args[0]

			//发行总量
			supplyNum, _ := sdk.NewIntFromString(args[1])

			margin, err := sdk.ParseCoin(args[2])
			if err != nil {
				return err
			}

			//精度
			a,_ := strconv.ParseUint(args[3],10,8)
			if err !=nil{
				return  err
			}
			precision  := uint8(a)

			website  := args[4]

			Description := args[5]

			msg := bank.NewMsgIssueToken(ownerAddr, outerName,supplyNum,margin,precision,website,Description)
			err = bank.CheckMsgIssueToken(msg)
			if err != nil{
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.MarkFlagRequired(client.FlagFrom)

	return cmd
}


//baccli tx bank  redeem 9999lhy-6c7 --from bac19qp38ktnphpy0v8883ht8yw56y70v788vgde9n --fees 2000000nbac
func GetRedeemTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeem [coin] --from mykey",
		Short: "redeem token",
		Long:  "baccli tx bank redeem {coin} --from mykey",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			fromAddr := cliCtx.GetFromAddress()
			amount, err := sdk.ParseCoin(args[0])
			if err != nil {
				return err
			}
			msg := bank.NewMsgRedeem(fromAddr,amount)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.MarkFlagRequired(client.FlagFrom)

	return cmd
}

func GetAddMaginTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-margin [inner_name] [amount] --from mykey",
		Short: "add-margin token",
		Long:  "baccli tx add-margin  {coin} {amount} --from mykey",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			fromAddr := cliCtx.GetFromAddress()

			innserName := args[0]
			amount, err := sdk.ParseCoin(args[1])
			if err != nil {
				return err
			}
			msg := bank.NewMsgAddMargin(fromAddr,innserName,amount)
			msg.ValidateBasic()
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.MarkFlagRequired(client.FlagFrom)

	return cmd
}



func GetSaveTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "save [utype] [msg] --from mykey",
		Short: "save edata",
		Long:  "baccli tx bank save  {utype}  {msg} --from mykey",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			fromAddr := cliCtx.GetFromAddress()
			a,err := strconv.ParseUint(args[0],10,8)
			if err != nil {
				return err
			}
			utype := uint8(a)
			data := args[1]


			msg := bank.NewMsgEdata(fromAddr,utype,data);
			err = msg.ValidateBasic()
			if err != nil{
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.MarkFlagRequired(client.FlagFrom)

	return cmd
}

