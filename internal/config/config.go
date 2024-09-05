package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	Locations       []Location `json:"locations"`
	TemperatureUnit string     `json:"temperature_unit"`
	ForecastInterval int       `json:"forecast_interval"`
	APIKey          string     `json:"api_key"`
}

// Location represents a saved location
type Location struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

var (
	defaultConfigFile = "config.json"
)

const (
	defaultTempUnit      = "C"
	defaultForecastHours = 24
)

// LoadConfig loads the configuration from the config file
func LoadConfig() (*Config, error) {
	configPath := defaultConfigFile
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return createDefaultConfig()
		}
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// Validate and set default values
	if err := config.validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

// validate checks the configuration and sets default values if needed
func (c *Config) validate() error {
	if c.TemperatureUnit == "" {
		c.TemperatureUnit = defaultTempUnit
	}
	if c.ForecastInterval <= 0 {
		c.ForecastInterval = defaultForecastHours
	}
	if c.APIKey == "" {
		return errors.New("API key is missing in the configuration")
	}
	if len(c.Locations) == 0 {
		return errors.New("no locations defined in the configuration")
	}
	return nil
}

// SaveConfig saves the configuration to the config file
func SaveConfig(config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	configPath := defaultConfigFile
	// If the directory does not exist, create it
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
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

// SetTemperatureUnit sets the temperature unit preference
func (c *Config) SetTemperatureUnit(unit string) {
	c.TemperatureUnit = unit
}

// SetForecastInterval sets the forecast interval
func (c *Config) SetForecastInterval(hours int) {
	c.ForecastInterval = hours
}

// SetAPIKey sets the OpenWeather API key
func (c *Config) SetAPIKey(key string) {
	c.APIKey = key
}

// createDefaultConfig creates a default configuration
func createDefaultConfig() (*Config, error) {
	config := &Config{
		TemperatureUnit:  defaultTempUnit,
		ForecastInterval: defaultForecastHours,
	}
	// If the directory does not exist, create it
	dir := filepath.Dir(defaultConfigFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	err := SaveConfig(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
