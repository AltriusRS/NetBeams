package http

import (
	"strconv"

	"github.com/altriusrs/netbeams/src/types"
)

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

func (p *Player) IntoPlayerEntity() types.Player {
	userId, _ := strconv.ParseInt(p.Uid, 10, 32)

	return types.Player{
		DisplayName: p.Name,
		Address:     nil,
		PlayerId:    0,
		Vehicles:    nil,
		Account: &types.Account{
			Name:        p.Name,
			Id:          p.Id,
			Guest:       p.Guest,
			Identifiers: p.Identifiers,
			Roles:       []string{p.Roles},
			UserId:      userId,
		},
		PublicKey: p.PublicKey,
	}
}
