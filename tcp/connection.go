package tcp

import (
	"net"
	"netbeams/config"
	"netbeams/environment"
	"netbeams/globals"
	"netbeams/http"
	"netbeams/logs"
	"strconv"
	"time"

	"github.com/Masterminds/semver/v3"
)

type TCPConnection struct {
	Address string
	Conn    net.Conn
	Parent  *Server
	Logger  logs.Logger
	State   globals.State
	Player  *http.Player
}

func NewTCPConnection(conn net.Conn, addr string, parent *Server, l *logs.Logger) TCPConnection {
	return TCPConnection{
		Address: addr,
		Conn:    conn,
		Parent:  parent,
		Logger:  l.Fork("TCP-" + addr),
		State:   globals.StateUnknown,
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

	defer c.Logger.Info("Connection closed")
	defer c.Close()

	defer c.Logger.Terminate()

	// Identify and authenticate the connection
	c.Identify()

	// Sync mod data and server info to the client
	c.SyncModData()

	for c.State == globals.StatePlaying {
		breakLoop := c.RuntimeLoop()
		if breakLoop {
			c.Kick("Connection closed by server")
			return
		}
	}

}

func (c *TCPConnection) Write(data Packet) {
	c.Logger.Debugf("Writing to connection %s - %d bytes", c.Address, data.Header)

	_, err := c.Conn.Write(data.Serialize())
	if err != nil {
		c.Logger.Error("Error writing to connection - Additional output below")
		c.Logger.Error(err.Error())
	}
}

func (c *TCPConnection) Identify() {
	c.SetState(globals.StateIdentify)

	// Read the first message
	sState := make([]byte, 1)
	_, err := c.Conn.Read(sState)
	if err != nil {
		c.Logger.Error("Error reading from connection - Additional output below")
		c.Logger.Error(err.Error())
		return
	}

	c.Logger.Debugf("Received state: %s", sState)

	switch sState[0] {
	case 'C':
		c.SetState(globals.StateAuthenticate)
		c.Authenticate()
	case 'D':
		c.SetState(globals.StateDownload)
	case 'P':
		c.Write(NewPacket("P"))
		c.SetState(globals.StatePingOnly)
	default:
		c.Logger.Error("Unknown starting state - Disconnecting - Additional output below")
		c.Logger.Errorf("Unknown starting state: %s", sState)
		c.Kick("Unknown starting state")
		return
	}
}

func (c *TCPConnection) Authenticate() {
	c.SetState(globals.StateAuthenticate)

	// Read the version information from the client
	packet, err := ReadPacket(c.Conn)

	if err != nil {
		c.Kick("Unable to read data")
		c.Logger.Error("Error authenticating - Additional output below")
		c.Logger.Fatal(err)
		return
	}

	text := packet.ToString()

	rawVersion := text[2:]

	// Parse the version provided by the client
	version, err := semver.NewVersion(rawVersion)

	if err != nil {
		c.Kick("Unable to parse version")
		c.Logger.Error("Error authenticating - Additional output below")
		c.Logger.Fatal(err)
		c.Logger.Error(rawVersion)
		return
	}

	if !environment.Context.SemverMaxClientVersion.Check(version) {
		c.Kick("Client version is too old")
		c.Logger.Error("Error authenticating - Additional output below")
		c.Logger.Fatal(err)
		return
	}

	if !environment.Context.SemverMinClientVersion.Check(version) {
		c.Kick("Client version is too new")
		c.Logger.Error("Error authenticating - Additional output below")
		c.Logger.Fatal(err)
		return
	}

	c.Logger.Debugf("Client version: %s - Continuing authentication", version)

	// The client version is valid, we can now read the authentication key
	c.Write(NewPacket("A"))

	packet, err = ReadPacket(c.Conn)

	if err != nil {
		c.Kick("Unable to read data")
		c.Logger.Error("Error authenticating - Additional output below")
		c.Logger.Fatal(err)
		return
	}

	key := packet.ToString()

	if len(key) > globals.MaxAuthKeyLength {
		c.Kick("Authentication key is too long")
		c.Logger.Error("Error authenticating - Additional output below")
		c.Logger.Fatal(err)
		return
	}

	c.Logger.Debugf("Authentication key: %s", key)

	player, err := globals.App.GetService("BeamMP API").(*http.API).AuthenticatePlayer(key)

	c.Player = player

	if err != nil {
		c.Kick("Unable to authenticate player")
		c.Logger.Error("Error authenticating - Additional output below")
		c.Logger.Fatal(err)
		return
	}

	if player == nil {
		c.Kick("Unable to authenticate player")
		c.Logger.Error("Error authenticating - Additional output below")
		c.Logger.Fatal(err)
		return
	}

	c.Logger.Debugf("Player: %s", player.Name)
	c.Logger.Debugf("UID: %s", player.Uid)
	c.Logger.Debugf("Roles: %s", player.Roles)
	c.Logger.Debugf("Identifiers: %s", player.Identifiers)
	c.Logger.Debugf("Is Guest?: %t", player.Guest)
	c.Logger.Infof("Changing logger ID to %s", player.Name)
	c.Logger.Module = player.Name

	if config.Configuration.General.Password != "" {
		success := c.HandlePassword()

		if !success {
			c.Kick("Unable to authenticate player")
			c.Logger.Error("Error authenticating - Failed to send valid password")
			return
		}
	}

}

// HandlePassword handles the password authentication
// TODO: Add password authentication support when that is better understood
func (c *TCPConnection) HandlePassword() bool {
	// Since passwords are not supported, we can merely return false,
	// to indicate that the password was not accepted
	return false

	// c.SetState(globals.StatePassword)

	// c.Logger.Debug("Sending password request")
	// c.Write(NewPacket("S"))
	// c.SetState(globals.Password)

	// // Read the password from the client
	// packet, err := ReadPacket(c.Conn)

	// if err != nil {
	// 	c.Kick("Unable to read data")
	// 	c.SetStatus(globals.Errored)
	// 	c.Logger.Error("Error authenticating - Additional output below")
	// 	c.Logger.Fatal(err)
	// 	return false
	// }

	// password := packet.ToString()

	// c.Logger.Debugf("Password: %s", password)

	// // if password != config.Configuration.General.Password {
	// // 	c.Kick("Invalid password")
	// // 	c.SetStatus(globals.Errored)
	// // 	c.Logger.Error("Error authenticating - Invalid password")
	// // 	return false
	// // }

	// // wait 10 secodns before closing
	// time.Sleep(10 * time.Second)

	// return false
}

func (c *TCPConnection) SyncModData() {
	c.Logger.Debug("Client is preparing to sync mod data")

	playerCount := len(c.Parent.Connections) - 1

	c.SetState(globals.StateDownload)

	c.Write(NewPacket("P" + strconv.Itoa(playerCount)))

	pauseStart := time.Now()

	for {
		packet, err := ReadPacket(c.Conn)
		if err != nil {
			if time.Since(pauseStart) > 5*time.Second {
				c.Kick("Unable to read data")
				c.Logger.Error("Error reading from connection - Additional output below")
				c.Logger.Error(err.Error())
				return
			} else {
				c.Kick("Unable to read data")
				c.Logger.Error("Error reading from connection - Additional output below")
				c.Logger.Error(err.Error())
				return
			}

		}

		c.Logger.Debugf("Received packet: %v", packet)

		if packet.IsEmpty() {
			c.Logger.Error("Failed to read packet from client - Malformed?")
			break
		} else if packet.Code(0) == 'f' {
			// The client is requesting a file
		} else if packet.Code(0) == 'S' {
			if packet.Code(1) == 'R' {
				// the client is requesting mod data
				// Since we do not support mods, we send an empty mod list
				c.Write(NewPacket("-"))
			} else {
				c.Logger.Error("The client sent an unknown request.")
				c.Kick("The client sent an unknown request.")
				break
			}
		} else if packet.String() == "Done" {
			c.Logger.Debug("Client mod list synced")
			c.SetState(globals.StateMapLoad)
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	if c.State != globals.StateMapLoad {
		c.Kick("Unable to sync mod data")
		return
	}

	c.Logger.Debug("Sending map files")

	c.Write(NewPacket("M" + config.Configuration.General.Map))

	packet, err := ReadPacket(c.Conn)

	if err != nil {
		c.Logger.Error("Error reading from connection - Additional output below")
		c.Logger.Error(err.Error())
		return
	}

	if packet.Code(0) == 'H' {
		c.Logger.Info("Client is connected and loaded")
		c.SetState(globals.StatePlaying)
	} else {
		c.Logger.Warn("Client may not be loaded - Unrecognized map load response")
	}
}

func (c *TCPConnection) Close() {
	c.Logger.Info("Closing connection")
	if c.Conn != nil {
		c.Kick("Server shutting down")
		c.Conn.Close()
		c.SetState(globals.StateDisconnected)
	}
	delete(c.Parent.Connections, c.Address)
}

// Kick a connection with a given message
func (c *TCPConnection) Kick(msg string) {
	c.Logger.Infof("Kicking connection %s", c.Address)
	c.Logger.Infof("Reason: %s", msg)

	c.Write(NewPacket("K" + msg)) // Kick the connection
	// c.Close()
}

// Main gamemplay loop for the connection
func (c *TCPConnection) RuntimeLoop() bool {
	c.Conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	packet, err := ReadPacket(c.Conn)

	if packet.IsEmpty() {
		if err != nil {
			switch err.Error() {

			// The packet is empty, because the connection was EOF
			case "EOF":
				// Sleep the thread for a second to allow the client to begin transmitting
				time.Sleep(time.Second)
				return false

			// The connection timed out
			case "i/o timeout":
				c.Logger.Debug("I/O Timeout Err - Ignoring")
				return false

			// Handle uncaught error cases
			default:
				c.Logger.Error("Error reading from connection - Additional output below")
				c.Logger.Error(err.Error())
				return true
			}
		}

	}

	c.GameplayParser(packet)

	return false
}

func (c *TCPConnection) GameplayParser(Packet Packet) {
	c.Logger.Debugf("Received packet: %s", Packet.data)
}
