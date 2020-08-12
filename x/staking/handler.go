package staking

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"

	sdk "github.com/bitcv-chain/bitcv-chain/types"
	"github.com/bitcv-chain/bitcv-chain/x/staking/keeper"
	"github.com/bitcv-chain/bitcv-chain/x/staking/tags"
	"github.com/bitcv-chain/bitcv-chain/x/staking/types"
	"time"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		// NOTE msg already has validate basic run
		switch msg := msg.(type) {
		case types.MsgCreateValidator:
			return handleMsgCreateValidator(ctx, msg, k)
		case types.MsgEditValidator:
			return handleMsgEditValidator(ctx, msg, k)
		case types.MsgDelegate:
			return handleMsgDelegate(ctx, msg, k)
		case types.MsgBeginRedelegate:
			return handleMsgBeginRedelegate(ctx, msg, k)
		case types.MsgUndelegate:
			return handleMsgUndelegate(ctx, msg, k)
		case types.MsgBurnBcvToEnergy:
			return handleMsgBurnBcvToEnergy(ctx, msg, k)
		case types.MsgBurnBcvToStake:
			return handleMsgBurnBcvToBcvstake(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("invalid message parse in staking module").Result()
		}
	}
}

// Called every block, update validator set
func EndBlocker(ctx sdk.Context, k keeper.Keeper) ([]abci.ValidatorUpdate, sdk.Tags) {
	resTags := sdk.NewTags()

	// Calculate validator set changes.
	//
	// NOTE: ApplyAndReturnValidatorSetUpdates has to come before
	// UnbondAllMatureValidatorQueue.
	// This fixes a bug when the unbonding period is instant (is the case in
	// some of the tests). The test expected the validator to be completely
	// unbonded after the Endblocker (go from Bonded -> Unbonding during
	// ApplyAndReturnValidatorSetUpdates and then Unbonding -> Unbonded during
	// UnbondAllMatureValidatorQueue).
	validatorUpdates := k.ApplyAndReturnValidatorSetUpdates(ctx)

	// Unbond all mature validators from the unbonding queue.
	k.UnbondAllMatureValidatorQueue(ctx)

	// Remove all mature unbonding delegations from the ubd queue.
	matureUnbonds := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	for _, dvPair := range matureUnbonds {
		err := k.CompleteUnbonding(ctx, dvPair.DelegatorAddress, dvPair.ValidatorAddress)
		if err != nil {
			continue
		}

		resTags = resTags.AppendTags(sdk.NewTags(
			tags.Action, tags.ActionCompleteUnbonding,
			tags.Delegator, dvPair.DelegatorAddress.String(),
			tags.SrcValidator, dvPair.ValidatorAddress.String(),
		))
	}

	// Remove all mature redelegations from the red queue.
	matureRedelegations := k.DequeueAllMatureRedelegationQueue(ctx, ctx.BlockHeader().Time)
	for _, dvvTriplet := range matureRedelegations {
		err := k.CompleteRedelegation(ctx, dvvTriplet.DelegatorAddress,
			dvvTriplet.ValidatorSrcAddress, dvvTriplet.ValidatorDstAddress)
		if err != nil {
			continue
		}

		resTags = resTags.AppendTags(sdk.NewTags(
			tags.Action, tags.ActionCompleteRedelegation,
			tags.Delegator, dvvTriplet.DelegatorAddress.String(),
			tags.SrcValidator, dvvTriplet.ValidatorSrcAddress.String(),
			tags.DstValidator, dvvTriplet.ValidatorDstAddress.String(),
		))
	}

	return validatorUpdates, resTags
}

// These functions assume everything has been authenticated,
// now we just perform action and save

func handleMsgCreateValidator(ctx sdk.Context, msg types.MsgCreateValidator, k keeper.Keeper) sdk.Result {
	// check to see if the pubkey or sender has been registered before
	if _, found := k.GetValidator(ctx, msg.ValidatorAddress); found {
		return ErrValidatorOwnerExists(k.Codespace()).Result()
	}

	if _, found := k.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(msg.PubKey)); found {
		return ErrValidatorPubKeyExists(k.Codespace()).Result()
	}
	if msg.Value.Amount.LT(sdk.NewInt(100000000000)){
		return ErrCreateValidatorAmountInvalid(k.Codespace(),"too small to create validator").Result()
	}

	//验证人使用bcv创建，抵押成bcvstake
	if msg.Value.Denom != sdk.DefaultBCVDemon {
			return ErrBadDenom(k.Codespace()).Result()
	}


	if _, err := msg.Description.EnsureLength(); err != nil {
		return err.Result()
	}

	if ctx.ConsensusParams() != nil {
		tmPubKey := tmtypes.TM2PB.PubKey(msg.PubKey)
		if !common.StringInSlice(tmPubKey.Type, ctx.ConsensusParams().Validator.PubKeyTypes) {
			return ErrValidatorPubKeyTypeUnsupported(k.Codespace(),
				tmPubKey.Type,
				ctx.ConsensusParams().Validator.PubKeyTypes).Result()
		}
	}

	validator := NewValidator(msg.ValidatorAddress, msg.PubKey, msg.Description)
	//设置validator的CanUnbondToken
	validator = validator.SetCanUnbondToken(msg.Value.Amount)

	commission := NewCommissionWithTime(
		msg.Commission.Rate, msg.Commission.MaxRate,
		msg.Commission.MaxChangeRate, ctx.BlockHeader().Time,
	)
	validator, err := validator.SetInitialCommission(commission)
	if err != nil {
		return err.Result()
	}

	//给bcvstakeaddr 转移资产
	err = k.ValidatorPledgeBcvToBcvstake(ctx,msg.DelegatorAddress,msg.Value.Amount)
	if err != nil{
		return err.Result()
	}

	validator.MinSelfDelegation = msg.MinSelfDelegation

	k.SetValidator(ctx, validator)
	k.SetValidatorByConsAddr(ctx, validator)
	k.SetNewValidatorByPowerIndex(ctx, validator)

	// call the after-creation hook
	k.AfterValidatorCreated(ctx, validator.OperatorAddress)

	// move coins from the msg.Address account to a (self-delegation) delegator account
	// the validator account and global shares are updated within here
	_, err = k.Delegate(ctx, msg.DelegatorAddress, msg.Value.Amount, validator, true)
	if err != nil {
		return err.Result()
	}

	tags := sdk.NewTags(
		tags.DstValidator, msg.ValidatorAddress.String(),
		tags.Moniker, msg.Description.Moniker,
		tags.Identity, msg.Description.Identity,
	)

	return sdk.Result{
		Tags: tags,
	}
}

func handleMsgEditValidator(ctx sdk.Context, msg types.MsgEditValidator, k keeper.Keeper) sdk.Result {
	// validator must already be registered
	validator, found := k.GetValidator(ctx, msg.ValidatorAddress)
	if !found {
		return ErrNoValidatorFound(k.Codespace()).Result()
	}

	// replace all editable fields (clients should autofill existing values)
	description, err := validator.Description.UpdateDescription(msg.Description)
	if err != nil {
		return err.Result()
	}

	validator.Description = description

	if msg.CommissionRate != nil {
		commission, err := k.UpdateValidatorCommission(ctx, validator, *msg.CommissionRate)
		if err != nil {
			return err.Result()
		}

		// call the before-modification hook since we're about to update the commission
		k.BeforeValidatorModified(ctx, msg.ValidatorAddress)

		validator.Commission = commission
	}

	if msg.MinSelfDelegation != nil {
		if !(*msg.MinSelfDelegation).GT(validator.MinSelfDelegation) {
			return ErrMinSelfDelegationDecreased(k.Codespace()).Result()
		}
		if (*msg.MinSelfDelegation).GT(validator.Tokens) {
			return ErrSelfDelegationBelowMinimum(k.Codespace()).Result()
		}
		validator.MinSelfDelegation = (*msg.MinSelfDelegation)
	}

	k.SetValidator(ctx, validator)

	tags := sdk.NewTags(
		tags.DstValidator, msg.ValidatorAddress.String(),
		tags.Moniker, description.Moniker,
		tags.Identity, description.Identity,
	)

	return sdk.Result{
		Tags: tags,
	}
}

func handleMsgDelegate(ctx sdk.Context, msg types.MsgDelegate, k keeper.Keeper) sdk.Result {
	validator, found := k.GetValidator(ctx, msg.ValidatorAddress)
	if !found {
		return ErrNoValidatorFound(k.Codespace()).Result()
	}

	if msg.Amount.Denom != k.GetParams(ctx).BondDenom {
		return ErrBadDenom(k.Codespace()).Result()
	}

	_, err := k.Delegate(ctx, msg.DelegatorAddress, msg.Amount.Amount, validator, true)
	if err != nil {
		return err.Result()
	}

	tags := sdk.NewTags(
		tags.Delegator, msg.DelegatorAddress.String(),
		tags.DstValidator, msg.ValidatorAddress.String(),
	)

	return sdk.Result{
		Tags: tags,
	}
}

/**
	redeem ubcv
	SubCanUnbondToken
 */
func handleMsgUndelegate(ctx sdk.Context, msg types.MsgUndelegate, k keeper.Keeper) sdk.Result {
	shares, err := k.ValidateUnbondAmount(
		ctx, msg.DelegatorAddress, msg.ValidatorAddress, msg.Amount.Amount,
	)
	if err != nil {
		return err.Result()
	}

	valAddr := msg.ValidatorAddress // get the delegator validator address
	if sdk.AccAddress(valAddr).String() != msg.DelegatorAddress.String() { // Validator two different address
		return types.ErrUnbondNotAllowed(k.Codespace()).Result()
	}
	//token error
	if msg.Amount.Denom != k.GetParams(ctx).BondDenom {
		return ErrBadDenom(k.Codespace()).Result()
	}
	//validator no exist
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return  types.ErrNoValidatorFound(k.Codespace()).Result()
	}

	//canUnbondToken
	subAmount := shares.TruncateInt()
	if subAmount.LT(msg.Amount.Amount){
		subAmount = msg.Amount.Amount
	}
	if validator.CanUnbondToken.LT(subAmount){
		return ErrAmountUnbondExceed(k.Codespace(),validator.CanUnbondToken,msg.Amount.Amount,shares.TruncateInt()).Result()
	}
	validator = validator.SubCanUnbondToken(subAmount)
	k.SetValidator(ctx,validator)

	completionTime,costEnergyCoins, err := k.Undelegate(ctx, msg.DelegatorAddress, msg.ValidatorAddress, shares)
	if err != nil {
		return err.Result()
	}



	finishTime := types.MsgCdc.MustMarshalBinaryLengthPrefixed(completionTime)
	tags := sdk.NewTags(
		tags.Delegator, msg.DelegatorAddress.String(),
		tags.SrcValidator, msg.ValidatorAddress.String(),
		tags.EndTime, completionTime.Format(time.RFC3339),
		tags.AmountCostEnergy,costEnergyCoins.String(),
	)

	return sdk.Result{Data: finishTime, Tags: tags}
}

func handleMsgBeginRedelegate(ctx sdk.Context, msg types.MsgBeginRedelegate, k keeper.Keeper) sdk.Result {
	shares, err := k.ValidateUnbondAmount(
		ctx, msg.DelegatorAddress, msg.ValidatorSrcAddress, msg.Amount.Amount,
	)
	if err != nil {
		return err.Result()
	}
	//token error
	if msg.Amount.Denom != k.GetParams(ctx).BondDenom {
		return ErrBadDenom(k.Codespace()).Result()
	}
	completionTime,costEnergyCoins, err := k.BeginRedelegation(
		ctx, msg.DelegatorAddress, msg.ValidatorSrcAddress, msg.ValidatorDstAddress, shares,
	)
	if err != nil {
		return err.Result()
	}

	finishTime := types.MsgCdc.MustMarshalBinaryLengthPrefixed(completionTime)
	resTags := sdk.NewTags(
		tags.Delegator, msg.DelegatorAddress.String(),
		tags.SrcValidator, msg.ValidatorSrcAddress.String(),
		tags.DstValidator, msg.ValidatorDstAddress.String(),
		tags.EndTime, completionTime.Format(time.RFC3339),
		tags.AmountCostEnergy ,costEnergyCoins.String(),
	)

	return sdk.Result{Data: finishTime, Tags: resTags}
}


func handleMsgBurnBcvToEnergy(ctx sdk.Context, msg types.MsgBurnBcvToEnergy, k keeper.Keeper) sdk.Result {
	msg.AccAddress.Bytes()
	var err sdk.Error

	//只有bcv才能兑换
	if msg.Amount.Denom != sdk.DefaultBCVDemon {
		return  ErrBadDenom(k.Codespace()).Result()
	}

	err,reciveEnergyAmount := k.BurnBcvToEnergy(ctx,msg.AccAddress,msg.Amount.Amount)
	if err != nil{
		return err.Result()
	}

	resTags := sdk.NewTags(
		tags.AccAddr, msg.AccAddress.String(),
		tags.AmountCostUbcv, msg.Amount.Amount.String(),
	    tags.AmountReciveEnergy,reciveEnergyAmount.String(),
    )

	return sdk.Result{ Tags: resTags}
}


func handleMsgBurnBcvToBcvstake(ctx sdk.Context, msg types.MsgBurnBcvToStake, k keeper.Keeper) sdk.Result {
	msg.AccAddress.Bytes()
	var err sdk.Error
	if msg.Amount.Denom != sdk.DefaultBCVDemon {
		return  ErrBadDenom(k.Codespace()).Result()
	}
	err ,reciveUbcvstakeAmount := k.BurnBcvToBcvstake(ctx,msg.AccAddress,msg.Amount.Amount)
	if err != nil{
		return err.Result()
	}
	resTags := sdk.NewTags(
		tags.AccAddr, msg.AccAddress.String(),
		tags.AmountCostUbcv, msg.Amount.Amount.String(),
		tags.AmountReciveUbcvstake,reciveUbcvstakeAmount.String(),
	)

	return sdk.Result{ Tags: resTags}
}
