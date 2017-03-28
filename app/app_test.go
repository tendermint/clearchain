package app

import (
	"testing"
	
	bscoin "github.com/tendermint/basecoin/types"
	"github.com/tendermint/clearchain/types"
	
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
		{"validPath_generic", args{"/resource/object"}, "resource", "object", false},
		{"validPath_no_resource", args{"/resource"}, "resource", "", false},
		{"validPath_legal_entity_all", args{"/legal_entity"}, "legal_entity", "", false},
		{"validPath_account_id", args{"/account/1d2df1ae-accb-11e6-bbbb-00ff5244ae7f"}, "account", "1d2df1ae-accb-11e6-bbbb-00ff5244ae7f", false},
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
	
	makeNewClient :=  bscoin.NewMemKVStore()
	s := state.NewState(makeNewClient)
	accountIndex := types.NewAccountIndex()
	accountIndex.Add("testId")
	s.SetAccountIndex(accountIndex)
	
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes abci.ResponseQuery
	}{
		{"OK",
			fields{
				nil, s,
				s, bctypes.NewPlugins(),
			},
			args{
				abci.RequestQuery{Path: "/account"}},
			abci.ResponseQuery{Code:abci.CodeType_OK},
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
			if gotRes := app.executeQuery(tt.args.req); gotRes.Code != abci.CodeType_OK {
				t.Errorf("Ledger.executeQuery() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
