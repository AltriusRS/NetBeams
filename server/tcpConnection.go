package server

import (
	"encoding/binary"
	"fmt"
	"net"
	"netbeams/logs"
	"time"
)

type TCPConnection struct {
	Address string
	Conn    net.Conn
	Parent  *TCPServer
	Logger  logs.Logger
	Status  Status
	State   State
}

func NewTCPConnection(conn net.Conn, addr string, parent *TCPServer, l *logs.Logger) TCPConnection {
	return TCPConnection{
		Address: addr,
		Conn:    conn,
		Parent:  parent,
		Logger:  l.Fork("TCP-" + addr),
		Status:  Healthy,
	}
}

func (c *TCPConnection) SetStatus(status Status) {
	if c.Status != status {
		c.Logger.Infof("Connection %s status changed from %s to %s", c.Address, c.Status, status)
		c.Status = status
	}
}

func (c *TCPConnection) SetState(state State) {
	if c.State != state {
		c.Logger.Infof("Connection %s state changed from %s to %s", c.Address, c.State, state)
		c.State = state
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
		c.SetStatus(Errored)
		return
	}

	for {
		switch c.Status {
		case Kicked:
			c.Logger.Info("Connection is kicked")
			return

		case Closed:
			c.Logger.Info("Connection is closed")
			return

		case Errored:
			c.Logger.Info("Connection is errored")
			return
		}

		switch c.State {
		case Unknown:
			c.SetState(Identify)

		case Identify:
			c.Logger.Info("Identifying")
			err := c.Identify()
			if err != nil {
				c.Kick("Unable to identify")
				return
			}

		case Authenticate:
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
		c.SetStatus(Errored)
		c.Logger.Error("Error writing to connection - Additional output below")
		c.Logger.Error(err.Error())
	}
}

// Read from the connection
func (c *TCPConnection) Read() ([]byte, error) {
	c.Conn.SetDeadline(time.Now().Add(time.Second))
	// HeaderBytes := make([]byte, 4)
	// _, err := c.Conn.Read(HeaderBytes)

	// if err != nil {
	// 	if strings.HasSuffix(err.Error(), "i/o timeout") {
	// 		c.Logger.Debug("Message polling timed out")
	// 		return nil, nil
	// 	}
	// 	c.Logger.Error("Error reading message header - Additional output below")
	// 	c.Logger.Fatal(err)
	// 	return nil, err
	// }

	// c.Logger.Debugf("Received message header %x", HeaderBytes)

	// // Read the header
	// Header := binary.LittleEndian.Uint32(HeaderBytes)

	// if Header > MaxHeaderSize {
	// 	return nil, fmt.Errorf("header size limit exceeded")
	// }

	// c.Logger.Debugf("Received message header %d", Header)

	// if Header >= MaxHeaderSize {
	// 	return nil, fmt.Errorf("header size limit exceeded")
	// }

	// BodyPayload := make([]byte, Header)

	// bytesRead, err := c.Conn.Read(BodyPayload)
	// if err != nil {
	// 	c.Logger.Error("Error reading message body - Additional output below")
	// 	c.Logger.Fatal(err)

	// 	return nil, err
	// }

	// Body := string(BodyPayload)

	// c.Logger.Debug("Received message body")
	// c.Logger.Debug(Body + "|END")

	// if bytesRead < int(Header) {
	// 	c.Logger.Warnf("Message body too short - %d bytes read out of %d", bytesRead, Header)
	// 	return nil, fmt.Errorf("message body too short")
	// }

	packet, err := ReadPacket(c.Conn)

	c.Logger.Debugf("Received packet: %s", packet.data)

	return nil, err
}

func (c *TCPConnection) Identify() error {
	c.SetState(Authenticate)
	return nil
}

func (c *TCPConnection) Authenticate() error {
	c.SetState(Authenticate)

	payload, err := c.Read()

	if err != nil {
		c.Kick("Unable to read data")
		c.SetStatus(Errored)
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
	c.Status = Closed
	delete(c.Parent.Connections, c.Address)
}

// Kick a connection with a given message
func (c *TCPConnection) Kick(msg string) {
	if c.Status != Healthy {
		c.Logger.Warn("Tried to kick a connection which is not healthy")
		return
	}

	c.SetStatus(Kicked)

	c.Logger.Infof("Kicking connection %s", c.Address)
	c.Logger.Infof("Reason: %s", msg)

	c.Write([]byte("K" + msg)) // Kick the connection
	c.Close()
}
