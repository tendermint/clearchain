package types

import (
	uuid "github.com/satori/go.uuid"
	abci "github.com/tendermint/abci/types"
	common "github.com/tendermint/go-common"
	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"
)

const (
	// TxTypeCreateLegalEntity defines CreateLegalEntityTx's code
	TxTypeCreateLegalEntity = byte(0x03)
)

// CreateLegalEntityTx defines the attributes of a legal entity create.
type CreateLegalEntityTx struct {
	Address   []byte           `json:"address"`   // Hash of the user's PubKey
	EntityID  string           `json:"entity_id"` // ID of the new legal entity
	ParentID  string           `json:"parent_id"` // ID of the new legal entity's parent
	Type      byte             `json:"type"`      // Mandatory
	Name      string           `json:"name"`      // Could be empty
	Signature crypto.Signature `json:"signature"`
}

// SignTx signs the transaction if its address and the privateKey's one match.
func (tx *CreateLegalEntityTx) SignTx(privateKey crypto.PrivKey, chainID string) error {
	sig, err := SignTx(tx.SignBytes(chainID), tx.Address, privateKey)
	if err != nil {
		return err
	}
	tx.Signature = sig
	return nil
}

// TxType returns the byte type of CreateLegalEntityTx
func (tx *CreateLegalEntityTx) TxType() byte {
	return TxTypeCreateLegalEntity
}

// SignBytes generates a byte-to-byte signature
func (tx *CreateLegalEntityTx) SignBytes(chainID string) []byte {
	signBytes := wire.BinaryBytes(chainID)
	sig := tx.Signature
	tx.Signature = nil
	signBytes = append(signBytes, wire.BinaryBytes(tx)...)
	tx.Signature = sig
	return signBytes
}

// ValidateBasic performs basic validation on the Tx.
func (tx *CreateLegalEntityTx) ValidateBasic() abci.Result {
	if len(tx.Address) != 20 {
		return abci.ErrBaseInvalidInput.AppendLog("Invalid address length")
	}
	if tx.Signature == nil {
		return abci.ErrBaseInvalidSignature.AppendLog("The transaction must be signed")
	}
	if _, err := uuid.FromString(tx.EntityID); err != nil {
		return abci.ErrBaseInvalidInput.AppendLog(common.Fmt("Invalid entity_id: %s", err))
	}
	if _, err := uuid.FromString(tx.ParentID); err != nil {
		return abci.ErrBaseInvalidInput.AppendLog(common.Fmt("Invalid parent_id: %s", err))
	}

	if !IsValidEntityType(tx.Type) {
		return abci.ErrBaseInvalidInput.AppendLog(common.Fmt("Invalid Type: %s", tx.Type))
	}
	return abci.OK
}

func (tx *CreateLegalEntityTx) String() string {
	return common.Fmt("CreateLegalEntityTx{%x,%q,%x,%s,%v}", tx.Address, tx.EntityID, tx.Type, tx.Name, tx.ParentID)
}
