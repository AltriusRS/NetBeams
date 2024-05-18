package server

import (
	"fmt"
	"netbeams/globals"
	"netbeams/http"
	"netbeams/logs"
	"netbeams/tcp"
)

type Player struct {
	Name          string
	Uid           string
	Guest         bool
	Status        globals.Status
	State         globals.State
	Logger        logs.Logger
	API           *http.API
	TCPConnection *tcp.TCPConnection
}

func PlayerFromApi(p *http.Player, l *logs.Logger, api *http.API, conn *tcp.TCPConnection) Player {
	return Player{
		Name:          p.Name,
		Uid:           p.Uid,
		Guest:         p.Guest,
		Status:        globals.Healthy,
		State:         globals.Identify,
		Logger:        l.Fork(fmt.Sprintf("Player %s", p.Name)),
		API:           api,
		TCPConnection: conn,
	}
}

func (p *Player) Kick(Reason string) bool {
	if p.Status != globals.Healthy {
		return false
	}

	p.Logger.Infof("Kicking player %s", p.Name)
	p.Logger.Infof("Reason: %s", Reason)

	p.TCPConnection.Kick(Reason)

	return true
}

func (p *Player) Shutdown() bool {
	if p.Status != globals.Healthy {
		return false
	}

	p.Logger.Info("Shutting down player")

	p.TCPConnection.Kick("Server is shutting down")

	return true
}
