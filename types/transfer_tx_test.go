package types

import (
	"bytes"
	"testing"

	"fmt"

	"github.com/satori/go.uuid"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire"
	tmsp "github.com/tendermint/tmsp/types"
)

func TestTransferTx_SignBytes(t *testing.T) {
	chainID := "test_chain_id"
	privKey := crypto.GenPrivKeyEd25519()
	tx := &TransferTx{
		Sender: TxTransferSender{
			Address:   []byte("test_sender_address"),
			AccountID: uuid.NewV4().String(),
			Amount:    99999999,
			Currency:  "USD",
			Sequence:  1,
			Signature: privKey.Sign(crypto.CRandBytes(128)),
		},
		CounterSigners: []TxTransferCounterSigner{
			TxTransferCounterSigner{
				Address:   []byte("test_countersigner_address"),
				Signature: privKey.Sign(crypto.CRandBytes(128)),
			},
		},
		Recipient: TxTransferRecipient{
			AccountID: uuid.NewV4().String(),
		},
	}
	signedBytes := tx.SignBytes(chainID)
	tx.Sender.Signature = nil
	tx.CounterSigners[0].Signature = nil
	expected := wire.BinaryBytes(chainID)
	expected = append(expected, wire.BinaryBytes(tx)...)
	if !bytes.Equal(signedBytes, expected) {
		t.Errorf("SignBytes() return %v, expected: %v", signedBytes, expected)
	}
}

func TestTransferTx_SetSignature(t *testing.T) {
	senderAddr := []byte("test_sender_address0")
	senderSignature := crypto.GenPrivKeyEd25519().Sign(crypto.CRandBytes(128))
	addresses := [][]byte{
		[]byte("test_account_address1"),
		[]byte("test_account_address2"),
	}
	signatures := []crypto.Signature{
		crypto.GenPrivKeyEd25519().Sign(crypto.CRandBytes(128)),
		crypto.GenPrivKeyEd25519().Sign(crypto.CRandBytes(128)),
	}
	tx := &TransferTx{
		Sender: TxTransferSender{
			Address:   senderAddr,
			AccountID: uuid.NewV4().String(),
			Amount:    99999999,
			Currency:  "USD",
			Sequence:  1,
		},
		CounterSigners: []TxTransferCounterSigner{
			TxTransferCounterSigner{
				Address: addresses[0],
			},
			TxTransferCounterSigner{
				Address: addresses[1],
			},
			TxTransferCounterSigner{
				Address: []byte("test_account_address3"),
			},
		},
	}
	if b := tx.SetSignature(senderAddr, senderSignature); !b {
		t.Errorf("SetSignature() on the sender return %v, expected %v", b, !b)
	}
	for i := range addresses {
		if b := tx.SetSignature(addresses[i], signatures[i]); !b {
			t.Errorf("%d:%s SetSignature() return %v, expected %v",
				i, addresses[i], b, !b)
		}
	}
	for i := range addresses {
		for _, sender := range tx.CounterSigners {
			if bytes.Equal(sender.Address, addresses[i]) &&
				sender.Signature != signatures[i] {
				t.Errorf("%d: found signature %v for %s, expected %v",
					i, sender.Signature, sender.Address, signatures[i])
			}
		}
	}
	nonExistentAddr := []byte("non_existent_addr")
	if b := tx.SetSignature(nonExistentAddr,
		crypto.GenPrivKeyEd25519().Sign(
			crypto.CRandBytes(128))); b {
		t.Errorf("SetSignature() on %s return %v, expected %v", nonExistentAddr, b, !b)
	}
}

func TestTransferTx_ValidateBasic(t *testing.T) {
	signature := crypto.GenPrivKeyEd25519().Sign([]byte("test_content"))
	type fields struct {
		Sender         TxTransferSender
		Recipient      TxTransferRecipient
		CounterSigners []TxTransferCounterSigner
	}
	tests := []struct {
		name   string
		fields fields
		want   tmsp.Result
	}{
		{"emptySender", fields{Sender: TxTransferSender{}}, tmsp.ErrBaseInvalidInput},
		{"invalidSender", fields{Sender: TxTransferSender{AccountID: uuid.NewV4().String()}}, tmsp.ErrBaseInvalidInput},
		{"invalidSequence", fields{Sender: TxTransferSender{
			AccountID: uuid.NewV4().String(),
			Address:   crypto.CRandBytes(20),
			Currency:  "USD",
			Amount:    100,
			Signature: signature,
		}}, tmsp.ErrBaseInvalidSequence},
		{"emptyRecipient", fields{Sender: TxTransferSender{
			AccountID: uuid.NewV4().String(),
			Address:   crypto.CRandBytes(20),
			Currency:  "USD",
			Amount:    100,
			Signature: signature,
			Sequence:  1}}, tmsp.ErrBaseInvalidOutput},
		{"validWithoutConterSigners", fields{
			Sender: TxTransferSender{
				AccountID: uuid.NewV4().String(),
				Address:   crypto.CRandBytes(20),
				Currency:  "USD",
				Amount:    100,
				Signature: signature,
				Sequence:  1}, Recipient: TxTransferRecipient{AccountID: uuid.NewV4().String()}}, tmsp.OK},
		{"invalidCurrency", fields{Sender: TxTransferSender{
			AccountID: uuid.NewV4().String(),
			Address:   crypto.CRandBytes(20),
			Currency:  "invalid",
			Amount:    100,
			Signature: signature,
			Sequence:  1}, Recipient: TxTransferRecipient{AccountID: uuid.NewV4().String()}}, tmsp.ErrBaseInvalidInput},
		{"emptyCounterSigner", fields{
			Sender: TxTransferSender{
				AccountID: uuid.NewV4().String(),
				Address:   crypto.CRandBytes(20),
				Currency:  "USD",
				Amount:    100,
				Signature: signature,
				Sequence:  1,
			},
			CounterSigners: []TxTransferCounterSigner{TxTransferCounterSigner{}},
			Recipient:      TxTransferRecipient{AccountID: uuid.NewV4().String()}}, tmsp.ErrBaseInvalidInput},
		{"invalidCounterSignature", fields{
			Sender: TxTransferSender{
				AccountID: uuid.NewV4().String(),
				Address:   crypto.CRandBytes(20),
				Currency:  "USD",
				Amount:    100,
				Signature: signature,
				Sequence:  3,
			},
			CounterSigners: []TxTransferCounterSigner{TxTransferCounterSigner{Address: crypto.CRandBytes(20)}},
			Recipient:      TxTransferRecipient{AccountID: uuid.NewV4().String()}}, tmsp.ErrBaseInvalidSignature},
		{"validWithCounterSignatures", fields{
			Sender: TxTransferSender{
				AccountID: uuid.NewV4().String(),
				Address:   crypto.CRandBytes(20),
				Currency:  "USD",
				Amount:    100,
				Signature: signature,
				Sequence:  3,
			}, CounterSigners: []TxTransferCounterSigner{TxTransferCounterSigner{
				Address:   crypto.CRandBytes(20),
				Signature: crypto.GenPrivKeyEd25519().Sign([]byte("test_content")),
			}}, Recipient: TxTransferRecipient{AccountID: uuid.NewV4().String()}}, tmsp.OK},
	}
	for _, tt := range tests {
		tx := TransferTx{
			Sender:         tt.fields.Sender,
			Recipient:      tt.fields.Recipient,
			CounterSigners: tt.fields.CounterSigners,
		}
		if got := tx.ValidateBasic(); got.Code != tt.want.Code {
			t.Errorf("%q. TransferTx.ValidateBasic() = %v, want: %v", tt.name, got, tt.want)
		}
	}
}

func TestTxTransferSender_ValidateBasic(t *testing.T) {
	signature := crypto.GenPrivKeyEd25519().Sign([]byte("test_content"))
	tests := []struct {
		name string
		tx   TxTransferSender
		want tmsp.Result
	}{
		{
			"emptyInput", TxTransferSender{}, tmsp.ErrBaseInvalidInput,
		},
		{
			"emptyAddress", TxTransferSender{AccountID: uuid.NewV4().String()}, tmsp.ErrBaseInvalidInput,
		},
		{
			"invalidAddress", TxTransferSender{Address: []byte{}}, tmsp.ErrBaseInvalidInput,
		},
		{
			"invalidAmount", TxTransferSender{
				AccountID: uuid.NewV4().String(),
				Address:   crypto.CRandBytes(20),
				Amount:    0,
				Signature: signature,
			}, tmsp.ErrBaseInvalidInput,
		},
		{
			"invalidCurrency", TxTransferSender{
				AccountID: uuid.NewV4().String(),
				Address:   crypto.CRandBytes(20),
				Amount:    10,
				Currency:  "Invalid currency",
				Signature: signature,
			}, tmsp.ErrBaseInvalidInput,
		},
		{
			"invalidSignature", TxTransferSender{
				AccountID: uuid.NewV4().String(),
				Address:   crypto.CRandBytes(20),
				Amount:    100,
				Currency:  "USD",
				Sequence:  0,
			}, tmsp.ErrBaseInvalidSignature,
		},
		{
			"invalidCurrency", TxTransferSender{
				AccountID: uuid.NewV4().String(),
				Address:   crypto.CRandBytes(20),
				Amount:    100,
				Currency:  "XXX",
				Sequence:  1,
				Signature: signature,
			}, tmsp.ErrBaseInvalidInput,
		},
		{
			"invalidAccount", TxTransferSender{
				Address:   crypto.CRandBytes(20),
				Amount:    100,
				Currency:  "USD",
				Sequence:  1,
				Signature: signature,
			}, tmsp.ErrBaseInvalidInput,
		},
		{
			"validInput", TxTransferSender{
				AccountID: uuid.NewV4().String(),
				Address:   crypto.CRandBytes(20),
				Amount:    100,
				Currency:  "USD",
				Sequence:  1,
				Signature: signature,
			}, tmsp.OK,
		},
		{
			"invalidSequence", TxTransferSender{
				AccountID: uuid.NewV4().String(),
				Address:   crypto.CRandBytes(20),
				Amount:    100,
				Currency:  "USD",
				Sequence:  0,
				Signature: signature,
			}, tmsp.ErrBaseInvalidSequence,
		},
	}

	for _, tc := range tests {
		if v := tc.tx.ValidateBasic(); tc.want.Code != v.Code {
			t.Errorf("%q. TxTransferSender.ValidateBasic() got = %v, want %v", tc.name, v, tc.want)
		}
	}
}

func TestTxTransferRecipient_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		tx   TxTransferRecipient
		want tmsp.Result
	}{
		{"emptyTx", TxTransferRecipient{}, tmsp.ErrBaseInvalidOutput},
		{"validUUID", TxTransferRecipient{AccountID: uuid.NewV4().String()}, tmsp.OK},
		{"invalidUUID", TxTransferRecipient{AccountID: "invalid"}, tmsp.ErrBaseInvalidOutput},
	}

	for _, tt := range tests {
		if v := tt.tx.ValidateBasic(); tt.want.Code != v.Code {
			t.Errorf("%q. TxTransferRecipient.ValidateBasic() got = %v, want %v", tt.name, v, tt.want)
		}
	}
}

func TestTxTransferCounterSigner_ValidateBasic(t *testing.T) {
	genAddr := func() []byte {
		return crypto.CRandBytes(20)
	}
	genSignature := func() crypto.Signature {
		return crypto.GenPrivKeyEd25519().Sign(crypto.CRandBytes(120))
	}
	tests := []struct {
		name string
		tx   TxTransferCounterSigner
		want tmsp.Result
	}{
		{"invalidAddress", TxTransferCounterSigner{}, tmsp.ErrBaseInvalidInput},
		{"invalidSignature", TxTransferCounterSigner{Address: genAddr()}, tmsp.ErrBaseInvalidSignature},
		{"validCounterSigner", TxTransferCounterSigner{Address: genAddr(), Signature: genSignature()}, tmsp.OK},
	}

	for _, tt := range tests {
		got := tt.tx.ValidateBasic()
		if got.Code != tt.want.Code {
			t.Errorf("%q. TxTransferCounterSigner.ValidateBasic() got = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestTxTransferCounterSigner_String(t *testing.T) {
	addr := crypto.CRandBytes(20)
	sign := crypto.GenPrivKeyEd25519().Sign(crypto.CRandBytes(120))
	tests := []struct {
		name string
		tx   TxTransferCounterSigner
		want string
	}{
		{"stringDescription", TxTransferCounterSigner{Address: addr, Signature: sign}, fmt.Sprintf("TxTransferCounterSigner{%x,%v}", addr, sign)},
	}

	for _, tt := range tests {
		got := tt.tx.String()
		if got != tt.want {
			t.Errorf("%q. TxTransferCounterSigner.String() got = %v, want %v", tt.name, got, tt.want)
		}
	}
}
