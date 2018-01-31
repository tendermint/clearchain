package types

import (
	wire "github.com/tendermint/go-wire"
	crypto "github.com/tendermint/go-crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)
var cdc = wire.NewCodec()

func init() {
	crypto.RegisterWire(cdc)	
	// Must register message interface to parse sdk.StdTx
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	RegisterWire(cdc)
}

func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(DepositMsg{},
		"com.tendermint.clearchain.DepositMsg", nil)
	cdc.RegisterConcrete(SettleMsg{},
		"com.tendermint.clearchain.SettleMsg", nil)
	cdc.RegisterConcrete(WithdrawMsg{},
		"com.tendermint.clearchain.WithdrawMsg", nil)
	cdc.RegisterConcrete(CreateAccountMsg{},
		"com.tendermint.clearchain.CreateAccountMsg", nil)
}
