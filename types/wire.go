package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"
)

var cdc = MakeTxCodec()

// RegisterWire is the functions that registers application's
// messages types to a wire.Codec.
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(DepositMsg{},
		"com.tendermint.clearchain.DepositMsg", nil)
	cdc.RegisterConcrete(SettleMsg{},
		"com.tendermint.clearchain.SettleMsg", nil)
	cdc.RegisterConcrete(WithdrawMsg{},
		"com.tendermint.clearchain.WithdrawMsg", nil)
	cdc.RegisterConcrete(CreateOperatorMsg{},
		"com.tendermint.clearchain.CreateOperatorMsg", nil)
	cdc.RegisterConcrete(CreateAdminMsg{},
		"com.tendermint.clearchain.CreateAdminMsg", nil)
	cdc.RegisterConcrete(CreateAssetAccountMsg{},
		"com.tendermint.clearchain.CreateAssetAccountMsg", nil)
	cdc.RegisterConcrete(FreezeOperatorMsg{},
		"com.tendermint.clearchain.FreezeOperatorMsg", nil)
	cdc.RegisterConcrete(FreezeAdminMsg{},
		"com.tendermint.clearchain.FreezeAdminMsg", nil)
}

// MakeTxCodec instantiate a wire.Codec and register
// all application's types; it returns the new codec.
func MakeTxCodec() (cdc *wire.Codec) {
	cdc = wire.NewCodec()

	// Register crypto.[PubKey,PrivKey,Signature] types.
	crypto.RegisterWire(cdc)

	// Register clearchain types.
	RegisterWire(cdc)

	// Must register message interface to parse sdk.StdTx
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)

	return
}
