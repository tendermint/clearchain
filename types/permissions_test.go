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
		{"validTxBytes", args{[]byte{TxTypeTransfer, TxTypeQueryAccount}}, 3},
	}
	for _, tt := range tests {
		if got := NewPermByTxType(tt.args.bs...); got != tt.want {
			t.Error(PermAccountQueryTx)
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
		{"hasNone", PermAccountQueryTx, args{PermNone}, false},
		{"hasNot", PermNone, args{PermAccountQueryTx}, false},
		{"has", NewPermByTxType(TxTypeQueryAccount, TxTypeTransfer), args{PermTransferTx}, true},
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
		{"addNone", PermAccountQueryTx, args{PermNone}, PermAccountQueryTx},
		{"addToNone", PermNone, args{PermAccountQueryTx}, PermAccountQueryTx},
		{"addPerm", PermAccountQueryTx, args{PermTransferTx}, 3},
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
		{"clearNone", PermAccountQueryTx, args{PermNone}, PermAccountQueryTx},
		{"clearOnNone", PermNone, args{PermAccountQueryTx}, PermNone},
		{"clearPerm", NewPermByTxType(TxTypeQueryAccount, TxTypeTransfer), args{PermTransferTx}, PermAccountQueryTx},
	}
	for _, tt := range tests {
		if got := tt.p.Clear(tt.args.perms); got != tt.want {
			t.Errorf("%q. Perm.Clear() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
