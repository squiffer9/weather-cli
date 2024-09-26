package cli

import (
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    *ParsedArgs
		wantErr bool
	}{
		{
			name: "Get weather for location",
			args: []string{"weather", "Tokyo"},
			want: &ParsedArgs{
				Command:  CommandGetWeather,
				Location: "Tokyo",
			},
			wantErr: false,
		},
		{
			name: "Add location",
			args: []string{"weather", "-i", "35.6895", "139.6917", "Tokyo"},
			want: &ParsedArgs{
				Command:   CommandAddLocation,
				Latitude:  35.6895,
				Longitude: 139.6917,
				Name:      "Tokyo",
			},
			wantErr: false,
		},
		{
			name: "Remove location",
			args: []string{"weather", "-r", "Tokyo"},
			want: &ParsedArgs{
				Command: CommandRemoveLocation,
				Name:    "Tokyo",
			},
			wantErr: false,
		},
		{
			name: "Set temperature unit to Celsius",
			args: []string{"weather", "--unit", "C"},
			want: &ParsedArgs{
				Command: CommandSetUnit,
				Unit:    "C",
			},
			wantErr: false,
		},
		{
			name: "Set temperature unit to Fahrenheit",
			args: []string{"weather", "--unit", "F"},
			want: &ParsedArgs{
				Command: CommandSetUnit,
				Unit:    "F",
			},
			wantErr: false,
		},
		{
			name: "Set forecast interval",
			args: []string{"weather", "--interval", "12"},
			want: &ParsedArgs{
				Command:  CommandSetInterval,
				Interval: 12,
			},
			wantErr: false,
		},
		{
			name: "List locations",
			args: []string{"weather", "--list"},
			want: &ParsedArgs{
				Command: CommandListLocations,
			},
			wantErr: false,
		},
		{
			name: "Show help",
			args: []string{"weather", "--help"},
			want: &ParsedArgs{
				Command:  CommandHelp,
				ShowHelp: true,
			},
			wantErr: false,
		},
		{
			name:    "Invalid add location (missing arguments)",
			args:    []string{"weather", "-i", "35.6895", "139.6917"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid temperature unit",
			args:    []string{"weather", "--unit", "K"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid latitude",
			args:    []string{"weather", "-i", "invalid", "139.6917", "Tokyo"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "No arguments",
			args:    []string{"weather"},
			want:    &ParsedArgs{Command: CommandHelp},
			wantErr: false,
		},
		// New test case for setting API key
		{
			name: "Set API key",
			args: []string{"weather", "--set-api-key", "abcdef123456"},
			want: &ParsedArgs{
				Command: CommandSetAPIKey,
				APIKey:  "abcdef123456",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
