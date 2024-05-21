package logs

import (
	"time"

	"github.com/altriusrs/netbeams/src/keyval"
)

// LogStream represents a log stream
type LogStream struct {
	net *keyval.KeyvalClient
}

// NewLogStream creates a new log stream
func NewLogStream() (*LogStream, error) {
	client, err := keyval.NewKeyvalClient()
	if err != nil {
		return nil, err
	}

	return &LogStream{
		net: client,
	}, nil
}

func ForkLogStream(client *keyval.KeyvalClient) *LogStream {
	return &LogStream{
		net: client,
	}
}

func (l *LogStream) Log(entry Log) error {
	if l.net == nil {
		return nil
	}

	response := l.net.Xadd("logstream", "*",
		[][]string{
			{"level", entry.Level.String()},
			{"module", entry.Module},
			{"message", entry.Message},
			{"time", entry.Time.Format(time.RFC3339)},
			{"machineid", entry.MachineID},
			{"shortid", entry.ShortId},
			{"hostname", entry.Hostname},
		},
	)

	return response.Error()
}

func (l *LogStream) Terminate() {
	l.net.Close()
}
