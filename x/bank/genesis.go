package bank

import (
	sdk "github.com/bitcv-chain/bitcv-chain/types"
	bacv1"github.com/bitcv-chain/bitcv-chain/bacchain/v1_00"
)

// GenesisState is the bank state that must be provided at genesis.
type GenesisState struct {
	SendEnabled bool `json:"send_enabled"`
	Tokens []IssueToken `json:"tokens"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(sendEnabled bool,tokens []IssueToken) GenesisState {
	return GenesisState{SendEnabled: sendEnabled,Tokens:tokens}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState { return GenesisState{SendEnabled:true,Tokens:nil} }

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetSendEnabled(ctx, data.SendEnabled)
	if data.Tokens != nil {
		for _, token := range data.Tokens {
			keeper.SetTokenByInnerName(ctx,token.InnerName,token)
		}
	}

}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {

	return NewGenesisState(
		bacv1.CheckSendEnable(keeper.GetSendEnabled(ctx),ctx.BlockHeight()),
		keeper.GetTokens(ctx))
}

// ValidateGenesis performs basic validation of bank genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error { return nil }
