package logs

import (
	"fmt"
	"netbeams/environment"
	"os"
	"strings"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/valkey-io/valkey-go"
)

// Logger represents a logger
type Logger struct {
	// The lowest level of messages that will be logged
	Level LogLevel

	// The module that the logger is for
	Module string

	// The valkey client
	Net *LogStream

	// The valkey client
	MachineID string

	// The hostname of the machine
	Hostname string
}

var internal_log_level *LogLevel

// Creates a new logger without a valkey client instance
func OfflineLogger(module string) Logger {
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

	if internal_log_level != nil {
		return Logger{
			Level:     *internal_log_level,
			Module:    module,
			MachineID: machineId,
			Hostname:  hostname,
		}
	}

	return Logger{
		Level:     LogLevelDebug,
		Module:    module,
		MachineID: machineId,
		Hostname:  hostname,
	}
}

// Creates a new logger with a valkey client instance
func NetLogger(module string) Logger {
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

	l.Debug("Machine ID: " + machineId)
	l.Debug("Fetching hostname")

	hostname, err := os.Hostname()

	if err != nil {
		l.Error("Failed to get hostname - Unable to initialise")
		panic(err)
	}

	l.Debug("Hostname: " + hostname)

	init_address, present := os.LookupEnv("VALKEY_URI")
	if !present || init_address == "" {
		l.Error("Failed to get VALKEY_URI environment variable - Falling back to offline logging")
		return l
	}

	l.Debug("Valkey URI is present")

	valkey_password, present := os.LookupEnv("VALKEY_PASSWORD")
	if !present || valkey_password == "" {
		l.Error("Failed to get VALKEY_PASSWORD environment variable - Falling back to offline logging")
		return l
	}
	l.Debug("Valkey Password is present")

	l.Debug("Creating valkey client")
	net := NewLogStream(valkey.ClientOption{InitAddress: []string{init_address}, Password: valkey_password, SelectDB: 0})

	return Logger{
		Module:    module,
		MachineID: machineId,
		Hostname:  hostname,
		Net:       &net,
		Level:     level,
	}
}

// Log a message
func (l *Logger) Log(level LogLevel, message string) {
	log := Log{
		Level:     level,
		Module:    l.Module,
		Message:   message,
		Time:      time.Now(),
		MachineID: l.MachineID,
		Hostname:  l.Hostname,
	}

	// If the log level is lower than the current log level, do nothing
	if level < l.Level {
		return
	}

	fmt.Println(log.String())

	if l.Net != nil {
		l.Net.Log(log)
	}
}

// logs a formatted message
func (l *Logger) Logf(level LogLevel, format string, args ...interface{}) {
	l.Log(level, fmt.Sprintf(format, args...))
}

// Log a debug message
func (l *Logger) Debug(message string) {
	l.Logf(LogLevelDebug, message)
}

// Log an info message
func (l *Logger) Info(message string) {
	l.Logf(LogLevelInfo, message)
}

// Log a warning message
func (l *Logger) Warn(message string) {
	l.Logf(LogLevelWarn, message)
}

// Log an error message
func (l *Logger) Error(message string) {
	l.Logf(LogLevelError, message)
}

// Log a fatal message
func (l *Logger) Fatal(message string) {
	l.Logf(LogLevelFatal, message)
}

func (l *Logger) Terminate() {
	if l.Net != nil {
		l.Net.Terminate()
	}
}
