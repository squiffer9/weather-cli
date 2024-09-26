package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var defaultConfigFile = "config.json"
var defaultTempUnit = "C"
var defaultForecastHours = 24

// Config represents the application configuration
type Config struct {
	Locations        []Location `json:"locations"`
	TemperatureUnit  string     `json:"temperature_unit"`
	ForecastInterval int        `json:"forecast_interval"`
	APIKey           string     `json:"api_key"`
}

// Location represents a saved location
type Location struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// LoadConfig loads the configuration from the config file
func LoadConfig() (*Config, error) {
	if _, err := os.Stat(defaultConfigFile); os.IsNotExist(err) {
		return createDefaultConfig()
	}

	file, err := os.ReadFile(defaultConfigFile)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the configuration to the config file
func SaveConfig(cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	if err := os.WriteFile(defaultConfigFile, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

// createDefaultConfig creates a default configuration
func createDefaultConfig() (*Config, error) {
	cfg := &Config{
		TemperatureUnit:  defaultTempUnit,
		ForecastInterval: defaultForecastHours,
	}

	if err := SaveConfig(cfg); err != nil {
		return nil, fmt.Errorf("error saving default config: %w", err)
	}

	return cfg, nil
}

// AddLocation adds a new location to the configuration
func (c *Config) AddLocation(name string, lat, lon float64) {
	c.Locations = append(c.Locations, Location{
		Name:      name,
		Latitude:  lat,
		Longitude: lon,
	})
}

// RemoveLocation removes a location from the configuration
func (c *Config) RemoveLocation(name string) error {
	for i, loc := range c.Locations {
		if loc.Name == name {
			c.Locations = append(c.Locations[:i], c.Locations[i+1:]...)
			return nil
		}
	}
	return errors.New("location not found")
}

// SetTemperatureUnit sets the temperature unit in the configuration
func (c *Config) SetTemperatureUnit(unit string) {
	c.TemperatureUnit = unit
}

// SetForecastInterval sets the forecast interval in the configuration
func (c *Config) SetForecastInterval(hours int) {
	c.ForecastInterval = hours
}

// SetAPIKey sets the API key in the configuration
func (c *Config) SetAPIKey(apiKey string) {
	c.APIKey = apiKey
}

// GetConfigDir returns the directory where the config file is stored
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".weather-cli")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("error creating config directory: %w", err)
	}

	return configDir, nil
}

func init() {
	configDir, err := GetConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting config directory: %v\n", err)
		os.Exit(1)
	}

	defaultConfigFile = filepath.Join(configDir, "config.json")
}
