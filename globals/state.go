package globals

// State represents the state of the client
type State int

const (
	StateUnknown      State = iota // When the state is unknown (Like on connection)
	StateIdentify                  // When the client needs identification
	StateAuthenticate              // When the client is authenticating
	StateDownload                  // When the client is trying to download files
	StatePingOnly                  // When the client needs to ping the server
	StatePassword                  // When the client needs to enter a password
	StateMapLoad                   // When the client is loading a map
	StatePlaying                   // When the client is playing
	StateDisconnected              // When the client is disconnected
)

func (s State) String() string {
	switch s {
	case StateUnknown:
		return "Unknown"
	case StateIdentify:
		return "Identify"
	case StateAuthenticate:
		return "Authenticate"
	case StateDownload:
		return "Download"
	case StatePingOnly:
		return "PingOnly"
	case StatePassword:
		return "Password"
	case StateMapLoad:
		return "MapLoad"
	case StatePlaying:
		return "Playing"
	default:
		return "Unknown"
	}
}
