package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func setupTestEnvironment(t *testing.T) func() {
	tempDir, err := os.MkdirTemp("", "weather-cli-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	oldConfigFile := defaultConfigFile
	defaultConfigFile = filepath.Join(tempDir, "config.json")

	return func() {
		os.RemoveAll(tempDir)
		defaultConfigFile = oldConfigFile
	}
}

func TestLoadAndSaveConfig(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	testConfig := &Config{
		Locations: []Location{
			{Name: "Tokyo", Latitude: 35.6895, Longitude: 139.6917},
		},
		TemperatureUnit:  "C",
		ForecastInterval: 24,
		APIKey:           "test-api-key",
	}

	err := SaveConfig(testConfig)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	loadedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if !reflect.DeepEqual(testConfig, loadedConfig) {
		t.Errorf("Loaded config does not match saved config.\nExpected: %+v\nGot: %+v", testConfig, loadedConfig)
	}
}

func TestAddAndRemoveLocation(t *testing.T) {
	config := &Config{}

	config.AddLocation("New York", 40.7128, -74.0060)
	if len(config.Locations) != 1 || config.Locations[0].Name != "New York" {
		t.Errorf("Failed to add location")
	}

	err := config.RemoveLocation("New York")
	if err != nil {
		t.Errorf("Failed to remove existing location: %v", err)
	}

	if len(config.Locations) != 0 {
		t.Errorf("Location was not removed")
	}

	err = config.RemoveLocation("London")
	if err == nil {
		t.Errorf("Removing non-existent location should return an error")
	}
}

func TestSetConfigOptions(t *testing.T) {
	config := &Config{}

	config.SetTemperatureUnit("F")
	if config.TemperatureUnit != "F" {
		t.Errorf("Failed to set temperature unit")
	}

	config.SetForecastInterval(48)
	if config.ForecastInterval != 48 {
		t.Errorf("Failed to set forecast interval")
	}

	config.SetAPIKey("new-api-key")
	if config.APIKey != "new-api-key" {
		t.Errorf("Failed to set API key")
	}
}

func TestCreateDefaultConfig(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	config, err := createDefaultConfig()
	if err != nil {
		t.Fatalf("Failed to create default config: %v", err)
	}

	if config.TemperatureUnit != defaultTempUnit {
		t.Errorf("Default temperature unit is incorrect")
	}

	if config.ForecastInterval != defaultForecastHours {
		t.Errorf("Default forecast interval is incorrect")
	}

	if len(config.Locations) != 0 {
		t.Errorf("Default config should have no locations")
	}

	if config.APIKey != "" {
		t.Errorf("Default config should have an empty API key")
	}
}

func TestLoadNonExistentConfig(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load non-existent config: %v", err)
	}

	if config.TemperatureUnit != defaultTempUnit {
		t.Errorf("Default temperature unit is incorrect")
	}

	if config.ForecastInterval != defaultForecastHours {
		t.Errorf("Default forecast interval is incorrect")
	}
}

func TestLoadInvalidConfig(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	err := os.WriteFile(defaultConfigFile, []byte("{invalid json}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid config file: %v", err)
	}

	_, err = LoadConfig()
	if err == nil {
		t.Errorf("Expected error when loading invalid config, but got nil")
	}
}

func TestSaveConfigError(t *testing.T) {
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Get the path of the temporary directory
	tempDir := filepath.Dir(defaultConfigFile)

	err := os.Chmod(tempDir, 0555)
	if err != nil {
		t.Fatalf("Failed to change directory permissions: %v", err)
	}

	config := &Config{}
	err = SaveConfig(config)
	if err == nil {
		t.Errorf("Expected error when saving config to unwritable location, but got nil")
	}

	// Directory permissions are restored after the test
	os.Chmod(tempDir, 0755)
}
