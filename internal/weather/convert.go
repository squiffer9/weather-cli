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

// ConvertTemperature converts temperature from the source unit to the target unit
func ConvertTemperature(temp float64, sourceUnit, targetUnit string) float64 {
    if sourceUnit == targetUnit {
        return RoundTemperature(temp)
    }

    switch {
    case sourceUnit == "C" && targetUnit == "F":
        return RoundTemperature(CelsiusToFahrenheit(temp))
    case sourceUnit == "F" && targetUnit == "C":
        return RoundTemperature(FahrenheitToCelsius(temp))
    default:
        return RoundTemperature(temp) // Default to no conversion if units are invalid
    }
}
