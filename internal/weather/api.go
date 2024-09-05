package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"weather-cli/internal/config"
)

var BaseURL = "https://api.openweathermap.org/data/2.5/forecast"

// WeatherData represents the structure of the OpenWeather API response
type WeatherData struct {
	Cod     string  `json:"cod"`
	Message float64 `json:"message"`
	Cnt     int     `json:"cnt"`
	List    []struct {
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
	} `json:"list"`
	City struct {
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
	} `json:"city"`
}

// WeatherService インターフェースを定義
type WeatherService interface {
	GetWeatherForecast(cfg *config.Config, location config.Location) (*WeatherData, error)
}

// 実際のWeatherServiceの実装
type RealWeatherService struct{}

func (s *RealWeatherService) GetWeatherForecast(cfg *config.Config, location config.Location) (*WeatherData, error) {
    client := &http.Client{Timeout: 10 * time.Second}

    url := fmt.Sprintf("%s?lat=%f&lon=%f&appid=%s&units=metric", BaseURL, location.Latitude, location.Longitude, cfg.APIKey)
	
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to OpenWeather API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenWeather API returned non-OK status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var weatherData WeatherData
	if err := json.Unmarshal(body, &weatherData); err != nil {
		return nil, fmt.Errorf("error unmarshaling weather data: %w", err)
	}

	return &weatherData, nil
}

// デフォルトのサービスインスタンス
var DefaultWeatherService WeatherService = &RealWeatherService{}

// GetWeatherForecast は DefaultWeatherService を使用
func GetWeatherForecast(cfg *config.Config, location config.Location) (*WeatherData, error) {
	return DefaultWeatherService.GetWeatherForecast(cfg, location)
}
