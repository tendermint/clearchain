package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// Base SDK reserves 0 ~ 99.

	CodeInvalidAmount      sdk.CodeType = 100
	CodeInvalidAddress     sdk.CodeType = 101
	CodeInvalidPubKey      sdk.CodeType = 102
	CodeInvalidAccount     sdk.CodeType = 103
	CodeWrongSigner        sdk.CodeType = 104
	CodeWrongMessageFormat sdk.CodeType = 105
)

func ErrInvalidAmount(typ string) sdk.Error {
	return sdk.NewError(CodeInvalidAmount, fmt.Sprintf("invalid amount: %s", typ))
}

func ErrInvalidAddress(typ string) sdk.Error {
	return sdk.NewError(CodeInvalidAddress, fmt.Sprintf("invalid address: %s", typ))
}

func ErrInvalidPubKey(typ string) sdk.Error {
	return sdk.NewError(CodeInvalidPubKey, fmt.Sprintf("invalid pub key: %s", typ))
}

func ErrInvalidAccount(typ string) sdk.Error {
	return sdk.NewError(CodeInvalidAccount, fmt.Sprintf("invalid account: %s", typ))
}

func ErrWrongSigner(typ string) sdk.Error {
	return sdk.NewError(CodeWrongSigner, fmt.Sprintf("wrong signer: %s", typ))
}
func ErrWrongMsgFormat(typ string) sdk.Error {
	return sdk.NewError(CodeWrongMessageFormat, fmt.Sprintf("wrong message format: %s", typ))
}
