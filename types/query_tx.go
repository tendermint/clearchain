package types

import (
	uuid "github.com/satori/go.uuid"
	abci "github.com/tendermint/abci/types"
	common "github.com/tendermint/go-common"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire"
)

const (
	TxTypeQueryBase    = byte(0x11)
	TxTypeQueryObjects = byte(0x12)
)

// BaseQueryTx defines the attribute of an accounts query
type BaseQueryTx struct {
	Address   []byte           `json:"address"` // Hash of the user's PubKey
	Signature crypto.Signature `json:"signature"`
}

// ObjectsQueryTx defines the attribute of an accounts query
type ObjectsQueryTx struct {
	BaseQueryTx
	Ids []string `json:"ids"`
}

//--------------------------------------------

// TxType returns the byte type of TransferTx.
func (tx BaseQueryTx) TxType() byte {
	return TxTypeQueryBase
}

// SignBytes generates a byte-to-byte signature.
func (tx BaseQueryTx) SignBytes(chainID string) []byte {
	signBytes := wire.BinaryBytes(chainID)
	sig := tx.Signature
	tx.Signature = nil
	signBytes = append(signBytes, wire.BinaryBytes(tx)...)
	tx.Signature = sig
	return signBytes
}

// ValidateBasic performs basic validation.
func (tx BaseQueryTx) ValidateBasic() abci.Result {
	if len(tx.Address) != 20 {
		return abci.ErrBaseInvalidInput.AppendLog(common.Fmt("Invalid address length: %v", len(tx.Address)))
	}
	if tx.Signature == nil {
		return abci.ErrBaseInvalidSignature.AppendLog("The query must be signed")
	}
	return abci.OK
}

// SignTx signs the transaction if its address and the privateKey's one match.
func (tx *BaseQueryTx) SignTx(privateKey crypto.PrivKey, chainID string) error {
	sig, err := SignTx(tx.SignBytes(chainID), tx.Address, privateKey)
	if err != nil {
		return err
	}
	tx.Signature = sig
	return nil
}

//--------------------------------------------

// TxType returns the byte type of TransferTx.
func (tx ObjectsQueryTx) TxType() byte {
	return TxTypeQueryObjects
}

func (tx ObjectsQueryTx) String() string {
	return common.Fmt("ObjectsQueryTx{%x %s %v}", tx.Address, tx.Ids, tx.Signature)
}

// ValidateBasic performs basic validation.
func (tx ObjectsQueryTx) ValidateBasic() abci.Result {
	if res := tx.BaseQueryTx.ValidateBasic(); res.IsErr() {
		return res
	}
	if len(tx.Ids) == 0 {
		return abci.ErrBaseInvalidInput.AppendLog("Ids must not be empty")
	}
	for _, id := range tx.Ids {
		if _, err := uuid.FromString(id); err != nil {
			return abci.ErrBaseInvalidInput.AppendLog(common.Fmt("Invalid UUID %q: %s", id, err))
		}
	}
	return abci.OK
}

// SignBytes generates a byte-to-byte signature.
func (tx ObjectsQueryTx) SignBytes(chainID string) []byte {
	signBytes := wire.BinaryBytes(chainID)
	sig := tx.Signature
	tx.Signature = nil
	signBytes = append(signBytes, wire.BinaryBytes(tx)...)
	tx.Signature = sig
	return signBytes
}
