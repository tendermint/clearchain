package types

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/satori/go.uuid"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire"
)

func TestTransferTx_SignBytes(t *testing.T) {
	chainID := "test_chain_id"
	privKey := crypto.GenPrivKeyEd25519()
	tx := &TransferTx{
		Committer: TxTransferCommitter{
			Address:   []byte("test_sender_address"),
			Signature: privKey.Sign(crypto.CRandBytes(128)),
		},
		Sender: TxTransferSender{
			AccountID: uuid.NewV4().String(),
			Amount:    99999999,
			Currency:  "USD",
			Sequence:  1,
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
	tx.Committer.Signature = nil
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
		Committer: TxTransferCommitter{
			Address: senderAddr,
		},
		Sender: TxTransferSender{
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
		Committer      TxTransferCommitter
		Sender         TxTransferSender
		Recipient      TxTransferRecipient
		CounterSigners []TxTransferCounterSigner
	}
	tests := []struct {
		name   string
		fields fields
		want   abci.Result
	}{
		{"emptySender", fields{Sender: TxTransferSender{}}, abci.ErrBaseInvalidInput},
		{"invalidSender", fields{Sender: TxTransferSender{AccountID: uuid.NewV4().String()}}, abci.ErrBaseInvalidInput},
		{"invalidSequence", fields{
			Sender: TxTransferSender{
				AccountID: uuid.NewV4().String(),
				Currency:  "USD",
				Amount:    100,
			},
			Committer: TxTransferCommitter{Address: crypto.CRandBytes(20), Signature: signature},
		}, abci.ErrBaseInvalidSequence},
		{"emptyRecipient", fields{
			Sender: TxTransferSender{
				AccountID: uuid.NewV4().String(),
				Currency:  "USD",
				Amount:    100,
				Sequence:  1},
			Committer: TxTransferCommitter{Address: crypto.CRandBytes(20), Signature: signature},
		}, abci.ErrBaseInvalidOutput},
		{"validWithoutConterSigners", fields{
			Sender: TxTransferSender{
				AccountID: uuid.NewV4().String(),
				Currency:  "USD",
				Amount:    100,
				Sequence:  1},
			Committer: TxTransferCommitter{Address: crypto.CRandBytes(20), Signature: signature},
			Recipient: TxTransferRecipient{AccountID: uuid.NewV4().String()}}, abci.OK},
		{"invalidCurrency", fields{
			Sender: TxTransferSender{
				AccountID: uuid.NewV4().String(),
				Currency:  "invalid",
				Amount:    100,
				Sequence:  1},
			Committer: TxTransferCommitter{Address: crypto.CRandBytes(20), Signature: signature},
			Recipient: TxTransferRecipient{AccountID: uuid.NewV4().String()}}, abci.ErrBaseInvalidInput},
		{"emptyCounterSigner", fields{
			Sender: TxTransferSender{
				AccountID: uuid.NewV4().String(),
				Currency:  "USD",
				Amount:    100,
				Sequence:  1,
			},
			Committer:      TxTransferCommitter{Address: crypto.CRandBytes(20), Signature: signature},
			CounterSigners: []TxTransferCounterSigner{TxTransferCounterSigner{}},
			Recipient:      TxTransferRecipient{AccountID: uuid.NewV4().String()}}, abci.ErrBaseInvalidInput},
		{"invalidCounterSignature", fields{
			Sender: TxTransferSender{
				AccountID: uuid.NewV4().String(),
				Currency:  "USD",
				Amount:    100,
				Sequence:  3,
			},
			Committer:      TxTransferCommitter{Address: crypto.CRandBytes(20), Signature: signature},
			CounterSigners: []TxTransferCounterSigner{TxTransferCounterSigner{Address: crypto.CRandBytes(20)}},
			Recipient:      TxTransferRecipient{AccountID: uuid.NewV4().String()}}, abci.ErrBaseInvalidSignature},
		{"validWithCounterSignatures", fields{
			Sender: TxTransferSender{
				AccountID: uuid.NewV4().String(),
				Currency:  "USD",
				Amount:    100,
				Sequence:  3,
			},
			Committer: TxTransferCommitter{Address: crypto.CRandBytes(20), Signature: signature},
			CounterSigners: []TxTransferCounterSigner{TxTransferCounterSigner{
				Address:   crypto.CRandBytes(20),
				Signature: crypto.GenPrivKeyEd25519().Sign([]byte("test_content")),
			}}, Recipient: TxTransferRecipient{AccountID: uuid.NewV4().String()}}, abci.OK},
	}
	for _, tt := range tests {
		tx := TransferTx{
			Committer:      tt.fields.Committer,
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
	tests := []struct {
		name string
		tx   TxTransferSender
		want abci.Result
	}{
		{
			"emptyInput", TxTransferSender{}, abci.ErrBaseInvalidInput,
		},
		{
			"invalidAmount", TxTransferSender{
				AccountID: uuid.NewV4().String(),
				Amount:    0,
			}, abci.ErrBaseInvalidInput,
		},
		{
			"invalidCurrency", TxTransferSender{
				AccountID: uuid.NewV4().String(),
				Amount:    10,
				Currency:  "Invalid currency",
			}, abci.ErrBaseInvalidInput,
		},
		{
			"invalidCurrency", TxTransferSender{
				AccountID: uuid.NewV4().String(),
				Amount:    100,
				Currency:  "XXX",
				Sequence:  1,
			}, abci.ErrBaseInvalidInput,
		},
		{
			"invalidAccount", TxTransferSender{
				Amount:   100,
				Currency: "USD",
				Sequence: 1,
			}, abci.ErrBaseInvalidInput,
		},
		{
			"validInput", TxTransferSender{
				AccountID: uuid.NewV4().String(),
				Amount:    100,
				Currency:  "USD",
				Sequence:  1,
			}, abci.OK,
		},
		{
			"invalidSequence", TxTransferSender{
				AccountID: uuid.NewV4().String(),
				Amount:    100,
				Currency:  "USD",
				Sequence:  0,
			}, abci.ErrBaseInvalidSequence,
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
		want abci.Result
	}{
		{"emptyTx", TxTransferRecipient{}, abci.ErrBaseInvalidOutput},
		{"validUUID", TxTransferRecipient{AccountID: uuid.NewV4().String()}, abci.OK},
		{"invalidUUID", TxTransferRecipient{AccountID: "invalid"}, abci.ErrBaseInvalidOutput},
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
		want abci.Result
	}{
		{"invalidAddress", TxTransferCounterSigner{}, abci.ErrBaseInvalidInput},
		{"invalidSignature", TxTransferCounterSigner{Address: genAddr()}, abci.ErrBaseInvalidSignature},
		{"validCounterSigner", TxTransferCounterSigner{Address: genAddr(), Signature: genSignature()}, abci.OK},
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

func TestTxTransferCommitter_ValidateBasic(t *testing.T) {
	type fields struct {
		Address   []byte
		Signature crypto.Signature
	}
	tests := []struct {
		name   string
		fields fields
		want   abci.Result
	}{
		{
			"emptyInput", fields{}, abci.ErrBaseInvalidInput,
		},
		{
			"invalidAddress", fields{Address: []byte{}}, abci.ErrBaseInvalidInput,
		},
		{
			"invalidSignature", fields{
				Address: crypto.CRandBytes(20),
			}, abci.ErrBaseInvalidSignature,
		},
		{
			"validInput", fields{
				Address:   crypto.CRandBytes(20),
				Signature: crypto.GenPrivKeyEd25519().Sign([]byte("test_content")),
			}, abci.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := TxTransferCommitter{
				Address:   tt.fields.Address,
				Signature: tt.fields.Signature,
			}
			if got := tx.ValidateBasic(); got.Code != tt.want.Code {
				t.Errorf("TxTransferCommitter.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransferTx_SignTx(t *testing.T) {
	privKey := crypto.GenPrivKeyEd25519()
	addr := privKey.PubKey().Address()
	type fields struct {
		Committer      TxTransferCommitter
		Sender         TxTransferSender
		Recipient      TxTransferRecipient
		CounterSigners []TxTransferCounterSigner
	}
	type args struct {
		privateKey crypto.PrivKey
		chainID    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"addressMismatch",
			fields{TxTransferCommitter{Address: []byte{}}, TxTransferSender{}, TxTransferRecipient{}, nil},
			args{privKey, "test"},
			true,
		},
		{
			"validSignature",
			fields{TxTransferCommitter{Address: addr}, TxTransferSender{}, TxTransferRecipient{}, nil},
			args{privKey, "test"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &TransferTx{
				Committer:      tt.fields.Committer,
				Sender:         tt.fields.Sender,
				Recipient:      tt.fields.Recipient,
				CounterSigners: tt.fields.CounterSigners,
			}
			if err := tx.SignTx(tt.args.privateKey, tt.args.chainID); (err != nil) != tt.wantErr {
				t.Errorf("TransferTx.SignTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTxTransferCommitter_SignTx(t *testing.T) {
	privKey := crypto.GenPrivKeyEd25519()
	addr := privKey.PubKey().Address()
	type fields struct {
		Address   []byte
		Signature crypto.Signature
	}
	type args struct {
		privateKey crypto.PrivKey
		chainID    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"addressMismatch",
			fields{[]byte{}, nil},
			args{privKey, "test"},
			true,
		},
		{
			"validSignature",
			fields{addr, nil},
			args{privKey, "test"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &TxTransferCommitter{
				Address:   tt.fields.Address,
				Signature: tt.fields.Signature,
			}
			if err := tx.SignTx(tt.args.privateKey, tt.args.chainID); (err != nil) != tt.wantErr {
				t.Errorf("TxTransferCommitter.SignTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTxTransferCounterSigner_SignTx(t *testing.T) {
	type fields struct {
		Address   []byte
		Signature crypto.Signature
	}
	type args struct {
		privateKey crypto.PrivKey
		chainID    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"alwaysOK", fields{}, args{crypto.GenPrivKeyEd25519(), "test"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &TxTransferCounterSigner{
				Address:   tt.fields.Address,
				Signature: tt.fields.Signature,
			}
			if err := tx.SignTx(tt.args.privateKey, tt.args.chainID); (err != nil) != tt.wantErr {
				t.Errorf("TxTransferCounterSigner.SignTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
