package types

import (
	uuid "github.com/satori/go.uuid"
	common "github.com/tendermint/go-common"
	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"
	tmsp "github.com/tendermint/tmsp/types"
)

const (
	// TxTypeCreateLegalEntity defines CreateLegalEntityTx's code
	TxTypeCreateLegalEntity = byte(0x03)
)

// CreateLegalEntityTx defines the attributes of a legal entity create.
type CreateLegalEntityTx struct {
	Address   []byte           `json:"address"`   // Hash of the user's PubKey
	EntityID  string           `json:"entity_id"` // ID of the new legal entity
	Type      byte             `json:"type"`      // Mandatory
	Name      string           `json:"name"`      // Could be empty
	Signature crypto.Signature `json:"signature"`
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
func (tx *CreateLegalEntityTx) ValidateBasic() tmsp.Result {
	if len(tx.Address) != 20 {
		return tmsp.ErrBaseInvalidInput.AppendLog("Invalid address length")
	}
	if tx.Signature == nil {
		return tmsp.ErrBaseInvalidSignature.AppendLog("The transaction must be signed")
	}
	if _, err := uuid.FromString(tx.EntityID); err != nil {
		return tmsp.ErrBaseInvalidInput.AppendLog(common.Fmt("Invalid entity_id: %s", err))
	}
	if !IsValidEntityType(tx.Type) {
		return tmsp.ErrBaseInvalidInput.AppendLog(common.Fmt("Invalid Type: %s", tx.Type))
	}
	return tmsp.OK
}

func (tx *CreateLegalEntityTx) String() string {
	return common.Fmt("CreateLegalEntityTx{%x,%q,%x,%s}", tx.Address, tx.EntityID, tx.Type, tx.Name)
}
