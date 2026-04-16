package config

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

var envVarPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// Load reads and parses a YAML config file, expanding environment variables.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	// Expand environment variables in the form ${VAR}
	expanded := envVarPattern.ReplaceAllFunc(data, func(match []byte) []byte {
		varName := envVarPattern.FindSubmatch(match)[1]
		if val, ok := os.LookupEnv(string(varName)); ok {
			return []byte(val)
		}
		return match
	})

	cfg := Default()
	if err := yaml.Unmarshal(expanded, cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if err := Validate(cfg); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	return cfg, nil
}

// ResolveEnvValue resolves a config value that may reference an environment variable.
// If the key ends with "_env", the value is treated as an environment variable name.
func ResolveEnvValue(config map[string]any, key string) string {
	// Check for direct value first
	if val, ok := config[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}

	// Check for _env variant
	envKey := key + "_env"
	if val, ok := config[envKey]; ok {
		if envVar, ok := val.(string); ok {
			return os.Getenv(envVar)
		}
	}

	return ""
}
