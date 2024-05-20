package globals

// Status represents the status of a service in the cluster
type Status int

const (
	StatusIdle       Status = iota // When the service is idle
	StatusStarting                 // When the service is starting
	StatusHealthy                  // When the service is healthy
	StatusStopping                 // When the service is stopping
	StatusStopped                  // When the service is stopped
	StatusShutdown                 // When the service is shutting down
	StatusErrored                  // When the service has encountered an error
	StatusRestarting               // When the service is restarting
)

func (s Status) String() string {
	switch s {
	case StatusIdle:
		return "Idle"
	case StatusStarting:
		return "Starting"
	case StatusHealthy:
		return "Healthy"
	case StatusStopping:
		return "Stopping"
	case StatusStopped:
		return "Stopped"
	case StatusShutdown:
		return "Shutdown"
	case StatusErrored:
		return "Errored"
	case StatusRestarting:
		return "Restarting"
	default:
		return "Unknown"
	}
}
