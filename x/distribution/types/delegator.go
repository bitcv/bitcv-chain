package types

import (
	sdk "github.com/bitcv-chain/bitcv-chain/types"
	"encoding/json"
)

// starting info for a delegator reward period
// tracks the previous validator period, the delegation's amount
// of staking token, and the creation height (to check later on
// if any slashes have occurred)
// NOTE that even though validators are slashed to whole staking tokens, the
// delegators within the validator may be left with less than a full token,
// thus sdk.Dec is used
type DelegatorStartingInfo struct {
	PreviousPeriod uint64  `json:"previous_period"` // period at which the delegation should withdraw starting from
	Stake          sdk.Dec `json:"stake"`           // amount of staking token delegated
	Height         uint64  `json:"height"`          // height at which delegation was created
}

// create a new DelegatorStartingInfo
func NewDelegatorStartingInfo(previousPeriod uint64, stake sdk.Dec, height uint64) DelegatorStartingInfo {
	return DelegatorStartingInfo{
		PreviousPeriod: previousPeriod,
		Stake:          stake,
		Height:         height,
	}
}


/*
 * 查看奖励能量消耗
 */
type RewardsDetailInfo struct{
	OperAddr	string					`json:"oper_addr"`
	StartingInfo  DelegatorStartingInfo `json:"starting_info"`
	Rewards	sdk.DecCoins				`json:"rewards"`
	CostEnergy sdk.Int					`json:"cost_energy"`
	EndHeight int64						`json:"end_height"`
}


func (r RewardsDetailInfo) String() string {
	data,_ := json.Marshal(r)
	return string(data)
}


/*
 * 查看奖励能量消耗
 */
type TotalRewardsDetailInfo struct{
	RewardList []RewardsDetailInfo `json:"reward_detail_list"`
	TotalRewards sdk.DecCoins `json:"total_rewards"`
}


func (r TotalRewardsDetailInfo) String() string {
	data,_ := json.Marshal(r)
	return string(data)
}