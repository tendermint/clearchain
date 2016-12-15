package types

import (
	common "github.com/tendermint/go-common"
	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"
	tmsp "github.com/tendermint/tmsp/types"
)

const (
	// TxTypeCreateUser defines CreateUserTx's code
	TxTypeCreateUser = byte(0x04)
)

// CreateUserTx defines the attributes of a user create.
type CreateUserTx struct {
	Address   []byte           `json:"address"`    // Hash of the user's PubKey
	Name      string           `json:"name"`       // Human-readable identifier, mandatory
	PubKey    crypto.PubKey    `json:"pub_key"`    // New user's public key
	CanCreate bool             `json:"can_create"` // Whether the user is a super user or not
	Signature crypto.Signature `json:"signature"`
}

func (tx *CreateUserTx) SignTx(privateKey crypto.PrivKey, chainID string) {
	tx.Signature = privateKey.Sign(tx.SignBytes(chainID))
}

// TxType returns the byte type of CreateUserTx
func (tx *CreateUserTx) TxType() byte {
	return TxTypeCreateUser
}

// SignBytes generates a byte-to-byte signature
func (tx *CreateUserTx) SignBytes(chainID string) []byte {
	signBytes := wire.BinaryBytes(chainID)
	sig := tx.Signature
	tx.Signature = nil
	signBytes = append(signBytes, wire.BinaryBytes(tx)...)
	tx.Signature = sig
	return signBytes
}

// ValidateBasic performs basic validation on the Tx.
func (tx *CreateUserTx) ValidateBasic() tmsp.Result {
	if len(tx.Address) != 20 {
		return tmsp.ErrBaseInvalidInput.AppendLog("Invalid address length")
	}
	if tx.PubKey == nil {
		return tmsp.ErrBaseInvalidPubKey.AppendLog("PubKey can't be nil")
	}
	if len(tx.Name) == 0 {
		return tmsp.ErrBaseInvalidInput.AppendLog("Name cannot be empty")
	}
	if tx.Signature == nil {
		return tmsp.ErrBaseInvalidSignature.AppendLog("The transaction must be signed")
	}
	return tmsp.OK
}

func (tx *CreateUserTx) String() string {
	return common.Fmt(
		"CreateUserTx{%x,%s,%t,%v}", tx.Address, tx.Name, tx.CanCreate, tx.PubKey)
}
