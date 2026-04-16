package gcp

import "github.com/deepnative/engine/internal/config"

// AuthConfig holds GCP authentication settings.
type AuthConfig struct {
	Project           string
	CredentialsFile   string
	ServiceAccountKey string
}

func resolveAuth(cfg map[string]any) *AuthConfig {
	return &AuthConfig{
		Project:           config.ResolveEnvValue(cfg, "project"),
		CredentialsFile:   config.ResolveEnvValue(cfg, "credentials_file"),
		ServiceAccountKey: config.ResolveEnvValue(cfg, "service_account_key"),
	}
}
