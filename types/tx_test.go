package types

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/tendermint/clearchain/testutil/mocks/mock_tx"
	"github.com/tendermint/go-crypto"
)

func TestCanExecTx(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockExecutor := mock_tx.NewMockTxExecutor(mockCtrl)
	mockTx := mock_tx.NewMockTx(mockCtrl)
	type args struct {
		executor TxExecutor
		tx       Tx
	}
	tests := []struct {
		name string
		args args
	}{
		{"canExecTx", args{mockExecutor, mockTx}},
	}
	for _, tt := range tests {
		mockTx.EXPECT().TxType().Return(byte(0))
		mockExecutor.EXPECT().CanExecTx(byte(0))
		CanExecTx(tt.args.executor, tt.args.tx)
	}
}

func TestSignTx(t *testing.T) {
	randBytes := crypto.CRandBytes(20)
	privKey := crypto.GenPrivKeyEd25519()
	type args struct {
		signedBytes []byte
		addr        []byte
		privKey     crypto.PrivKey
	}
	tests := []struct {
		name    string
		args    args
		want    crypto.Signature
		wantErr bool
	}{
		{"validSignature", args{randBytes, privKey.PubKey().Address(), privKey}, privKey.Sign(randBytes), false},
		{"invalidSignature", args{randBytes, crypto.GenPrivKeyEd25519().PubKey().Address(), privKey}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SignTx(tt.args.signedBytes, tt.args.addr, tt.args.privKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SignTx() = %v, want %v", got, tt.want)
			}
		})
	}
}
