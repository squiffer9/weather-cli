package weather

import "math"

// CelsiusToFahrenheit converts temperature from Celsius to Fahrenheit
func CelsiusToFahrenheit(celsius float64) float64 {
	return celsius*9/5 + 32
}

// FahrenheitToCelsius converts temperature from Fahrenheit to Celsius
func FahrenheitToCelsius(fahrenheit float64) float64 {
	return (fahrenheit - 32) * 5 / 9
}

// RoundTemperature rounds the temperature to one decimal place
func RoundTemperature(temp float64) float64 {
	return math.Round(temp*10) / 10
}

// ConvertTemperature converts temperature based on the unit preference
func ConvertTemperature(temp float64, sourceUnit string, targetUnit string) float64 {
	if sourceUnit == targetUnit {
		return RoundTemperature(temp)
	}

	if sourceUnit == "C" && targetUnit == "F" {
		return RoundTemperature(CelsiusToFahrenheit(temp))
	}

	if sourceUnit == "F" && targetUnit == "C" {
		return RoundTemperature(FahrenheitToCelsius(temp))
	}

	// If we reach here, something is wrong, return the original temperature
	return RoundTemperature(temp)
}
