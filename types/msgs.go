package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
)

// Defines all the messages (requests) supported by the app

// message types definitions
const (
	DepositType       = "deposit"
	SettlementType    = "settlement"
	WithdrawType      = "withdraw"
	CreateAccountType = "createAccount"
)

// DepositMsg defines the properties of an asset transfer
type DepositMsg struct {
	Sender    crypto.Address
	Recipient crypto.Address
	Amount    sdk.Coin
}

// ensure DepositMsg implements the sdk.Msg interface
var _ sdk.Msg = (*DepositMsg)(nil)

//Called by SDk automatically
func (d DepositMsg) ValidateBasic() sdk.Error {
	if d.Amount.Amount <= 0 {
		return ErrInvalidAmount("negative amount")
	}
	if d.Amount.Denom == "" {
		return ErrInvalidAmount("empty denom")
	}
	if err := validateAddress(d.Sender); err != nil {
		return err
	}
	if err := validateAddress(d.Recipient); err != nil {
		return err
	}

	if bytes.Equal(d.Sender, d.Recipient) {
		return ErrInvalidAddress("sender and recipient have the same address")
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

//Called by SDk automatically
func (s SettleMsg) ValidateBasic() sdk.Error {
	// amount may be negative
	if s.Amount.Amount == 0 {
		return ErrInvalidAmount("empty or 0 amount not allowed")
	}
	if s.Amount.Denom == "" {
		return ErrInvalidAmount("empty denom")
	}
	if err := validateAddress(s.Sender); err != nil {
		return err
	}
	if err := validateAddress(s.Recipient); err != nil {
		return err
	}

	if bytes.Equal(s.Sender, s.Recipient) {
		return ErrInvalidAddress("sender and recipient have the same address")
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

// Called by SDk automatically
func (w WithdrawMsg) ValidateBasic() sdk.Error {
	if w.Amount.Amount <= 0 {
		return ErrInvalidAmount("negative or 0 amount not allowed")
	}
	if w.Amount.Denom == "" {
		return ErrInvalidAmount("empty denom")
	}
	if err := validateAddress(w.Sender); err != nil {
		return err
	}
	if err := validateAddress(w.Recipient); err != nil {
		return err
	}
	if err := validateAddress(w.Operator); err != nil {
		return err
	}

	if bytes.Equal(w.Sender, w.Recipient) {
		return ErrInvalidAddress("sender and recipient have the same address")
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

//*********************

type CreateAccountMsg struct {
	Creator     crypto.Address
	PubKey      crypto.PubKey
	AccountType string
}

var _ sdk.Msg = CreateAccountMsg{}

//Called by SDk automatically
func (msg CreateAccountMsg) ValidateBasic() sdk.Error {
	if err := validateAddress(msg.Creator); err != nil {
		return err
	}
	if msg.PubKey == nil {
		return ErrInvalidPubKey("pub key is nil")
	}
	if bytes.Equal(msg.Creator, msg.PubKey.Address()) {
		return ErrInvalidAddress("creator and new account have the same address")
	}
	if !IsValidEntityType(msg.AccountType) {
		return ErrInvalidAccount("unrecognized entity type")
	}
	return nil
}

// Type returns the message type.
// Must be alphanumeric or empty.
func (msg CreateAccountMsg) Type() string {
	return CreateAccountType
}

// Get returns some property of the Msg.
func (msg CreateAccountMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// GetSignBytes returns the canonical byte representation of the Msg.
func (msg CreateAccountMsg) GetSignBytes() []byte {
	bz, err := cdc.MarshalBinary(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

// GetSigners returns the addrs of signers that must sign.
// CONTRACT: All signatures must be present to be valid.
// CONTRACT: Returns addrs in some deterministic order.
func (msg CreateAccountMsg) GetSigners() []crypto.Address {
	return []crypto.Address{msg.Creator}
}

//******************************* helper methods *****************************

func validateAddress(addr crypto.Address) sdk.Error {
	if addr == nil {
		return ErrInvalidAddress("address is nil")
	}
	if len(addr) != 20 {
		return ErrInvalidAddress("invalid address length")
	}
	return nil
}
