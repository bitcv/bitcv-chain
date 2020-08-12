package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"

	bac "github.com/bitcv-chain/bitcv-chain/app"
	sdk "github.com/bitcv-chain/bitcv-chain/types"
	"github.com/bitcv-chain/bitcv-chain/x/auth"
	bacinit "github.com/bitcv-chain/bitcv-chain/cmd/bacinit"
	"github.com/bitcv-chain/bitcv-chain/codec"
	"github.com/bitcv-chain/bitcv-chain/x/mint"
	bacv1 "github.com/bitcv-chain/bitcv-chain/bacchain/v1_00"

)

func init() {

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	rootCmd.AddCommand(txCmd)
	rootCmd.AddCommand(pubkeyCmd)
	rootCmd.AddCommand(encodePubkeyCmd)
	rootCmd.AddCommand(addrCmd)
	rootCmd.AddCommand(decodeBech32AccAddrCmd)
	rootCmd.AddCommand(encodeBech32AccAddrCmd)
	rootCmd.AddCommand(hackCmd)
	rootCmd.AddCommand(rawBytesCmd)
	rootCmd.AddCommand(parseGensisCmd)
	rootCmd.AddCommand(printSystemAddrCmd)
	rootCmd.AddCommand(blockCmd)
}

var rootCmd = &cobra.Command{
	Use:          "bacdebug",
	Short:        "bacchain debug tool",
	SilenceUsage: true,
}

var txCmd = &cobra.Command{
	Use:   "tx",
	Short: "Decode a bac tx from hex or base64",
	RunE:  runTxCmd,
}

var pubkeyCmd = &cobra.Command{
	Use:   "pubkey",
	Short: "Decode a pubkey from hex, base64, or bech32",
	RunE:  runPubKeyCmd,
}
var encodePubkeyCmd = &cobra.Command{
	Use:   "encode_pubkey",
	Short: "Decode a pubkey from hex, base64, or bech32",
	RunE:  runEncodePubkeyCmd,
}
var addrCmd = &cobra.Command{
	Use:   "addr",
	Short: "Convert an address between hex and bech32",
	RunE:  runAddrCmd,
}

var decodeBech32AccAddrCmd = &cobra.Command{
	Use:   "decode_bech32",
	Short: "Convert an address between hex and bech32",
	RunE:  runDecodeBech32AccAddrCmd,
}

var encodeBech32AccAddrCmd = &cobra.Command{
	Use:   "encode_bech32",
	Short: "Convert an address to bech32",
	RunE:  runEncodeBech32AccAddrCmd,
}

var hackCmd = &cobra.Command{
	Use:   "hack",
	Short: "Boilerplate to Hack on an existing state by scripting some Go...",
	RunE:  runHackCmd,
}

var blockCmd = &cobra.Command{
	Use:   "block",
	Short: "traversal block..",
	RunE:  runBlockCmd,
}

var rawBytesCmd = &cobra.Command{
	Use:   "raw-bytes",
	Short: "Convert raw bytes output (eg. [10 21 13 255]) to hex",
	RunE:  runRawBytesCmd,
}

var parseGensisCmd = &cobra.Command{
	Use:"parse_genesis",
	Short:"parse_genesis gensis.json",
	RunE : runParseGenesis,
}


var printSystemAddrCmd = &cobra.Command{
	Use:"print_system_addr",
	Short:"parse_genesis gensis.json",
	RunE : runPrintSystemAddr,
}


func runRawBytesCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expected single arg")
	}
	stringBytes := args[0]
	stringBytes = strings.Trim(stringBytes, "[")
	stringBytes = strings.Trim(stringBytes, "]")
	spl := strings.Split(stringBytes, " ")

	byteArray := []byte{}
	for _, s := range spl {
		b, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		byteArray = append(byteArray, byte(b))
	}
	fmt.Printf("%X\n", byteArray)
	return nil
}

func runPubKeyCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expected single arg")
	}

	pubkeyString := args[0]
	var pubKeyI crypto.PubKey

	// try hex, then base64, then bech32
	pubkeyBytes, err := hex.DecodeString(pubkeyString)
	if err != nil {
		var err2 error
		pubkeyBytes, err2 = base64.StdEncoding.DecodeString(pubkeyString)
		if err2 != nil {
			var err3 error
			pubKeyI, err3 = sdk.GetAccPubKeyBech32(pubkeyString)
			if err3 != nil {
				var err4 error
				pubKeyI, err4 = sdk.GetValPubKeyBech32(pubkeyString)

				if err4 != nil {
					var err5 error
					pubKeyI, err5 = sdk.GetConsPubKeyBech32(pubkeyString)
					if err5 != nil {
						return fmt.Errorf(`Expected hex, base64, or bech32. Got errors:
								hex: %v,
								base64: %v
								bech32 Acc: %v
								bech32 Val: %v
								bech32 Cons: %v`,
							err, err2, err3, err4, err5)
					}

				}
			}

		}
	}

	var pubKey ed25519.PubKeyEd25519
	if pubKeyI == nil {
		copy(pubKey[:], pubkeyBytes)
	} else {
		pubKey = pubKeyI.(ed25519.PubKeyEd25519)
		pubkeyBytes = pubKey[:]
	}

	cdc := bac.MakeCodec()
	pubKeyJSONBytes, err := cdc.MarshalJSON(pubKey)
	if err != nil {
		return err
	}
	accPub, err := sdk.Bech32ifyAccPub(pubKey)
	if err != nil {
		return err
	}
	valPub, err := sdk.Bech32ifyValPub(pubKey)
	if err != nil {
		return err
	}

	consenusPub, err := sdk.Bech32ifyConsPub(pubKey)
	if err != nil {
		return err
	}
	fmt.Println("Address:", pubKey.Address())
	fmt.Printf("Hex: %X\n", pubkeyBytes)
	fmt.Println("JSON (base64):", string(pubKeyJSONBytes))
	fmt.Println("Bech32 Acc:", accPub)
	fmt.Println("Bech32 Validator Operator:", valPub)
	fmt.Println("Bech32 Validator Consensus:", consenusPub)
	return nil
}


func runEncodePubkeyCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expected single arg")
	}

	pubkeyString := args[0]
	var pubKeyI crypto.PubKey
	fmt.Println("pubkey:",pubkeyString)

	var pubkeyBytes []byte
	//pubkeyBytes, err := hex.DecodeString(pubkeyString)
	pubKeyI, err5 := sdk.GetConsPubKeyBech32(pubkeyString)
	if err5 != nil {
		return fmt.Errorf(`Expected hex, base64, or bech32. Got errors:bech32 Cons: %v`,
			 err5)
	}

	var pubKey ed25519.PubKeyEd25519
	if pubKeyI == nil {
		copy(pubKey[:], pubkeyBytes)
	} else {
		fmt.Println("copy succ...")
		pubKey = pubKeyI.(ed25519.PubKeyEd25519)
		pubkeyBytes = pubKey[:]
	}

	cdc := bac.MakeCodec()
	pubKeyJSONBytes, err := cdc.MarshalJSON(pubKey)
	if err != nil {
		return err
	}
	accPub, err := sdk.Bech32ifyAccPub(pubKey)
	if err != nil {
		return err
	}
	valPub, err := sdk.Bech32ifyValPub(pubKey)
	if err != nil {
		return err
	}

	consenusPub, err := sdk.Bech32ifyConsPub(pubKey)
	if err != nil {
		return err
	}
	fmt.Println("validator_address,Address:在blocks接口中使用", pubKey.Address())
	fmt.Printf("Hex: %X\n", pubkeyBytes)
	fmt.Println("JSON (base64):", string(pubKeyJSONBytes))
	fmt.Println("Bech32 Acc:", accPub)
	fmt.Println("Bech32 Validator Operator:", valPub)
	fmt.Println("Bech32 Validator Consensus:", consenusPub)
	return nil
}

func runAddrCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expected single arg")
	}

	addrString := args[0]
	var addr []byte

	// try hex, then bech32
	var err error
	addr, err = hex.DecodeString(addrString)
	if err != nil {
		var err2 error
		addr, err2 = sdk.AccAddressFromBech32(addrString)
		if err2 != nil {
			var err3 error
			addr, err3 = sdk.ValAddressFromBech32(addrString)

			if err3 != nil {
				return fmt.Errorf(`Expected hex or bech32. Got errors:
			hex: %v,
			bech32 acc: %v
			bech32 val: %v
			`, err, err2, err3)

			}
		}
	}

	accAddr := sdk.AccAddress(addr)
	valAddr := sdk.ValAddress(addr)
	conAddr := sdk.ConsAddress(addr)
	fmt.Println("Address:", addr)
	fmt.Printf("Address (hex): %X\n", addr)
	fmt.Printf("Bech32 Acc: %s\n", accAddr)
	fmt.Printf("Bech32 Val: %s\n", valAddr)
	fmt.Printf("Bech32 Con: %s\n", conAddr)

	return nil
}


func  runDecodeBech32AccAddrCmd (cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expected single arg")
	}

	addrString := args[0]
	addr, err2 := sdk.AccAddressFromBech32(addrString)
	if err2 != nil {
			return fmt.Errorf(`Expected hex or bech32. Got errors:bech32 acc: %v`,err2 )
	}

	fmt.Println("bench32Addr:",addrString)
	accAddr := sdk.AccAddress(addr)
	valAddr := sdk.ValAddress(addr)
	conAddr := sdk.ConsAddress(addr)

	fmt.Println("accAddr",accAddr)
	fmt.Println("accAddr",valAddr)
	fmt.Println("accAddr",conAddr)

	fmt.Println("Address:", addr)
	fmt.Printf("Address (hex): %X\n", addr)
	fmt.Printf("Bech32 Acc: %s\n", accAddr)
	fmt.Printf("Bech32 Val: %s\n", valAddr)
	fmt.Printf("Bech32 Con: %s\n", conAddr)

	fmt.Println(accAddr.String())

	return nil
}

func  runEncodeBech32AccAddrCmd (cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expected single arg")
	}

	addrString := args[0]
	fmt.Println(addrString)
	accAddr2 := sdk.AccAddress(addrString)//error
	accAddr ,_:= sdk.AccAddressFromHex(addrString)//ok
	fmt.Println(accAddr)
	fmt.Printf("Bech32 Acc: %s\n", accAddr)
	fmt.Println(accAddr2)
	return nil
}

func runTxCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expected single arg")
	}

	txString := args[0]

	// try hex, then base64
	txBytes, err := hex.DecodeString(txString)
	if err != nil {
		var err2 error
		txBytes, err2 = base64.StdEncoding.DecodeString(txString)
		if err2 != nil {
			return fmt.Errorf(`Expected hex or base64. Got errors:
			hex: %v,
			base64: %v
			`, err, err2)
		}
	}

	var tx = auth.StdTx{}
	cdc := bac.MakeCodec()

	err = cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
	if err != nil {
		return err
	}

	bz, err := cdc.MarshalJSON(tx)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer([]byte{})
	err = json.Indent(buf, bz, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(buf.String())
	return nil
}

func runParseGenesis(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expected filename arg")
	}

	fmt.Println("充电地址",sdk.AccAddrEnergyPool.String())
	fmt.Println("耗电地址",sdk.AccAddrEnergyBurn.String())
	fmt.Println("burn地址",sdk.AccAddrBcvBurnFromBuyBcvStake.String())
	fmt.Println("stake地址",sdk.AccAddrBcvstakePool.String())

	fileName := args[0]

	//获取高度 bcv_export_{ip}_{height}
	height,_:= strconv.ParseInt(strings.Split(strings.Split(fileName,"_")[3],".")[0],10,64)
	//计算现在应该出都BAC数量
	diff := mint.FirstReduceHeight  - bacv1.StartParamInitHeight
	var bacNumShoudProduce  int64
	if height < diff{
		bacNumShoudProduce = height  * 10
	}else{
		bacNumShoudProduce = diff * 10 + (height - diff) * 5
	}
	genDoc, err := bacinit.LoadGenesisDoc(codec.Cdc, fileName)
	if err != nil {
		fmt.Println(err)
	}

	genesisState := bac.GenesisState{}
	cdc := bac.MakeCodec()
	if err:=cdc.UnmarshalJSON(genDoc.AppState, &genesisState);err != nil{
		fmt.Println(err)
	}

	//用户总资产
	var totalAccountCoins  sdk.Coins
	for _,account := range genesisState.Accounts{
		totalAccountCoins = totalAccountCoins.Add(account.Coins)
	}
	fmt.Println(".................")
	fmt.Println("所有账户coins:",totalAccountCoins)
	fmt.Println("genesis BAC用户持有:",totalAccountCoins.AmountOf(sdk.DEFAULT_FEE_COIN))
	fmt.Println("genesis BAC烧毁:",genesisState.AuthData.BacManagePool.AlreadyBurn)
	fmt.Println("genesis BAC社区:",	genesisState.DistrData.FeePool.CommunityPool)
	fmt.Println("总计:",totalAccountCoins.AmountOf(sdk.DEFAULT_FEE_COIN).Add(genesisState.AuthData.BacManagePool.AlreadyBurn.AmountOf(sdk.DEFAULT_FEE_COIN)).Add(genesisState.DistrData.FeePool.CommunityPool.AmountOf(sdk.DEFAULT_FEE_COIN).TruncateInt()))
	fmt.Println("计算总产生BAC 1:",bacNumShoudProduce * 1000000000 + 10000000000000000)
	fmt.Println("genesis BAC社区中总产生:",genesisState.AuthData.BacManagePool.GenerateTokenSupply())
	fmt.Println("计算总产生BAC 2:",bacNumShoudProduce * 1000000000)
	fmt.Println(".................")




	return nil
}



func runPrintSystemAddr(cmd *cobra.Command, args []string) error {
	//fmt.Println("销毁代币地址",sdk.AccAddrBurn.String())
	//fmt.Println("购买矿机销毁BCV地址",sdk.AccAddrBcvBurn.String())
	//fmt.Println("stake地址",sdk.AccAddrBcvstakePool.String())
	//fmt.Println("充电地址",sdk.AccAddrEnergyPool.String())
	//fmt.Println("耗电地址",sdk.AccAddrEnergyBurn.String())


	fmt.Println("销毁代币地址",sdk.AccAddrBurn.String())
	fmt.Println("购买矿机销毁BCV地址",sdk.AccAddrBcvBurnFromBuyBcvStake.String())
	fmt.Println("stake地址",sdk.AccAddrBcvstakePool.String())
	fmt.Println("充电地址",sdk.AccAddrEnergyPool.String())
	fmt.Println("耗电地址",sdk.AccAddrEnergyBurn.String())

	return nil
}
func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
