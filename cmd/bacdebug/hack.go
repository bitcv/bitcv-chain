package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path"

	"github.com/bitcv-chain/bitcv-chain/store"

	"github.com/bitcv-chain/bitcv-chain/baseapp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"

	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	bam "github.com/bitcv-chain/bitcv-chain/baseapp"
	sdk "github.com/bitcv-chain/bitcv-chain/types"

	"github.com/bitcv-chain/bitcv-chain/codec"
	"github.com/bitcv-chain/bitcv-chain/x/auth"
	"github.com/bitcv-chain/bitcv-chain/x/bank"
	"github.com/bitcv-chain/bitcv-chain/x/params"
	"github.com/bitcv-chain/bitcv-chain/x/slashing"
	"github.com/bitcv-chain/bitcv-chain/x/staking"

	bac "github.com/bitcv-chain/bitcv-chain/app"
	"github.com/bitcv-chain/bitcv-chain/x/mint"
	distr"github.com/bitcv-chain/bitcv-chain/x/distribution"
	distrtypes"github.com/bitcv-chain/bitcv-chain/x/distribution/types"

	"github.com/bitcv-chain/bitcv-chain/x/gov"
	"github.com/bitcv-chain/bitcv-chain/x/crisis"
	"github.com/bitcv-chain/bitcv-chain/x/staking/keeper"
	"github.com/tendermint/tendermint/blockchain"
	"github.com/tendermint/tendermint/types"
)

//./bacdebug hack /Users/liuhaoyang/.bacd



func runHackCmd(cmd *cobra.Command, args []string) error {

	if len(args) != 1 {
		return fmt.Errorf("Expected 1 arg")
	}

	dataDir := args[0]
	dataDir = path.Join(dataDir, "data")

	// load the app
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	db, err := sdk.NewLevelDB("application", dataDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	app := NewBacApp(logger, db,true, baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))))

	// print some info
	id := app.LastCommitID()
	lastBlockHeight := app.LastBlockHeight()
	fmt.Println("ID", id)
	fmt.Println("LastBlockHeight", lastBlockHeight)

	//----------------------------------------------------
	// XXX: start hacking!
	//----------------------------------------------------
	// eg. bac-6001 testnet bug
	// We paniced when iterating through the "bypower" keys.
	// The following powerKey was there, but the corresponding "trouble" validator did not exist.
	// So here we do a binary search on the past states to find when the powerKey first showed up ...

	// operator of the validator the bonds, gets jailed, later unbonds, and then later is still found in the bypower store
	trouble := hexToBytes("D3DC0FF59F7C3B548B7AFA365561B87FD0208AF8")
	// this is his "bypower" key

	powerKey := hexToBytes("05303030303030303030303033FFFFFFFFFFFF4C0C0000FFFED3DC0FF59F7C3B548B7AFA365561B87FD0208AF8")


	fmt.Println("powerKey",powerKey)
	topHeight := lastBlockHeight
	bottomHeight := int64(0)
	checkHeight := topHeight

	fmt.Println(trouble,powerKey,bottomHeight,checkHeight)
	for {

		//// load the given version of the state
		app := NewBacApp(logger, db,false, baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))))
		err = app.LoadVersion(checkHeight, app.keyMain)
		checkHeight --
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		ctx := app.NewContext(true, abci.Header{})

		//// check for the powerkey and the validator from the store
		//validatorList  := app.stakingKeeper.GetAllValidators(ctx)
		iterator := app.stakingKeeper.ValidatorsPowerStoreIterator(ctx)
		for ;iterator.Valid();iterator.Next(){
			validator ,isFound := app.stakingKeeper.GetValidator(ctx, sdk.ValAddress(iterator.Value()))
			if !isFound{
				continue
			}
			b := keeper.GetValidatorsByPowerIndexKey(validator)
			fmt.Println(hex.EncodeToString(b))

		}
		iterator.Close()



		store := ctx.KVStore(app.keyStaking)
		iterator2 := store.Iterator([]byte{0x00},[]byte{0x50})
		for ; iterator2.Valid();iterator2.Next(){
			val := iterator2.Value()
			fmt.Println("##key:",hex.EncodeToString(iterator2.Key()),"##val:",hex.EncodeToString(val))
		}
		iterator2.Close()

		//validator_historical_rewards
		fmt.Println("begin_print_validator_historical_rewards.............")
		PrintHistoricalRewards(ctx,app)
		fmt.Println("end_print_validator_historical_rewards.............")



		//GetRedelegationsFromValidator
		fmt.Println("begin_print_redelegation.............")
		stakingKeeper  :=  app.stakingKeeper
		operatorAddress,_:= sdk.ValAddressFromBech32("bacvaloper10tyhju9pfpfkt7hrd2zqr0vjn4k5sfrrl5kav9")
		val  := stakingKeeper.GetRedelegationsFromValidator(ctx, operatorAddress)
		fmt.Println("staking keeper:",val)
		fmt.Println("end_print_redelegation.............")
		break;


	}


	return nil
}


//遍历block
func runBlockCmd(cmd *cobra.Command, args []string) error {

	if len(args) != 1 {
		return fmt.Errorf("Expected 1 arg")
	}

	dataDir := args[0]
	dataDir = path.Join(dataDir, "data")

	// load the app
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	applicationDb, err := sdk.NewLevelDB("application", dataDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	blockstoreDb, err := sdk.NewLevelDB("blockstore", dataDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app := NewBacApp(logger, applicationDb,true, baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))))

	// print some info
	id := app.LastCommitID()
	lastBlockHeight := app.LastBlockHeight()

	fmt.Println("ID", id)
	fmt.Println("LastBlockHeight", lastBlockHeight)
	bottomHeight := int64(0)
	for {

		//// load the given version of the state
		app := NewBacApp(logger, applicationDb,false, baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))))
		err = app.LoadVersion(bottomHeight, app.keyMain)
		bottomHeight ++
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if(bottomHeight > lastBlockHeight){
			os.Exit(1)
		}

		//pool
		//resp = app.Query(abci.RequestQuery{
		//	Path: "/custom/staking/pool",
		//	Data: []byte{},
		//})

		//blok
		//resp = app.Query(abci.RequestQuery{
		//	Path: "/store/block/2",
		//	Data: []byte{},
		//})
		//viper.Set(client.FlagTrustNode, true)
		//cliContext := context.NewCLIContext()
		//cliContext.Height = 4
		//node,_:= cliContext.GetNode()

		//var a  *int64
		//var  b = int64(2)
		//a = &b
		//c ,_:= node.Block(a)
		//fmt.Println(c)

		blockStore := blockchain.NewBlockStore(blockstoreDb)
		b := blockStore.LoadBlock(bottomHeight)
		fmt.Println("###b.size",b.Size())

		var tx interface{}
		txs :=   b.Data.Txs
		for tx,_   = range txs {
			a,ok:= tx.(types.Tx)
			fmt.Println(a)
			fmt.Println(ok)
		}
		fmt.Println("###b.size",b.NumTxs)
	}


	return nil
}

func base64ToPub(b64 string) ed25519.PubKeyEd25519 {
	data, _ := base64.StdEncoding.DecodeString(b64)
	var pubKey ed25519.PubKeyEd25519
	copy(pubKey[:], data)
	return pubKey

}

func hexToBytes(h string) []byte {
	trouble, _ := hex.DecodeString(h)
	return trouble

}

//--------------------------------------------------------------------------------
// NOTE: This is all copied from bac/app/app.go
// so we can access internal fields!

const (
	appName = "bacchain"
)

// default home directories for expected binaries
var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.baccli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.bacd")
)

// Extended ABCI application
type BacApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// keys to access the substores
	keyMain          *sdk.KVStoreKey
	keyAccount       *sdk.KVStoreKey
	keyStaking       *sdk.KVStoreKey
	tkeyStaking      *sdk.TransientStoreKey
	keySlashing      *sdk.KVStoreKey
	keyMint          *sdk.KVStoreKey
	keyDistr         *sdk.KVStoreKey
	tkeyDistr        *sdk.TransientStoreKey
	keyGov           *sdk.KVStoreKey
	keyFeeCollection *sdk.KVStoreKey
	keyParams        *sdk.KVStoreKey
	tkeyParams       *sdk.TransientStoreKey

	// Manage getting and setting accounts
	accountKeeper       auth.AccountKeeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	bankKeeper          bank.Keeper
	stakingKeeper       staking.Keeper
	slashingKeeper      slashing.Keeper
	mintKeeper          mint.Keeper
	distrKeeper         distr.Keeper
	govKeeper           gov.Keeper
	crisisKeeper        crisis.Keeper
	paramsKeeper        params.Keeper
}

func NewBacApp(logger log.Logger, db dbm.DB, loadLatest bool,baseAppOptions ...func(*bam.BaseApp)) *BacApp {

	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	invCheckPeriod := uint(100)
	var app = &BacApp{
		BaseApp:          bApp,
		cdc:              cdc,
		invCheckPeriod:   invCheckPeriod,
		keyMain:          sdk.NewKVStoreKey(bam.MainStoreKey),
		keyAccount:       sdk.NewKVStoreKey(auth.StoreKey),
		keyStaking:       sdk.NewKVStoreKey(staking.StoreKey),
		tkeyStaking:      sdk.NewTransientStoreKey(staking.TStoreKey),
		keyMint:          sdk.NewKVStoreKey(mint.StoreKey),
		keyDistr:         sdk.NewKVStoreKey(distr.StoreKey),
		tkeyDistr:        sdk.NewTransientStoreKey(distr.TStoreKey),
		keySlashing:      sdk.NewKVStoreKey(slashing.StoreKey),
		keyGov:           sdk.NewKVStoreKey(gov.StoreKey),
		keyFeeCollection: sdk.NewKVStoreKey(auth.FeeStoreKey),
		keyParams:        sdk.NewKVStoreKey(params.StoreKey),
		tkeyParams:       sdk.NewTransientStoreKey(params.TStoreKey),
	}

	app.paramsKeeper = params.NewKeeper(app.cdc, app.keyParams, app.tkeyParams)

	// define the accountKeeper
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		app.keyAccount,
		app.paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)

	// add handlers
	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		app.paramsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(
		app.cdc,
		app.keyFeeCollection,
	)
	stakingKeeper := staking.NewKeeper(
		app.cdc,
		app.keyStaking, app.tkeyStaking,
		app.bankKeeper, app.paramsKeeper.Subspace(staking.DefaultParamspace),
		staking.DefaultCodespace,
	)
	app.mintKeeper = mint.NewKeeper(app.cdc, app.keyMint,
		app.paramsKeeper.Subspace(mint.DefaultParamspace),
		&stakingKeeper, app.feeCollectionKeeper,
	)
	app.distrKeeper = distr.NewKeeper(
		app.cdc,
		app.keyDistr,
		app.paramsKeeper.Subspace(distr.DefaultParamspace),
		app.bankKeeper, &stakingKeeper, app.feeCollectionKeeper,
		distr.DefaultCodespace,
	)

	app.slashingKeeper = slashing.NewKeeper(
		app.cdc,
		app.keySlashing,
		&stakingKeeper, app.paramsKeeper.Subspace(slashing.DefaultParamspace),
		slashing.DefaultCodespace,
	)
	app.govKeeper = gov.NewKeeper(
		app.cdc,
		app.keyGov,
		app.paramsKeeper, app.paramsKeeper.Subspace(gov.DefaultParamspace), app.bankKeeper, &stakingKeeper,
		gov.DefaultCodespace,
	)
	app.crisisKeeper = crisis.NewKeeper(
		app.paramsKeeper.Subspace(crisis.DefaultParamspace),
		app.distrKeeper,
		app.bankKeeper,
		app.feeCollectionKeeper,
	)

	// register the staking hooks
	// NOTE: The stakingKeeper above is passed by reference, so that it can be
	// modified like below:
	app.stakingKeeper = *stakingKeeper.SetHooks(
		NewStakingHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()),
	)
	/*stakingKeeper*/
	app.stakingKeeper = *stakingKeeper.SetDistrKeeper(app.distrKeeper)
	/***/
	// register the crisis routes
	bank.RegisterInvariants(&app.crisisKeeper, app.accountKeeper)
	distr.RegisterInvariants(&app.crisisKeeper, app.distrKeeper, app.stakingKeeper)
	staking.RegisterInvariants(&app.crisisKeeper, app.stakingKeeper, app.feeCollectionKeeper, app.distrKeeper, app.accountKeeper)

	// register message routes
	app.Router().
		AddRoute(bank.RouterKey, bank.NewHandler(app.bankKeeper)).
		AddRoute(staking.RouterKey, staking.NewHandler(app.stakingKeeper)).
		AddRoute(distr.RouterKey, distr.NewHandler(app.distrKeeper)).
		AddRoute(slashing.RouterKey, slashing.NewHandler(app.slashingKeeper)).
		AddRoute(gov.RouterKey, gov.NewHandler(app.govKeeper)).
		AddRoute(crisis.RouterKey, crisis.NewHandler(app.crisisKeeper))

	app.QueryRouter().
		AddRoute(auth.QuerierRoute, auth.NewQuerier(app.accountKeeper)).
		AddRoute(distr.QuerierRoute, distr.NewQuerier(app.distrKeeper)).
		AddRoute(gov.QuerierRoute, gov.NewQuerier(app.govKeeper)).
		AddRoute(slashing.QuerierRoute, slashing.NewQuerier(app.slashingKeeper, app.cdc)).
		AddRoute(staking.QuerierRoute, staking.NewQuerier(app.stakingKeeper, app.cdc)).
		AddRoute(mint.QuerierRoute, mint.NewQuerier(app.mintKeeper))

	// initialize BaseApp
	app.MountStores(app.keyMain, app.keyAccount, app.keyStaking, app.keyMint, app.keyDistr,
		app.keySlashing, app.keyGov, app.keyFeeCollection, app.keyParams,
		app.tkeyParams, app.tkeyStaking, app.tkeyDistr,
	)
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.feeCollectionKeeper))
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		err := app.LoadLatestVersion(app.keyMain)
		if err != nil {
			cmn.Exit(err.Error())
		}
	}

	return app
}

// custom tx codec
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	slashing.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	cdc.Seal()
	return cdc
}

// application updates every end block
func (app *BacApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	tags := slashing.BeginBlocker(ctx, req, app.slashingKeeper)

	return abci.ResponseBeginBlock{
		Tags: tags.ToKVPairs(),
	}
}

// application updates every end block
// nolint: unparam
func (app *BacApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	validatorUpdates, tags := staking.EndBlocker(ctx, app.stakingKeeper)

	return abci.ResponseEndBlock{
		ValidatorUpdates: validatorUpdates,
		Tags:             tags,
	}
}

// custom logic for bac initialization
func (app *BacApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes
	// TODO is this now the whole genesis file?

	var genesisState bac.GenesisState
	err := app.cdc.UnmarshalJSON(stateJSON, &genesisState)
	if err != nil {
		panic(err) // TODO https://github.com/bitcv-chain/bitcv-chain/issues/468 // return sdk.ErrGenesisParse("").TraceCause(err, "")
	}

	// load the accounts
	for _, gacc := range genesisState.Accounts {
		acc := gacc.ToAccount()
		app.accountKeeper.SetAccount(ctx, acc)
	}

	// load the initial staking information
	validators, err := staking.InitGenesis(ctx, app.stakingKeeper, genesisState.StakingData)
	if err != nil {
		panic(err) // TODO https://github.com/bitcv-chain/bitcv-chain/issues/468 // return sdk.ErrGenesisParse("").TraceCause(err, "")
	}

	slashing.InitGenesis(ctx, app.slashingKeeper, genesisState.SlashingData, genesisState.StakingData.Validators.ToSDKValidators())

	return abci.ResponseInitChain{
		Validators: validators,
	}
}



var _ sdk.StakingHooks = StakingHooks{}

// StakingHooks contains combined distribution and slashing hooks needed for the
// staking module.
type StakingHooks struct {
	dh distr.Hooks
	sh slashing.Hooks
}

func NewStakingHooks(dh distr.Hooks, sh slashing.Hooks) StakingHooks {
	return StakingHooks{dh, sh}
}

// nolint
func (h StakingHooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
	h.dh.AfterValidatorCreated(ctx, valAddr)
	h.sh.AfterValidatorCreated(ctx, valAddr)
}
func (h StakingHooks) BeforeValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress) {
	h.dh.BeforeValidatorModified(ctx, valAddr)
	h.sh.BeforeValidatorModified(ctx, valAddr)
}
func (h StakingHooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.dh.AfterValidatorRemoved(ctx, consAddr, valAddr)
	h.sh.AfterValidatorRemoved(ctx, consAddr, valAddr)
}
func (h StakingHooks) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.dh.AfterValidatorBonded(ctx, consAddr, valAddr)
	h.sh.AfterValidatorBonded(ctx, consAddr, valAddr)
}
func (h StakingHooks) AfterValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.dh.AfterValidatorBeginUnbonding(ctx, consAddr, valAddr)
	h.sh.AfterValidatorBeginUnbonding(ctx, consAddr, valAddr)
}
func (h StakingHooks) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.dh.BeforeDelegationCreated(ctx, delAddr, valAddr)
	h.sh.BeforeDelegationCreated(ctx, delAddr, valAddr)
}
func (h StakingHooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.dh.BeforeDelegationSharesModified(ctx, delAddr, valAddr)
	h.sh.BeforeDelegationSharesModified(ctx, delAddr, valAddr)
}
func (h StakingHooks) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.dh.BeforeDelegationRemoved(ctx, delAddr, valAddr)
	h.sh.BeforeDelegationRemoved(ctx, delAddr, valAddr)
}
func (h StakingHooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.dh.AfterDelegationModified(ctx, delAddr, valAddr)
	h.sh.AfterDelegationModified(ctx, delAddr, valAddr)
}
func (h StakingHooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
	h.dh.BeforeValidatorSlashed(ctx, valAddr, fraction)
	h.sh.BeforeValidatorSlashed(ctx, valAddr, fraction)
}




func PrintHistoricalRewards(ctx sdk.Context,app *BacApp)  {
	app.distrKeeper.IterateValidatorHistoricalRewards(ctx, func(val sdk.ValAddress, period uint64, rewards distrtypes.ValidatorHistoricalRewards) (stop bool) {
		fmt.Println("period:",period,"rewards",rewards);
		return false;
	})
}