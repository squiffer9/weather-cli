![Build Status](https://github.com/squiffer9/weather-cli/actions/workflows/go.yaml/badge.svg)
[![codecov](https://codecov.io/github/squiffer9/weather-cli/graph/badge.svg?token=5C9YAY6ARU)](https://codecov.io/github/squiffer9/weather-cli)

# Weather CLI Application

## Overview

This Weather CLI Application is a command-line tool that provides weather forecasts for specified locations. It utilizes the **OpenWeather API (version 2.5)** to fetch weather data and displays it in a user-friendly format, including ASCII art representations of weather conditions.

## Features

- Fetch and display weather forecasts for specified locations
- Support for both latitude/longitude and named location inputs
- Temperature display in both Celsius and Fahrenheit
- Precipitation information
- ASCII art representation of weather conditions
- Customizable forecast interval
- Location management (add, remove, list)

## Prerequisites

- Go 1.22 or later
- OpenWeather API key (sign up at [OpenWeather](https://openweathermap.org/api) to get your API key)

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/squiffer9/weather-cli.git
   ```
2. Navigate to the project directory:
   ```
   cd weather-cli
   ```
3. Build the application:
   ```
   go build -o weather cmd/weather/main.go
   ```

## Configuration

The application will create and manage its configuration file automatically. You don't need to create a config.json file manually. Instead, you should set your OpenWeather API key using the CLI command after installation.

## Usage

### Setting up the API Key

Before using the application for the first time, set your OpenWeather API key:

```
./weather --set-api-key your_openweather_api_key_here
```

Replace `your_openweather_api_key_here` with your actual OpenWeather API key.

### Basic Usage

```
./weather <location>
```

Replace `<location>` with either a named location you've added or latitude and longitude coordinates.

### Commands

- Get weather for a location:
  ```
  ./weather tokyo
  ./weather 35.6895 139.6917
  ```

- Add a new location:
  ```
  ./weather -i <latitude> <longitude> <name>
  ```

- Remove a location:
  ```
  ./weather -r <name>
  ```

- Set temperature unit:
  ```
  ./weather --unit <C|F>
  ```

- Set forecast interval:
  ```
  ./weather --interval <hours>
  ```

- List saved locations:
  ```
  ./weather --list
  ```

- Show help:
  ```
  ./weather --help
  ```

## Development

### Project Structure

```
weather-cli/
├── cmd/
│   └── weather/
│       └── main.go
├── internal/
│   ├── config/
│   ├── weather/
│   ├── location/
│   └── cli/
├── test/
├── mock/
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

### Running Tests

To run the tests, use the following command:

```
go test ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.

## Acknowledgements

This application uses the OpenWeather API to fetch weather data. You can find more information about the API at [OpenWeather API Documentation](https://openweathermap.org/api).
