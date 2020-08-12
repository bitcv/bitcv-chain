package app

import (
	"os"
	"testing"

	"github.com/bitcv-chain/bitcv-chain/x/bank"
	"github.com/bitcv-chain/bitcv-chain/x/crisis"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/bitcv-chain/bitcv-chain/codec"
	"github.com/bitcv-chain/bitcv-chain/x/auth"
	distr "github.com/bitcv-chain/bitcv-chain/x/distribution"
	"github.com/bitcv-chain/bitcv-chain/x/gov"
	"github.com/bitcv-chain/bitcv-chain/x/mint"
	"github.com/bitcv-chain/bitcv-chain/x/slashing"
	"github.com/bitcv-chain/bitcv-chain/x/staking"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/bitcv-chain/bitcv-chain/types"
	"github.com/tendermint/tendermint/crypto"
)

func setGenesis(gapp *BacApp, accs ...*auth.BaseAccount) error {
	genaccs := make([]GenesisAccount, len(accs))
	for i, acc := range accs {
		genaccs[i] = NewGenesisAccount(acc)
	}

	genesisState := NewGenesisState(
		genaccs,
		[]GenesisAccountEdatas{},
		auth.DefaultGenesisState(),
		bank.DefaultGenesisState(),
		staking.DefaultGenesisState(),
		mint.DefaultGenesisState(),
		distr.DefaultGenesisState(),
		gov.DefaultGenesisState(),
		crisis.DefaultGenesisState(),
		slashing.DefaultGenesisState(),
	)

	stateBytes, err := codec.MarshalJSONIndent(gapp.cdc, genesisState)
	if err != nil {
		return err
	}

	// Initialize the chain
	vals := []abci.ValidatorUpdate{}
	gapp.InitChain(abci.RequestInitChain{Validators: vals, AppStateBytes: stateBytes})
	gapp.Commit()

	return nil
}


func TestBacdExport(t *testing.T) {
	db := db.NewMemDB()
	gapp := NewBacApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)


	acc := &auth.BaseAccount{
		Address: types.AccAddress(crypto.AddressHash([]byte("test"))),
		Coins:   types.Coins{
					types.Coin{Denom:"nbac",Amount:types.MustNewIntFromString("29016000000000000")},
					types.Coin{Denom:"ubcv",Amount:types.MustNewIntFromString("1200000000000000")},

		},

	}
	setGenesis(gapp,acc)

	// Making a new app object with the db, so that initchain hasn't been called
	newGapp := NewBacApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)
	_, _, err := newGapp.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}





