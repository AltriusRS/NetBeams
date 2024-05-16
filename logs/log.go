package logs

import (
	"fmt"
	"strings"
	"time"
)

// Log represents a log entry
type Log struct {
	// The level of the log entry
	Level LogLevel

	// The module that generated the entry
	Module string

	// The message passed
	Message string

	// The time the entry was generated
	Time time.Time

	// The machine ID - This is used to identify the machine that generated the log entry
	MachineID string

	// The hostname of the machine that generated the log entry
	Hostname string
}

func (l *Log) String() string {
	timeString := l.Time.Local().Format("15:04:05.000")
	upperLevel := strings.ToUpper(l.Level.String())
	logColor := ""

	switch l.Level {
	case LogLevelDebug:
		logColor = "\x1b[36m"
	case LogLevelInfo:
		logColor = "\x1b[34m"
	case LogLevelWarn:
		logColor = "\x1b[33m"
	case LogLevelError:
		logColor = "\x1b[31m"
	case LogLevelFatal:
		logColor = "\x1b[41m\x1b[37m"
	}

	switch l.Level {
	case LogLevelInfo, LogLevelWarn:
		return fmt.Sprintf("\x1b[35m[%s]\x1b[0m %s[ %s]\x1b[0m \x1b[36m[%s]\x1b[0m %s", timeString, logColor, upperLevel, l.Module, l.Message)
	default:
		return fmt.Sprintf("\x1b[35m[%s]\x1b[0m %s[%s]\x1b[0m \x1b[36m[%s]\x1b[0m %s", timeString, logColor, upperLevel, l.Module, l.Message)
	}
}
