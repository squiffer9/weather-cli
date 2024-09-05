package cli

import (
	"bytes"
	"io"
	"os"
	"testing"
	"weather-cli/internal/config"
	"weather-cli/internal/weather"
)

func TestExecuteGetWeather(t *testing.T) {
	cfg := &config.Config{
		Locations: []config.Location{
			{Name: "Tokyo", Latitude: 35.6895, Longitude: 139.6917},
		},
		TemperatureUnit:  "C",
		ForecastInterval: 24,
	}

	args := &ParsedArgs{
		Command:  CommandGetWeather,
		Location: "Tokyo",
	}

	mockService := &MockWeatherService{
		GetWeatherForecastFunc: func(*config.Config, config.Location) (*weather.WeatherData, error) {
			return &weather.WeatherData{
				City: struct {
					ID    int    `json:"id"`
					Name  string `json:"name"`
					Coord struct {
						Lat float64 `json:"lat"`
						Lon float64 `json:"lon"`
					} `json:"coord"`
					Country    string `json:"country"`
					Population int    `json:"population"`
					Timezone   int    `json:"timezone"`
					Sunrise    int    `json:"sunrise"`
					Sunset     int    `json:"sunset"`
				}{
					Name:    "Tokyo",
					Country: "JP",
				},
				List: []struct {
					Dt   int64 `json:"dt"`
					Main struct {
						Temp      float64 `json:"temp"`
						FeelsLike float64 `json:"feels_like"`
						TempMin   float64 `json:"temp_min"`
						TempMax   float64 `json:"temp_max"`
						Pressure  int     `json:"pressure"`
						SeaLevel  int     `json:"sea_level"`
						GrndLevel int     `json:"grnd_level"`
						Humidity  int     `json:"humidity"`
						TempKf    float64 `json:"temp_kf"`
					} `json:"main"`
					Weather []struct {
						ID          int    `json:"id"`
						Main        string `json:"main"`
						Description string `json:"description"`
						Icon        string `json:"icon"`
					} `json:"weather"`
					Clouds struct {
						All int `json:"all"`
					} `json:"clouds"`
					Wind struct {
						Speed float64 `json:"speed"`
						Deg   int     `json:"deg"`
						Gust  float64 `json:"gust"`
					} `json:"wind"`
					Visibility int     `json:"visibility"`
					Pop        float64 `json:"pop"`
					Sys        struct {
						Pod string `json:"pod"`
					} `json:"sys"`
					DtTxt string `json:"dt_txt"`
					Rain  struct {
						ThreeH float64 `json:"3h"`
					} `json:"rain,omitempty"`
					Snow struct {
						ThreeH float64 `json:"3h"`
					} `json:"snow,omitempty"`
				}{
					{
						Dt: 1625097600,
						Main: struct {
							Temp      float64 `json:"temp"`
							FeelsLike float64 `json:"feels_like"`
							TempMin   float64 `json:"temp_min"`
							TempMax   float64 `json:"temp_max"`
							Pressure  int     `json:"pressure"`
							SeaLevel  int     `json:"sea_level"`
							GrndLevel int     `json:"grnd_level"`
							Humidity  int     `json:"humidity"`
							TempKf    float64 `json:"temp_kf"`
						}{
							Temp:      25.5,
							FeelsLike: 26.1,
							Humidity:  60,
						},
						Weather: []struct {
							ID          int    `json:"id"`
							Main        string `json:"main"`
							Description string `json:"description"`
							Icon        string `json:"icon"`
						}{
							{
								ID:          800,
								Main:        "Clear",
								Description: "clear sky",
								Icon:        "01d",
							},
						},
						Wind: struct {
							Speed float64 `json:"speed"`
							Deg   int     `json:"deg"`
							Gust  float64 `json:"gust"`
						}{
							Speed: 2.5,
							Deg:   180,
						},
						DtTxt: "2023-07-01 12:00:00",
					},
				},
			}, nil
		},
	}

	// Set the test service
	originalService := weather.DefaultWeatherService
	weather.DefaultWeatherService = mockService
	defer func() { weather.DefaultWeatherService = originalService }()

	// Capture standard output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := executeGetWeather(args, cfg)

	// Restore standard output
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Errorf("executeGetWeather returned an error: %v", err)
	}

	// Read the captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	capturedOutput := buf.String()

	if len(capturedOutput) == 0 {
		t.Errorf("executeGetWeather didn't produce any output")
	}

	// Check if the output contains expected strings
	expectedOutputs := []string{"Tokyo", "JP", "25.5°C"}
	for _, expected := range expectedOutputs {
		if !bytes.Contains(buf.Bytes(), []byte(expected)) {
			t.Errorf("executeGetWeather output didn't contain expected string: %s", expected)
		}
	}
}

func TestExecuteAddLocation(t *testing.T) {
	cfg := &config.Config{}

	args := &ParsedArgs{
		Command:   CommandAddLocation,
		Name:      "New York",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}

	err := executeAddLocation(args, cfg)

	if err != nil {
		t.Errorf("executeAddLocation returned an error: %v", err)
	}

	if len(cfg.Locations) != 1 || cfg.Locations[0].Name != "New York" {
		t.Errorf("Location was not added correctly")
	}
}

func TestExecuteRemoveLocation(t *testing.T) {
	cfg := &config.Config{
		Locations: []config.Location{
			{Name: "Tokyo", Latitude: 35.6895, Longitude: 139.6917},
		},
	}

	args := &ParsedArgs{
		Command: CommandRemoveLocation,
		Name:    "Tokyo",
	}

	err := executeRemoveLocation(args, cfg)

	if err != nil {
		t.Errorf("executeRemoveLocation returned an error: %v", err)
	}

	if len(cfg.Locations) != 0 {
		t.Errorf("Location was not removed correctly")
	}
}

func TestExecuteSetUnit(t *testing.T) {
	cfg := &config.Config{TemperatureUnit: "C"}

	args := &ParsedArgs{
		Command: CommandSetUnit,
		Unit:    "F",
	}

	err := executeSetUnit(args, cfg)

	if err != nil {
		t.Errorf("executeSetUnit returned an error: %v", err)
	}

	if cfg.TemperatureUnit != "F" {
		t.Errorf("Temperature unit was not set correctly")
	}
}

func TestExecuteSetInterval(t *testing.T) {
	cfg := &config.Config{ForecastInterval: 24}

	args := &ParsedArgs{
		Command:  CommandSetInterval,
		Interval: 48,
	}

	err := executeSetInterval(args, cfg)

	if err != nil {
		t.Errorf("executeSetInterval returned an error: %v", err)
	}

	if cfg.ForecastInterval != 48 {
		t.Errorf("Forecast interval was not set correctly")
	}
}

func TestExecuteListLocations(t *testing.T) {
	cfg := &config.Config{
		Locations: []config.Location{
			{Name: "Tokyo", Latitude: 35.6895, Longitude: 139.6917},
			{Name: "New York", Latitude: 40.7128, Longitude: -74.0060},
		},
	}

	// Capture standard output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := executeListLocations(cfg)

	// Restore standard output
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Errorf("executeListLocations returned an error: %v", err)
	}

	// Read the captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	capturedOutput := buf.String()

	if len(capturedOutput) == 0 {
		t.Errorf("executeListLocations didn't produce any output")
	}

	// Check if the output contains expected strings
	expectedOutputs := []string{"Tokyo", "New York"}
	for _, expected := range expectedOutputs {
		if !bytes.Contains(buf.Bytes(), []byte(expected)) {
			t.Errorf("executeListLocations output didn't contain expected string: %s", expected)
		}
	}
}

func TestExecuteCommand(t *testing.T) {
	cfg := &config.Config{
		Locations: []config.Location{
			{Name: "Tokyo", Latitude: 35.6895, Longitude: 139.6917},
		},
		TemperatureUnit:  "C",
		ForecastInterval: 24,
	}

	mockService := &MockWeatherService{
		GetWeatherForecastFunc: func(*config.Config, config.Location) (*weather.WeatherData, error) {
			return &weather.WeatherData{
				City: struct {
					ID    int    `json:"id"`
					Name  string `json:"name"`
					Coord struct {
						Lat float64 `json:"lat"`
						Lon float64 `json:"lon"`
					} `json:"coord"`
					Country    string `json:"country"`
					Population int    `json:"population"`
					Timezone   int    `json:"timezone"`
					Sunrise    int    `json:"sunrise"`
					Sunset     int    `json:"sunset"`
				}{
					Name:    "Tokyo",
					Country: "JP",
				},
				List: []struct {
					Dt   int64 `json:"dt"`
					Main struct {
						Temp      float64 `json:"temp"`
						FeelsLike float64 `json:"feels_like"`
						TempMin   float64 `json:"temp_min"`
						TempMax   float64 `json:"temp_max"`
						Pressure  int     `json:"pressure"`
						SeaLevel  int     `json:"sea_level"`
						GrndLevel int     `json:"grnd_level"`
						Humidity  int     `json:"humidity"`
						TempKf    float64 `json:"temp_kf"`
					} `json:"main"`
					Weather []struct {
						ID          int    `json:"id"`
						Main        string `json:"main"`
						Description string `json:"description"`
						Icon        string `json:"icon"`
					} `json:"weather"`
					Clouds struct {
						All int `json:"all"`
					} `json:"clouds"`
					Wind struct {
						Speed float64 `json:"speed"`
						Deg   int     `json:"deg"`
						Gust  float64 `json:"gust"`
					} `json:"wind"`
					Visibility int     `json:"visibility"`
					Pop        float64 `json:"pop"`
					Sys        struct {
						Pod string `json:"pod"`
					} `json:"sys"`
					DtTxt string `json:"dt_txt"`
					Rain  struct {
						ThreeH float64 `json:"3h"`
					} `json:"rain,omitempty"`
					Snow struct {
						ThreeH float64 `json:"3h"`
					} `json:"snow,omitempty"`
				}{
					{
						Dt: 1625097600,
						Main: struct {
							Temp      float64 `json:"temp"`
							FeelsLike float64 `json:"feels_like"`
							TempMin   float64 `json:"temp_min"`
							TempMax   float64 `json:"temp_max"`
							Pressure  int     `json:"pressure"`
							SeaLevel  int     `json:"sea_level"`
							GrndLevel int     `json:"grnd_level"`
							Humidity  int     `json:"humidity"`
							TempKf    float64 `json:"temp_kf"`
						}{
							Temp:      25.5,
							FeelsLike: 26.1,
							Humidity:  60,
						},
						Weather: []struct {
							ID          int    `json:"id"`
							Main        string `json:"main"`
							Description string `json:"description"`
							Icon        string `json:"icon"`
						}{
							{
								ID:          800,
								Main:        "Clear",
								Description: "clear sky",
								Icon:        "01d",
							},
						},
						Wind: struct {
							Speed float64 `json:"speed"`
							Deg   int     `json:"deg"`
							Gust  float64 `json:"gust"`
						}{
							Speed: 2.5,
							Deg:   180,
						},
						DtTxt: "2023-07-01 12:00:00",
					},
				},
			}, nil
		},
	}

	// Set the test service
	originalService := weather.DefaultWeatherService
	weather.DefaultWeatherService = mockService
	defer func() { weather.DefaultWeatherService = originalService }()

	testCases := []struct {
		name    string
		args    *ParsedArgs
		wantErr bool
	}{
		{
			name: "Get Weather",
			args: &ParsedArgs{
				Command:  CommandGetWeather,
				Location: "Tokyo",
			},
			wantErr: false,
		},
		{
			name: "Add Location",
			args: &ParsedArgs{
				Command:   CommandAddLocation,
				Name:      "New York",
				Latitude:  40.7128,
				Longitude: -74.0060,
			},
			wantErr: false,
		},
		{
			name: "Remove Location",
			args: &ParsedArgs{
				Command: CommandRemoveLocation,
				Name:    "Tokyo",
			},
			wantErr: false,
		},
		{
			name: "Set Unit",
			args: &ParsedArgs{
				Command: CommandSetUnit,
				Unit:    "F",
			},
			wantErr: false,
		},
		{
			name: "Set Interval",
			args: &ParsedArgs{
				Command:  CommandSetInterval,
				Interval: 48,
			},
			wantErr: false,
		},
		{
			name: "List Locations",
			args: &ParsedArgs{
				Command: CommandListLocations,
			},
			wantErr: false,
		},
		{
			name: "Help",
			args: &ParsedArgs{
				Command: CommandHelp,
			},
			wantErr: false,
		},
		{
			name: "Unknown Command",
			args: &ParsedArgs{
				Command: Command(999), // Unknown command 
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Capture standard output
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := ExecuteCommand(tc.args, cfg)

			// Restore standard output
			w.Close()
			os.Stdout = oldStdout

			if (err != nil) != tc.wantErr {
				t.Errorf("ExecuteCommand() error = %v, wantErr %v", err, tc.wantErr)
			}

			// Read the captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)

			// Verify the output (e.g. for the "Get Weather" command)
			if tc.name == "Get Weather" {
				expectedOutputs := []string{"Tokyo", "JP", "25.5°C"}
				for _, expected := range expectedOutputs {
					if !bytes.Contains(buf.Bytes(), []byte(expected)) {
						t.Errorf("ExecuteCommand() output for %s didn't contain expected string: %s", tc.name, expected)
					}
				}
			}
		})
	}
}
