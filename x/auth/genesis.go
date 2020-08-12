package auth

import (
	"fmt"
	sdk "github.com/bitcv-chain/bitcv-chain/types"
)

// GenesisState - all auth state that must be provided at genesis
type GenesisState struct {
	BacManagePool   BacManagePool `json:"bac_manage_pool"`
	CollectedFees sdk.Coins `json:"collected_fees"`
	Params        Params    `json:"params"`
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(coinNewPool BacManagePool,collectedFees sdk.Coins, params Params) GenesisState {
	return GenesisState{
		BacManagePool:  coinNewPool,
		Params:        params,
		CollectedFees: collectedFees,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(InitialBacManagePool(), sdk.NewCoins(), DefaultParams())
}

// InitGenesis - Init store state from genesis data
func InitGenesis(ctx sdk.Context, ak AccountKeeper, fck FeeCollectionKeeper, data GenesisState) {
	ak.SetParams(ctx, data.Params)
	fck.setCollectedFees(ctx, data.CollectedFees)
	fck.SetBacManagePool(ctx, data.BacManagePool)
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, ak AccountKeeper, fck FeeCollectionKeeper) GenesisState {
	collectedFees := fck.GetCollectedFees(ctx)
	params := ak.GetParams(ctx)
	bacManagePool := fck.GetBacManagePool(ctx)
	return NewGenesisState(bacManagePool, collectedFees, params)
}

// ValidateGenesis performs basic validation of auth genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	if data.Params.TxSigLimit == 0 {
		return fmt.Errorf("invalid tx signature limit: %d", data.Params.TxSigLimit)
	}
	if data.Params.SigVerifyCostED25519 == 0 {
		return fmt.Errorf("invalid ED25519 signature verification cost: %d", data.Params.SigVerifyCostED25519)
	}
	if data.Params.SigVerifyCostSecp256k1 == 0 {
		return fmt.Errorf("invalid SECK256k1 signature verification cost: %d", data.Params.SigVerifyCostSecp256k1)
	}
	if data.Params.MaxMemoCharacters == 0 {
		return fmt.Errorf("invalid max memo characters: %d", data.Params.MaxMemoCharacters)
	}
	if data.Params.TxSizeCostPerByte == 0 {
		return fmt.Errorf("invalid tx size cost per byte: %d", data.Params.TxSizeCostPerByte)
	}
	return nil
}
