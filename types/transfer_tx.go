package types

import (
	"bytes"

	"github.com/satori/go.uuid"
	"github.com/tendermint/go-common"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire"
	tmsp "github.com/tendermint/tmsp/types"
)

const (
	// TxTypeTransfer defines TrasferTx's code
	TxTypeTransfer = byte(0x01)
)

// TransferTx defines the attributes of transfer transaction
type TransferTx struct {
	Sender         TxTransferSender          `json:"sender"`
	Recipient      TxTransferRecipient       `json:"recipient"`
	CounterSigners []TxTransferCounterSigner `json:"counter_signers"`
}

// TxType returns the byte type of TransferTx
func (tx *TransferTx) TxType() byte {
	return TxTypeTransfer
}

// SignBytes generates a byte-to-byte signature
func (tx *TransferTx) SignBytes(chainID string) []byte {
	signBytes := wire.BinaryBytes(chainID)
	senderSig := tx.Sender.Signature
	tx.Sender.Signature = nil
	sigz := make([]crypto.Signature, len(tx.CounterSigners))
	for i, sender := range tx.CounterSigners {
		sigz[i] = sender.Signature
		tx.CounterSigners[i].Signature = nil
	}
	signBytes = append(signBytes, wire.BinaryBytes(tx)...)
	tx.Sender.Signature = senderSig
	for i := range tx.CounterSigners {
		tx.CounterSigners[i].Signature = sigz[i]
	}
	return signBytes
}

// SetSignature sets account's signature to the relevant TxInputTransfer
func (tx *TransferTx) SetSignature(addr []byte, sig crypto.Signature) bool {
	if bytes.Equal(tx.Sender.Address, addr) {
		tx.Sender.Signature = sig
		return true
	}
	for i, input := range tx.CounterSigners {
		if bytes.Equal(input.Address, addr) {
			tx.CounterSigners[i].Signature = sig
			return true
		}
	}
	return false
}

func (tx *TransferTx) String() string {
	return common.Fmt("TransferTx{%v->%v, %v}", tx.Sender, tx.Recipient, tx.CounterSigners)
}

// ValidateBasic validates Tx basic structure.
func (tx *TransferTx) ValidateBasic() (res tmsp.Result) {
	if res := tx.Sender.ValidateBasic(); res.IsErr() {
		// Check the sender
		return res
	}
	// Check the recipient
	if res := tx.Recipient.ValidateBasic(); res.IsErr() {
		return res
	}
	// Check the countersigners
	for _, in := range tx.CounterSigners {
		if res := in.ValidateBasic(); res.IsErr() {
			return res
		}
	}
	return tmsp.OK
}

// SignTx signs the transaction if its address and the privateKey's one match.
func (tx *TransferTx) SignTx(privateKey crypto.PrivKey, chainID string) error {
	for i := 0; i < len(tx.CounterSigners); i++ {
		err := tx.CounterSigners[i].SignTx(privateKey, chainID)
		if err != nil {
			return err
		}
	}

	sig, err := SignTx(tx.SignBytes(chainID), tx.Sender.Address, privateKey)
	if err != nil {
		return err
	}
	tx.Sender.Signature = sig

	return nil
}

// TxTransferSender defines the attributes of a transfer's sender
type TxTransferSender struct {
	Address   []byte           `json:"address"`    // Hash of the user's PubKey
	AccountID string           `json:"account_id"` // Sender's Account ID
	Amount    int64            `json:"amount"`
	Currency  string           `json:"currency"` //3-letter ISO 4217 code or ""
	Sequence  int              `json:"sequence"` // Must be 1 greater than the last committed TxInput or 0 if the account is just a countersigner
	Signature crypto.Signature `json:"signature"`
}

// ValidateBasic performs basic validation on a TxInputTransfer
func (t TxTransferSender) ValidateBasic() tmsp.Result {
	if len(t.Address) != 20 {
		return tmsp.ErrBaseInvalidInput.AppendLog("Invalid address length")
	}
	if t.Signature == nil {
		return tmsp.ErrBaseInvalidSignature.AppendLog("TxTransferSender transaction must be signed")
	}
	if _, err := uuid.FromString(t.AccountID); err != nil {
		return tmsp.ErrBaseInvalidInput.AppendLog(common.Fmt("Invalid account_id: %s", err))
	}
	if t.Amount <= 0 {
		return tmsp.ErrBaseInvalidInput.AppendLog(common.Fmt("Amount must be non-negative: %q", t.Amount))
	}
	if len(t.Currency) != 3 {
		return tmsp.ErrBaseInvalidInput.AppendLog(common.Fmt("Currency must be either empty or a 3-letter ISO 4217 compliant code: %q", t.Currency))
	}
	// Validate currency and amount
	currency, ok := Currencies[t.Currency]
	if !ok {
		return tmsp.ErrBaseInvalidInput.AppendLog(common.Fmt("Unsupported currency: %q", t.Currency))
	}
	// Validate the amount against the currency
	if !currency.ValidateAmount(t.Amount) {
		return tmsp.ErrBaseInvalidInput.AppendLog(
			common.Fmt("Invalid amount %d for currency %s", t.Amount, currency.Symbol()))
	}
	if t.Sequence <= 0 {
		return tmsp.ErrBaseInvalidSequence.AppendLog(common.Fmt("Sequence must be greater than 0"))
	}
	return tmsp.OK
}

// String returns a string representation of TxInputTransfer
func (t TxTransferSender) String() string {
	return common.Fmt("TxTransferSender{%x,%v,%v,%v, %v}", t.Address, t.Amount, t.Currency, t.Sequence, t.Signature)
}

// SignBytes generates a byte-to-byte signature.
func (tx TxTransferSender) SignBytes(chainID string) []byte {
	signBytes := wire.BinaryBytes(chainID)
	sig := tx.Signature
	tx.Signature = nil
	signBytes = append(signBytes, wire.BinaryBytes(tx)...)
	tx.Signature = sig
	return signBytes
}

// SignTx signs the transaction if its address and the privateKey's one match.
func (tx *TxTransferSender) SignTx(privateKey crypto.PrivKey, chainID string) error {

	sig, err := SignTx(tx.SignBytes(chainID), tx.Address, privateKey)
	if err != nil {
		return err
	}
	tx.Signature = sig
	return nil
}

//-----------------------------------------------------------------------------

// TxTransferRecipient defines the attributes of a transfer's recipient
type TxTransferRecipient struct {
	AccountID string `json:"account_id"` // Recipient's Account ID
}

// ValidateBasic performs basic validation on a TxTransferRecipient
func (t TxTransferRecipient) ValidateBasic() tmsp.Result {
	if _, err := uuid.FromString(t.AccountID); err != nil {
		return tmsp.ErrBaseInvalidOutput.AppendLog(common.Fmt("Invalid account_id: %s", err))
	}
	return tmsp.OK
}

// String returns a string representation of TxInputTransfer
func (t TxTransferRecipient) String() string {
	return common.Fmt("TxTransferRecipient{%s}", t.AccountID)
}

// TxTransferCounterSigner defines the attributes of a transfer's counter signer
type TxTransferCounterSigner struct {
	Address   []byte           `json:"address"` // Hash of the user's PubKey
	Signature crypto.Signature `json:"signature"`
}

// ValidateBasic performs basic validation on a TxTransferCounterSigner
func (t TxTransferCounterSigner) ValidateBasic() tmsp.Result {
	if len(t.Address) != 20 {
		return tmsp.ErrBaseInvalidInput.AppendLog("Invalid address length")
	}
	if t.Signature == nil {
		return tmsp.ErrBaseInvalidSignature.AppendLog("TxTransferCounterSigner transaction must be signed")
	}
	return tmsp.OK
}

// String returns a string representation of TxTransferCounterSigner
func (t TxTransferCounterSigner) String() string {
	return common.Fmt("TxTransferCounterSigner{%x,%v}", t.Address, t.Signature)
}

// SignBytes generates a byte-to-byte signature.
func (tx TxTransferCounterSigner) SignBytes(chainID string) []byte {
	signBytes := wire.BinaryBytes(chainID)
	sig := tx.Signature
	tx.Signature = nil
	signBytes = append(signBytes, wire.BinaryBytes(tx)...)
	tx.Signature = sig
	return signBytes
}

// SignTx signs the transaction if its address and the privateKey's one match.
func (tx *TxTransferCounterSigner) SignTx(privateKey crypto.PrivKey, chainID string) error {

	sig, err := SignTx(tx.SignBytes(chainID), tx.Address, privateKey)
	if err != nil {
		return err
	}
	tx.Signature = sig
	return nil
}
