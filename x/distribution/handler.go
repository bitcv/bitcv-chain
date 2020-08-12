package distribution

import (
	sdk "github.com/bitcv-chain/bitcv-chain/types"
	"github.com/bitcv-chain/bitcv-chain/x/distribution/keeper"
	"github.com/bitcv-chain/bitcv-chain/x/distribution/tags"
	"github.com/bitcv-chain/bitcv-chain/x/distribution/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		// NOTE msg already has validate basic run
		switch msg := msg.(type) {
		case types.MsgSetWithdrawAddress:
			return handleMsgModifyWithdrawAddress(ctx, msg, k)
		case types.MsgWithdrawDelegatorReward:
			return handleMsgWithdrawDelegatorReward(ctx, msg, k)
		case types.MsgWithdrawValidatorCommission:
			return handleMsgWithdrawValidatorCommission(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("invalid message parse in distribution module").Result()
		}
	}
}

// These functions assume everything has been authenticated (ValidateBasic passed, and signatures checked)

func handleMsgModifyWithdrawAddress(ctx sdk.Context, msg types.MsgSetWithdrawAddress, k keeper.Keeper) sdk.Result {

	err := k.SetWithdrawAddr(ctx, msg.DelegatorAddress, msg.WithdrawAddress)
	if err != nil {
		return err.Result()
	}

	tags := sdk.NewTags(
		tags.Delegator, []byte(msg.DelegatorAddress.String()),
	)
	return sdk.Result{
		Tags: tags,
	}
}

func handleMsgWithdrawDelegatorReward(ctx sdk.Context, msg types.MsgWithdrawDelegatorReward, k keeper.Keeper) sdk.Result {
	rewards,costEnergyCoins,err := k.WithdrawDelegationRewards(ctx, msg.DelegatorAddress, msg.ValidatorAddress)
	if err != nil {
		return err.Result()
	}

	addr := k.GetDelegatorWithdrawAddr(ctx,msg.DelegatorAddress).String()
	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Rewards, rewards.String(),
			tags.Delegator, []byte(msg.DelegatorAddress.String()),
			tags.Validator, []byte(msg.ValidatorAddress.String()),
			tags.AmountCostEnergy,costEnergyCoins.String(),
			tags.WithdrawDstAddress,[]byte(addr),
		),
	}
}

func handleMsgWithdrawValidatorCommission(ctx sdk.Context, msg types.MsgWithdrawValidatorCommission, k keeper.Keeper) sdk.Result {
	commission, err := k.WithdrawValidatorCommission(ctx, msg.ValidatorAddress)
	if err != nil {
		return err.Result()
	}

	//addr := k.GetDelegatorWithdrawAddr(ctx,sdk.AccAddress(msg.ValidatorAddress.Bytes())).String()
	addr := sdk.AccAddress(msg.ValidatorAddress.Bytes()).String()
	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Commission, commission.String(),
			tags.Validator, []byte(msg.ValidatorAddress.String()),
			tags.WithdrawDstAddress,[]byte(addr),
		),
	}
}
