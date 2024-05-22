package types

import "net"

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
	DisplayName string     // The name of the player (can be changed by the client through plugins)
	Address     net.Addr   // The IP address of the player
	PlayerId    int        // The ID of the player within the game's ID system
	Vehicles    []*Vehicle // The vehicles actively owned by the player in the current session
	Account     *Account   // The account information for the player from the BeamMP API
	PublicKey   string     // The public key of the player for their current session
}

// AuthKey = "84ce26b9-5302-42a1-b3c2-e5d09dea4bd8"
