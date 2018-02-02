package types

import "testing"

func TestIsCreatableEntity(t *testing.T) {
	type args struct {
		entityType string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"can't create ch", args{"ch"}, false},
		{"bad entity", args{""}, false},
		{"bad entity", args{"bad"}, false},
		{"valid entity", args{"gcm"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCreatableEntity(tt.args.entityType); got != tt.want {
				t.Errorf("IsCreatableEntity() = %v, want %v", got, tt.want)
			}
		})
	}
}
