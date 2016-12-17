package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tendermint/clearchain/cmd/ledgerctl/cmd"
)

func init() {
	log.SetPrefix("")
	log.SetFlags(0)
}

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
