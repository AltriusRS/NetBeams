package http

// Player is a wrapper around a player object that is returned from the API
type Player struct {
	V           int64    `json:"__v"`
	Id          string   `json:"_id"`
	Name        string   `json:"username"`
	Guest       bool     `json:"guest"`
	CreatedAt   string   `json:"createdAt"`
	Identifiers []string `json:"identifiers"`
	PublicKey   string   `json:"public_key"`
	Roles       string   `json:"roles"`
	Uid         string   `json:"uid"`
}
