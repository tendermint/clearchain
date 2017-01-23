package cmd

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tendermint/clearchain/client"
	"github.com/tendermint/go-crypto"
)

func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use:   "transfer_money",
		Short: "creates money transfer enty on blockchain",
		Run: func(cmd *cobra.Command, args []string) {

			var chainID, serverAddress, privateKeyParam, senderID, recipientID, counterSignerParam, amountParam, currency string

			if len(args) == 8 {
//ledgerctl transfer_money test_chain_id tcp://127.0.0.1:46658 ATRXWwlJ6bvNRcNRT/EMmymjZvAGsLZp5a95t9HL5NRhhDh4uTLuSQikLSS//AOeuN+s1DQMgzQjEGgglAR/r6s= 1d2df1ae-accb-11e6-bbbb-00ff5244ae7f 6b6d3a08-5527-4955-b4fd-f5ba7e083548 ASrNVL489e9TlRNmIqC+vRs96+ntDRkAi1+jWnf89Nrdc4YgmMK2CzG5yTgMPvNyEq4+b5F41q79tR0MImWtYJA= 10000 EUR
				
				chainID = args[0]
				serverAddress = args[1]
				privateKeyParam = args[2]
				senderID = args[3]
				recipientID = args[4]
				counterSignerParam = args[5] // `-` as value indicates no counter signers. 
				amountParam = args[6]
				currency = args[7]
			} else {
				chainID = readParameter("chainID")
				serverAddress = readParameter("serverAddress")
				privateKeyParam = readParameter("privateKey")
				senderID = readParameter("senderID")
				recipientID = readParameter("recipientID")
				counterSignerParam = readParameter("counterSignerAddresses (comma separated list)")
				amountParam = readParameter("amount")
				currency = readParameter("currency")
			}

			privateKey, err := crypto.PrivKeyFromBytes(client.Decode(privateKeyParam))
			if err != nil {
				panic(err)
			}

			counterSignerAddresses := [][]byte{}
			if len(counterSignerParam) > 1 {
				splitCSA := strings.Split(counterSignerParam, ",")

				counterSignerAddresses = make([][]byte, len(splitCSA))
				for i, cs := range splitCSA {
					counterSignerAddresses[i] = client.Decode(cs)
				}
			}

			amount, err := strconv.Atoi(amountParam)
			if err != nil {
				panic(err)
			}

			client.SetChainID(chainID)
			client.StartClient(serverAddress)
			client.TransferMoney(privateKey, senderID, recipientID, counterSignerAddresses, int64(amount), currency)
		},
	})
}
