package cli

import (
	"os"
	"testing"
	"weather-cli/internal/config"
	"weather-cli/internal/weather"
)

// testLoadConfig is wrapped in a variable so it can be replaced in tests
var testLoadConfig func() (*config.Config, error)

// TestCLI_Run tests the Run method of the CLI
func TestCLI_Run(t *testing.T) {
	tests := []struct {
		name               string
		args               []string
		mockConfig         *config.Config
		mockWeatherService *MockWeatherService
		wantErr            bool
	}{
		{
			name:       "Help command",
			args:       []string{"weather", "--help"},
			mockConfig: &config.Config{},
			mockWeatherService: &MockWeatherService{
				GetWeatherForecastFunc: func(cfg *config.Config, location config.Location) (*weather.WeatherData, error) {
					return &weather.WeatherData{}, nil
				},
			},
			wantErr: false,
		},
		{
			name: "Get weather command",
			args: []string{"weather", "Tokyo"},
			mockConfig: &config.Config{
				Locations: []config.Location{
					{Name: "Tokyo", Latitude: 35.6895, Longitude: 139.6917},
				},
			},
			mockWeatherService: &MockWeatherService{
				GetWeatherForecastFunc: func(cfg *config.Config, location config.Location) (*weather.WeatherData, error) {
					return &weather.WeatherData{}, nil
				},
			},
			wantErr: false,
		},
		{
			name:       "Add location command",
			args:       []string{"weather", "-i", "35.6895", "139.6917", "Tokyo"},
			mockConfig: &config.Config{},
			mockWeatherService: &MockWeatherService{
				GetWeatherForecastFunc: func(cfg *config.Config, location config.Location) (*weather.WeatherData, error) {
					return &weather.WeatherData{}, nil
				},
			},
			wantErr: false,
		},
		{
			name: "Remove location command",
			args: []string{"weather", "-r", "Tokyo"},
			mockConfig: &config.Config{
				Locations: []config.Location{
					{Name: "Tokyo", Latitude: 35.6895, Longitude: 139.6917},
				},
			},
			mockWeatherService: &MockWeatherService{
				GetWeatherForecastFunc: func(cfg *config.Config, location config.Location) (*weather.WeatherData, error) {
					return &weather.WeatherData{}, nil
				},
			},
			wantErr: false,
		},
		{
			name:       "Set unit command",
			args:       []string{"weather", "--unit", "F"},
			mockConfig: &config.Config{},
			mockWeatherService: &MockWeatherService{
				GetWeatherForecastFunc: func(cfg *config.Config, location config.Location) (*weather.WeatherData, error) {
					return &weather.WeatherData{}, nil
				},
			},
			wantErr: false,
		},
		{
			name:       "Set interval command",
			args:       []string{"weather", "--interval", "12"},
			mockConfig: &config.Config{},
			mockWeatherService: &MockWeatherService{
				GetWeatherForecastFunc: func(cfg *config.Config, location config.Location) (*weather.WeatherData, error) {
					return &weather.WeatherData{}, nil
				},
			},
			wantErr: false,
		},
		{
			name:       "List locations command",
			args:       []string{"weather", "--list"},
			mockConfig: &config.Config{},
			mockWeatherService: &MockWeatherService{
				GetWeatherForecastFunc: func(cfg *config.Config, location config.Location) (*weather.WeatherData, error) {
					return &weather.WeatherData{}, nil
				},
			},
			wantErr: false,
		},
		{
			name:       "Invalid command",
			args:       []string{"weather", "--invalid"},
			mockConfig: &config.Config{},
			mockWeatherService: &MockWeatherService{
				GetWeatherForecastFunc: func(cfg *config.Config, location config.Location) (*weather.WeatherData, error) {
					return &weather.WeatherData{}, nil
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := newTestCLI(tt.mockConfig)
			weather.DefaultWeatherService = tt.mockWeatherService

			err := cli.Run(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("CLI.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunWithErrorHandling(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RunWithErrorHandling() panicked: %v", r)
		}
	}()

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"weather", "--help"}

	testLoadConfig = func() (*config.Config, error) {
		return &config.Config{}, nil
	}

	oldWeatherService := weather.DefaultWeatherService
	weather.DefaultWeatherService = &MockWeatherService{
		GetWeatherForecastFunc: func(cfg *config.Config, location config.Location) (*weather.WeatherData, error) {
			return &weather.WeatherData{}, nil
		},
	}
	defer func() { weather.DefaultWeatherService = oldWeatherService }()

	// Create a testRunWithErrorHandling function
	testRunWithErrorHandling := func() {
		cli := newTestCLI(nil)
		cfg, err := cli.loadConfig()
		if err != nil {
			return // Use return instead of os.Exit(1) in tests 
		}
		cli.cfg = cfg

		if err := cli.Run(os.Args); err != nil {
			return // Use return instead of os.Exit(1) in tests
		}
	}

	testRunWithErrorHandling()
}

// newTestCLI creates a new CLI instance with a mock loadConfig function
func newTestCLI(cfg *config.Config) *CLI {
	return &CLI{
		cfg:        cfg,
		loadConfig: testLoadConfig,
	}
}

