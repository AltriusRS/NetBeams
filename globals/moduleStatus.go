package globals

type Status int

const (
	Idle Status = iota
	Starting
	Healthy
	Stopping
	Stopped
	Shutdown
	Errored
	Closed
	Kicked
	Disconnected
	Disconnecting
)

func (s Status) String() string {
	switch s {
	case Idle:
		return "Idle"
	case Starting:
		return "Starting"
	case Healthy:
		return "Healthy"
	case Stopping:
		return "Stopping"
	case Stopped:
		return "Stopped"
	case Shutdown:
		return "Shutdown"
	case Errored:
		return "Errored"
	case Closed:
		return "Closed"
	case Kicked:
		return "Kicked"
	case Disconnected:
		return "Disconnected"
	case Disconnecting:
		return "Disconnecting"
	default:
		return "Unknown"
	}
}
