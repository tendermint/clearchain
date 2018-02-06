package types

import (
	"bytes"
	"strings"

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

// ValidateBasic is called by the SDK automatically.
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

// Type returns the message type.
// Must be alphanumeric or empty.
func (d DepositMsg) Type() string {
	return DepositType
}

// Get some property of the Msg.
func (d DepositMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// GetSignBytes returns the canonical byte representation of the Msg.
func (d DepositMsg) GetSignBytes() []byte {
	bz, err := cdc.MarshalBinary(d)
	if err != nil {
		panic(err)
	}
	return bz
}

// GetSigners returns the addrs of signers that must sign.
// CONTRACT: All signatures must be present to be valid.
// CONTRACT: Returns addrs in some deterministic order.
func (d DepositMsg) GetSigners() []crypto.Address {
	return []crypto.Address{d.Sender}
}

// SettleMsg defines the properties of a settle transaction.
type SettleMsg struct {
	Sender    crypto.Address
	Recipient crypto.Address
	Amount    sdk.Coin
}

var _ sdk.Msg = SettleMsg{}

// ValidateBasic is called by the SDK automatically.
func (msg SettleMsg) ValidateBasic() sdk.Error {
	// amount may be negative
	if msg.Amount.Amount == 0 {
		return ErrInvalidAmount("empty or 0 amount not allowed")
	}
	if msg.Amount.Denom == "" {
		return ErrInvalidAmount("empty denom")
	}
	if err := validateAddress(msg.Sender); err != nil {
		return err
	}
	if err := validateAddress(msg.Recipient); err != nil {
		return err
	}

	if bytes.Equal(msg.Sender, msg.Recipient) {
		return ErrInvalidAddress("sender and recipient have the same address")
	}

	return nil
}

// Type returns the message type.
// Must be alphanumeric or empty.
func (msg SettleMsg) Type() string {
	return SettlementType
}

// Get some property of the Msg.
func (msg SettleMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// GetSignBytes returns the canonical byte representation of the Msg.
func (msg SettleMsg) GetSignBytes() []byte {
	bz, err := cdc.MarshalBinary(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

// GetSigners returns the addrs of signers that must sign.
// CONTRACT: All signatures must be present to be valid.
// CONTRACT: Returns addrs in some deterministic order.
func (msg SettleMsg) GetSigners() []crypto.Address {
	return []crypto.Address{msg.Sender}
}

// WithdrawMsg defines the properties of a withdraw transaction.
type WithdrawMsg struct {
	Sender    crypto.Address
	Recipient crypto.Address
	Operator  crypto.Address
	Amount    sdk.Coin
}

var _ sdk.Msg = WithdrawMsg{}

// ValidateBasic is called by the SDK automatically.
func (msg WithdrawMsg) ValidateBasic() sdk.Error {
	if msg.Amount.Amount <= 0 {
		return ErrInvalidAmount("negative or 0 amount not allowed")
	}
	if msg.Amount.Denom == "" {
		return ErrInvalidAmount("empty denom")
	}
	if err := validateAddress(msg.Sender); err != nil {
		return err
	}
	if err := validateAddress(msg.Recipient); err != nil {
		return err
	}
	if err := validateAddress(msg.Operator); err != nil {
		return err
	}
	if bytes.Equal(msg.Sender, msg.Recipient) {
		return ErrInvalidAddress("sender and recipient have the same address")
	}
	return nil
}

// Type returns the message type.
// Must be alphanumeric or empty.
func (msg WithdrawMsg) Type() string {
	return WithdrawType
}

// Get some property of the Msg.
func (msg WithdrawMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// GetSignBytes returns the canonical byte representation of the Msg.
func (msg WithdrawMsg) GetSignBytes() []byte {
	bz, err := cdc.MarshalBinary(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

// GetSigners returns the addrs of signers that must sign.
// CONTRACT: All signatures must be present to be valid.
// CONTRACT: Returns addrs in some deterministic order.
func (msg WithdrawMsg) GetSigners() []crypto.Address {
	return []crypto.Address{msg.Sender, msg.Operator}
}

// CreateAccountMsg defines the property of a create account transaction.
type CreateAccountMsg struct {
	Creator         crypto.Address
	PubKey          crypto.PubKey
	AccountType     string
	LegalEntityName string
	IsAdmin         bool
}

var _ sdk.Msg = CreateAccountMsg{}

// ValidateBasic is called by the SDK automatically.
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
	if len(strings.TrimSpace(msg.LegalEntityName)) == 0 {
		return ErrInvalidAccount("legal entity name must be non-nil")
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

// Auxiliary functions, might be undocumented

func validateAddress(addr crypto.Address) sdk.Error {
	if addr == nil {
		return ErrInvalidAddress("address is nil")
	}
	if len(addr) != 20 {
		return ErrInvalidAddress("invalid address length")
	}
	return nil
}
