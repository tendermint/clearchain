package types

import "testing"

func TestIsValidEntityType(t *testing.T) {
	type args struct {
		entityType string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"valid entity", args{"ch"}, true},
		{"bad entity", args{""}, false},
		{"bad entity", args{"bad"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidEntityType(tt.args.entityType); got != tt.want {
				t.Errorf("IsValidEntityType() = %v, want %v", got, tt.want)
			}
		})
	}
}
