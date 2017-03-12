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
		{"validTxBytes", args{[]byte{TxTypeTransfer, TxTypeQueryBase}}, 3},
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
		{"hasNone", PermBaseQueryTx, args{PermNone}, false},
		{"hasNot", PermNone, args{PermObjectsQueryTx}, false},
		{"has", NewPermByTxType(TxTypeQueryObjects, TxTypeTransfer), args{PermTransferTx}, true},
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
		{"addNone", PermObjectsQueryTx, args{PermNone}, PermObjectsQueryTx},
		{"addToNone", PermNone, args{PermObjectsQueryTx}, PermObjectsQueryTx},
		{"addPerm", PermBaseQueryTx, args{PermTransferTx}, 3},
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
		{"clearNone", PermObjectsQueryTx, args{PermNone}, PermObjectsQueryTx},
		{"clearOnNone", PermNone, args{PermObjectsQueryTx}, PermNone},
		{"clearPerm", NewPermByTxType(TxTypeQueryObjects, TxTypeTransfer), args{PermTransferTx}, PermObjectsQueryTx},
	}
	for _, tt := range tests {
		if got := tt.p.Clear(tt.args.perms); got != tt.want {
			t.Errorf("%q. Perm.Clear() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
