package tags

import (
	sdk "github.com/bitcv-chain/bitcv-chain/types"
)

// Governance tags
var (
	ActionProposalDropped  = "proposal_dropped"
	ActionProposalPassed   = "proposal_passed"
	ActionProposalRejected = "proposal_rejected"

	Action            = sdk.TagAction
	Proposer          = "proposer"
	ProposalID        = "proposal_id"
	VotingPeriodStart = "voting_period_start"
	Depositor         = "depositor"
	Voter             = "voter"
	ProposalResult    = "proposal_result"
)
