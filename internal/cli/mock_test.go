package cli

import (
	"weather-cli/internal/config"
	"weather-cli/internal/weather"
)

// MockWeatherService is a mock implementation of WeatherService for testing
type MockWeatherService struct {
	GetWeatherForecastFunc func(cfg *config.Config, location config.Location) (*weather.WeatherData, error)
}

// GetWeatherForecast calls the mock function
func (m *MockWeatherService) GetWeatherForecast(cfg *config.Config, location config.Location) (*weather.WeatherData, error) {
	return m.GetWeatherForecastFunc(cfg, location)
}
