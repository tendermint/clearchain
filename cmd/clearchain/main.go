package main

import (
	"fmt"

	"github.com/tendermint/clearchain/app"
)

func main() {
	app := app.NewClearchainApp()
	fmt.Println("Running forever on :46658")
	app.RunForever()
}
