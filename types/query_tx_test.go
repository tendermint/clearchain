package types

import (
	"reflect"
	"testing"

	"github.com/satori/go.uuid"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
)

func TestBaseQueryTx_TxType(t *testing.T) {
	type fields struct {
		Address   []byte
		Signature crypto.Signature
	}
	tests := []struct {
		name   string
		fields fields
		want   byte
	}{
		{"queryType", fields{}, TxTypeQueryBase},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := BaseQueryTx{
				Address:   tt.fields.Address,
				Signature: tt.fields.Signature,
			}
			if got := tx.TxType(); got != tt.want {
				t.Errorf("BaseQueryTx.TxType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseQueryTx_SignBytes(t *testing.T) {
	chainID := "test_chain_id"
	privKey := crypto.GenPrivKeyEd25519()
	pubKeyAddr := privKey.PubKey().Address()
	signBytes := func(addr []byte) []byte {
		return BaseQueryTx{Address: addr}.SignBytes(chainID)
	}
	type fields struct {
		Address   []byte
		Signature crypto.Signature
	}
	type args struct {
		chainID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{"signed", fields{pubKeyAddr, nil}, args{chainID}, signBytes(pubKeyAddr)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := BaseQueryTx{
				Address:   tt.fields.Address,
				Signature: tt.fields.Signature,
			}
			if got := tx.SignBytes(tt.args.chainID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseQueryTx.SignBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseQueryTx_ValidateBasic(t *testing.T) {
	privKey := crypto.GenPrivKeyEd25519()
	pubKeyAddr := privKey.PubKey().Address()
	signature := privKey.Sign(pubKeyAddr)
	type fields struct {
		Address   []byte
		Signature crypto.Signature
	}
	tests := []struct {
		name   string
		fields fields
		want   abci.Result
	}{
		{"invalidAddress", fields{[]byte("addr"), nil}, abci.ErrBaseInvalidInput},
		{"invalidSignature", fields{pubKeyAddr, nil}, abci.ErrBaseInvalidSignature},
		{"valid", fields{pubKeyAddr, signature}, abci.OK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := BaseQueryTx{
				Address:   tt.fields.Address,
				Signature: tt.fields.Signature,
			}
			if got := tx.ValidateBasic(); got.Code != tt.want.Code {
				t.Errorf("BaseQueryTx.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseQueryTx_SignTx(t *testing.T) {
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
		{"validSignature", fields{addr, nil}, args{privKey, "chainID"}, false},
		{"invalidSignature", fields{addr, nil}, args{crypto.GenPrivKeyEd25519(), "chainID"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &BaseQueryTx{
				Address:   tt.fields.Address,
				Signature: tt.fields.Signature,
			}
			if err := tx.SignTx(tt.args.privateKey, tt.args.chainID); (err != nil) != tt.wantErr {
				t.Errorf("BaseQueryTx.SignTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestObjectsQueryTx_TxType(t *testing.T) {
	type fields struct {
		BaseQueryTx BaseQueryTx
		Ids         []string
	}
	tests := []struct {
		name   string
		fields fields
		want   byte
	}{
		{"queryType", fields{}, TxTypeQueryObjects},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := ObjectsQueryTx{
				BaseQueryTx: tt.fields.BaseQueryTx,
				Ids:         tt.fields.Ids,
			}
			if got := tx.TxType(); got != tt.want {
				t.Errorf("ObjectsQueryTx.TxType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObjectsQueryTx_String(t *testing.T) {
	signature, _ := crypto.SignatureFromBytes([]byte{1, 100, 140, 5, 246, 69, 107, 210, 41, 250, 189, 162, 44, 49, 6, 222, 185, 227, 247, 12, 213, 215, 246, 182, 66, 0, 233, 54, 215, 124, 175, 172, 235, 72, 151, 154, 26, 65, 145, 127, 121, 223, 4, 233, 210, 18, 188, 144, 72, 18, 63, 80, 158, 68, 221, 110, 82, 249, 26, 46, 202, 154, 43, 1, 13})
	type fields struct {
		BaseQueryTx BaseQueryTx
		Ids         []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"stringRepr", fields{BaseQueryTx{Address: []byte{byte(0x01)}, Signature: signature}, []string{"account_1", "account_2"}}, "ObjectsQueryTx{01 [account_1 account_2] /648C05F6456B.../}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := ObjectsQueryTx{
				BaseQueryTx: tt.fields.BaseQueryTx,
				Ids:         tt.fields.Ids,
			}
			if got := tx.String(); got != tt.want {
				t.Errorf("ObjectsQueryTx.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObjectsQueryTx_ValidateBasic(t *testing.T) {
	privKey := crypto.GenPrivKeyEd25519()
	pubKeyAddr := privKey.PubKey().Address()
	signature := privKey.Sign(pubKeyAddr)
	type fields struct {
		BaseQueryTx BaseQueryTx
		Ids         []string
	}
	tests := []struct {
		name   string
		fields fields
		want   abci.Result
	}{
		{"invalidAddress", fields{BaseQueryTx{[]byte("addr"), nil}, nil}, abci.ErrBaseInvalidInput},
		{"invalidSignature", fields{BaseQueryTx{pubKeyAddr, nil}, nil}, abci.ErrBaseInvalidSignature},
		{"emptyIds", fields{BaseQueryTx{pubKeyAddr, signature}, nil}, abci.ErrBaseInvalidInput},
		{"invalidIds", fields{BaseQueryTx{pubKeyAddr, signature}, []string{"test"}}, abci.ErrBaseInvalidInput},
		{"valid", fields{BaseQueryTx{pubKeyAddr, signature}, []string{uuid.NewV4().String()}}, abci.OK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := ObjectsQueryTx{
				BaseQueryTx: tt.fields.BaseQueryTx,
				Ids:         tt.fields.Ids,
			}
			if got := tx.ValidateBasic(); got.Code != tt.want.Code {
				t.Errorf("ObjectsQueryTx.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObjectsQueryTx_SignBytes(t *testing.T) {
	chainID := "test_chain_id"
	privKey := crypto.GenPrivKeyEd25519()
	pubKeyAddr := privKey.PubKey().Address()
	signBytes := func(addr []byte) []byte {
		return BaseQueryTx{Address: addr}.SignBytes(chainID)
	}
	type fields struct {
		BaseQueryTx BaseQueryTx
		Ids         []string
	}
	type args struct {
		chainID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{"signed", fields{BaseQueryTx{Address: pubKeyAddr}, nil}, args{chainID},
			append(signBytes(pubKeyAddr), 0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := ObjectsQueryTx{
				BaseQueryTx: tt.fields.BaseQueryTx,
				Ids:         tt.fields.Ids,
			}
			if got := tx.SignBytes(tt.args.chainID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ObjectsQueryTx.SignBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
