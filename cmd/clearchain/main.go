package main

import (
	"fmt"

	"github.com/tendermint/clearchain/app"
)

// Entry point of the Go app

func main() {
	app := app.NewClearchainApp(app.AppName, "cc")
	fmt.Println("Clearchain app started. Running forever on :46658")
	app.RunForever()
}
