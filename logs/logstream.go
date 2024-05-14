package logs

import (
	"context"
	"time"

	"github.com/valkey-io/valkey-go"
)

// LogStream represents a log stream
type LogStream struct {
	net valkey.Client
}

// NewLogStream creates a new log stream
func NewLogStream(opts valkey.ClientOption) LogStream {
	client, err := valkey.NewClient(opts)
	if err != nil {
		panic(err)
	}

	return LogStream{
		net: client,
	}
}

func (l *LogStream) Log(entry Log) {
	response := l.net.Do(
		context.Background(),
		l.net.B().Xadd().Key("logstream").Id("*").FieldValue().FieldValue("level", entry.Level.String()).FieldValue("module", entry.Module).FieldValue("message", entry.Message).FieldValue("time", entry.Time.Format(time.RFC3339)).FieldValue("machineid", entry.MachineID).FieldValue("hostname", entry.Hostname).Build(),
	)

	if response.Error() != nil {
		panic(response.Error())
	}
}

func (l *LogStream) Terminate() {
	l.net.Close()
}
