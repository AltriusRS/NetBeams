package globals

type State int

const (
	Unknown State = iota
	Identify
	Authenticate
	Download
	PingOnly
	Password
)

func (s State) String() string {
	switch s {
	case Unknown:
		return "Unknown"
	case Identify:
		return "Identify"
	case Authenticate:
		return "Authenticate"
	case Download:
		return "Download"
	case PingOnly:
		return "PingOnly"
	case Password:
		return "Password"
	default:
		return "Unknown"
	}
}
