package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
)

// Defines all the messages (requests) supported by the app

// message types definitions
const (
	DepositType            = "deposit"
	SettlementType         = "settlement"
	WithdrawType           = "withdraw"
	CreateUserAccountType  = "createUser"
	CreateAssetAccountType = "createAsset"
)

// DepositMsg defines the properties of an asset transfer
type DepositMsg struct {
	Operator  crypto.Address
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
	if err := validateAddress(d.Operator); err != nil {
		return err
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
func (d DepositMsg) Type() string { return DepositType }

// Get some property of the Msg.
func (d DepositMsg) Get(key interface{}) (value interface{}) { return nil }

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
func (d DepositMsg) GetSigners() []crypto.Address { return []crypto.Address{d.Operator} }

// SettleMsg defines the properties of a settle transaction.
type SettleMsg struct {
	Operator  crypto.Address
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
	if err := validateAddress(msg.Operator); err != nil {
		return err
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
func (msg SettleMsg) Type() string { return SettlementType }

// Get some property of the Msg.
func (msg SettleMsg) Get(key interface{}) (value interface{}) { return nil }

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
func (msg SettleMsg) GetSigners() []crypto.Address { return []crypto.Address{msg.Operator} }

// WithdrawMsg defines the properties of a withdraw transaction.
type WithdrawMsg struct {
	Operator  crypto.Address
	Sender    crypto.Address
	Recipient crypto.Address
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
	if err := validateAddress(msg.Operator); err != nil {
		return err
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
func (msg WithdrawMsg) Type() string { return WithdrawType }

// Get some property of the Msg.
func (msg WithdrawMsg) Get(key interface{}) (value interface{}) { return nil }

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
	return []crypto.Address{msg.Operator}
}

// CreateUserAccountMsg defines the property of a create user transaction.
type CreateUserAccountMsg struct {
	Creator         crypto.Address
	PubKey          crypto.PubKey
	LegalEntityType string
	LegalEntityName string
	IsAdmin         bool
}

var _ sdk.Msg = CreateUserAccountMsg{}

// ValidateBasic performs basic validation checks and it's
// called by the SDK automatically.
func (msg CreateUserAccountMsg) ValidateBasic() sdk.Error {
	if err := validateAddress(msg.Creator); err != nil {
		return err
	}
	if msg.PubKey == nil {
		return ErrInvalidPubKey("pub key is nil")
	}
	if bytes.Equal(msg.Creator, msg.PubKey.Address()) {
		return ErrInvalidPubKey("creator and new account have the same address")
	}
	if err := validateEntity(msg.LegalEntityName, msg.LegalEntityType); err != nil {
		return ErrInvalidAccount(err.Error())
	}
	return nil
}

// Type returns the message type.
// Must be alphanumeric or empty.
func (msg CreateUserAccountMsg) Type() string { return CreateUserAccountType }

// Get returns some property of the Msg.
func (msg CreateUserAccountMsg) Get(key interface{}) (value interface{}) { return nil }

// GetSignBytes returns the canonical byte representation of the Msg.
func (msg CreateUserAccountMsg) GetSignBytes() []byte {
	bz, err := cdc.MarshalBinary(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

// GetSigners returns the addrs of signers that must sign.
// CONTRACT: All signatures must be present to be valid.
// CONTRACT: Returns addrs in some deterministic order.
func (msg CreateUserAccountMsg) GetSigners() []crypto.Address { return []crypto.Address{msg.Creator} }

// CreateAssetAccountMsg defines the property of a create user transaction.
type CreateAssetAccountMsg struct {
	Creator crypto.Address
	PubKey  crypto.PubKey
}

var _ sdk.Msg = CreateAssetAccountMsg{}

// Auxiliary functions, might be undocumented

// ValidateBasic performs basic validation checks and it's
// called by the SDK automatically.
func (msg CreateAssetAccountMsg) ValidateBasic() sdk.Error {
	if err := validateAddress(msg.Creator); err != nil {
		return err
	}
	if msg.PubKey == nil {
		return ErrInvalidPubKey("pub key is nil")
	}
	if bytes.Equal(msg.Creator, msg.PubKey.Address()) {
		return ErrInvalidPubKey("creator and new account have the same address")
	}
	return nil
}

// Type returns the message type.
// Must be alphanumeric or empty.
func (msg CreateAssetAccountMsg) Type() string { return CreateAssetAccountType }

// Get returns some property of the Msg.
func (msg CreateAssetAccountMsg) Get(key interface{}) (value interface{}) { return nil }

// GetSignBytes returns the canonical byte representation of the Msg.
func (msg CreateAssetAccountMsg) GetSignBytes() []byte {
	bz, err := cdc.MarshalBinary(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

// GetSigners returns the addrs of signers that must sign.
// CONTRACT: All signatures must be present to be valid.
// CONTRACT: Returns addrs in some deterministic order.
func (msg CreateAssetAccountMsg) GetSigners() []crypto.Address { return []crypto.Address{msg.Creator} }

func validateAddress(addr crypto.Address) sdk.Error {
	if addr == nil {
		return ErrInvalidAddress("address is nil")
	}
	if len(addr) != 20 {
		return ErrInvalidAddress("invalid address length")
	}
	return nil
}
