package bank
import (

	sdk "github.com/bitcv-chain/bitcv-chain/types"
	"github.com/bitcv-chain/bitcv-chain/codec"
	"encoding/hex"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/crypto"
	"time"
	"math/big"
)
var _ Keeper = (*BaseKeeper)(nil)


type IssueToken struct{
	OwnerAddress   	    sdk.AccAddress   `json:"owner_address"`
	ExchangeAddress     sdk.AccAddress	 `json:"exchange_address"`
	OuterName 			string			 `json:"outer_name"`
	InnerName           string			 `json:"inner_name"`
	SupplyNum           sdk.Int			 `json:"supply_num"` 	//发行总量
	Margin              sdk.Coin          `json:"margin"`
	Precision           uint8			 `json:"precision"`
	Website             string			 `json:"website"`
	Description         string 			 `json:"description"`
	ExchangeRate        sdk.Dec			 `json:"exchange_rate"`
	Timestamp			string			 `json:"timestamp"`
}

func (it IssueToken) String() string {
	return it.String()
}

func NewIssueToken(ownerAddress sdk.AccAddress, outerName string ,supplyNum sdk.Int,margin sdk.Coin,
	precision uint8 ,website string,description string) IssueToken {
	return IssueToken{
		OwnerAddress: ownerAddress,
		OuterName:outerName,
		SupplyNum: supplyNum,
		Margin:margin,
		Precision:precision,
		Website:website,
		Description:description,
	}
}

type IssueTokenKeeper struct{
	key           sdk.StoreKey
    cdc          *codec.Codec
	codespace sdk.CodespaceType
}



/**
*	发币keeper
 */
func NewIssueTokenKeeper(key sdk.StoreKey,cdc *codec.Codec,	 codespace sdk.CodespaceType) IssueTokenKeeper {
	return IssueTokenKeeper{
		key:	key,
		cdc:             cdc,
		codespace:     codespace,
	}
}

func (itk IssueTokenKeeper) GetTokenByInnerName(ctx sdk.Context,innerNameByte  []byte) (exist bool, token IssueToken) {
	store := ctx.KVStore(itk.key)
	b := store.Get(itk.getIssueTokenKeyName(innerNameByte))
	if b == nil {
		return  false,IssueToken{}
	}
	itk.cdc.MustUnmarshalBinaryLengthPrefixed(b, &token)
	return true ,token
}

func (itk IssueTokenKeeper) SetTokenByInnerName(ctx sdk.Context,innerNameByte  []byte,token IssueToken) {
	store := ctx.KVStore(itk.key)
	b := itk.cdc.MustMarshalBinaryLengthPrefixed(token)
	store.Set(itk.getIssueTokenKeyName(innerNameByte), b)
}

//赎回
func (itk IssueTokenKeeper) Redeem(ctx sdk.Context, account sdk.AccAddress,issueCoin sdk.Coin,bk BaseKeeper)(sdk.Tags,sdk.Error) {
	isExist,token := itk.GetTokenByInnerName(ctx,[]byte(issueCoin.Denom))
	if !isExist{
		return nil,ErrInvaidCoin(itk.codespace,issueCoin.Denom)
	}

	tmpTags, err := bk.SendCoins(ctx, account, token.ExchangeAddress,sdk.Coins{issueCoin})
	if err != nil {
		return nil,err
	}

	exchangeAmount := token.ExchangeRate.MulInt(issueCoin.Amount).TruncateInt()
	exchangeCoin := sdk.Coin{Denom:sdk.DefaultBCVDemon,Amount:exchangeAmount}
	if exchangeAmount.GT(sdk.ZeroInt()){
		_, err = bk.SendCoins(ctx,token.ExchangeAddress,account,sdk.Coins{exchangeCoin})
		if err != nil {
			return  nil ,err
		}
	}

	tags := tmpTags
	tags = tags.AppendTag(TagRedeemExchangeAddr,token.ExchangeAddress.String())
	tags = tags.AppendTag(TagRedeemExchangeAddrAddAmonut,exchangeCoin.String())
	return  tags,nil
}

/**
 * 存储发行代币信息
 */
func (itk IssueTokenKeeper) getIssueTokenKeyName(innerNameByte []byte) []byte {
	return append(IssueTokenStoreKeyPrefix, innerNameByte...)
}

func (itk IssueTokenKeeper) IssueToken(ctx sdk.Context,  token IssueToken,bk Keeper) (sdk.Tags, sdk.Error,IssueToken)  {
	var tags  sdk.Tags
	suffix := hex.EncodeToString(tmhash.Sum(ctx.TxBytes()))[0:3]
	innerName  := token.OuterName+"-"+suffix
	token.InnerName = innerName
	token.SupplyNum = token.SupplyNum.Mul(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(token.Precision)), nil)))
	innerNameByte := []byte(innerName)
	exist ,_ := bk.GetTokenByInnerName(ctx,innerName)
	if exist {
		return nil,ErrIssueTokenExist(bk.Codespace()),IssueToken{}
	}

	exchangeAddress := sdk.AccAddress(crypto.AddressHash(innerNameByte))
	token.ExchangeAddress = exchangeAddress
	token.ExchangeRate  = token.Margin.Amount.ToDec().QuoTruncate(token.SupplyNum.ToDec())
	token.Timestamp     = ctx.BlockHeader().Time.Format(time.RFC3339)

	if token.Margin.Amount.LT(sdk.NewInt(int64(sdk.ISSUE_TOKEN_MIN_MARGIN_NUM)).Mul(sdk.NewInt(1000000))){
		return  nil, ErrIssueTokenMarginNumErr(DefaultCodespace),IssueToken{}
	}

	tmpTags,err := bk.SendCoins(ctx,token.OwnerAddress,exchangeAddress,sdk.Coins{token.Margin})
	if err != nil{
		return  nil ,err,IssueToken{}
	}
	tags = tags.AppendTags(tmpTags)

	//add new token
	supplyCoins := sdk.Coins{sdk.Coin{Denom:innerName,Amount:token.SupplyNum}}
	_,tmpTags,err = bk.AddCoins(ctx, token.OwnerAddress,supplyCoins)
	if err != nil{
		return tmpTags,err,IssueToken{}
	}

	//设置issueToken
	bk.SetTokenByInnerName(ctx,innerName,token)

	tags = tags.AppendTags(tmpTags)
	tags = tags.AppendTag(TagIssueTokenMargin,token.Margin.String())
	tags = tags.AppendTag(TagIssueTokenSupplyNum,sdk.Coin{Amount:token.SupplyNum,Denom:innerName}.String())
	tags = tags.AppendTag(TagIssueTokenExchangeRate,token.ExchangeRate.String())
	tags = tags.AppendTag(TagIssueExchangeAddr,token.ExchangeAddress.String())
	tags = tags.AppendTag(TagIssueTokenInnerName,token.InnerName)
	return tags,nil,token
}


func (itk IssueTokenKeeper) GetTokens(ctx sdk.Context) (tokens []IssueToken) {
	store := ctx.KVStore(itk.key)
	iterator := sdk.KVStorePrefixIterator(store, IssueTokenStoreKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var token IssueToken
		itk.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &token)
		tokens = append(tokens, token)
	}
	return tokens
}


func (itk IssueTokenKeeper) AddMarginByInnerName(ctx sdk.Context,account sdk.AccAddress,innerName string,amount sdk.Coin,bk Keeper) (sdk.Tags, sdk.Error,IssueToken)  {
	var tags  sdk.Tags
	exist ,token := bk.GetTokenByInnerName(ctx,innerName)
	if !exist {
		return nil,ErrIssueNoExist(bk.Codespace()),IssueToken{}
	}

	//转移资产
	tmpTags,err := bk.SendCoins(ctx,account,token.ExchangeAddress,sdk.Coins{amount})
	tags = tags.AppendTags(tmpTags)

	if err != nil{
		return tmpTags,err,token;
	}

	token.Margin = token.Margin.Add(amount)
	token.ExchangeRate  = token.Margin.Amount.ToDec().QuoTruncate(token.SupplyNum.ToDec())
	itk.SetTokenByInnerName(ctx,[]byte(innerName),token)

	tags = tags.AppendTag(TagIssueTokenAddMargin,amount.String())
	return tags,nil,token;
}
