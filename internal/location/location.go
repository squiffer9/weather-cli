package location

import (
	"errors"
	"fmt"
	"weather-cli/internal/config"
)

// Manager handles location-related operations
type Manager struct {
	cfg *config.Config
}

// NewManager creates a new location manager
func NewManager(cfg *config.Config) *Manager {
	return &Manager{cfg: cfg}
}

// AddLocation adds a new location to the configuration
func (m *Manager) AddLocation(name string, lat, lon float64) error {
	for _, loc := range m.cfg.Locations {
		if loc.Name == name {
			return fmt.Errorf("location with name '%s' already exists", name)
		}
	}

	m.cfg.AddLocation(name, lat, lon)
	return config.SaveConfig(m.cfg)
}

// RemoveLocation removes a location from the configuration
func (m *Manager) RemoveLocation(name string) error {
	err := m.cfg.RemoveLocation(name)
	if err != nil {
		return err
	}
	return config.SaveConfig(m.cfg)
}

// GetLocation retrieves a location by name
func (m *Manager) GetLocation(name string) (*config.Location, error) {
	for _, loc := range m.cfg.Locations {
		if loc.Name == name {
			return &loc, nil
		}
	}
	return nil, errors.New("location not found")
}

// ListLocations returns all saved locations
func (m *Manager) ListLocations() []config.Location {
	return m.cfg.Locations
}

// UpdateLocation updates an existing location
func (m *Manager) UpdateLocation(name string, lat, lon float64) error {
	for i, loc := range m.cfg.Locations {
		if loc.Name == name {
			m.cfg.Locations[i].Latitude = lat
			m.cfg.Locations[i].Longitude = lon
			return config.SaveConfig(m.cfg)
		}
	}
	return errors.New("location not found")
}
