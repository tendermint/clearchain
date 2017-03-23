package app

import (
	"reflect"
	"testing"

	abci "github.com/tendermint/abci/types"
	bctypes "github.com/tendermint/basecoin/types"
	"github.com/tendermint/clearchain/state"
	eyes "github.com/tendermint/merkleeyes/client"
)

func Test_splitQueryPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{"validPath", args{"/resource/object"}, "resource", "object", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := splitQueryPath(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitQueryPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("splitQueryPath() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("splitQueryPath() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestLedger_executeQuery(t *testing.T) {
	type fields struct {
		eyesCli    *eyes.Client
		state      *state.State
		cacheState *state.State
		plugins    *bctypes.Plugins
	}
	type args struct {
		req abci.RequestQuery
	}
	makeNewClient := func() *eyes.Client {
		v, _ := eyes.NewClient("dummy")
		return v
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes abci.Result
	}{
		{"OK",
			fields{
				makeNewClient(), state.NewState(makeNewClient),
				state.NewState(eyes.NewClient("dummy")), bctypes.NewPlugins(),
			},
			args{
				abci.RequestQuery{Path: "ciao"}}, abci.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Ledger{
				eyesCli:    tt.fields.eyesCli,
				state:      tt.fields.state,
				cacheState: tt.fields.cacheState,
				plugins:    tt.fields.plugins,
			}
			if gotRes := app.executeQuery(tt.args.req); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Ledger.executeQuery() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
