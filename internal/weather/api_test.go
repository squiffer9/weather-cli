package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"weather-cli/internal/config"
)

func TestGetWeatherForecast(t *testing.T) {
	tests := []struct {
		name           string
		location       config.Location
		apiResponse    WeatherData
		expectedError  bool
		errorContains  string
		checkResponse  func(*testing.T, *WeatherData)
		serverBehavior func(http.ResponseWriter, *http.Request)
	}{
		{
			name:     "Successful API call",
			location: config.Location{Name: "Tokyo", Latitude: 35.6895, Longitude: 139.6917},
			apiResponse: WeatherData{
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
							Temp: 25.5,
						},
						Weather: []struct {
							ID          int    `json:"id"`
							Main        string `json:"main"`
							Description string `json:"description"`
							Icon        string `json:"icon"`
						}{
							{ID: 800, Description: "clear sky"},
						},
					},
				},
			},
			checkResponse: func(t *testing.T, wd *WeatherData) {
				if wd.City.Name != "Tokyo" {
					t.Errorf("Expected city name Tokyo, got %s", wd.City.Name)
				}
				if wd.List[0].Main.Temp != 25.5 {
					t.Errorf("Expected temperature 25.5, got %f", wd.List[0].Main.Temp)
				}
			},
		},
		{
			name:          "API Error",
			location:      config.Location{Name: "Invalid", Latitude: 0, Longitude: 0},
			expectedError: true,
			errorContains: "API returned non-OK status",
			serverBehavior: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "Bad Request")
			},
		},
		{
			name:          "Timeout Error",
			location:      config.Location{Name: "Timeout", Latitude: 1, Longitude: 1},
			expectedError: true,
			errorContains: "context deadline exceeded",
			serverBehavior: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(500 * time.Millisecond)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.serverBehavior != nil {
					tt.serverBehavior(w, r)
					return
				}
				json.NewEncoder(w).Encode(tt.apiResponse)
			}))
			defer server.Close()

			originalBaseURL := BaseURL
			BaseURL = server.URL
			defer func() { BaseURL = originalBaseURL }()

			cfg := &config.Config{
				APIKey: "test_api_key",
			}

			service := &RealWeatherService{}

			// Create a context with a timeout
			ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
			defer cancel()

			weatherData, err := service.GetWeatherForecastWithContext(ctx, cfg, tt.location)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				} else if tt.errorContains != "" && !containsString(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', but got '%s'", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tt.checkResponse != nil {
					tt.checkResponse(t, weatherData)
				}
			}
		})
	}
}

func TestGetWeatherForecastQueryParameters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Get("lat") != "35.689500" || query.Get("lon") != "139.691700" {
			t.Errorf("Incorrect latitude or longitude: lat=%s, lon=%s", query.Get("lat"), query.Get("lon"))
		}
		if query.Get("appid") != "test_api_key" {
			t.Errorf("Incorrect API key: %s", query.Get("appid"))
		}
		if query.Get("units") != "metric" {
			t.Errorf("Incorrect units: %s", query.Get("units"))
		}
		json.NewEncoder(w).Encode(WeatherData{})
	}))
	defer server.Close()

	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	cfg := &config.Config{
		APIKey: "test_api_key",
	}
	location := config.Location{
		Name:      "Tokyo",
		Latitude:  35.6895,
		Longitude: 139.6917,
	}

	service := &RealWeatherService{}
	_, err := service.GetWeatherForecast(cfg, location)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}

// GetWeatherForecastWithContext is a new method that accepts a context
func (s *RealWeatherService) GetWeatherForecastWithContext(ctx context.Context, cfg *config.Config, location config.Location) (*WeatherData, error) {
	client := &http.Client{}
	
	url := fmt.Sprintf("%s?lat=%f&lon=%f&appid=%s&units=metric", BaseURL, location.Latitude, location.Longitude, cfg.APIKey)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request to OpenWeather API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenWeather API returned non-OK status: %s", resp.Status)
	}

	var weatherData WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return nil, fmt.Errorf("error unmarshaling weather data: %w", err)
	}

	return &weatherData, nil
}
