package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.API.Address != ":8080" {
		t.Errorf("expected :8080, got %s", cfg.API.Address)
	}
	if !cfg.Insights.Enabled {
		t.Error("expected insights to be enabled by default")
	}
	if cfg.Insights.AnomalyZScore != 2.0 {
		t.Errorf("expected z-score 2.0, got %f", cfg.Insights.AnomalyZScore)
	}
}

func TestValidateEmptyConfig(t *testing.T) {
	cfg := Default()
	if err := Validate(cfg); err != nil {
		t.Errorf("default config should validate: %v", err)
	}
}

func TestValidateMissingName(t *testing.T) {
	cfg := Default()
	cfg.Providers = []ProviderConfig{
		{Provider: "aws", Kind: "cloud"},
	}
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected validation error for missing name")
	}
}

func TestValidateMissingProvider(t *testing.T) {
	cfg := Default()
	cfg.Providers = []ProviderConfig{
		{Name: "test", Kind: "cloud"},
	}
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected validation error for missing provider")
	}
}

func TestValidateMissingKind(t *testing.T) {
	cfg := Default()
	cfg.Providers = []ProviderConfig{
		{Name: "test", Provider: "aws"},
	}
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected validation error for missing kind")
	}
}

func TestValidateInvalidKind(t *testing.T) {
	cfg := Default()
	cfg.Providers = []ProviderConfig{
		{Name: "test", Provider: "aws", Kind: "invalid"},
	}
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected validation error for invalid kind")
	}
}

func TestValidateDuplicateName(t *testing.T) {
	cfg := Default()
	cfg.Providers = []ProviderConfig{
		{Name: "dup", Provider: "aws", Kind: "cloud"},
		{Name: "dup", Provider: "azure", Kind: "cloud"},
	}
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected validation error for duplicate name")
	}
}

func TestValidateValidProviders(t *testing.T) {
	cfg := Default()
	cfg.Providers = []ProviderConfig{
		{Name: "aws-prod", Provider: "aws", Kind: "cloud"},
		{Name: "argocd", Provider: "argocd", Kind: "pipeline"},
		{Name: "pagerduty", Provider: "pagerduty", Kind: "sre"},
	}
	if err := Validate(cfg); err != nil {
		t.Errorf("expected valid config, got: %v", err)
	}
}

func TestLoadWithEnvExpansion(t *testing.T) {
	t.Setenv("TEST_TOKEN", "my-secret-token")

	yaml := `
engine:
  log_level: debug
providers:
  - name: test
    provider: argocd
    kind: pipeline
    config:
      token: "${TEST_TOKEN}"
api:
  address: ":9090"
`
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.API.Address != ":9090" {
		t.Errorf("expected :9090, got %s", cfg.API.Address)
	}
	if cfg.Providers[0].Config["token"] != "my-secret-token" {
		t.Errorf("expected env var expansion, got %v", cfg.Providers[0].Config["token"])
	}
}

func TestResolveEnvValue(t *testing.T) {
	t.Setenv("MY_TOKEN", "secret123")

	config := map[string]any{
		"server":    "https://example.com",
		"token_env": "MY_TOKEN",
	}

	if got := ResolveEnvValue(config, "server"); got != "https://example.com" {
		t.Errorf("expected direct value, got %s", got)
	}
	if got := ResolveEnvValue(config, "token"); got != "secret123" {
		t.Errorf("expected env value, got %s", got)
	}
	if got := ResolveEnvValue(config, "missing"); got != "" {
		t.Errorf("expected empty string, got %s", got)
	}
}
