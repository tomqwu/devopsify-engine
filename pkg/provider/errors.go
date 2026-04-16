package provider

import "errors"

var (
	// ErrProviderNotFound is returned when a provider is not registered.
	ErrProviderNotFound = errors.New("provider not found")

	// ErrProviderInit is returned when a provider fails to initialize.
	ErrProviderInit = errors.New("provider initialization failed")

	// ErrUnsupportedOperation is returned when a provider does not support a requested operation.
	ErrUnsupportedOperation = errors.New("unsupported operation")

	// ErrProviderUnhealthy is returned when a provider health check fails.
	ErrProviderUnhealthy = errors.New("provider unhealthy")
)
