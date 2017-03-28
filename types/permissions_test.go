package types

import "testing"

func TestNewPermByTxType(t *testing.T) {
	type args struct {
		bs []byte
	}
	tests := []struct {
		name string
		args args
		want Perm
	}{
		{"invalidTxBytes", args{}, 0},
		{"validTxBytes", args{[]byte{TxTypeTransfer, TxTypeCreateAccount}}, 3},
	}
	for _, tt := range tests {
		if got := NewPermByTxType(tt.args.bs...); got != tt.want {
			t.Errorf("%q. NewPermByTxType() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestPerm_Has(t *testing.T) {
	type args struct {
		perms Perm
	}
	tests := []struct {
		name string
		p    Perm
		args args
		want bool
	}{
		{"hasNone", PermTransferTx, args{PermNone}, false},
		{"hasNot", PermNone, args{PermTransferTx}, false},
		{"has", NewPermByTxType(TxTypeTransfer), args{PermTransferTx}, true},
	}
	for _, tt := range tests {
		if got := tt.p.Has(tt.args.perms); got != tt.want {
			t.Errorf("%q. Perm.Clear() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestPerm_Add(t *testing.T) {
	type args struct {
		perms Perm
	}
	tests := []struct {
		name string
		p    Perm
		args args
		want Perm
	}{
		{"addNone", PermTransferTx, args{PermNone}, PermTransferTx},
		{"addToNone", PermNone, args{PermTransferTx}, PermTransferTx},
		{"addPerm", PermCreateAccountTx, args{PermCreateLegalEntityTx}, 6},
	}
	for _, tt := range tests {
		if got := tt.p.Add(tt.args.perms); got != tt.want {
			t.Errorf("%q. Perm.Add() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestPerm_Clear(t *testing.T) {
	type args struct {
		perms Perm
	}
	tests := []struct {
		name string
		p    Perm
		args args
		want Perm
	}{
		{"clearNone", PermTransferTx, args{PermNone}, PermTransferTx},
		{"clearOnNone", PermNone, args{PermTransferTx}, PermNone},
		{"clearPerm", NewPermByTxType(TxTypeTransfer, TxTypeCreateAccount), args{PermTransferTx}, PermCreateAccountTx},
	}
	for _, tt := range tests {
		if got := tt.p.Clear(tt.args.perms); got != tt.want {
			t.Errorf("%q. Perm.Clear() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
