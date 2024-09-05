package cli

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

func TestParseOptions(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected *Options
		wantErr  bool
	}{
		{
			name:     "Get weather for location",
			args:     []string{"weather", "Tokyo"},
			expected: &Options{Location: "Tokyo"},
			wantErr:  false,
		},
		{
			name:     "Add location",
			args:     []string{"weather", "-i", "35.6895", "139.6917", "Tokyo"},
			expected: &Options{AddLocation: true, Latitude: 35.6895, Longitude: 139.6917, Name: "Tokyo"},
			wantErr:  false,
		},
		{
			name:     "Remove location",
			args:     []string{"weather", "-r", "Tokyo"},
			expected: &Options{RemoveName: "Tokyo"},
			wantErr:  false,
		},
		{
			name:     "Set temperature unit",
			args:     []string{"weather", "--unit", "F"},
			expected: &Options{Unit: "F"},
			wantErr:  false,
		},
		{
			name:     "Set forecast interval",
			args:     []string{"weather", "--interval", "12"},
			expected: &Options{Interval: 12},
			wantErr:  false,
		},
		{
			name:     "List locations",
			args:     []string{"weather", "--list"},
			expected: &Options{List: true},
			wantErr:  false,
		},
		{
			name:     "Show help",
			args:     []string{"weather", "--help"},
			expected: &Options{Help: true},
			wantErr:  false,
		},
		{
			name:     "Invalid temperature unit",
			args:     []string{"weather", "--unit", "K"},
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "Invalid forecast interval",
			args:     []string{"weather", "--interval", "-1"},
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "Invalid add location arguments",
			args:     []string{"weather", "-i", "35.6895", "Tokyo"},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original args and restore them after the test
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			// Set up test args
			os.Args = tt.args

			// Reset flags for each test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			got, err := ParseOptions()

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParseOptions() = %v, want %v", got, tt.expected)
			}
		})
	}
}
func TestValidateOptions(t *testing.T) {
	tests := []struct {
		name    string
		opts    *Options
		wantErr bool
	}{
		{
			name:    "Valid options",
			opts:    &Options{Location: "Tokyo"},
			wantErr: false,
		},
		{
			name:    "Invalid temperature unit",
			opts:    &Options{Unit: "K"},
			wantErr: true,
		},
		{
			name:    "Invalid forecast interval",
			opts:    &Options{Interval: -1},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOptions(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
