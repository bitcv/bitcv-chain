package bank
import (
	"github.com/tendermint/tendermint/crypto/tmhash"
	sdk "github.com/bitcv-chain/bitcv-chain/types"
	"github.com/bitcv-chain/bitcv-chain/codec"
	"encoding/hex"
	"github.com/bitcv-chain/bitcv-chain/version"
	"strconv"
)
var _ Keeper = (*BaseKeeper)(nil)


type Edata struct{
	Utype     uint8  `json:"utype"`
	Data      string	 `json:"data"`
	Location  string			 `json:"location"`
}

type EdataKeeper struct{
	key           sdk.StoreKey
    cdc          *codec.Codec
	codespace sdk.CodespaceType
}



/**
*	书记存储keeper
 */
func NewEdataKeeper(key sdk.StoreKey,cdc *codec.Codec,	 codespace sdk.CodespaceType) EdataKeeper {
	return EdataKeeper{
		key:	    key,
		cdc:        cdc,
		codespace:  codespace,
	}
}

type  AccountEdata struct{
	account sdk.AccAddress
	edatas  []Edata
}


func (k EdataKeeper) SetEdataByAccount(ctx sdk.Context,account sdk.AccAddress,edata Edata,bk BaseKeeper) (sdk.Tags,sdk.Error){
	//组装location
	bigVersion := version.GetBigVersion(version.Version)
	height := ctx.BlockHeader().Height
	suffix := hex.EncodeToString(tmhash.Sum(ctx.TxBytes()))[0:4]
	edata.Location =  bigVersion+"." + strconv.FormatInt(height,10) + "." + suffix

	//销毁bac
	costCoin := sdk.Coins{sdk.Coin{Denom:sdk.DefaultMintDenom,Amount:sdk.NewInt(int64(len(edata.Data)  * 20000000)) }}
	tmpTags ,err := bk.SendCoins(ctx,account,sdk.AccAddrBurnFromEdataSave,costCoin)
	if err != nil{
		return  tmpTags,err
	}
	//存储数据
	found,edatas := k.GetEdatasByAccount(ctx,account)
	if found{
		edatas = append(edatas,edata )
	}else{
		edatas = []Edata{edata}
	}
	store := ctx.KVStore(k.key)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(edatas)
	store.Set(k.GetEdatasKeyByaccount(ctx,account), bz)

	tags := sdk.Tags{}
	tags = tags.AppendTag(TagEdataCost,costCoin.String())
	tags = tags.AppendTag(TagEdataLocation,edata.Location)

	return tags,nil
}


func (k EdataKeeper) GetEdatasByAccount(ctx sdk.Context, account sdk.AccAddress) ( found bool,edatas []Edata,) {
	store := ctx.KVStore(k.key)
	key := k.GetEdatasKeyByaccount(ctx,account)
	value := store.Get(key)
	if value == nil {
		return  false,edatas
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(value,&edatas)
	return  true,edatas
}

func (k EdataKeeper) GetEdatasKeyByaccount(ctx sdk.Context,account sdk.AccAddress)  []byte {
	return append(EdataStoreKeyPrefix, account...)
}

