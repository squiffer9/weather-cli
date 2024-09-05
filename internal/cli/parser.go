package cli

import (
	"errors"
	"flag"
	"strconv"
	"strings"
)

// Command represents the different commands available in the CLI
type Command int

const (
	CommandGetWeather Command = iota
	CommandAddLocation
	CommandRemoveLocation
	CommandSetUnit
	CommandSetInterval
	CommandListLocations
	CommandHelp
)

// ParsedArgs holds the parsed command-line arguments
type ParsedArgs struct {
	Command    Command
	Location   string
	Latitude   float64
	Longitude  float64
	Name       string
	Unit       string
	Interval   int
	ShowHelp   bool
}

// ParseArgs parses the command-line arguments and returns a ParsedArgs struct
func ParseArgs(args []string) (*ParsedArgs, error) {
	if len(args) < 2 {
		return &ParsedArgs{Command: CommandHelp}, nil
	}

	parsed := &ParsedArgs{}
	flagSet := flag.NewFlagSet("weather", flag.ContinueOnError)

	// Define flags
	flagSet.BoolVar(&parsed.ShowHelp, "help", false, "Show help message")
	addLocation := flagSet.Bool("i", false, "Add a new location")
	removeLocation := flagSet.String("r", "", "Remove a location")
	setUnit := flagSet.String("unit", "", "Set temperature unit (C or F)")
	setInterval := flagSet.Int("interval", 0, "Set forecast interval in hours")
	listLocations := flagSet.Bool("list", false, "List saved locations")

	// Parse flags
	err := flagSet.Parse(args[1:])
	if err != nil {
		return nil, err
	}

	// Handle different commands
	switch {
	case *addLocation:
		return handleAddLocation(parsed, flagSet.Args())
	case *removeLocation != "":
		parsed.Command = CommandRemoveLocation
		parsed.Name = *removeLocation
	case *setUnit != "":
		return handleSetUnit(parsed, *setUnit)
	case *setInterval != 0:
		parsed.Command = CommandSetInterval
		parsed.Interval = *setInterval
	case *listLocations:
		parsed.Command = CommandListLocations
	case parsed.ShowHelp:
		parsed.Command = CommandHelp
	default:
		// If no flags are set, assume it's a get weather command
		return handleGetWeather(parsed, flagSet.Args())
	}

	return parsed, nil
}

// handleAddLocation parses the arguments for adding a new location
func handleAddLocation(parsed *ParsedArgs, args []string) (*ParsedArgs, error) {
	if len(args) != 3 {
		return nil, errors.New("invalid arguments for adding location. Use: -i <latitude> <longitude> <name>")
	}

	lat, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		return nil, errors.New("invalid latitude")
	}

	lon, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return nil, errors.New("invalid longitude")
	}

	parsed.Command = CommandAddLocation
	parsed.Latitude = lat
	parsed.Longitude = lon
	parsed.Name = args[2]

	return parsed, nil
}

// handleSetUnit parses the arguments for setting the temperature unit
func handleSetUnit(parsed *ParsedArgs, unit string) (*ParsedArgs, error) {
	unit = strings.ToUpper(unit)
	if unit != "C" && unit != "F" {
		return nil, errors.New("invalid temperature unit. Use C or F")
	}

	parsed.Command = CommandSetUnit
	parsed.Unit = unit

	return parsed, nil
}

// handleGetWeather parses the arguments for getting the weather forecast
func handleGetWeather(parsed *ParsedArgs, args []string) (*ParsedArgs, error) {
	if len(args) == 0 {
		return nil, errors.New("location is required for getting weather")
	}

	parsed.Command = CommandGetWeather
	parsed.Location = strings.Join(args, " ")

	return parsed, nil
}
