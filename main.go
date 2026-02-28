package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nemessiah/network-tool/entrypoint/api"
	"github.com/nemessiah/network-tool/entrypoint/interactive"
	"github.com/nemessiah/network-tool/internal"
)

func main() {
	mode := flag.String("mode", "cli", "Run mode (e.g. dev, prod, test)")
	flag.Parse()

	err := internal.GenerateConfig()
	if err != nil {
		log.Fatalf("Error (%T): %+v", err, err)
	}

	switch *mode {
	case "":
		fmt.Println("Error: --mode is required")
		os.Exit(1)
	case "cli":
		interactive.CliEntry()
	case "api":
		fmt.Println("Running in mode:", *mode)
		api.ApiEntry()
	default:
		fmt.Println("Invalid mode:", *mode)
		os.Exit(1)
	}
}
