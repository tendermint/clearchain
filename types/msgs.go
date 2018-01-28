package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
)

const (
	DepositType    = "deposit"
	SettlementType = "settlement"
	WithdrawType   = "withdraw"
)

// DepositMsg defines the properties of an asset transfer
type DepositMsg struct {
	Sender    crypto.Address
	Recipient crypto.Address
	Amount    sdk.Coin
}

var _ sdk.Msg = DepositMsg{}
var _ sdk.Msg = (*DepositMsg)(nil)

func (d DepositMsg) ValidateBasic() sdk.Error {
	if d.Amount.Amount <= 0 {
		return sdk.NewError(CodeInvalidAmount, "negative amount")
	}
	if d.Amount.Denom == "" {
		return sdk.NewError(CodeInvalidAmount, "invalid denom")
	}
	if len(d.Sender) != 20 {
		return sdk.NewError(CodeInvalidAddress, "invalid sender address")
	}
	if len(d.Recipient) != 20 {
		return sdk.NewError(CodeInvalidAddress, "invalid recipient address")
	}
	if bytes.Equal(d.Sender, d.Recipient) {
		return sdk.NewError(CodeSameAddress, "same addresses")
	}

	return nil
}

// Return the message type.
// Must be alphanumeric or empty.
func (d DepositMsg) Type() string {
	return DepositType
}

// Get some property of the Msg.
func (d DepositMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// Get the canonical byte representation of the Msg.
func (d DepositMsg) GetSignBytes() []byte {
	bz, err := cdc.MarshalBinary(d)
	if err != nil {
		panic(err)
	}
	return bz
}

// Signers returns the addrs of signers that must sign.
// CONTRACT: All signatures must be present to be valid.
// CONTRACT: Returns addrs in some deterministic order.
func (d DepositMsg) GetSigners() []crypto.Address {
	return []crypto.Address{d.Sender}
}

type SettleMsg struct {
	Sender    crypto.Address
	Recipient crypto.Address
	Amount    sdk.Coin
}

var _ sdk.Msg = SettleMsg{}

// amount may be negative
func (s SettleMsg) ValidateBasic() sdk.Error {

	if s.Amount.Amount == 0 {
		return sdk.NewError(CodeInvalidAmount, "invalid amount")
	}
	if s.Amount.Denom == "" {
		return sdk.NewError(CodeInvalidAmount, "invalid denom")
	}
	if len(s.Sender) != 20 {
		return sdk.NewError(CodeInvalidAddress, "invalid sender address")
	}
	if len(s.Recipient) != 20 {
		return sdk.NewError(CodeInvalidAddress, "invalid recipient address")
	}
	if bytes.Equal(s.Sender, s.Recipient) {
		return sdk.NewError(CodeSameAddress, "same addresses")
	}

	return nil
}

// Return the message type.
// Must be alphanumeric or empty.
func (s SettleMsg) Type() string {
	return SettlementType
}

// Get some property of the Msg.
func (s SettleMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// Get the canonical byte representation of the Msg.
func (s SettleMsg) GetSignBytes() []byte {
	bz, err := cdc.MarshalBinary(s)
	if err != nil {
		panic(err)
	}
	return bz
}

// Signers returns the addrs of signers that must sign.
// CONTRACT: All signatures must be present to be valid.
// CONTRACT: Returns addrs in some deterministic order.
func (s SettleMsg) GetSigners() []crypto.Address {
	return []crypto.Address{s.Sender}
}

type WithdrawMsg struct {
	Sender    crypto.Address
	Recipient crypto.Address
	Operator  crypto.Address
	Amount    sdk.Coin
}

var _ sdk.Msg = WithdrawMsg{}

func (w WithdrawMsg) ValidateBasic() sdk.Error {
	if w.Amount.Amount <= 0 {
		return sdk.NewError(CodeInvalidAmount, "invalid amount")
	}
	if w.Amount.Denom == "" {
		return sdk.NewError(CodeInvalidAmount, "invalid denom")
	}
	if len(w.Sender) != 20 {
		return sdk.NewError(CodeInvalidAddress, "invalid sender address")
	}
	if len(w.Recipient) != 20 {
		return sdk.NewError(CodeInvalidAddress, "invalid recipient address")
	}
	if len(w.Operator) != 20 {
		return sdk.NewError(CodeInvalidAddress, "invalid operator address")
	}
	if bytes.Equal(w.Sender, w.Recipient) {
		return sdk.NewError(CodeSameAddress, "same addresses")
	}

	return nil
}

// Return the message type.
// Must be alphanumeric or empty.
func (w WithdrawMsg) Type() string {
	return WithdrawType
}

// Get some property of the Msg.
func (w WithdrawMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// Get the canonical byte representation of the Msg.
func (w WithdrawMsg) GetSignBytes() []byte {
	bz, err := cdc.MarshalBinary(w)
	if err != nil {
		panic(err)
	}
	return bz
}

// Signers returns the addrs of signers that must sign.
// CONTRACT: All signatures must be present to be valid.
// CONTRACT: Returns addrs in some deterministic order.
func (w WithdrawMsg) GetSigners() []crypto.Address {
	return []crypto.Address{w.Sender, w.Operator}
}
