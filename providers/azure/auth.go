package azure

import "github.com/deepnative/engine/internal/config"

// AuthConfig holds Azure authentication settings.
type AuthConfig struct {
	SubscriptionID string
	TenantID       string
	ClientID       string
	ClientSecret   string
}

func resolveAuth(cfg map[string]any) *AuthConfig {
	return &AuthConfig{
		SubscriptionID: config.ResolveEnvValue(cfg, "subscription_id"),
		TenantID:       config.ResolveEnvValue(cfg, "tenant_id"),
		ClientID:       config.ResolveEnvValue(cfg, "client_id"),
		ClientSecret:   config.ResolveEnvValue(cfg, "client_secret"),
	}
}
