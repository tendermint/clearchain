package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodeInvalidAmount  sdk.CodeType = 100
	CodeInvalidAddress sdk.CodeType = 101
	CodeSameAddress    sdk.CodeType = 102
	CodeInvalidPubKey  sdk.CodeType = 103
	CodeInvalidAccount sdk.CodeType = 104
	CodeWrongSigner    sdk.CodeType = 105
)

func ErrWrongSigner(typ string) sdk.Error {
	return sdk.NewError(CodeWrongSigner, fmt.Sprintf("wrong signer: %s", typ))
}

func ErrInvalidAddress(typ string) sdk.Error {
	return sdk.NewError(CodeInvalidAddress, fmt.Sprintf("invalid address: %s", typ))
}

func ErrInvalidAccount(typ string) sdk.Error {
	return sdk.NewError(CodeInvalidAccount, fmt.Sprintf("invalid account: %s", typ))
}
