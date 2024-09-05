package cli

import (
	"flag"
	"fmt"
)

// Options represents the command-line options for the weather CLI
type Options struct {
	Location    string
	AddLocation bool
	Latitude    float64
	Longitude   float64
	Name        string
	RemoveName  string
	Unit        string
	Interval    int
	List        bool
	Help        bool
}

// ParseOptions parses the command-line arguments and returns an Options struct
func ParseOptions() (*Options, error) {
	opts := &Options{}

	flag.BoolVar(&opts.AddLocation, "i", false, "Add a new location")
	flag.StringVar(&opts.RemoveName, "r", "", "Remove a location by name")
	flag.StringVar(&opts.Unit, "unit", "", "Set temperature unit (C or F)")
	flag.IntVar(&opts.Interval, "interval", 0, "Set forecast interval in hours")
	flag.BoolVar(&opts.List, "list", false, "List saved locations")
	flag.BoolVar(&opts.Help, "help", false, "Show help message")

	flag.Usage = usage
	flag.Parse()

	// Handle different cases based on flags and arguments
	args := flag.Args()
	if opts.AddLocation {
		if len(args) != 3 {
			return nil, fmt.Errorf("invalid arguments for adding location. Use: -i <latitude> <longitude> <name>")
		}
		var err error
		opts.Latitude, err = parseFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("invalid latitude: %v", err)
		}
		opts.Longitude, err = parseFloat(args[1])
		if err != nil {
			return nil, fmt.Errorf("invalid longitude: %v", err)
		}
		opts.Name = args[2]
	} else if len(args) > 0 {
		opts.Location = args[0]
	}

	return opts, validateOptions(opts)
}

// parseFloat parses a string to a float64 value
func parseFloat(s string) (float64, error) {
	var v float64
	_, err := fmt.Sscanf(s, "%f", &v)
	return v, err
}

// validateOptions checks if the provided options are valid
func validateOptions(opts *Options) error {
	if opts.Unit != "" && opts.Unit != "C" && opts.Unit != "F" {
		return fmt.Errorf("invalid temperature unit. Use C or F")
	}

	if opts.Interval < 0 {
		return fmt.Errorf("invalid forecast interval. Must be a positive number")
	}

	return nil
}

// usage prints the usage message for the weather CLI
func usage() {
	fmt.Println("Weather CLI Application Usage:")
	fmt.Println("  weather <location>                        Get weather for a location")
	fmt.Println("  weather -i <latitude> <longitude> <name>  Add a new location")
	fmt.Println("  weather -r <name>                         Remove a location")
	fmt.Println("  weather --unit <C|F>                      Set temperature unit")
	fmt.Println("  weather --interval <hours>                Set forecast interval")
	fmt.Println("  weather --list                            List saved locations")
	fmt.Println("  weather --help                            Show this help message")
}
