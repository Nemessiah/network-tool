package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nemessiah/network-tool/entrypoint/api"
	"github.com/nemessiah/network-tool/entrypoint/interactive"
	"github.com/nemessiah/network-tool/internal"
)

func main() {
	mode := flag.String("mode", "cli", "Run mode (e.g. cli, api)")
	config := flag.String("config", "local", "Where to generate config from")
	flag.Parse()

	switch *config {
	case "":
		fmt.Println("Error: --config is required")
		os.Exit(1)
	case "local":
		fmt.Println("Loading config from home directory")
		err := internal.GenerateConfig()
		if err != nil {
			log.Fatalf("Error (%T): %+v", err, err)
		}
	default:
		fmt.Println("Invalid config:", *config)
		os.Exit(1)
	}

	switch *mode {
	case "":
		fmt.Println("Error: --mode is required")
		os.Exit(1)
	case "cli":
		interactive.CliEntry()
	case "api":
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		dbpool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
			os.Exit(1)
		}
		defer dbpool.Close()

		fmt.Println("Running in mode:", *mode)
		api.ApiEntry(context, dbpool)
	default:
		fmt.Println("Invalid mode:", *mode)
		os.Exit(1)
	}
}
