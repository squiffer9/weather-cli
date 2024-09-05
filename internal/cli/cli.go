package cli

import (
	"fmt"
	"os"
	"weather-cli/internal/config"
	"weather-cli/internal/weather"
)

// CLI represents the command-line interface for the weather application
type CLI struct {
	cfg        *config.Config
	loadConfig func() (*config.Config, error)
}

// NewCLI creates a new CLI instance
func NewCLI(cfg *config.Config) *CLI {
	return &CLI{
		cfg:        cfg,
		loadConfig: config.LoadConfig,
	}
}

// Run executes the CLI application
func (c *CLI) Run(args []string) error {
	// Parse command-line arguments
	parsedArgs, err := ParseArgs(args)
	if err != nil {
		return fmt.Errorf("error parsing arguments: %w", err)
	}

	// If help is requested, display help and exit
	if parsedArgs.ShowHelp {
		weather.DisplayHelp()
		return nil
	}

	// Execute the appropriate command
	err = ExecuteCommand(parsedArgs, c.cfg)
	if err != nil {
		return fmt.Errorf("error executing command: %w", err)
	}

	return nil
}

// RunWithErrorHandling runs the CLI application and handles errors
func RunWithErrorHandling() {
	// Create a new CLI instance
	cli := NewCLI(nil)

	// Load configuration
	cfg, err := cli.loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}
	cli.cfg = cfg

	// Run CLI
	if err := cli.Run(os.Args); err != nil {
		weather.DisplayError(err)
		os.Exit(1)
	}
}
