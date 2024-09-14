package main

import (
	"fmt"
	"os"

	"github.com/haqqer/keuanganku/cmd"
)

func main() {
	var argsRaw = os.Args
	if len(argsRaw) <= 1 {
		fmt.Println("use `run` or `migrate` arguments")
		os.Exit(1)
	}
	arg := argsRaw[1]
	switch arg {
	case "serve":
		cmd.Serve()
	case "migrate":
		cmd.Migrate()
	}
}
