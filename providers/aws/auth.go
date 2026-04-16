package aws

import (
	"fmt"

	"github.com/deepnative/engine/internal/config"
)

// AuthMethod represents the AWS authentication method.
type AuthMethod string

const (
	AuthMethodAccessKey   AuthMethod = "access_key"
	AuthMethodProfile     AuthMethod = "profile"
	AuthMethodRole        AuthMethod = "role"
	AuthMethodEnvironment AuthMethod = "environment"
)

// AuthConfig holds AWS authentication settings.
type AuthConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Region          string
	Profile         string
	RoleARN         string
	Method          AuthMethod
}

func resolveAuth(cfg map[string]any) (*AuthConfig, error) {
	auth := &AuthConfig{
		Region: getStringOr(cfg, "region", "us-east-1"),
	}

	// Check for access key auth
	accessKey := config.ResolveEnvValue(cfg, "access_key_id")
	secretKey := config.ResolveEnvValue(cfg, "secret_access_key")
	if accessKey != "" && secretKey != "" {
		auth.AccessKeyID = accessKey
		auth.SecretAccessKey = secretKey
		auth.SessionToken = config.ResolveEnvValue(cfg, "session_token")
		auth.Method = AuthMethodAccessKey
		return auth, nil
	}

	// Check for profile auth
	if profile, ok := cfg["profile"].(string); ok && profile != "" {
		auth.Profile = profile
		auth.Method = AuthMethodProfile
		return auth, nil
	}

	// Check for role assumption
	if roleARN, ok := cfg["role_arn"].(string); ok && roleARN != "" {
		auth.RoleARN = roleARN
		auth.Method = AuthMethodRole
		return auth, nil
	}

	// Default to environment-based auth
	auth.Method = AuthMethodEnvironment
	return auth, nil
}

func validateAuth(auth *AuthConfig) error {
	if auth.Region == "" {
		return fmt.Errorf("region is required")
	}
	if auth.Method == AuthMethodAccessKey {
		if auth.AccessKeyID == "" || auth.SecretAccessKey == "" {
			return fmt.Errorf("access_key_id and secret_access_key are required for access_key auth")
		}
	}
	return nil
}

func getStringOr(cfg map[string]any, key, defaultVal string) string {
	if v, ok := cfg[key].(string); ok && v != "" {
		return v
	}
	return defaultVal
}
