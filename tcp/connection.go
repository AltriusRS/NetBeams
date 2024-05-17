package tcp

import (
	"encoding/binary"
	"fmt"
	"net"
	"netbeams/globals"
	"netbeams/logs"
	"time"
)

type TCPConnection struct {
	Address string
	Conn    net.Conn
	Parent  *Server
	Logger  logs.Logger
	Status  globals.Status
	State   globals.State
}

func NewTCPConnection(conn net.Conn, addr string, parent *Server, l *logs.Logger) TCPConnection {
	return TCPConnection{
		Address: addr,
		Conn:    conn,
		Parent:  parent,
		Logger:  l.Fork("TCP-" + addr),
		Status:  globals.Healthy,
	}
}

func (c *TCPConnection) SetStatus(status globals.Status) {
	if c.Status != status {
		c.Logger.Infof("Connection %s status changed from %s to %s", c.Address, c.Status, status)
		c.Status = status
	}
}

func (c *TCPConnection) SetState(state globals.State) {
	if c.State != state {
		c.Logger.Infof("Connection %s state changed from %s to %s", c.Address, c.State, state)
		c.State = state
	}
}

func (c *TCPConnection) Listen() {
	c.Logger.Info("Listening for messages")

	if c.Status != globals.Healthy {
		c.Logger.Info("Connection is not healthy")
		return
	}

	defer c.Logger.Info("Connection closed")

	defer c.Logger.Terminate()

	if c.Conn == nil {
		c.Logger.Error("Connection is nil")
		c.SetStatus(globals.Errored)
		return
	}

	c.Identify()

	for {

		switch c.Status {
		case globals.Kicked:
			c.Logger.Info("Connection is kicked")
			return

		case globals.Closed:
			c.Logger.Info("Connection is closed")
			return

		case globals.Errored:
			c.Logger.Info("Connection is errored")
			return
		}

		switch c.State {
		case globals.Unknown:
			c.SetState(globals.Identify)

		case globals.Identify:
			err := c.Identify()
			if err != nil {
				c.Kick("Unable to identify")
				return
			}

		case globals.Authenticate:
			c.Logger.Info("Authenticating")
			err := c.Authenticate()
			if err != nil {
				c.Kick("Unable to authenticate")
				return
			}

		default:
			c.Logger.Warnf("Unknown state: %s", c.State)
			c.Kick("Unknown state")
			return
		}

	}
}

func (c *TCPConnection) Write(data []byte) {
	c.Logger.Debugf("Writing to connection %s - %d bytes", c.Address, len(data))

	header := int32(len(data))
	headerBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(headerBytes, uint32(header))

	packet := append(headerBytes, data...)

	_, err := c.Conn.Write(packet)
	if err != nil {
		c.SetStatus(globals.Errored)
		c.Logger.Error("Error writing to connection - Additional output below")
		c.Logger.Error(err.Error())
	}
}

// Read from the connection
func (c *TCPConnection) Read() ([]byte, error) {
	c.Conn.SetDeadline(time.Now().Add(time.Second))

	packet, err := ReadPacket(c.Conn)

	c.Logger.Debugf("Received packet: %s", packet.data)

	return nil, err
}

func (c *TCPConnection) Identify() error {
	// Allow the connection to hang for 5 seconds to allow for latency issues on startup
	c.Conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Read the first message
	sState := make([]byte, 1)
	_, err := c.Conn.Read(sState)
	if err != nil {
		c.SetStatus(globals.Errored)
		c.Logger.Error("Error reading from connection - Additional output below")
		c.Logger.Error(err.Error())
		return err
	}

	c.Logger.Debugf("Received state: %s", sState)

	switch sState[0] {
	case 'C':
		c.SetState(globals.Authenticate)
	case 'D':
		c.SetState(globals.Download)
	case 'P':
		c.Write([]byte("P"))
		c.SetState(globals.PingOnly)
	default:
		c.Kick("Unknown starting state")
		return fmt.Errorf("unknown starting state")
	}

	return nil
}

func (c *TCPConnection) Authenticate() error {
	c.SetState(globals.Authenticate)

	payload, err := c.Read()

	if err != nil {
		c.Kick("Unable to read data")
		c.SetStatus(globals.Errored)
		c.Logger.Error("Error authenticating - Additional output below")
		c.Logger.Fatal(err)
		return err
	}

	c.Logger.Debugf("Received payload: %s", payload)

	return fmt.Errorf("not implemented")

	return nil
}

func (c *TCPConnection) Close() {
	c.Logger.Info("Closing connection")
	if c.Conn != nil {
		c.Conn.Close()
	}
	c.Status = globals.Closed
	delete(c.Parent.Connections, c.Address)
}

// Kick a connection with a given message
func (c *TCPConnection) Kick(msg string) {
	if c.Status != globals.Healthy {
		c.Logger.Warn("Tried to kick a connection which is not healthy")
		return
	}

	c.SetStatus(globals.Kicked)

	c.Logger.Infof("Kicking connection %s", c.Address)
	c.Logger.Infof("Reason: %s", msg)

	c.Write([]byte("K" + msg)) // Kick the connection
	c.Close()
}
