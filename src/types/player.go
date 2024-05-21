package types

import "net"

// A BeamMP Account object that is returned from the API in a format which can be modified by the client here
type Account struct {
	Name        string   //
	PublicKey   string   //
	Id          int64    //
	Guest       bool     //
	Identifiers []string //
	Roles       []string //
	UserId      string   //
}

// A structure representing a player in the game itself
type Player struct {
	DisplayName string     // The name of the player (can be changed by the client through plugins)
	Address     net.Addr   // The IP address of the player
	PlayerId    int        // The ID of the player within the game's ID system
	Vehicles    []*Vehicle // The vehicles actively owned by the player in the current session
	Account     *Account   // The account information for the player from the BeamMP API
}
