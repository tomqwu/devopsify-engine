package config

import (
	"fmt"
	"strings"
)

var validKinds = map[string]bool{
	"cloud":    true,
	"pipeline": true,
	"sre":      true,
}

// Validate checks the config for correctness.
func Validate(cfg *Config) error {
	seen := make(map[string]bool)

	for i, p := range cfg.Providers {
		if p.Name == "" {
			return fmt.Errorf("provider[%d]: name is required", i)
		}
		if p.Provider == "" {
			return fmt.Errorf("provider[%d] (%s): provider type is required", i, p.Name)
		}
		if p.Kind == "" {
			return fmt.Errorf("provider[%d] (%s): kind is required", i, p.Name)
		}
		if !validKinds[strings.ToLower(p.Kind)] {
			return fmt.Errorf("provider[%d] (%s): invalid kind %q (must be cloud, pipeline, or sre)", i, p.Name, p.Kind)
		}
		if seen[p.Name] {
			return fmt.Errorf("provider[%d] (%s): duplicate provider name", i, p.Name)
		}
		seen[p.Name] = true
	}

	return nil
}
