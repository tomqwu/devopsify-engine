package config

// Config is the top-level configuration for the Deep Native Engine.
type Config struct {
	Engine    EngineConfig     `yaml:"engine"`
	Providers []ProviderConfig `yaml:"providers"`
	Insights  InsightsConfig   `yaml:"insights"`
	API       APIConfig        `yaml:"api"`
}

// EngineConfig holds engine-level settings.
type EngineConfig struct {
	LogLevel string `yaml:"log_level"`
}

// ProviderConfig holds configuration for a single provider instance.
type ProviderConfig struct {
	Name     string         `yaml:"name"`
	Provider string         `yaml:"provider"`
	Kind     string         `yaml:"kind"`
	Config   map[string]any `yaml:"config"`
}

// InsightsConfig controls the insights engine behavior.
type InsightsConfig struct {
	Enabled           bool    `yaml:"enabled"`
	CostThreshold     float64 `yaml:"cost_threshold"`
	DriftCheckEnabled bool    `yaml:"drift_check_enabled"`
	AnomalyZScore     float64 `yaml:"anomaly_z_score"`
}

// APIConfig holds API server configuration.
type APIConfig struct {
	Address         string `yaml:"address"`
	ReadTimeout     int    `yaml:"read_timeout"`
	WriteTimeout    int    `yaml:"write_timeout"`
	ShutdownTimeout int    `yaml:"shutdown_timeout"`
}

// Default returns a default configuration.
func Default() *Config {
	return &Config{
		Engine: EngineConfig{
			LogLevel: "info",
		},
		Providers: []ProviderConfig{},
		Insights: InsightsConfig{
			Enabled:           true,
			CostThreshold:     100.0,
			DriftCheckEnabled: true,
			AnomalyZScore:     2.0,
		},
		API: APIConfig{
			Address:         ":8080",
			ReadTimeout:     30,
			WriteTimeout:    30,
			ShutdownTimeout: 15,
		},
	}
}
