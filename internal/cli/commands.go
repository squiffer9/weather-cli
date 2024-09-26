package cli

import (
	"fmt"
	"weather-cli/internal/config"
	"weather-cli/internal/location"
	"weather-cli/internal/weather"
)

// ExecuteCommand executes the appropriate command based on the parsed arguments
func ExecuteCommand(args *ParsedArgs, cfg *config.Config) error {
	switch args.Command {
	case CommandGetWeather:
		return executeGetWeather(args, cfg)
	case CommandAddLocation:
		return executeAddLocation(args, cfg)
	case CommandRemoveLocation:
		return executeRemoveLocation(args, cfg)
	case CommandSetUnit:
		return executeSetUnit(args, cfg)
	case CommandSetInterval:
		return executeSetInterval(args, cfg)
	case CommandListLocations:
		return executeListLocations(cfg)
	case CommandHelp:
		weather.DisplayHelp()
		return nil
	case CommandSetAPIKey:
		return executeSetAPIKey(args, cfg)
	default:
		return fmt.Errorf("unknown command")
	}
}

// executeGetWeather fetches and displays weather data for a given location
func executeGetWeather(args *ParsedArgs, cfg *config.Config) error {
	locationManager := location.NewManager(cfg)
	loc, err := locationManager.GetLocation(args.Location)
	if err != nil {
		return fmt.Errorf("failed to get location: %w", err)
	}

	weatherData, err := weather.GetWeatherForecast(cfg, *loc)
	if err != nil {
		return fmt.Errorf("failed to fetch weather data: %w", err)
	}

	weather.DisplayWeather(weatherData, cfg)
	return nil
}

// executeAddLocation adds a new location to the configuration
func executeAddLocation(args *ParsedArgs, cfg *config.Config) error {
	locationManager := location.NewManager(cfg)
	if err := locationManager.AddLocation(args.Name, args.Latitude, args.Longitude); err != nil {
		return fmt.Errorf("failed to add location: %w", err)
	}
	fmt.Printf("Location '%s' added successfully.\n", args.Name)
	return nil
}

// executeRemoveLocation removes a location from the configuration
func executeRemoveLocation(args *ParsedArgs, cfg *config.Config) error {
	locationManager := location.NewManager(cfg)
	if err := locationManager.RemoveLocation(args.Name); err != nil {
		return fmt.Errorf("failed to remove location: %w", err)
	}
	fmt.Printf("Location '%s' removed successfully.\n", args.Name)
	return nil
}

// executeSetUnit sets the temperature unit in the configuration
func executeSetUnit(args *ParsedArgs, cfg *config.Config) error {
	cfg.SetTemperatureUnit(args.Unit)
	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}
	fmt.Printf("Temperature unit set to %s.\n", args.Unit)
	return nil
}

// executeSetInterval sets the forecast interval in the configuration
func executeSetInterval(args *ParsedArgs, cfg *config.Config) error {
	cfg.SetForecastInterval(args.Interval)
	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}
	fmt.Printf("Forecast interval set to %d hours.\n", args.Interval)
	return nil
}

// executeListLocations displays the list of saved locations
func executeListLocations(cfg *config.Config) error {
	locationManager := location.NewManager(cfg)
	locations := locationManager.ListLocations()
	weather.DisplayLocationList(locations)
	return nil
}

// executeSetAPIKey sets the API key in the configuration
func executeSetAPIKey(args *ParsedArgs, cfg *config.Config) error {
	cfg.SetAPIKey(args.APIKey)
	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}
	fmt.Println("API key has been set successfully.")
	return nil
}
