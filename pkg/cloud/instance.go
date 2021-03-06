package cloud

import "context"

type ProviderStatus string

var (
	StatusRunning  ProviderStatus = "RUNNING"
	StatusStopped  ProviderStatus = "STOPPED"
	StatusStopping ProviderStatus = "STOPPING"
	StatusStarting ProviderStatus = "STARTING"
	StatusUnknown  ProviderStatus = "UNKNOWN"
)

type Provider interface {
	// Status fetches the status of a remote instance
	Status(ctx context.Context, instanceID string) (ProviderStatus, error)

	// Start starts a remote instance
	Start(ctx context.Context, instanceID string) error
}
