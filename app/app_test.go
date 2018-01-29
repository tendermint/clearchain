package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: init state
// TODO: query
func TestApp(t *testing.T) {
	cc := NewClearchainApp()
	junk := []byte("khekfhewgfsug")

	cc.BeginBlock(abci.RequestBeginBlock{})
	// garbage in, garbage out
	dres := cc.DeliverTx(junk)
	assert.EqualValues(t, sdk.CodeTxParse, dres.Code, dres.Log)

	cc.EndBlock(abci.RequestEndBlock{})
	// no data in db
	// cres := cc.Commit()
	// assert.Equal(t, 0, len(cres.Data))
}
