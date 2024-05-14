package logs

import (
	"fmt"
	"os"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/valkey-io/valkey-go"
)

// Logger represents a logger
type Logger struct {
	// The module that the logger is for
	Module string

	// The valkey client
	net LogStream

	// The valkey client
	MachineID string

	// The hostname of the machine
	Hostname string
}

// NewLogger creates a new logger
func NewLogger(module string) Logger {
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

	return Logger{
		Module:    module,
		net:       NewLogStream(valkey.ClientOption{InitAddress: []string{"127.0.0.1:6379"}, Password: "test_password", SelectDB: 0}),
		MachineID: machineId,
		Hostname:  hostname,
	}
}

func (l *Logger) Log(level string, message string) {
	log := Log{
		Level:     level,
		Module:    l.Module,
		Message:   message,
		Time:      time.Now(),
		MachineID: l.MachineID,
		Hostname:  l.Hostname,
	}

	l.net.Log(log)
	fmt.Println(log.String())
}

func (l *Logger) Logf(level string, format string, args ...interface{}) {
	l.Log(level, fmt.Sprintf(format, args...))
}

func (l *Logger) Debug(message string) {
	l.Logf("debug", message)
}

func (l *Logger) Info(message string) {
	l.Logf("info", message)
}

func (l *Logger) Warn(message string) {
	l.Logf("warn", message)
}

func (l *Logger) Error(message string) {
	l.Logf("error", message)
}

func (l *Logger) Fatal(message string) {
	l.Logf("fatal", message)
}

func (l *Logger) Terminate() {
	l.net.Terminate()
}
