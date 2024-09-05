package weather

import (
	"fmt"
	"math"
	"testing"
)

func TestCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		name     string
		celsius  float64
		expected float64
	}{
		{"Freezing point", 0, 32},
		{"Boiling point", 100, 212},
		{"Body temperature", 37, 98.6},
		{"Negative temperature", -40, -40},
		{"Absolute zero", -273.15, -459.67},
		{"High temperature", 1000, 1832},
		{"Decimal temperature", 23.5, 74.3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CelsiusToFahrenheit(tt.celsius)
			if math.Abs(result-tt.expected) > 0.1 {
				t.Errorf("CelsiusToFahrenheit(%f) = %f; want %f", tt.celsius, result, tt.expected)
			}
		})
	}
}

func TestFahrenheitToCelsius(t *testing.T) {
	tests := []struct {
		name       string
		fahrenheit float64
		expected   float64
	}{
		{"Freezing point", 32, 0},
		{"Boiling point", 212, 100},
		{"Body temperature", 98.6, 37},
		{"Negative temperature", -40, -40},
		{"Absolute zero", -459.67, -273.15},
		{"High temperature", 1832, 1000},
		{"Decimal temperature", 74.3, 23.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FahrenheitToCelsius(tt.fahrenheit)
			if math.Abs(result-tt.expected) > 0.1 {
				t.Errorf("FahrenheitToCelsius(%f) = %f; want %f", tt.fahrenheit, result, tt.expected)
			}
		})
	}
}

func TestRoundTemperature(t *testing.T) {
	tests := []struct {
		name     string
		temp     float64
		expected float64
	}{
		{"Round down", 23.34, 23.3},
		{"Round up", 23.36, 23.4},
		{"No rounding needed", 23.30, 23.3},
		{"Negative temperature round down", -23.34, -23.3},
		{"Negative temperature round up", -23.36, -23.4},
		{"Zero", 0, 0},
		{"Large number", 1234.56, 1234.6},
		{"Very small number", 0.001, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RoundTemperature(tt.temp)
			if result != tt.expected {
				t.Errorf("RoundTemperature(%f) = %f; want %f", tt.temp, result, tt.expected)
			}
		})
	}
}


func TestConvertTemperature(t *testing.T) {
	tests := []struct {
		name       string
		temp       float64
		sourceUnit string
		targetUnit string
		expected   float64
	}{
		{"Celsius to Fahrenheit", 25, "C", "F", 77},
		{"Fahrenheit to Celsius", 77, "F", "C", 25},
		{"Celsius as is", 25, "C", "C", 25},
		{"Fahrenheit as is", 77, "F", "F", 77},
		{"Negative Celsius to Fahrenheit", -10, "C", "F", 14},
		{"Negative Fahrenheit to Celsius", 14, "F", "C", -10},
		{"Zero Celsius to Fahrenheit", 0, "C", "F", 32},
		{"Zero Fahrenheit to Celsius", 32, "F", "C", 0},
		{"High temperature Celsius to Fahrenheit", 100, "C", "F", 212},
		{"High temperature Fahrenheit to Celsius", 212, "F", "C", 100},
		{"Low temperature Celsius to Fahrenheit", -100, "C", "F", -148},
		{"Low temperature Fahrenheit to Celsius", -148, "F", "C", -100},
		{"Invalid source unit", 25, "K", "C", 25}, // Should default to no conversion
		{"Invalid target unit", 25, "C", "K", 25}, // Should default to no conversion
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertTemperature(tt.temp, tt.sourceUnit, tt.targetUnit)
			if math.Abs(result-tt.expected) > 0.1 {
				t.Errorf("ConvertTemperature(%f, %s, %s) = %f; want %f", tt.temp, tt.sourceUnit, tt.targetUnit, result, tt.expected)
			}
		})
	}
}
func TestTemperatureConversionRoundTrip(t *testing.T) {
	temperatures := []float64{-100, -50, 0, 25, 50, 100, 200}

	for _, temp := range temperatures {
		t.Run(fmt.Sprintf("RoundTrip %.2f", temp), func(t *testing.T) {
			fahrenheit := CelsiusToFahrenheit(temp)
			backToCelsius := FahrenheitToCelsius(fahrenheit)
			if math.Abs(temp-backToCelsius) > 0.001 {
				t.Errorf("Round trip conversion failed: %.2f C -> %.2f F -> %.2f C", temp, fahrenheit, backToCelsius)
			}
		})
	}
}
