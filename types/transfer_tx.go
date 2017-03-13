package types

import (
	"bytes"

	"github.com/satori/go.uuid"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-common"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire"
)

const (
	// TxTypeTransfer defines TrasferTx's code
	TxTypeTransfer = byte(0x01)
)

// TransferTx defines the attributes of transfer transaction
type TransferTx struct {
	Committer      TxTransferCommitter       `json:"committer"`
	Sender         TxTransferSender          `json:"sender"`
	Recipient      TxTransferRecipient       `json:"recipient"`
	CounterSigners []TxTransferCounterSigner `json:"counter_signers"`
}

// TxTransferCommitter defines the attributes of a transfer's sender
type TxTransferCommitter struct {
	Address   []byte           `json:"address"` // Hash of the user's PubKey
	Signature crypto.Signature `json:"signature"`
}

// TxTransferSender defines the attributes of a transfer's sender
type TxTransferSender struct {
	AccountID string `json:"account_id"` // Sender's Account ID
	Amount    int64  `json:"amount"`
	Currency  string `json:"currency"` //3-letter ISO 4217 code or ""
	Sequence  int    `json:"sequence"` // Must be 1 greater than the last committed TxInput or 0 if the account is just a countersigner
}

// TxTransferRecipient defines the attributes of a transfer's recipient
type TxTransferRecipient struct {
	AccountID string `json:"account_id"` // Recipient's Account ID
}

// TxTransferCounterSigner defines the attributes of a transfer's counter signer
type TxTransferCounterSigner struct {
	Address   []byte           `json:"address"` // Hash of the user's PubKey
	Signature crypto.Signature `json:"signature"`
}

//-----------------------------------------------------------------------------

// TxType returns the byte type of TransferTx
func (tx *TransferTx) TxType() byte {
	return TxTypeTransfer
}

// SignBytes generates a byte-to-byte signature
func (tx *TransferTx) SignBytes(chainID string) []byte {
	signBytes := wire.BinaryBytes(chainID)
	commiterSig := tx.Committer.Signature
	tx.Committer.Signature = nil
	sigz := make([]crypto.Signature, len(tx.CounterSigners))
	for i, counterSig := range tx.CounterSigners {
		sigz[i] = counterSig.Signature
		tx.CounterSigners[i].Signature = nil
	}
	signBytes = append(signBytes, wire.BinaryBytes(tx)...)
	tx.Committer.Signature = commiterSig
	for i := range tx.CounterSigners {
		tx.CounterSigners[i].Signature = sigz[i]
	}
	return signBytes
}

// SetSignature sets account's signature to the relevant TxInputTransfer
func (tx *TransferTx) SetSignature(addr []byte, sig crypto.Signature) bool {
	if bytes.Equal(tx.Committer.Address, addr) {
		tx.Committer.Signature = sig
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
	return common.Fmt("TransferTx{%v: %v->%v, %v}", tx.Committer, tx.Sender, tx.Recipient, tx.CounterSigners)
}

// ValidateBasic validates Tx basic structure.
func (tx *TransferTx) ValidateBasic() (res abci.Result) {
	// Check the committer
	if res := tx.Committer.ValidateBasic(); res.IsErr() {
		return res
	}
	// Check the sender
	if res := tx.Sender.ValidateBasic(); res.IsErr() {
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
	return abci.OK
}

// SignTx signs the transaction if its address and the privateKey's one match.
func (tx *TransferTx) SignTx(privateKey crypto.PrivKey, chainID string) error {
	sig, err := SignTx(tx.SignBytes(chainID), tx.Committer.Address, privateKey)
	if err != nil {
		return err
	}
	tx.Committer.Signature = sig

	return nil
}

//-----------------------------------------------------------------------------

// ValidateBasic performs basic validation on a TxInputTransfer
func (t TxTransferCommitter) ValidateBasic() abci.Result {
	if len(t.Address) != 20 {
		return abci.ErrBaseInvalidInput.AppendLog("Invalid address length")
	}
	if t.Signature == nil {
		return abci.ErrBaseInvalidSignature.AppendLog("TxTransferCommitter transaction must be signed")
	}
	return abci.OK
}

// SignBytes generates a byte-to-byte signature.
func (t TxTransferCommitter) SignBytes(chainID string) []byte {
	signBytes := wire.BinaryBytes(chainID)
	sig := t.Signature
	t.Signature = nil
	signBytes = append(signBytes, wire.BinaryBytes(t)...)
	t.Signature = sig
	return signBytes
}

// SignTx signs the transaction if its address and the privateKey's one match.
func (t *TxTransferCommitter) SignTx(privateKey crypto.PrivKey, chainID string) error {
	sig, err := SignTx(t.SignBytes(chainID), t.Address, privateKey)
	if err != nil {
		return err
	}
	t.Signature = sig
	return nil
}

// String returns a string representation of TxTransferCommitter
func (t TxTransferCommitter) String() string {
	return common.Fmt("TxTransferCommitter{%x, %v}", t.Address, t.Signature)
}

//-----------------------------------------------------------------------------

// ValidateBasic performs basic validation on a TxInputTransfer
func (t TxTransferSender) ValidateBasic() abci.Result {
	if _, err := uuid.FromString(t.AccountID); err != nil {
		return abci.ErrBaseInvalidInput.AppendLog(common.Fmt("Invalid account_id: %s", err))
	}
	if t.Amount <= 0 {
		return abci.ErrBaseInvalidInput.AppendLog(common.Fmt("Amount must be non-negative: %q", t.Amount))
	}
	if len(t.Currency) != 3 {
		return abci.ErrBaseInvalidInput.AppendLog(common.Fmt("Currency must be either empty or a 3-letter ISO 4217 compliant code: %q", t.Currency))
	}
	// Validate currency and amount
	currency, ok := Currencies[t.Currency]
	if !ok {
		return abci.ErrBaseInvalidInput.AppendLog(common.Fmt("Unsupported currency: %q", t.Currency))
	}
	// Validate the amount against the currency
	if !currency.ValidateAmount(t.Amount) {
		return abci.ErrBaseInvalidInput.AppendLog(
			common.Fmt("Invalid amount %d for currency %s", t.Amount, currency.Symbol()))
	}
	if t.Sequence <= 0 {
		return abci.ErrBaseInvalidSequence.AppendLog(common.Fmt("Sequence must be greater than 0"))
	}
	return abci.OK
}

// String returns a string representation of TxTransferSender
func (t TxTransferSender) String() string {
	return common.Fmt("TxTransferSender{%v,%v,%v}", t.Amount, t.Currency, t.Sequence)
}

//-----------------------------------------------------------------------------

// ValidateBasic performs basic validation on a TxTransferRecipient
func (t TxTransferRecipient) ValidateBasic() abci.Result {
	if _, err := uuid.FromString(t.AccountID); err != nil {
		return abci.ErrBaseInvalidOutput.AppendLog(common.Fmt("Invalid account_id: %s", err))
	}
	return abci.OK
}

// String returns a string representation of TxInputTransfer
func (t TxTransferRecipient) String() string {
	return common.Fmt("TxTransferRecipient{%s}", t.AccountID)
}

//-----------------------------------------------------------------------------

// ValidateBasic performs basic validation on a TxTransferCounterSigner
func (t TxTransferCounterSigner) ValidateBasic() abci.Result {
	if len(t.Address) != 20 {
		return abci.ErrBaseInvalidInput.AppendLog("Invalid address length")
	}
	if t.Signature == nil {
		return abci.ErrBaseInvalidSignature.AppendLog("TxTransferCounterSigner transaction must be signed")
	}
	return abci.OK
}

// String returns a string representation of TxTransferCounterSigner
func (t TxTransferCounterSigner) String() string {
	return common.Fmt("TxTransferCounterSigner{%x,%v}", t.Address, t.Signature)
}

// SignBytes generates a byte-to-byte signature.
func (t TxTransferCounterSigner) SignBytes(chainID string) []byte {
	signBytes := wire.BinaryBytes(chainID)
	sig := t.Signature
	t.Signature = nil
	signBytes = append(signBytes, wire.BinaryBytes(t)...)
	t.Signature = sig
	return signBytes
}

// SignTx signs the transaction if its address and the privateKey's one match.
func (t *TxTransferCounterSigner) SignTx(privateKey crypto.PrivKey, chainID string) error {
	t.Signature = privateKey.Sign(t.SignBytes(chainID))
	return nil
}
