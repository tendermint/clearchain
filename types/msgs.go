package types

import (
	"bytes"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
)

// Defines all the messages (requests) supported by the app

// message types definitions
const (
	DepositType            = "deposit"
	SettlementType         = "settlement"
	WithdrawType           = "withdraw"
	CreateOperatorType     = "createOperator"
	CreateAdminType        = "createAdmin"
	CreateAssetAccountType = "createAsset"
	FreezeOperatorType     = "freezeOperator"
	FreezeAdminType        = "freezeAdmin"
)

const (
	// AddressLength represents the number of bytes that compose addresses.
	AddressLength = 20
)

// DepositMsg defines the properties of an asset transfer
type DepositMsg struct {
	Operator  sdk.Address
	Sender    sdk.Address
	Recipient sdk.Address
	Amount    sdk.Coin
}

// ensure DepositMsg implements the sdk.Msg interface
var _ sdk.Msg = DepositMsg{}
var _ sdk.Msg = (*DepositMsg)(nil)

// ValidateBasic is called by the SDK automatically.
func (msg DepositMsg) ValidateBasic() sdk.Error {
	if msg.Amount.Amount <= 0 {
		return ErrInvalidAmount("negative amount")
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
func (msg DepositMsg) Type() string { return DepositType }

// Get some property of the Msg.
func (msg DepositMsg) Get(key interface{}) (value interface{}) { return nil }

// GetSignBytes returns the canonical byte representation of the Msg.
func (msg DepositMsg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

// GetSigners returns the addrs of signers that must sign.
// CONTRACT: All signatures must be present to be valid.
// CONTRACT: Returns addrs in some deterministic order.
func (msg DepositMsg) GetSigners() []sdk.Address { return []sdk.Address{msg.Operator} }

// SettleMsg defines the properties of a settle transaction.
type SettleMsg struct {
	Operator  sdk.Address
	Sender    sdk.Address
	Recipient sdk.Address
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
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

// GetSigners returns the addrs of signers that must sign.
// CONTRACT: All signatures must be present to be valid.
// CONTRACT: Returns addrs in some deterministic order.
func (msg SettleMsg) GetSigners() []sdk.Address { return []sdk.Address{msg.Operator} }

// WithdrawMsg defines the properties of a withdraw transaction.
type WithdrawMsg struct {
	Operator  sdk.Address
	Sender    sdk.Address
	Recipient sdk.Address
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
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

// GetSigners returns the addrs of signers that must sign.
// CONTRACT: All signatures must be present to be valid.
// CONTRACT: Returns addrs in some deterministic order.
func (msg WithdrawMsg) GetSigners() []sdk.Address { return []sdk.Address{msg.Operator} }

// CreateAssetAccountMsg defines the property of a create user transaction.
type CreateAssetAccountMsg struct {
	Creator sdk.Address
	PubKey  crypto.PubKey
}

var _ sdk.Msg = CreateAssetAccountMsg{}

// ValidateBasic performs basic validation checks and it's
// called by the SDK automatically.
func (msg CreateAssetAccountMsg) ValidateBasic() sdk.Error {
	if err := validateAddress(msg.Creator); err != nil {
		return err
	}
	if msg.PubKey.Empty() {
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
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

// GetSigners returns the addrs of signers that must sign.
// CONTRACT: All signatures must be present to be valid.
// CONTRACT: Returns addrs in some deterministic order.
func (msg CreateAssetAccountMsg) GetSigners() []sdk.Address { return []sdk.Address{msg.Creator} }

// BaseCreateUserMsg defines the properties of a transaction
// that triggers the creation of a new generic user.
// Legal entitiy is inherited from the creator.
type BaseCreateUserMsg struct {
	Creator sdk.Address
	PubKey  crypto.PubKey
}

// ValidateBasic is called by the SDK automatically.
func (msg BaseCreateUserMsg) ValidateBasic() sdk.Error {
	if msg.PubKey.Empty() {
		return ErrInvalidPubKey("pub key is nil")
	}
	if err := validateAddress(msg.Creator); err != nil {
		return err
	}
	if bytes.Equal(msg.Creator, msg.PubKey.Address()) {
		return ErrSelfCreate(fmt.Sprintf("%v", msg.Creator))
	}
	return nil
}

// Type returns the message type.
// Must be alphanumeric or empty.
//func (msg BaseCreateUserMsg) Type() string { return CreateOperatorType }

// Get returns some property of the Msg.
func (msg BaseCreateUserMsg) Get(key interface{}) (value interface{}) { return nil }

// GetSignBytes returns the canonical byte representation of the Msg.
func (msg BaseCreateUserMsg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

// GetSigners returns the addrs of signers that must sign.
// CONTRACT: All signatures must be present to be valid.
// CONTRACT: Returns addrs in some deterministic order.
func (msg BaseCreateUserMsg) GetSigners() []sdk.Address { return []sdk.Address{msg.Creator} }

// CreateOperatorMsg defines the properties of a transaction
// that triggers the creation of a new unprivileged user.
// Legal entitiy is inherited from the creator.
type CreateOperatorMsg struct{ BaseCreateUserMsg }

var _ sdk.Msg = (*CreateOperatorMsg)(nil)

// Type returns the message type.
// Must be alphanumeric or empty.
func (msg CreateOperatorMsg) Type() string { return CreateOperatorType }

// CreateAdminMsg defines the properties of a transaction
// that triggers the cross-entity creation of privileged users
// The message must carry the legal entity.
// Only a clearing house can utilise this endpoint.
type CreateAdminMsg struct {
	BaseCreateUserMsg
	BaseLegalEntity
}

var _ sdk.Msg = (*CreateAdminMsg)(nil)

// ValidateBasic is called by the SDK automatically.
func (msg CreateAdminMsg) ValidateBasic() sdk.Error {
	if err := msg.BaseCreateUserMsg.ValidateBasic(); err != nil {
		return err
	}
	if err := ValidateLegalEntity(msg.BaseLegalEntity); err != nil {
		return ErrInvalidLegalEntity(err.Error())
	}
	return nil
}

// Type returns the message type.
// Must be alphanumeric or empty.
func (msg CreateAdminMsg) Type() string { return CreateAdminType }

// BaseFreezeAccountMsg defines the properties of a transaction
// that freezes user or asset accounts.
type BaseFreezeAccountMsg struct {
	Admin  sdk.Address
	Target sdk.Address
}

// ValidateBasic is called by the SDK automatically.
func (msg BaseFreezeAccountMsg) ValidateBasic() sdk.Error {
	if err := validateAddress(msg.Admin); err != nil {
		return err
	}
	if err := validateAddress(msg.Target); err != nil {
		return err
	}
	if bytes.Equal(msg.Admin, msg.Target) {
		return ErrSelfFreeze(fmt.Sprintf("%v", msg.Admin))
	}
	return nil
}

// Type returns the message type.
// Must be alphanumeric or empty.

// Get returns some property of the Msg.
func (msg BaseFreezeAccountMsg) Get(key interface{}) (value interface{}) { return nil }

// GetSignBytes returns the canonical byte representation of the Msg.
func (msg BaseFreezeAccountMsg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

// GetSigners returns the addrs of signers that must sign.
// CONTRACT: All signatures must be present to be valid.
// CONTRACT: Returns addrs in some deterministic order.
func (msg BaseFreezeAccountMsg) GetSigners() []sdk.Address { return []sdk.Address{msg.Admin} }

// FreezeOperatorMsg defines the properties of a transaction
// that freezes an operator. Admin accounts can freeze their
// own legal entity's operators.
type FreezeOperatorMsg struct{ BaseFreezeAccountMsg }

var _ sdk.Msg = (*FreezeOperatorMsg)(nil)

// Type returns the message type.
// Must be alphanumeric or empty.
func (msg FreezeOperatorMsg) Type() string { return FreezeOperatorType }

// FreezeAdminMsg defines the properties of a transaction
// that freezes an admin. Only clearing house Admin accounts
// can freeze other Admin accounts, regardless of the entity
// that own them.
type FreezeAdminMsg struct{ BaseFreezeAccountMsg }

var _ sdk.Msg = (*FreezeAdminMsg)(nil)

// Type returns the message type.
// Must be alphanumeric or empty.
func (msg FreezeAdminMsg) Type() string { return FreezeAdminType }

/* Constructors */

// NewCreateAdminMsg creates a new CreateAdminMsg.
func NewCreateAdminMsg(creator sdk.Address, pubkey crypto.PubKey,
	entityName, entityType string) (msg CreateAdminMsg) {
	msg.BaseCreateUserMsg.Creator = creator
	msg.BaseCreateUserMsg.PubKey = pubkey
	msg.BaseLegalEntity.EntityName = entityName
	msg.BaseLegalEntity.EntityType = entityType
	return
}

// NewCreateOperatorMsg creates a new CreateOperatorMsg.
func NewCreateOperatorMsg(creator sdk.Address, pubkey crypto.PubKey) (msg CreateOperatorMsg) {
	msg.BaseCreateUserMsg.Creator = creator
	msg.BaseCreateUserMsg.PubKey = pubkey
	return
}

// NewCreateAssetAccountMsg creates a new CreateAssetAccountMsg.
func NewCreateAssetAccountMsg(creator sdk.Address, pubkey crypto.PubKey) (msg CreateAssetAccountMsg) {
	msg.Creator = creator
	msg.PubKey = pubkey
	return
}

/* Auxiliary functions, could be undocumented */

func validateAddress(addr sdk.Address) sdk.Error {
	if addr == nil {
		return ErrInvalidAddress("address is nil")
	}
	if len(addr) != AddressLength {
		return ErrInvalidAddress("invalid address length")
	}
	return nil
}
