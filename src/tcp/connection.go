package tcp

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/altriusrs/netbeams/src/config"
	"github.com/altriusrs/netbeams/src/environment"
	"github.com/altriusrs/netbeams/src/http"
	"github.com/altriusrs/netbeams/src/logs"
	"github.com/altriusrs/netbeams/src/player_manager"
	"github.com/altriusrs/netbeams/src/types"
)

type TCPConnection struct {
	logs.Logger
	Address string
	Conn    net.Conn
	Parent  *Server
	State   types.State
	Player  types.Player
}

func NewTCPConnection(conn net.Conn, addr string, parent *Server) TCPConnection {
	return TCPConnection{
		Address: addr,
		Conn:    conn,
		Parent:  parent,
		Logger:  logs.NetLogger("TCP-" + addr),
		State:   types.StateUnknown,
	}
}

func (c *TCPConnection) SetState(state types.State) {
	if c.State != state {
		c.Infof("Connection %s state changed from %s to %s", c.Address, c.State, state)
		c.State = state
	}
}

func (c *TCPConnection) Listen() {
	c.Info("Listening for messages")

	defer c.Info("Connection closed")
	defer c.Close()

	defer c.Terminate()

	// Identify and authenticate the connection
	c.Identify()

	// Sync mod data and server info to the client
	c.SyncModData()

	for c.State == types.StatePlaying {
		breakLoop := c.RuntimeLoop()
		if breakLoop {
			c.Kick("Connection closed by server")
			return
		}
	}

}

func (c *TCPConnection) Write(data types.TcpPacket) {
	c.Debugf("Writing to connection %s - %d bytes", c.Address, data.Header)

	_, err := c.Conn.Write(data.Serialize())
	if err != nil {
		c.Error("Error writing to connection - Additional output below")
		c.Error(err.Error())
	}
}

func (c *TCPConnection) Identify() {
	c.SetState(types.StateIdentify)

	// Read the first message
	sState := make([]byte, 1)
	_, err := c.Conn.Read(sState)
	if err != nil {
		c.Error("Error reading from connection - Additional output below")
		c.Error(err.Error())
		return
	}

	c.Debugf("Received state: %s", sState)

	pm := types.App.GetService("Player Manager").(*player_manager.PlayerManager)
	_, err = pm.GetNextID()

	if err != nil {
		if err.Error() == "server is full" {
			c.Kick("Server is full")
			return
		} else {
			c.Kick("The server is experiencing an error - Please try again later")
			c.Error("Error authenticating - Additional output below")
			c.Error(err.Error())
			return
		}
	}

	switch sState[0] {
	case 'C':
		c.SetState(types.StateAuthenticate)
		c.Authenticate()
	case 'D':
		c.SetState(types.StateDownload)
	case 'P':
		c.Write(types.NewTcpPacket("P"))
		c.SetState(types.StatePingOnly)
	default:
		c.Error("Unknown starting state - Disconnecting - Additional output below")
		c.Errorf("Unknown starting state: %s", sState)
		c.Kick("Unknown starting state")
		return
	}
}

func (c *TCPConnection) Authenticate() {
	c.SetState(types.StateAuthenticate)

	// Read the version information from the client
	packet, err := types.ReadTcpPacket(c.Conn)

	if err != nil {
		c.Kick("Unable to read data")
		c.Error("Error authenticating - Additional output below")
		c.Fatal(err)
		return
	}

	text := packet.ToString()

	rawVersion := text[2:]

	// Parse the version provided by the client
	version, err := semver.NewVersion(rawVersion)

	if err != nil {
		c.Kick("Unable to parse version")
		c.Error("Error authenticating - Additional output below")
		c.Fatal(err)
		c.Error(rawVersion)
		return
	}

	if !environment.Context.SemverMaxClientVersion.Check(version) {
		c.Kick("Client version is too old")
		c.Error("Error authenticating - Additional output below")
		c.Fatal(err)
		return
	}

	if !environment.Context.SemverMinClientVersion.Check(version) {
		c.Kick("Client version is too new")
		c.Error("Error authenticating - Additional output below")
		c.Fatal(err)
		return
	}

	c.Debugf("Client version: %s - Continuing authentication", version)

	// The client version is valid, we can now read the authentication key
	c.Write(types.NewTcpPacket("A"))

	packet, err = types.ReadTcpPacket(c.Conn)

	if err != nil {
		c.Kick("Unable to read data")
		c.Error("Error authenticating - Additional output below")
		c.Fatal(err)
		return
	}

	key := packet.ToString()

	if len(key) > types.MaxAuthKeyLength {
		c.Kick("Authentication key is too long")
		c.Error("Error authenticating - Additional output below")
		c.Fatal(err)
		return
	}

	c.Debugf("Authentication key: %s", key)

	player, err := types.App.GetService("BeamMP API").(*http.API).AuthenticatePlayer(key)

	if err != nil {
		c.Kick("Unable to authenticate player")
		c.Error("Error authenticating - Additional output below")
		c.Fatal(err)
		return
	}

	c.Player = player.IntoPlayerEntity()

	pid, err := types.App.GetService("Player Manager").(*player_manager.PlayerManager).ReserveSlot(&c.Player)

	if err != nil {
		c.Kick("Unable to reserve slot")
		c.Error("Error authenticating - Additional output below")
		c.Fatal(err)
		return
	}

	c.Player.PlayerId = *pid

	if player == nil {
		c.Kick("Unable to authenticate player")
		c.Error("Error authenticating - Additional output below")
		c.Fatal(err)
		return
	}

	c.Debugf("Player: %s", player.Name)
	c.Debugf("UID: %s", player.Uid)
	c.Debugf("Roles: %s", player.Roles)
	c.Debugf("Identifiers: %s", player.Identifiers)
	c.Debugf("Is Guest?: %t", player.Guest)
	c.Infof("Changing logger ID to %s", player.Name)
	c.Module = player.Name

	if config.Configuration.General.Password != "" {
		success := c.HandlePassword()

		if !success {
			c.Kick("Unable to authenticate player")
			c.Error("Error authenticating - Failed to send valid password")
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

	// c.SetState(types.StatePassword)

	// c.Debug("Sending password request")
	// c.Write(types.NewTcpPacket("S"))
	// c.SetState(types.Password)

	// // Read the password from the client
	// packet, err := ReadPacket(c.Conn)

	// if err != nil {
	// 	c.Kick("Unable to read data")
	// 	c.SetStatus(types.Errored)
	// 	c.Error("Error authenticating - Additional output below")
	// 	c.Fatal(err)
	// 	return false
	// }

	// password := packet.ToString()

	// c.Debugf("Password: %s", password)

	// // if password != config.Configuration.General.Password {
	// // 	c.Kick("Invalid password")
	// // 	c.SetStatus(types.Errored)
	// // 	c.Error("Error authenticating - Invalid password")
	// // 	return false
	// // }

	// // wait 10 secodns before closing
	// time.Sleep(10 * time.Second)

	// return false
}

func (c *TCPConnection) SyncModData() {
	c.Debug("Client is preparing to sync mod data")

	playerCount := len(c.Parent.Connections) - 1

	c.SetState(types.StateDownload)

	c.Write(types.NewTcpPacket("P" + strconv.Itoa(playerCount)))

	pauseStart := time.Now()

	for {
		packet, err := types.ReadTcpPacket(c.Conn)
		if err != nil {
			if time.Since(pauseStart) > 5*time.Second {
				c.Kick("Unable to read data")
				c.Error("Error reading from connection - Additional output below")
				c.Error(err.Error())
				return
			} else {
				c.Kick("Unable to read data")
				c.Error("Error reading from connection - Additional output below")
				c.Error(err.Error())
				return
			}

		}

		c.Debugf("Received packet: %v", packet)

		if packet.IsEmpty() {
			c.Error("Failed to read packet from client - Malformed?")
			break
		} else if packet.Code(0) == 'f' {
			// The client is requesting a file
		} else if packet.Code(0) == 'S' {
			if packet.Code(1) == 'R' {
				// the client is requesting mod data
				// Since we do not support mods, we send an empty mod list
				c.Write(types.NewTcpPacket("-"))
			} else {
				c.Error("The client sent an unknown request.")
				c.Kick("The client sent an unknown request.")
				break
			}
		} else if packet.String() == "Done" {
			c.Debug("Client mod list synced")
			c.SetState(types.StateMapLoad)
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	if c.State != types.StateMapLoad {
		c.Kick("Unable to sync mod data")
		return
	}

	c.Debug("Sending map files")

	c.Write(types.NewTcpPacket("M" + config.Configuration.General.Map))

	packet, err := types.ReadTcpPacket(c.Conn)

	if err != nil {
		c.Error("Error reading from connection - Additional output below")
		c.Error(err.Error())
		return
	}

	if packet.Code(0) == 'H' {
		c.Info("Client is connected and loaded")
		c.SetState(types.StatePlaying)
	} else {
		c.Warn("Client may not be loaded - Unrecognized map load response")
	}
}

func (c *TCPConnection) Close() {
	c.Info("Closing connection")
	if c.Conn != nil {
		c.Kick("Server shutting down")
		_ = c.Conn.Close()
		c.SetState(types.StateDisconnected)
	}
	delete(c.Parent.Connections, c.Address)
}

// Kick a connection with a given message
func (c *TCPConnection) Kick(msg string) {
	c.Infof("Kicking connection %s", c.Address)
	c.Infof("Reason: %s", msg)

	c.Write(types.NewTcpPacket("K" + msg)) // Kick the connection
	// c.Close()
}

// Main gamemplay loop for the connection
func (c *TCPConnection) RuntimeLoop() bool {
	_ = c.Conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	packet, err := types.ReadTcpPacket(c.Conn)

	if packet.IsEmpty() {
		if err != nil {

			e := err.Error()

			if strings.HasSuffix(e, "i/o timeout") {
				c.Debug("I/O Timeout Err - Ignoring")
				return false
			} else if strings.HasSuffix(e, "EOF") {
				// Sleep the thread for a second to allow the client to begin transmitting
				time.Sleep(time.Second)
				return false
			} else {
				c.Error("Error reading from connection - Additional output below")
				c.Error(err.Error())
				return true
			}

		}

	}

	c.GameplayParser(packet)

	return false
}

func (c *TCPConnection) GameplayParser(Packet types.TcpPacket) {
	c.Debugf("Received packet: %s", Packet.Data)
}
