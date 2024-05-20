package heartbeat

import (
	"net/http"
	"net/url"
	"netbeams/config"
	"netbeams/globals"
	"strings"
	"time"
)

type Manager struct {
	globals.Service
	client *http.Client
}

func New(client *http.Client) *Manager {
	return &Manager{client: client}
}

func (h *Manager) StartHook() (*globals.Status, error) {

	go h.routine()

	return &globals.Healthy, nil
}

func (h *Manager) routine() {
	for {
		data := url.Values{}

		data.Set("uuid", config.Configuration.General.AuthKey)
		data.Set("players", config.Configuration.General.Players)

		h.client.Post(globals.BaseAPIURL+"/heartbeat", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
		time.Sleep(time.Second * 30)
	}
}
