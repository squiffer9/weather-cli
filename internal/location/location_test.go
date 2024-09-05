package location

import (
	"reflect"
	"testing"
	"weather-cli/internal/config"
)

// mockConfig is a helper function to create a mock config for testing
func mockConfig() *config.Config {
	return &config.Config{
		Locations: []config.Location{
			{Name: "Tokyo", Latitude: 35.6895, Longitude: 139.6917},
		},
	}
}

func TestNewManager(t *testing.T) {
	cfg := mockConfig()
	manager := NewManager(cfg)
	if manager.cfg != cfg {
		t.Errorf("NewManager() did not set the config correctly")
	}
}

func TestAddLocation(t *testing.T) {
	cfg := mockConfig()
	manager := NewManager(cfg)

	// Test adding a new location
	err := manager.AddLocation("New York", 40.7128, -74.0060)
	if err != nil {
		t.Errorf("AddLocation() failed: %v", err)
	}
	if len(manager.cfg.Locations) != 2 {
		t.Errorf("AddLocation() did not add the location")
	}

	// Test adding a duplicate location
	err = manager.AddLocation("Tokyo", 35.6895, 139.6917)
	if err == nil {
		t.Errorf("AddLocation() should fail when adding a duplicate location")
	}
}

func TestRemoveLocation(t *testing.T) {
	cfg := mockConfig()
	manager := NewManager(cfg)

	// Test removing an existing location
	err := manager.RemoveLocation("Tokyo")
	if err != nil {
		t.Errorf("RemoveLocation() failed: %v", err)
	}
	if len(manager.cfg.Locations) != 0 {
		t.Errorf("RemoveLocation() did not remove the location")
	}

	// Test removing a non-existent location
	err = manager.RemoveLocation("Non-existent")
	if err == nil {
		t.Errorf("RemoveLocation() should fail when removing a non-existent location")
	}
}

func TestGetLocation(t *testing.T) {
	cfg := mockConfig()
	manager := NewManager(cfg)

	// Test getting an existing location
	loc, err := manager.GetLocation("Tokyo")
	if err != nil {
		t.Errorf("GetLocation() failed: %v", err)
	}
	if loc.Name != "Tokyo" {
		t.Errorf("GetLocation() returned incorrect location")
	}

	// Test getting a non-existent location
	_, err = manager.GetLocation("Non-existent")
	if err == nil {
		t.Errorf("GetLocation() should fail when getting a non-existent location")
	}
}

func TestListLocations(t *testing.T) {
	cfg := mockConfig()
	manager := NewManager(cfg)

	locations := manager.ListLocations()
	if !reflect.DeepEqual(locations, cfg.Locations) {
		t.Errorf("ListLocations() returned incorrect locations")
	}
}

func TestUpdateLocation(t *testing.T) {
	cfg := mockConfig()
	manager := NewManager(cfg)

	// Test updating an existing location
	err := manager.UpdateLocation("Tokyo", 35.6762, 139.6503)
	if err != nil {
		t.Errorf("UpdateLocation() failed: %v", err)
	}
	loc, _ := manager.GetLocation("Tokyo")
	if loc.Latitude != 35.6762 || loc.Longitude != 139.6503 {
		t.Errorf("UpdateLocation() did not update the location correctly")
	}

	// Test updating a non-existent location
	err = manager.UpdateLocation("Non-existent", 0, 0)
	if err == nil {
		t.Errorf("UpdateLocation() should fail when updating a non-existent location")
	}
}
