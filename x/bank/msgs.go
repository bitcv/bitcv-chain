package bank

import (
	sdk "github.com/bitcv-chain/bitcv-chain/types"
	"strings"
	"regexp"
)

// RouterKey is they name of the bank module
const RouterKey = "bank"

// MsgSend - high level transaction of the coin module
type MsgSend struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address"`
	Amount      sdk.Coins      `json:"amount"`
}

var _ sdk.Msg = MsgSend{}

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgSend(fromAddr, toAddr sdk.AccAddress, amount sdk.Coins) MsgSend {
	return MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: amount}
}

// Route Implements Msg.
func (msg MsgSend) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSend) Type() string { return "send" }

// ValidateBasic Implements Msg.
func (msg MsgSend) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	if msg.ToAddress.Empty() {
		return sdk.ErrInvalidAddress("missing recipient address")
	}
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins("send amount is invalid: " + msg.Amount.String())
	}
	if !msg.Amount.IsAllPositive() {
		return sdk.ErrInsufficientCoins("send amount must be positive")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSend) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgSend) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

// MsgMultiSend - high level transaction of the coin module
type MsgMultiSend struct {
	Inputs  []Input  `json:"inputs"`
	Outputs []Output `json:"outputs"`
}

var _ sdk.Msg = MsgMultiSend{}

// NewMsgMultiSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgMultiSend(in []Input, out []Output) MsgMultiSend {
	return MsgMultiSend{Inputs: in, Outputs: out}
}

// Route Implements Msg
func (msg MsgMultiSend) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgMultiSend) Type() string { return "multisend" }

// ValidateBasic Implements Msg.
func (msg MsgMultiSend) ValidateBasic() sdk.Error {
	// this just makes sure all the inputs and outputs are properly formatted,
	// not that they actually have the money inside
	if len(msg.Inputs) == 0 {
		return ErrNoInputs(DefaultCodespace).TraceSDK("")
	}
	if len(msg.Outputs) == 0 {
		return ErrNoOutputs(DefaultCodespace).TraceSDK("")
	}

	return ValidateInputsOutputs(msg.Inputs, msg.Outputs)
}

// GetSignBytes Implements Msg.
func (msg MsgMultiSend) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgMultiSend) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.Inputs))
	for i, in := range msg.Inputs {
		addrs[i] = in.Address
	}
	return addrs
}

// Input models transaction input
type Input struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

// ValidateBasic - validate transaction input
func (in Input) ValidateBasic() sdk.Error {
	if len(in.Address) == 0 {
		return sdk.ErrInvalidAddress(in.Address.String())
	}
	if !in.Coins.IsValid() {
		return sdk.ErrInvalidCoins(in.Coins.String())
	}
	if !in.Coins.IsAllPositive() {
		return sdk.ErrInvalidCoins(in.Coins.String())
	}
	return nil
}

// NewInput - create a transaction input, used with MsgMultiSend
func NewInput(addr sdk.AccAddress, coins sdk.Coins) Input {
	return Input{
		Address: addr,
		Coins:   coins,
	}
}

// Output models transaction outputs
type Output struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

// ValidateBasic - validate transaction output
func (out Output) ValidateBasic() sdk.Error {
	if len(out.Address) == 0 {
		return sdk.ErrInvalidAddress(out.Address.String())
	}
	if !out.Coins.IsValid() {
		return sdk.ErrInvalidCoins(out.Coins.String())
	}
	if !out.Coins.IsAllPositive() {
		return sdk.ErrInvalidCoins(out.Coins.String())
	}
	return nil
}

// NewOutput - create a transaction output, used with MsgMultiSend
func NewOutput(addr sdk.AccAddress, coins sdk.Coins) Output {
	return Output{
		Address: addr,
		Coins:   coins,
	}
}

// ValidateInputsOutputs validates that each respective input and output is
// valid and that the sum of inputs is equal to the sum of outputs.
func ValidateInputsOutputs(inputs []Input, outputs []Output) sdk.Error {
	var totalIn, totalOut sdk.Coins

	for _, in := range inputs {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
		totalIn = totalIn.Add(in.Coins)
	}

	for _, out := range outputs {
		if err := out.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
		totalOut = totalOut.Add(out.Coins)
	}

	// make sure inputs and outputs match
	if !totalIn.IsEqual(totalOut) {
		return ErrInputOutputMismatch(DefaultCodespace)
	}

	return nil
}





// MsgIssueToken - high level transaction of the coin module
type MsgIssueToken struct {
	OwnerAddress     sdk.AccAddress  `json:"owner_address"`
	OuterName        string	         `json:"outer_name"`
	SupplyNum        sdk.Int         `json:"supply_num"`
	Margin        sdk.Coin         `json:"margin"`
	Precision        uint8           `json:"precision"`
	Website          string          `json:"website"`
	Description      string          `json:"description"`
}

var _ sdk.Msg = MsgIssueToken{}

// NewMsgIssueToken
func NewMsgIssueToken(ownerAddress sdk.AccAddress, outerName string ,supplyNum sdk.Int,margin sdk.Coin,
	precision uint8 ,website string,description string) MsgIssueToken {
	return MsgIssueToken{
		OwnerAddress: ownerAddress,
		OuterName:outerName,
		SupplyNum: supplyNum,
		Margin:margin,
		Precision:precision,
		Website:website,
		Description:description,
	}
}

// Route Implements Msg.
func (msg MsgIssueToken) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgIssueToken) Type() string { return "issue_token" }

// ValidateBasic Implements Msg.
func (msg MsgIssueToken) ValidateBasic() sdk.Error {
	return CheckMsgIssueToken(msg)
}

// GetSignBytes Implements Msg.
func (msg MsgIssueToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgIssueToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
}


func CheckMsgIssueToken(msg MsgIssueToken) sdk.Error{
	if msg.OwnerAddress.Empty() {
		return sdk.ErrInvalidAddress("missing owner address")
	}

	if len(msg.Website) > 128 {
		return ErrIssueTokenWebsiteLenErr(DefaultCodespace)
	}
	if len(msg.Description) > 512 {
		return ErrIssueTokenDescriptionLenErr(DefaultCodespace)
	}

	precision := sdk.NewInt(int64(msg.Precision))
	if precision.GT(sdk.NewInt(12)) ||  precision.LT(sdk.ZeroInt()){
		return ErrIssueTokenPrecisionErr(DefaultCodespace)
	}

	if msg.SupplyNum.LTE(sdk.ZeroInt()){
		return ErrIssueTokenSupplyNumErr(DefaultCodespace)
	}

	if msg.Margin.Amount.LT(sdk.NewInt(int64(sdk.ISSUE_TOKEN_MIN_MARGIN_NUM)).Mul(sdk.NewInt(1000000))){
		return ErrIssueTokenMarginNumErr(DefaultCodespace)
	}

	outerNameP :=`^[a-z]{2,15}$`
	reOuterName := regexp.MustCompile(outerNameP)

	if !reOuterName.MatchString(msg.OuterName){
		return ErrIssueTokenOutNameErr(DefaultCodespace)
	}

	if strings.HasPrefix(strings.ToLower(msg.OuterName),"test"){
		//if !msg.SupplyNum.Equal(sdk.ISSUE_TOKEN_TEST_SUPPLY_NUM){
		//	return ErrIssueTokenSupplyNumErr(DefaultCodespace)
		//}
		//
		//if !msg.MarginNum.Equal(sdk.ISSUE_TOKEN_TEST_MARGIN_NUM){
		//	return  ErrIssueTokenTestMarginNumErr(DefaultCodespace)
		//}
		//
		//if msg.Precision != uint8(sdk.ISSUE_TOKEN_TEST_PRECISION){
		//	return ErrIssueTokenTestPrecisionErr(DefaultCodespace)
		//}
	}

	return  nil
}




// MsgRedeem -
type MsgRedeem struct {
	Account sdk.AccAddress `json:"account"`
	Amount      sdk.Coin      `json:"amount"`
}

var _ sdk.Msg = MsgRedeem{}

func NewMsgRedeem(account sdk.AccAddress, amount sdk.Coin) MsgRedeem {
	return MsgRedeem{Account: account, Amount:amount}
}

func (msg MsgRedeem) Route() string { return RouterKey }

func (msg MsgRedeem) Type() string { return "redeem" }

func (msg MsgRedeem) ValidateBasic() sdk.Error {
	if msg.Account.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	if !msg.Amount.IsPositive() {
		return sdk.ErrInvalidCoins("amount is invalid: " + msg.Amount.String())
	}
	return nil
}

func (msg MsgRedeem) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgRedeem) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Account}
}



// MsgRedeem -
type MsgAddMargin struct {
	Account sdk.AccAddress `json:"account"`
	InnerName	string			`json:"inner_name"`
	Amount      sdk.Coin     `json:"amount"`
}

var _ sdk.Msg = MsgRedeem{}

func NewMsgAddMargin(account sdk.AccAddress, innerName string ,amount sdk.Coin) MsgAddMargin {
	return MsgAddMargin{Account: account,InnerName:innerName ,Amount:amount}
}

func (msg MsgAddMargin) Route() string { return RouterKey }

func (msg MsgAddMargin) Type() string { return "add_margin" }

func (msg MsgAddMargin) ValidateBasic() sdk.Error {
	if msg.Account.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	if !msg.Amount.IsPositive() {
		return sdk.ErrInvalidCoins("amount is invalid: " + msg.Amount.String())
	}
	if msg.Amount.Denom != sdk.DefaultBCVDemon{
		return ErrIssueTokenCoinErr(DefaultCodespace)
	}
	return nil
}

func (msg MsgAddMargin) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgAddMargin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Account}
}




//
type MsgEdata struct {
	Account sdk.AccAddress `json:"account"`
	Utype        uint8     `json:"utype"`
	Data		string 		`json:"data"`
}

var _ sdk.Msg = MsgEdata{}

func NewMsgEdata(account sdk.AccAddress, utype uint8,data string) MsgEdata {
	return MsgEdata{Account:account , Utype:utype,Data:data}
}

func (msg MsgEdata) Route() string { return RouterKey }

func (msg MsgEdata) Type() string { return "edata" }

func (msg MsgEdata) ValidateBasic() sdk.Error {
	if msg.Account.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	if len(msg.Data) <=0 || len(msg.Data) > 2048{
		return  ErrEdataLenghtErr(DefaultCodespace,len(msg.Data));
	}

	if msg.Utype <= 0 || msg.Utype > 255{
		return ErrEdataUtypeErr(DefaultCodespace)
	}
	return nil
}

func (msg MsgEdata) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

func (msg MsgEdata) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Account}
}
