package app

import (
	"os"
	"github.com/bitcv-chain/bitcv-chain/x/bank"
	"github.com/bitcv-chain/bitcv-chain/x/crisis"

	"github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/bitcv-chain/bitcv-chain/codec"
	"github.com/bitcv-chain/bitcv-chain/x/auth"
	distr "github.com/bitcv-chain/bitcv-chain/x/distribution"
	"github.com/bitcv-chain/bitcv-chain/x/gov"
	"github.com/bitcv-chain/bitcv-chain/x/mint"
	"github.com/bitcv-chain/bitcv-chain/x/slashing"
	"github.com/bitcv-chain/bitcv-chain/x/staking"

	abci "github.com/tendermint/tendermint/abci/types"
	"encoding/hex"
	"github.com/tendermint/tendermint/crypto"
	sdk "github.com/bitcv-chain/bitcv-chain/types"

	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/bitcv-chain/bitcv-chain/bacchain/v1_00"
)


var (
	pk1   = ed25519.GenPrivKey().PubKey()
	pk2   = ed25519.GenPrivKey().PubKey()
	pk3   = ed25519.GenPrivKey().PubKey()
	valAddr1 = sdk.ValAddress(pk1.Address())
	valAddr2 = sdk.ValAddress(pk2.Address())
	ValAddr3 = sdk.ValAddress(pk3.Address())

	emptyAddr   sdk.ValAddress
	emptyPubkey crypto.PubKey


	priv1 = secp256k1.GenPrivKey()
	addr1 = sdk.AccAddress(priv1.PubKey().Address())
	priv2 = secp256k1.GenPrivKey()
	addr2 = sdk.AccAddress(priv2.PubKey().Address())
	addr3 = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	priv4 = secp256k1.GenPrivKey()
	addr4 = sdk.AccAddress(priv4.PubKey().Address())

)

func setGenesis2(gapp *BacApp, accs ...*auth.BaseAccount) error {
	genaccs := make([]GenesisAccount, len(accs))
	for i, acc := range accs {
		genaccs[i] = NewGenesisAccount(acc)
	}

	stakingGenesis := staking.DefaultGenesisState()
	stakingGenesis.Pool = staking.Pool{BondedTokens:sdk.ZeroInt(),NotBondedTokens:sdk.GetChainParamBcvStakeAmount()}

	//set test param
	slashGenesisStatus := slashing.DefaultGenesisState()
	slashGenesisStatus.Params = slashing.DefaultTestParams()

	genesisState := NewGenesisState(
		genaccs,
		[]GenesisAccountEdatas{},
		auth.DefaultGenesisState(),
		bank.DefaultGenesisState(),
		stakingGenesis,
		mint.DefaultGenesisState(),
		distr.DefaultGenesisState(),
		gov.DefaultGenesisState(),
		crisis.DefaultGenesisState(),
		slashGenesisStatus,
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



/**
 整体测试
 */
var (

pks = []crypto.PubKey{
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB51"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB52"),
	}
	addrs = []sdk.ValAddress{
		sdk.ValAddress(pks[0].Address()),
		sdk.ValAddress(pks[1].Address()),
		sdk.ValAddress(pks[2].Address()),
	}
	initCoins = sdk.TokensFromTendermintPower(5000000)
)
func newPubKey(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	var pkEd ed25519.PubKeyEd25519
	copy(pkEd[:], pkBytes[:])
	return pkEd
}

func CreateTestInput(accs []auth.BaseAccount) (app *BacApp,ctx sdk.Context) {
	db := db.NewMemDB()
	app = NewBacApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)
	var accountList []*auth.BaseAccount

	//add energy
	accountEnergy := auth.NewBaseAccountWithAddress(sdk.AccAddrEnergyPool)
	accountEnergy.Coins = sdk.Coins{sdk.Coin{Denom: sdk.CHAIN_COIN_NAME_ENERGY, Amount: sdk.GetChainParamEnergyAmount()}}
	accountList = append(accountList, &accountEnergy)

	var account0 auth.BaseAccount
	account0 = auth.NewBaseAccountWithAddress(sdk.AccAddress(addrs[0]))
	account0.Coins = sdk.Coins{
		{sdk.DefaultBCVDemon, initCoins},
	}
	accountList = append(accountList,&account0)

	var account1 auth.BaseAccount
	account1 = auth.NewBaseAccountWithAddress(sdk.AccAddress(addrs[1]))
	account1.Coins = sdk.Coins{
		{sdk.DefaultBCVDemon, initCoins},
	}
	accountList = append(accountList,&account1)


	for k, _ := range accs {
		accountList = append(accountList, &accs[k])
	}

	//remaiNbac remainUbcv
	remainNbac := sdk.MustNewIntFromString(v1_00.StartParamBacAlreadyProduce)
	remainUbcv := sdk.GetChainParamBcvAmount()
	remainUbcvstake := sdk.GetChainParamBcvStakeAmount()
	for k, _ := range accountList {
		remainNbac = remainNbac.Sub(accountList[k].Coins.AmountOf(sdk.DefaultMintDenom))
		remainUbcv = remainUbcv.Sub(accountList[k].Coins.AmountOf(sdk.DefaultBCVDemon))
		remainUbcvstake = remainUbcvstake.Sub(accountList[k].Coins.AmountOf(sdk.DefaultBondDenom))
	}

	accountRemain := auth.NewBaseAccountWithAddress(sdk.AccAddress(addrs[2]))
	accountRemain.Coins = sdk.Coins{
		sdk.Coin{Denom: sdk.DEFAULT_FEE_COIN, Amount: remainNbac},
		sdk.Coin{Denom: sdk.DefaultBCVDemon, Amount: remainUbcv}}
	accountList = append(accountList, &accountRemain)

	//remainUbcvstake
	var accountPool= auth.NewBaseAccountWithAddress(sdk.AccAddrBcvstakePool)
	accountPool.Coins = sdk.Coins{sdk.Coin{Denom: sdk.CHAIN_COIN_NAME_BCVSTAKE, Amount: remainUbcvstake}}
	accountList = append(accountList, &accountPool)

	ctx = app.NewContext(true,abci.Header{Height: 1})
	setGenesis2(app,accountList ...)
	return app,ctx
}

func burnBcvToGetEnergy(ctx sdk.Context,app *BacApp,addr sdk.AccAddress, ubcvAmount sdk.Int) sdk.Result {
	sh := staking.NewHandler(app.stakingKeeper)
	got := sh(ctx,staking.NewMsgBurnBcvToEnergy(addr,sdk.Coin{Denom:"ubcv",Amount:ubcvAmount}))
	return  got
}
