package logs

type LogLevel uint8

const (
	LogLevelDebug LogLevel = 0 // Debug level
	LogLevelInfo  LogLevel = 1 //  Info level
	LogLevelWarn  LogLevel = 2 //  Warn level
	LogLevelError LogLevel = 3 // Error level
	LogLevelFatal LogLevel = 4 // Fatal level
)

func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "debug"
	case LogLevelInfo:
		return "info"
	case LogLevelWarn:
		return "warn"
	case LogLevelError:
		return "error"
	case LogLevelFatal:
		return "fatal"
	default:
		return "unknown"
	}
}

func FromString(l string) LogLevel {
	switch l {
	case "debug":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	case "warn":
		return LogLevelWarn
	case "error":
		return LogLevelError
	case "fatal":
		return LogLevelFatal
	default:
		return LogLevelInfo
	}
}
