package main

import (
	"os"
	"testing"
	"weather-cli/internal/config"
)

func TestRun(t *testing.T) {
	testConfig := &config.Config{
		Locations: []config.Location{
			{Name: "Tokyo", Latitude: 35.6895, Longitude: 139.6917},
		},
		TemperatureUnit:  "C",
		ForecastInterval: 24,
		APIKey:           "test_api_key",
	}
	err := config.SaveConfig(testConfig)
	if err != nil {
		t.Fatalf("Failed to save test config: %v", err)
	}
	defer os.Remove("config.json")

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"Help command", []string{"weather", "--help"}, false},
		{"Invalid command", []string{"weather", "--invalid"}, true},
		{"List locations", []string{"weather", "--list"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			err := run()
			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
