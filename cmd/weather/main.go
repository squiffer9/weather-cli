package main

import (
	"fmt"
	"os"
	"weather-cli/internal/cli"
	"weather-cli/internal/config"
)

// run is a helper function to run the CLI application
func run() error {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	// Create and run CLI
	weatherCLI := cli.NewCLI(cfg)
	return weatherCLI.Run(os.Args)
}

// main is the entry point for the CLI application
func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Please check your config.json file and ensure all required fields are properly set.")
		os.Exit(1)
	}
}
