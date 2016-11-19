package types

import (
	"reflect"
	"testing"
)

func TestAccountIndex_Has(t *testing.T) {
	type fields struct {
		Accounts []string
	}
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"empty", fields{[]string{}}, args{"string"}, false},
		{"hasIt", fields{[]string{"a", "b", "c"}}, args{"a"}, true},
		{"hasNotIt", fields{[]string{"a", "b", "c"}}, args{"d"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &AccountIndex{
				Accounts: tt.fields.Accounts,
			}
			if got := i.Has(tt.args.s); got != tt.want {
				t.Errorf("AccountIndex.Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccountIndex_ToStringSlice(t *testing.T) {
	type fields struct {
		Accounts []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{"empty", fields{[]string{}}, []string{}},
		{"nonEmpty", fields{[]string{"a", "b", "c"}}, []string{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &AccountIndex{
				Accounts: tt.fields.Accounts,
			}
			if got := i.ToStringSlice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AccountIndex.ToStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccountIndex_Add(t *testing.T) {
	type fields struct {
		Accounts []string
	}
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"addNewItem", fields{[]string{"a", "b", "c"}}, args{"d"}},
		{"addExistingItem", fields{[]string{"a", "b", "c"}}, args{"a"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &AccountIndex{
				Accounts: tt.fields.Accounts,
			}
			i.Add(tt.args.s)
			if !i.Has(tt.args.s) {
				t.Errorf("AccountIndex.Add() didn't add %s", tt.args.s)
			}
		})
	}
}
