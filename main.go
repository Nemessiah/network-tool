package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nemessiah/network-tool/entrypoint/api"
	"github.com/nemessiah/network-tool/entrypoint/interactive"
)

func main() {
	mode := flag.String("mode", "cli", "Run mode (e.g. dev, prod, test)")
	flag.Parse()

	switch *mode {
	case "":
		fmt.Println("Error: --mode is required")
		os.Exit(1)
	case "cli":
		interactive.Cli()
	case "api":
		fmt.Println("Running in mode:", *mode)
		api.Api()
	default:
		fmt.Println("Invalid mode:", *mode)
		os.Exit(1)
	}
}
