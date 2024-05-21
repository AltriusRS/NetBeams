package logs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/altriusrs/netbeams/src/environment"
	"github.com/denisbrodbeck/machineid"
)

var globalLogger *Logger

// Logger represents a logger
type Logger struct {
	// The lowest level of messages that will be logged
	Level LogLevel

	// The module that the logger is for
	Module string

	// The valkey client
	Net *LogStream

	// A unique identifier for the machine
	MachineID string

	// A short identifier based on the machine ID
	ShortId string

	// The hostname of the machine
	Hostname string
}

var internal_log_level *LogLevel

// Creates a new logger without a valkey client instance
func OfflineLogger(module string) Logger {

	if globalLogger != nil {
		return globalLogger.Fork(module)
	}

	if module == "" {
		panic("module cannot be empty")
	}

	machineId, err := machineid.ProtectedID("NetBeams")

	if err != nil {
		panic(err)
	}

	hostname, err := os.Hostname()

	if err != nil {
		panic(err)
	}

	shortId := machineId[5:15]

	if internal_log_level != nil {
		globalLogger = &Logger{
			Level:     *internal_log_level,
			Module:    module,
			MachineID: machineId,
			ShortId:   shortId,
			Hostname:  hostname,
		}
	} else {
		globalLogger = &Logger{
			Level:     LogLevelInfo,
			Module:    module,
			MachineID: machineId,
			ShortId:   shortId,
			Hostname:  hostname,
		}
	}

	return globalLogger.Fork(module)
}

// Creates a new logger with a valkey client instance
func NetLogger(module string) Logger {

	if globalLogger != nil {
		return globalLogger.Fork(module)
	}

	if module == "" {
		panic("module cannot be empty")
	}

	l := OfflineLogger("LogPrimer") // Since we're using this as a temporary logger
	l.Debug("Creating logger")

	if internal_log_level == nil {
		l.Debug("Loading environment variables")

		err := environment.LoadEnvFile(".dev.env")

		if err != nil {
			l.Warn(err.Error())
			l.Debug("Error loading .dev.env file - Checking production environment")
			err = environment.LoadEnvFile(".env")

			if err != nil {
				l.Warn(err.Error())
				l.Debug("Error loading .dev.env file - Checking production environment")
			}
		}
	}

	log_level, present := os.LookupEnv("NETBEAMS_LOG_LEVEL")

	var level LogLevel

	if !present || log_level == "" {
		if internal_log_level != nil {
			level = *internal_log_level
		} else {
			l.Warn("Failed to get log level - Falling back to INFO")
			level = LogLevelInfo
			internal_log_level = &level
		}
	} else {
		level = FromString(strings.ToLower(log_level))
		internal_log_level = &level
	}

	l.Level = level

	l.Debug("Log level set to " + level.String())

	l.Debug("Fetching machine ID")
	machineId, err := machineid.ProtectedID("NetBeams")

	if err != nil {
		l.Error("Failed to get machine ID - Unable to initialise")
		panic(err)
	}
	l.Debug("Fetching hostname")

	hostname, err := os.Hostname()

	if err != nil {
		l.Error("Failed to get hostname - Unable to initialise")
		panic(err)
	}

	shortId := machineId[5:15]

	net, err := NewLogStream()

	if err != nil {
		l.Error("Failed to create valkey client - Unable to initialise")
		return l
	}

	globalLogger = &Logger{
		Level:     level,
		Module:    module,
		MachineID: machineId,
		ShortId:   shortId,
		Hostname:  hostname,
		Net:       net,
	}

	return globalLogger.Fork(module)
}

// Log a message
func (l *Logger) Log(level LogLevel, message string) {
	log := Log{
		Level:     level,
		Module:    l.Module,
		Message:   message,
		Time:      time.Now(),
		MachineID: l.MachineID,
		ShortId:   l.ShortId,
		Hostname:  l.Hostname,
	}

	// If the log level is lower than the current log level, do nothing
	if level < l.Level && !environment.Context.IsDev {
		return
	}

	fmt.Println(log.String())

	if l.Net != nil {
		_ = l.Net.Log(log) // Ignore write errors to valkey as they are not critical and could cause spam
	}
}

// Fork a new logger from the existing instance (minimises the amount of data logged when creating a new logger)
func (l *Logger) Fork(module string) Logger {
	return Logger{
		Level:     l.Level,
		Module:    module,
		MachineID: l.MachineID,
		Hostname:  l.Hostname,
		ShortId:   l.ShortId,
		Net:       l.Net,
	}

}

// logs a formatted message
func (l *Logger) Logf(level LogLevel, format string, args ...any) {
	l.Log(level, fmt.Sprintf(format, args...))
}

// Log a debug message
func (l *Logger) Debug(message string) {
	l.Logf(LogLevelDebug, message)
}

// Log a formatted debug message
func (l *Logger) Debugf(format string, args ...any) {
	l.Logf(LogLevelDebug, format, args...)
}

// Log an info message
func (l *Logger) Info(message string) {
	l.Logf(LogLevelInfo, message)
}

// Log a formatted info message
func (l *Logger) Infof(format string, args ...any) {
	l.Logf(LogLevelInfo, format, args...)
}

// Log a warning message
func (l *Logger) Warn(message string) {
	l.Logf(LogLevelWarn, message)
}

// Log a formatted warning message
func (l *Logger) Warnf(format string, args ...any) {
	l.Logf(LogLevelWarn, format, args...)
}

// Log an error message
func (l *Logger) Error(message string) {
	l.Logf(LogLevelError, message)
}

// Log a formatted error message
func (l *Logger) Errorf(format string, args ...any) {
	l.Logf(LogLevelError, format, args...)
}

// Log a fatal message
func (l *Logger) Fatal(message error) {
	l.Logf(LogLevelFatal, message.Error())
}

func (l *Logger) Terminate() {
	if l.Net != nil {
		l.Net.Terminate()
	}
}
