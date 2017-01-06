package types

import (
	uuid "github.com/satori/go.uuid"
	common "github.com/tendermint/go-common"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire"
	tmsp "github.com/tendermint/tmsp/types"
)

const (
	// TxTypeQueryAccount defines AccountQueryTx's code
	TxTypeQueryAccount = byte(0x11)
	// TxTypeQueryAccountIndex defines AccountIndexQueryTx's code
	TxTypeQueryAccountIndex = byte(0x12)
	TxTypeQueryLegalEntityIndex = byte(0x13)
	TxTypeLegalEntity = byte(0x14)
)

// AccountQueryTx defines the attribute of an accounts query
type AccountQueryTx struct {
	Accounts  []string         `json:"accounts"`
	Address   []byte           `json:"address"` // Hash of the user's PubKey
	Signature crypto.Signature `json:"signature"`
}

func (tx AccountQueryTx) String() string {
	return common.Fmt("AccountQueryTx{%x %s %v}", tx.Address, tx.Accounts, tx.Signature)
}

// TxType returns the byte type of TransferTx.
func (tx AccountQueryTx) TxType() byte {
	return TxTypeQueryAccount
}

// ValidateBasic performs basic validation.
func (tx AccountQueryTx) ValidateBasic() tmsp.Result {
	if len(tx.Address) != 20 {
		return tmsp.ErrBaseInvalidInput.AppendLog(common.Fmt("Invalid address length: %v", len(tx.Address)))
	}
	if tx.Signature == nil {
		return tmsp.ErrBaseInvalidSignature.AppendLog("The query must be signed")
	}
	if len(tx.Accounts) == 0 {
		return tmsp.ErrBaseInvalidInput.AppendLog("Accounts must not be empty")
	}
	for _, accountID := range tx.Accounts {
		if _, err := uuid.FromString(accountID); err != nil {
			return tmsp.ErrBaseInvalidInput.AppendLog(common.Fmt("Invalid account_id %q: %s", accountID, err))
		}
	}
	return tmsp.OK
}

// SignBytes generates a byte-to-byte signature.
func (tx AccountQueryTx) SignBytes(chainID string) []byte {
	signBytes := wire.BinaryBytes(chainID)
	sig := tx.Signature
	tx.Signature = nil
	signBytes = append(signBytes, wire.BinaryBytes(tx)...)
	tx.Signature = sig
	return signBytes
}

// SignTx signs the transaction if its address and the privateKey's one match.
func (tx *AccountQueryTx) SignTx(privateKey crypto.PrivKey, chainID string) error {
	sig, err := SignTx(tx.SignBytes(chainID), tx.Address, privateKey)
	if err != nil {
		return err
	}
	tx.Signature = sig
	return nil
}

//--------------------------------------------

// AccountIndexQueryTx defines the attribute of a wildcard query
type AccountIndexQueryTx struct {
	Address   []byte           `json:"address"` // Hash of the user's PubKey
	Signature crypto.Signature `json:"signature"`
}

func (tx AccountIndexQueryTx) String() string {
	return common.Fmt("AccountIndexQueryTx{%x %v}", tx.Address, tx.Signature)
}

// TxType returns the byte type of TransferTx.
func (tx AccountIndexQueryTx) TxType() byte {
	return TxTypeQueryAccountIndex
}

// ValidateBasic performs basic validation.
func (tx AccountIndexQueryTx) ValidateBasic() tmsp.Result {
	if len(tx.Address) != 20 {
		return tmsp.ErrBaseInvalidInput.AppendLog(common.Fmt("Invalid address length: %v", len(tx.Address)))
	}
	if tx.Signature == nil {
		return tmsp.ErrBaseInvalidSignature.AppendLog("The query must be signed")
	}
	return tmsp.OK
}

// SignBytes generates a byte-to-byte signature.
func (tx AccountIndexQueryTx) SignBytes(chainID string) []byte {
	signBytes := wire.BinaryBytes(chainID)
	sig := tx.Signature
	tx.Signature = nil
	signBytes = append(signBytes, wire.BinaryBytes(tx)...)
	tx.Signature = sig
	return signBytes
}

// SignTx signs the transaction if its address and the privateKey's one match.
func (tx *AccountIndexQueryTx) SignTx(privateKey crypto.PrivKey, chainID string) error {
	sig, err := SignTx(tx.SignBytes(chainID), tx.Address, privateKey)
	if err != nil {
		return err
	}
	tx.Signature = sig
	return nil
}

type LegalEntityIndexQueryTx struct {
	Address   []byte           `json:"address"` // Hash of the user's PubKey
	Signature crypto.Signature `json:"signature"`
}

func (tx LegalEntityIndexQueryTx) String() string {
	return common.Fmt("LegalEntityIndexQueryTx{%x %v}", tx.Address, tx.Signature)
}

func (tx LegalEntityIndexQueryTx) TxType() byte {
	return TxTypeQueryLegalEntityIndex
}
func (tx LegalEntityIndexQueryTx) SignBytes(chainID string) []byte {
	signBytes := wire.BinaryBytes(chainID)
	sig := tx.Signature
	tx.Signature = nil
	signBytes = append(signBytes, wire.BinaryBytes(tx)...)
	tx.Signature = sig
	return signBytes
}

func (tx *LegalEntityIndexQueryTx) SignTx(privateKey crypto.PrivKey, chainID string) error {
	sig, err := SignTx(tx.SignBytes(chainID), tx.Address, privateKey)
	if err != nil {
		return err
	}
	tx.Signature = sig
	return nil
}

func (tx LegalEntityIndexQueryTx) ValidateBasic() tmsp.Result {
	if len(tx.Address) != 20 {
		return tmsp.ErrBaseInvalidInput.AppendLog(common.Fmt("Invalid address length: %v", len(tx.Address)))
	}
	if tx.Signature == nil {
		return tmsp.ErrBaseInvalidSignature.AppendLog("The query must be signed")
	}

	return tmsp.OK
}

type LegalEntityQueryTx struct {
	Ids  []string         `json:"ids"`
	Address   []byte           `json:"address"` // Hash of the user's PubKey
	Signature crypto.Signature `json:"signature"`
}

func (tx LegalEntityQueryTx) String() string {
	return common.Fmt("LegalEntityQueryTx{%x %s %v}", tx.Address, tx.Ids, tx.Signature)
}

// TxType returns the byte type of TransferTx.
func (tx LegalEntityQueryTx) TxType() byte {
	return TxTypeLegalEntity
}

// ValidateBasic performs basic validation.
func (tx LegalEntityQueryTx) ValidateBasic() tmsp.Result {
	if len(tx.Address) != 20 {
		return tmsp.ErrBaseInvalidInput.AppendLog(common.Fmt("Invalid address length: %v", len(tx.Address)))
	}
	if tx.Signature == nil {
		return tmsp.ErrBaseInvalidSignature.AppendLog("The query must be signed")
	}
	if len(tx.Ids) == 0 {
		return tmsp.ErrBaseInvalidInput.AppendLog("Ids must not be empty")
	}
	for _, accountID := range tx.Ids {
		if _, err := uuid.FromString(accountID); err != nil {
			return tmsp.ErrBaseInvalidInput.AppendLog(common.Fmt("Invalid id %q: %s", accountID, err))
		}
	}
	return tmsp.OK
}

// SignBytes generates a byte-to-byte signature.
func (tx LegalEntityQueryTx) SignBytes(chainID string) []byte {
	signBytes := wire.BinaryBytes(chainID)
	sig := tx.Signature
	tx.Signature = nil
	signBytes = append(signBytes, wire.BinaryBytes(tx)...)
	tx.Signature = sig
	return signBytes
}

// SignTx signs the transaction if its address and the privateKey's one match.
func (tx *LegalEntityQueryTx) SignTx(privateKey crypto.PrivKey, chainID string) error {
	sig, err := SignTx(tx.SignBytes(chainID), tx.Address, privateKey)
	if err != nil {
		return err
	}
	tx.Signature = sig
	return nil
}
