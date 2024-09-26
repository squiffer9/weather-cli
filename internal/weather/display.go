package weather

import (
	"fmt"
	"strings"
	"time"

	"weather-cli/internal/config"
)

// DisplayWeather formats and displays the weather forecast
func DisplayWeather(data *WeatherData, cfg *config.Config) {
	fmt.Printf("Weather forecast for %s, %s\n\n", data.City.Name, data.City.Country)

	for i, item := range data.List {
		if i >= cfg.ForecastInterval {
			break
		}

		date := time.Unix(item.Dt, 0)
		temp := ConvertTemperature(item.Main.Temp, "C", cfg.TemperatureUnit)
		feelsLike := ConvertTemperature(item.Main.FeelsLike, "C", cfg.TemperatureUnit)

		fmt.Printf("Date: %s\n", date.Format("2006-01-02 15:04:05"))
		fmt.Printf("Temperature: %.1f°%s (Feels like: %.1f°%s)\n", temp, cfg.TemperatureUnit, feelsLike, cfg.TemperatureUnit)
		fmt.Printf("Humidity: %d%%\n", item.Main.Humidity)
		fmt.Printf("Wind: %.1f m/s\n", item.Wind.Speed)
		fmt.Printf("Weather: %s\n", item.Weather[0].Description)

		// Display ASCII art for the weather condition
		fmt.Println(GetWeatherAscii(item.Weather[0].ID))

		// Display precipitation information if available
		if item.Rain.ThreeH > 0 {
			fmt.Printf("Rain: %.1f mm\n", item.Rain.ThreeH)
		}
		if item.Snow.ThreeH > 0 {
			fmt.Printf("Snow: %.1f mm\n", item.Snow.ThreeH)
		}

		fmt.Println(strings.Repeat("-", 40))
	}
}

// DisplayLocationList formats and displays the list of saved locations
func DisplayLocationList(locations []config.Location) {
	fmt.Println("Saved Locations:")
	for _, loc := range locations {
		fmt.Printf("- %s (Lat: %.4f, Lon: %.4f)\n", loc.Name, loc.Latitude, loc.Longitude)
	}
}

// DisplayHelp shows the usage information for the CLI
func DisplayHelp() {
	fmt.Println("Weather CLI Application Usage:")
	fmt.Println("  weather <location>                   Get weather for a location")
	fmt.Println("  weather -i <latitude> <longitude> <name>  Add a new location")
	fmt.Println("  weather -r <name>                    Remove a location")
	fmt.Println("  weather --unit <C|F>                 Set temperature unit")
	fmt.Println("  weather --interval <hours>           Set forecast interval")
	fmt.Println("  weather --list                       List saved locations")
	fmt.Println("  weather --set-api-key <api_key>      Set the OpenWeather API key")
	fmt.Println("  weather --help                       Show this help message")
}

// DisplayError formats and displays error messages
func DisplayError(err error) {
	if err == nil {
		fmt.Println("Error: <nil>")
	} else {
		fmt.Printf("Error: %s\n", err)
	}
}
