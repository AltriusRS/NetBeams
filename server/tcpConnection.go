package server

import (
	"net"
	"netbeams/logs"
	"strings"
	"time"
)

type TCPConnection struct {
	Conn   net.Conn
	Parent *TCPServer
	Logger logs.Logger
	Status Status
}

func NewTCPConnection(conn net.Conn, addr string, parent *TCPServer, l *logs.Logger) TCPConnection {
	return TCPConnection{
		Conn:   conn,
		Parent: parent,
		Logger: l.Fork("TCP-" + addr),
	}
}

func (c *TCPConnection) Listen() {
	c.Logger.Info("Listening for messages")

	if c.Status != Healthy {
		c.Logger.Info("Connection is not healthy")
		return
	}

	defer c.Logger.Info("Connection closed")

	defer c.Logger.Terminate()

	if c.Conn == nil {
		c.Logger.Error("Connection is nil")
		c.Status = Errored
		return
	}

	for {

		HeaderBytes := make([]byte, 4)

		c.Logger.Debug("Waiting for message")
		c.Conn.SetDeadline(time.Now().Add(time.Second))
		if c.Status != Healthy {
			c.Logger.Error("Connection is not healthy")
			return
		}
		_, err := c.Conn.Read(HeaderBytes)

		if err != nil {
			if strings.HasSuffix(err.Error(), "i/o timeout") {
				c.Logger.Debug("Message polling timed out")
				continue
			}
			c.Logger.Error("Error reading message header - Additional output below")
			c.Logger.Fatal(err)
			return
		}
	}
}
