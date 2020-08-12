// nolint
package tags

import (
	sdk "github.com/bitcv-chain/bitcv-chain/types"
)

var (
	ActionCompleteUnbonding    = "complete_unbonding"
	ActionCompleteRedelegation = "complete_redelegation"

	Action       = sdk.TagAction
	SrcValidator = sdk.TagSrcValidator
	DstValidator = sdk.TagDstValidator
	Delegator    = sdk.TagDelegator
	Moniker      = "moniker"
	Identity     = "identity"
	EndTime      = "end_time"
	AccAddr		 = "account_addr"
	Amount		 = "amount"

	AmountCostUbcv    = "amount_cost_ubcv"
	AmountReciveEnergy  = "amount_recive_energy"
	AmountReciveUbcvstake  = "amount_recive_ubcvstake"
	AmountCostEnergy = "amount_cost_energy"

)
