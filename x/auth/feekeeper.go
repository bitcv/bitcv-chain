package auth

import (
	"github.com/bitcv-chain/bitcv-chain/codec"
	sdk "github.com/bitcv-chain/bitcv-chain/types"
	"fmt"
	"github.com/bitcv-chain/bitcv-chain/bacchain/v1_00"
)

var (
	collectedFeesKey = []byte("collectedFees")
)
// FeeCollectionKeeper handles collection of fees in the anteHandler
// and setting of MinFees for different fee tokens
type FeeCollectionKeeper struct {

	// The (unexposed) key used to access the fee store from the Context.
	key sdk.StoreKey

	// The codec codec for binary encoding/decoding of accounts.
	cdc *codec.Codec
}

// NewFeeCollectionKeeper returns a new FeeCollectionKeeper
func NewFeeCollectionKeeper(cdc *codec.Codec, key sdk.StoreKey) FeeCollectionKeeper {
	return FeeCollectionKeeper{
		key: key,
		cdc: cdc,
	}
}

// GetCollectedFees - retrieves the collected fee pool
func (fck FeeCollectionKeeper) GetCollectedFees(ctx sdk.Context) sdk.Coins {
	store := ctx.KVStore(fck.key)
	bz := store.Get(collectedFeesKey)
	if bz == nil {
		return sdk.NewCoins()
	}

	emptyFees := sdk.NewCoins()
	feePool := &emptyFees
	fck.cdc.MustUnmarshalBinaryLengthPrefixed(bz, feePool)
	return *feePool
}

func (fck FeeCollectionKeeper) setCollectedFees(ctx sdk.Context, coins sdk.Coins) {
	bz := fck.cdc.MustMarshalBinaryLengthPrefixed(coins)
	store := ctx.KVStore(fck.key)
	store.Set(collectedFeesKey, bz)
}

// AddCollectedFees - add to the fee pool
func (fck FeeCollectionKeeper) AddCollectedFees(ctx sdk.Context, coins sdk.Coins) sdk.Coins {
	newCoins := fck.GetCollectedFees(ctx).Add(coins)
	fck.setCollectedFees(ctx, newCoins)

	return newCoins
}

// ClearCollectedFees - clear the fee pool
func (fck FeeCollectionKeeper) ClearCollectedFees(ctx sdk.Context) {
	fck.setCollectedFees(ctx, sdk.NewCoins())
}

/////////////////New Coin Manage/////////////////////////////////////
var (
	bacManagePoolKey = []byte("bacmanagepool") // TODO
	DefaultGenerateCoin = sdk.DEFAULT_FEE_COIN
)

// Record generate new coin
type BacManagePool struct {
	TotalGenerate sdk.Coins `json:"total_generate"`
	AlreadyBurn sdk.Coins   `json:"already_burn"`
}
// initial pool for testing
func InitialBacManagePool() BacManagePool {
	totalGenerate ,_ := sdk.NewIntFromString(v1_00.StartParamBacAlreadyProduce)
	return BacManagePool{
		TotalGenerate: sdk.Coins{sdk.NewCoin(DefaultGenerateCoin,totalGenerate)},
		AlreadyBurn:   sdk.Coins{sdk.NewCoin(DefaultGenerateCoin,sdk.ZeroInt())},
	}
}
// Generate total BAC supply
func (p BacManagePool) GenerateTokenSupply() sdk.Coins {
	return p.TotalGenerate
}
// BAC alreay burn
func (p BacManagePool) AlreadyBurnSupply() sdk.Coins {
	return p.AlreadyBurn
}

func (p BacManagePool) String() string {
	return fmt.Sprintf("supply %s;burn %s",p.TotalGenerate.String(),p.AlreadyBurn.String())
}

func (fck FeeCollectionKeeper) GetBacManagePool(ctx sdk.Context) (coinNewPool BacManagePool){
	store := ctx.KVStore(fck.key)
	b := store.Get(bacManagePoolKey)
	if b == nil {
		panic("Stored coinNewPool should not have been nil")
	}
	fck.cdc.MustUnmarshalBinaryLengthPrefixed(b, &coinNewPool)
	return
}

func (fck FeeCollectionKeeper) SetBacManagePool(ctx sdk.Context, coinNewPool BacManagePool) {
	store := ctx.KVStore(fck.key)
	b := fck.cdc.MustMarshalBinaryLengthPrefixed(coinNewPool)
	store.Set(bacManagePoolKey, b)
}

// when minting new tokens
func (fck FeeCollectionKeeper) NewCoinInflateSupply(ctx sdk.Context, newTokens sdk.Coins)  {
	p := fck.GetBacManagePool(ctx)
	p.TotalGenerate = p.TotalGenerate.Add(newTokens)
	fck.SetBacManagePool(ctx, p)
}
// Destroy transfer fee
func (fck FeeCollectionKeeper) BurnSupply(ctx sdk.Context, tokens sdk.Coins) {
	p := fck.GetBacManagePool(ctx)
	p.AlreadyBurn = p.AlreadyBurn.Add(tokens)
	fck.SetBacManagePool(ctx, p)
}
// get current bac total circulation
func (fck FeeCollectionKeeper) GetNewTokenCirculation(ctx sdk.Context) sdk.Coins {
	p := fck.GetBacManagePool(ctx)
	// total generate token - alreay burn token
	return p.TotalGenerate.Sub(p.AlreadyBurn)
}
/////////////////TODO/////////////////////////////////////

type BacManagePoolForApi struct {
	BacManagePool BacManagePool `json:"bac_manage_pool"`
	Height int64 `json:"height"`
}

func (fck FeeCollectionKeeper) GetBacManagePoolForApi(ctx sdk.Context) ( BacManagePoolForApi){
	store := ctx.KVStore(fck.key)
	b := store.Get(bacManagePoolKey)
	if b == nil {
		panic("Stored coinNewPool should not have been nil")
	}
	var bacManagePool BacManagePool
	fck.cdc.MustUnmarshalBinaryLengthPrefixed(b, &bacManagePool)

	bacManagePoolForApi := BacManagePoolForApi{
		BacManagePool:bacManagePool,
		Height:ctx.BlockHeight(),
	}

	return bacManagePoolForApi
}