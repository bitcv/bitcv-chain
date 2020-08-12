// nolint
package tags

import (
	sdk "github.com/bitcv-chain/bitcv-chain/types"
)

// Distribution tx tags
var (
	Rewards    = "rewards"
	Commission = "commission"

	Validator = sdk.TagSrcValidator
	Delegator = sdk.TagDelegator

	WithdrawDstAddress = sdk.TagWithdrawDstAddress

	AmountCostEnergy = "amount_cost_energy"
)
