package types

import (
	"testing"

	"github.com/tendermint/clearchain/testutil/mocks/mock_tx"
	"github.com/golang/mock/gomock"
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
