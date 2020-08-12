package bank

import (
	sdk "github.com/bitcv-chain/bitcv-chain/types"
	"fmt"
)

// Bank errors reserve 100 ~ 199.
const (
	DefaultCodespace sdk.CodespaceType = "bank"

	CodeSendDisabled         sdk.CodeType = 101
	CodeInvalidInputsOutputs sdk.CodeType = 102

	CodeIssueTokenNeedMoreBcvToMargin	     sdk.CodeType   = 197
	CodeIssueTokenExist	     sdk.CodeType   = 198
	CodeInvaidCoin 			 sdk.CodeType = 199

	CodeIssueTokenWebsiteLenErr sdk.CodeType = 120
	CodeIssueTokenDescriptionLenErr sdk.CodeType = 121
	CodeIssueTokenPrecisionErr sdk.CodeType = 122
	CodeIssueTokenSupplyNumErr sdk.CodeType = 123
	CodeIssueTokenSupplyMarginErr sdk.CodeType = 124
	CodeIssueTokenOutNameErr sdk.CodeType = 125
	CodeIssueTokenTestMarginNumErr sdk.CodeType = 126
	CodeIssueTokenTestSupplyNumErr sdk.CodeType = 127
	CodeIssueTokenTestPrecisonErr  sdk.CodeType = 128
	CodeIssueTokenCoinErr			sdk.CodeType = 129
	CodeIssueTokenNoExist	     sdk.CodeType   = 130


	CodeEdataLenghtErr sdk.CodeType = 151
	CodeEdataUtypeErr 	sdk.CodeType = 152


)

// ErrNoInputs is an error
func ErrNoInputs(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInputsOutputs, "no inputs to send transacction")
}

// ErrNoOutputs is an error
func ErrNoOutputs(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInputsOutputs, "no outputs to send transaction")
}

// ErrInputOutputMismatch is an error
func ErrInputOutputMismatch(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInputsOutputs, "sum inputs != sum outputs")
}

// ErrSendDisabled is an error
func ErrSendDisabled(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSendDisabled, "send transactions are currently disabled")
}


func ErrInvaidCoin(codespace sdk.CodespaceType,coins string) sdk.Error {
	return sdk.NewError(codespace, CodeInvaidCoin, fmt.Sprintf("coins %v can not tran",coins))
}

func ErrIssueTokenExist(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace,CodeIssueTokenExist,"issue token already exist")
}

func ErrIssueNoExist(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace,CodeIssueTokenNoExist,"issue token is not  exist")
}

//发币的时候需要更多的BCV
func ErrIssueTokenNeedMoreBcvToMargin(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace,CodeIssueTokenNeedMoreBcvToMargin,"more margin required")
}



func ErrIssueTokenWebsiteLenErr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace,CodeIssueTokenWebsiteLenErr,"bad Website length ")
}

func ErrIssueTokenDescriptionLenErr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace,CodeIssueTokenDescriptionLenErr,"bad description length ")
}

func ErrIssueTokenPrecisionErr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace,CodeIssueTokenPrecisionErr,"bad precision,must big tran 0 and small than 6 ")
}

func ErrIssueTokenSupplyNumErr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace,CodeIssueTokenSupplyNumErr,"bad  supply num ")
}

func ErrIssueTokenMarginNumErr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace,CodeIssueTokenSupplyMarginErr,fmt.Sprintf("bad  margin num,min margin is %v bcv",sdk.ISSUE_TOKEN_MIN_MARGIN_NUM))
}

func ErrIssueTokenOutNameErr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace,CodeIssueTokenOutNameErr,"outer name error ")
}

func ErrIssueTokenTestMarginNumErr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace,CodeIssueTokenTestMarginNumErr,fmt.Sprintf("issue test token,margin num must be  %v",sdk.ISSUE_TOKEN_TEST_MARGIN_NUM))
}

func ErrIssueTokenTestSupplyNumErr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace,CodeIssueTokenTestSupplyNumErr,fmt.Sprintf("issue test token,supply num must be  %v",sdk.ISSUE_TOKEN_TEST_SUPPLY_NUM))
}

func ErrIssueTokenTestPrecisionErr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace,CodeIssueTokenTestPrecisonErr,fmt.Sprintf("issue test token,margin num must be  %v",sdk.ISSUE_TOKEN_TEST_PRECISION))
}

func ErrIssueTokenCoinErr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace,CodeIssueTokenCoinErr,"margin coin error,must be bcv")
}

func ErrEdataLenghtErr(codespace sdk.CodespaceType,len int) sdk.Error {
	return sdk.NewError(codespace,CodeEdataLenghtErr,fmt.Sprintf("data lenght must big than 0 and  small than 2048,now  tran after hex is  %v",len))
}

func ErrEdataUtypeErr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace,CodeEdataUtypeErr,"utype error,must 1-255")
}

