package weather

import (
	"errors"
	"fmt"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"weather-cli/internal/config"
)

// createMockWeatherData creates a mock WeatherData for testing
func createMockWeatherData() *WeatherData {
	return &WeatherData{
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
				Dt: time.Now().Unix(),
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
					FeelsLike: 26.0,
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
						Description: "clear sky",
					},
				},
				Wind: struct {
					Speed float64 `json:"speed"`
					Deg   int     `json:"deg"`
					Gust  float64 `json:"gust"`
				}{
					Speed: 3.5,
				},
				Rain: struct {
					ThreeH float64 `json:"3h"`
				}{
					ThreeH: 0.5,
				},
			},
		},
	}
}

func TestDisplayWeather(t *testing.T) {
	mockWeatherData := createMockWeatherData()

	testCases := []struct {
		name     string
		config   *config.Config
		expected []string
	}{
		{
			name: "Celsius display",
			config: &config.Config{
				TemperatureUnit:  "C",
				ForecastInterval: 1,
			},
			expected: []string{"Tokyo, JP", "Temperature: 25.5°C", "Humidity: 60%", "Wind: 3.5 m/s", "clear sky", "Rain: 0.5 mm"},
		},
		{
			name: "Fahrenheit display",
			config: &config.Config{
				TemperatureUnit:  "F",
				ForecastInterval: 1,
			},
			expected: []string{"Tokyo, JP", "Temperature: 77.9°F", "Humidity: 60%", "Wind: 3.5 m/s", "clear sky", "Rain: 0.5 mm"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			DisplayWeather(mockWeatherData, tc.config)

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			for _, expected := range tc.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', but it didn't.\nActual output:\n%s", expected, output)
				}
			}
		})
	}
}

func TestDisplayLocationList(t *testing.T) {
	locations := []config.Location{
		{Name: "Tokyo", Latitude: 35.6895, Longitude: 139.6917},
		{Name: "New York", Latitude: 40.7128, Longitude: -74.0060},
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	DisplayLocationList(locations)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	expectedOutputs := []string{
		"Saved Locations:",
		"- Tokyo (Lat: 35.6895, Lon: 139.6917)",
		"- New York (Lat: 40.7128, Lon: -74.0060)",
	}

	for _, expected := range expectedOutputs {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain '%s', but it didn't.\nActual output:\n%s", expected, output)
		}
	}
}

func TestDisplayHelp(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	DisplayHelp()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	expectedOutputs := []string{
		"Weather CLI Application Usage:",
		"weather <location>",
		"weather -i <latitude> <longitude> <name>",
		"weather -r <name>",
		"weather --unit <C|F>",
		"weather --interval <hours>",
		"weather --list",
		"weather --help",
	}

	for _, expected := range expectedOutputs {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain '%s', but it didn't.\nActual output:\n%s", expected, output)
		}
	}
}

func TestDisplayError(t *testing.T) {
	testCases := []struct {
		name          string
		err           error
		expectedError string
	}{
		{
			name:          "Simple error",
			err:           errors.New("Test error message"),
			expectedError: "Error: Test error message\n",
		},
		{
			name:          "Empty error",
			err:           errors.New(""),
			expectedError: "Error: \n",
		},
		{
			name:          "Nil error",
			err:           nil,
			expectedError: "Error: <nil>\n",
		},
		{
			name:          "Formatted error",
			err:           fmt.Errorf("formatted error: %d", 42),
			expectedError: "Error: formatted error: 42\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			DisplayError(tc.err)

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			if output != tc.expectedError {
				t.Errorf("Expected output '%s', but got '%s'", tc.expectedError, output)
			}
		})
	}
}
