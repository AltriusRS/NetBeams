package types

import (
	"net"
)

// A BeamMP Account object that is returned from the API in a format which can be modified by the client here
type Account struct {
	Name        string   // The player's name on the forum
	Id          string   // The player's ID internally to BeamMP
	Guest       bool     // Whether the player is a guest account or not
	Identifiers []string // The player's identifiers
	Roles       []string // The roles of the player
	UserId      int64    // The user ID of the player on the forum
}

// A structure representing a player in the game itself
type Player struct {
	Status      PlayerStatus            // The status of the player
	DisplayName string                  // The name of the player (can be changed by the client through plugins)
	Address     net.Addr                // The IP address of the player
	PlayerId    int                     // The ID of the player within the game's ID system
	Vehicles    []*Vehicle              // The vehicles actively owned by the player in the current session
	Account     *Account                // The account information for the player from the BeamMP API
	PublicKey   string                  // The public key of the player for their current session
	Permissions PlayerPermissionsConfig // The permissions of the player, this is a struct based off of the config file, but skips any fields which are handled before the player loads to this state
}

type PlayerStatus int

const (
	PlayerStatusUnknown    PlayerStatus = iota // A player status that is unknown
	PlayerStatusConnecting                     // A player that is connecting to the server
	PlayerStatusLoading                        // A player that is loading into the server (maps or mods)
	PlayerStatusConnected                      // A player that is currently connected to the server
	PlayerStatusIdle                           // A player that is idle on the server
	PlayerStatusPlaying                        // A player that is currently playing the game
)

// A struct representing the permissions for each player
type PlayerPermissionsConfig struct {
	// If the player may bypass idle timeouts
	BypassIdle bool `toml:"BypassIdle" comment:"Whether idle connections are allowed to join the server"`

	// If the player may bypass online quota limits
	BypassOnline bool `toml:"BypassOnline" comment:"Whether online connections are allowed to join the server"`

	// If the player may bypass vehicle limiting
	BypassVehicles bool `toml:"BypassVehicles" comment:"Whether vehicles are allowed to join the server"`

	// If the player may hide their name in the leaderboard
	HideName bool `toml:"HideName" comment:"Whether the player's name is hidden in the leaderboard"`

	// If the player may kick other players
	KickPlayers bool `toml:"KickPlayers" comment:"Whether the user can kick other players"`

	// If the player may ban other players
	BanPlayers bool `toml:"BanPlayers" comment:"Whether the user can ban other players"`

	// If the player may mute other players
	MutePlayers bool `toml:"MutePlayers" comment:"Whether the user can mute other players"`
}
