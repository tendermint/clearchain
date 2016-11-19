package types

import (
	"bytes"
	"testing"

	uuid "github.com/satori/go.uuid"
	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"
	tmsp "github.com/tendermint/tmsp/types"
)

func TestCreateAccountTx_TxType(t *testing.T) {
	type fields struct {
		Address   []byte
		AccountID string
		Signature crypto.Signature
	}
	tests := []struct {
		name   string
		fields fields
		want   byte
	}{
		{"default", fields{}, TxTypeCreateAccount},
	}
	for _, tt := range tests {
		tx := &CreateAccountTx{
			Address:   tt.fields.Address,
			AccountID: tt.fields.AccountID,
			Signature: tt.fields.Signature,
		}
		if got := tx.TxType(); got != tt.want {
			t.Errorf("%q. CreateAccountTx.TxType() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCreateAccountTx_SignBytes(t *testing.T) {
	chainID := "chainID"
	privKey := crypto.GenPrivKeyEd25519()
	tx := &CreateAccountTx{
		Address:   privKey.PubKey().Address(),
		AccountID: "account_id",
		Signature: nil,
	}
	signedBytes := tx.SignBytes(chainID)
	expected := append(wire.BinaryBytes(chainID), wire.BinaryBytes(tx)...)
	if !bytes.Equal(signedBytes, expected) {
		t.Errorf("CreateAccountTx.SignBytes() = %v, want: %v", signedBytes, expected)
	}
}

func TestCreateAccountTx_ValidateBasic(t *testing.T) {
	type fields struct {
		Address   []byte
		AccountID string
		Signature crypto.Signature
	}
	tests := []struct {
		name   string
		fields fields
		want   tmsp.Result
	}{
		{"emptyTx", fields{}, tmsp.ErrBaseInvalidInput},
		{"invalidAddress", fields{Address: []byte("")}, tmsp.ErrBaseInvalidInput},
		{"invalidSignature", fields{crypto.CRandBytes(20), uuid.NewV4().String(), nil}, tmsp.ErrBaseInvalidSignature},
		{"invalidAccountID", fields{crypto.CRandBytes(20), "", crypto.GenPrivKeyEd25519().Sign(crypto.CRandBytes(20))}, tmsp.ErrBaseInvalidInput},
		{"valid", fields{crypto.CRandBytes(20), uuid.NewV4().String(), crypto.GenPrivKeyEd25519().Sign(crypto.CRandBytes(20))}, tmsp.OK},
	}
	for _, tt := range tests {
		tx := &CreateAccountTx{
			Address:   tt.fields.Address,
			AccountID: tt.fields.AccountID,
			Signature: tt.fields.Signature,
		}
		if got := tx.ValidateBasic(); got.Code != tt.want.Code {
			t.Errorf("%q. CreateAccountTx.ValidateBasic() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCreateAccountTx_String(t *testing.T) {
	type fields struct {
		Address   []byte
		AccountID string
		Signature crypto.Signature
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"empty", fields{}, "CreateAccountTx{,\"\"}"},
		{"stringRepr", fields{[]byte{0}, "account_id", nil}, "CreateAccountTx{00,\"account_id\"}"},
	}
	for _, tt := range tests {
		tx := &CreateAccountTx{
			Address:   tt.fields.Address,
			AccountID: tt.fields.AccountID,
			Signature: tt.fields.Signature,
		}
		if got := tx.String(); got != tt.want {
			t.Errorf("%q. CreateAccountTx.String() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
