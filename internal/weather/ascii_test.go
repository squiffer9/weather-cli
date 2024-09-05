package weather

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetWeatherAscii(t *testing.T) {
	tests := []struct {
		name        string
		conditionID int
		wantContain string
		description string
	}{
		{
			name:        "Clear sky",
			conditionID: 800,
			wantContain: `\   /`,
			description: "Should contain sun representation",
		},
		{
			name:        "Few clouds",
			conditionID: 801,
			wantContain: `/"".-.`,
			description: "Should contain small cloud representation",
		},
		{
			name:        "Scattered clouds",
			conditionID: 802,
			wantContain: "(    )",
			description: "Should contain medium cloud representation",
		},
		{
			name:        "Broken clouds",
			conditionID: 803,
			wantContain: "*   *",
			description: "Should contain larger cloud representation with stars",
		},
		{
			name:        "Shower rain",
			conditionID: 500,
			wantContain: "' ' ' '",
			description: "Should contain light rain representation",
		},
		{
			name:        "Rain",
			conditionID: 501,
			wantContain: "‚'‚'‚'‚'",
			description: "Should contain heavier rain representation",
		},
		{
			name:        "Thunderstorm",
			conditionID: 200,
			wantContain: "⚡''⚡''",
			description: "Should contain lightning representation",
		},
		{
			name:        "Snow",
			conditionID: 600,
			wantContain: "*  *  *",
			description: "Should contain snowflake representation",
		},
		{
			name:        "Mist",
			conditionID: 701,
			wantContain: "_ - _ - _",
			description: "Should contain mist representation",
		},
		{
			name:        "Unknown condition",
			conditionID: 999,
			wantContain: "?????",
			description: "Should return default unknown weather representation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetWeatherAscii(tt.conditionID)
			if !strings.Contains(got, tt.wantContain) {
				t.Errorf("GetWeatherAscii(%d) = %v, want it to contain %v", tt.conditionID, got, tt.wantContain)
			}
			if len(got) == 0 {
				t.Errorf("GetWeatherAscii(%d) returned an empty string", tt.conditionID)
			}
		})
	}
}

func TestWeatherAsciiArtCompleteness(t *testing.T) {
	expectedConditions := []int{800, 801, 802, 803, 500, 501, 200, 600, 701}

	for _, condition := range expectedConditions {
		t.Run(fmt.Sprintf("Condition %d", condition), func(t *testing.T) {
			if _, exists := WeatherAsciiArt[condition]; !exists {
				t.Errorf("WeatherAsciiArt is missing ASCII art for condition ID %d", condition)
			}
		})
	}

	// Test for unexpected conditions
	for condition := range WeatherAsciiArt {
		t.Run(fmt.Sprintf("Unexpected condition %d", condition), func(t *testing.T) {
			if !containsInt(expectedConditions, condition) {
				t.Errorf("WeatherAsciiArt contains unexpected condition ID %d", condition)
			}
		})
	}
}

/// Helper function to check if a slice contains a value
func containsInt(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}


func TestGetWeatherAsciiEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		conditionID int
		description string
	}{
		{
			name:        "Negative condition ID",
			conditionID: -1,
			description: "Should return default unknown weather representation",
		},
		{
			name:        "Very large condition ID",
			conditionID: 10000,
			description: "Should return default unknown weather representation",
		},
		{
			name:        "Zero condition ID",
			conditionID: 0,
			description: "Should return default unknown weather representation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetWeatherAscii(tt.conditionID)
			if !strings.Contains(got, "?????") {
				t.Errorf("GetWeatherAscii(%d) = %v, want it to contain default unknown representation", tt.conditionID, got)
			}
		})
	}
}
