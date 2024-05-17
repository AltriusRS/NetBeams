package globals

type State int

const (
	Unknown State = iota
	Identify
	Authenticate
	Download
	PingOnly
)

func (s State) String() string {
	switch s {
	case Unknown:
		return "Unknown"
	case Identify:
		return "Identify"
	case Authenticate:
		return "Authenticate"
	default:
		return "Unknown"
	}
}
