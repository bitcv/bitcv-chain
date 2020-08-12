package bank

import (
	sdk "github.com/bitcv-chain/bitcv-chain/types"
	bacv1"github.com/bitcv-chain/bitcv-chain/bacchain/v1_00"
	"encoding/hex"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSend:
			return handleMsgSend(ctx, k, msg)
		//case MsgMultiSend:
		//	return handleMsgMultiSend(ctx, k, msg)

		case MsgIssueToken:
			return handleMsgIssueToken(ctx,k,msg)

		case MsgRedeem:
			return handleMsgRedeem(ctx,k,msg)

		case MsgAddMargin:
			return handleMsgAddMargin(ctx,k,msg)

		case MsgEdata:
			return handleMsgEdata(ctx,k,msg)

		default:
			errMsg := "Unrecognized bank Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSend.
func handleMsgSend(ctx sdk.Context, k Keeper, msg MsgSend) sdk.Result {
	if !bacv1.CheckSendEnable(k.GetSendEnabled(ctx),ctx.BlockHeight()){
		return ErrSendDisabled(k.Codespace()).Result()
	}
	if !checkCanTran(msg.Amount){
		return  ErrInvaidCoin(k.Codespace(),msg.Amount.String()).Result()
	}
	tags, err := k.SendCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}

// Handle MsgMultiSend.
func handleMsgMultiSend(ctx sdk.Context, k Keeper, msg MsgMultiSend) sdk.Result {
	// NOTE: totalIn == totalOut should already have been checked
	if !bacv1.CheckSendEnable(k.GetSendEnabled(ctx),ctx.BlockHeight()){
		return ErrSendDisabled(k.Codespace()).Result()
	}
	for _,input  := range msg.Inputs{
		if !checkCanTran(input.Coins){
			return  ErrInvaidCoin(k.Codespace(),input.Coins.String()).Result()
		}
	}
	tags, err := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}

func checkCanTran(coins sdk.Coins)  bool {
	for _,coin := range coins {
		if coin.Denom == sdk.DefaultBondDenom || coin.Denom == sdk.CHAIN_COIN_NAME_ENERGY{
			return false;
		}
	}
	return true
}



// Handle MsgSend.
func handleMsgIssueToken(ctx sdk.Context, k Keeper, msg MsgIssueToken) sdk.Result {
	err := CheckMsgIssueToken(msg)
	if err != nil {
		return err.Result()
	}

	var token IssueToken

	token.OuterName = msg.OuterName
	token.OwnerAddress   = msg.OwnerAddress
	token.SupplyNum = msg.SupplyNum
	token.Margin = msg.Margin
	token.Precision = msg.Precision
	token.Website = msg.Website
	token.Description = msg.Description


	tags, err ,_:= k.IssueToken(ctx,token)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}


// Handle MsgSend.
func handleMsgRedeem(ctx sdk.Context, k Keeper, msg MsgRedeem) sdk.Result {
	tags, err := k.Redeem(ctx,msg.Account,msg.Amount)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func handleMsgAddMargin(ctx sdk.Context, k Keeper, msg MsgAddMargin) sdk.Result {
	err := msg.ValidateBasic()
	if err!= nil{
		return err.Result()
	}
	tags,err,_:= k.AddMarginByInnerName(ctx,msg.Account,msg.InnerName,msg.Amount)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

// Handle MsgSend.
func handleMsgEdata(ctx sdk.Context, k Keeper, msg MsgEdata) sdk.Result {
	edata := Edata{Utype:msg.Utype,Data:hex.EncodeToString([]byte(msg.Data))}
	tags ,err := k.SetEdataByAccount(ctx,msg.Account,edata)
	if err != nil{
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}
