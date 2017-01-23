package types

import (
	abci "github.com/tendermint/abci/types"
	common "github.com/tendermint/go-common"
	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"
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

// SignTx signs the transaction if its address and the privateKey's one match.
func (tx *CreateUserTx) SignTx(privateKey crypto.PrivKey, chainID string) error {
	sig, err := SignTx(tx.SignBytes(chainID), tx.Address, privateKey)
	if err != nil {
		return err
	}
	tx.Signature = sig
	return nil
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
func (tx *CreateUserTx) ValidateBasic() abci.Result {
	if len(tx.Address) != 20 {
		return abci.ErrBaseInvalidInput.AppendLog("Invalid address length")
	}
	if tx.PubKey == nil {
		return abci.ErrBaseInvalidPubKey.AppendLog("PubKey can't be nil")
	}
	if len(tx.Name) == 0 {
		return abci.ErrBaseInvalidInput.AppendLog("Name cannot be empty")
	}
	if tx.Signature == nil {
		return abci.ErrBaseInvalidSignature.AppendLog("The transaction must be signed")
	}
	return abci.OK
}

func (tx *CreateUserTx) String() string {
	return common.Fmt(
		"CreateUserTx{%x,%s,%t,%v}", tx.Address, tx.Name, tx.CanCreate, tx.PubKey)
}
