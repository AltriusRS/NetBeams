package heartbeat

import (
	"net/http"

	"github.com/altriusrs/netbeams/src/types"
)

type Manager struct {
	types.Service
	client *http.Client
}

func Service() *Manager {
	client := http.DefaultClient

	manager := &Manager{
		Service: types.SpinUp("Heartbeat"),
		client:  client,
	}

	manager.RegisterServiceHooks(manager.StartHook, manager.StopHook, manager.CleanupHook)

	return manager
}

func (h *Manager) StartHook() (types.Status, error) {

	// go h.routine()

	return types.StatusHealthy, nil
}

func (h *Manager) StopHook() (types.Status, error) {
	return types.StatusStopped, nil
}

func (h *Manager) CleanupHook() (types.Status, error) {
	return types.StatusStopped, nil
}

// func (h *Manager) routine() {
// 	for {
// 		data := url.Values{}

// 		data.Set("uuid", config.Configuration.General.AuthKey)
// 		// data.Set("players", string(rune(config.Configuration.General.MaxPlayers)))

// 		h.client.Post(types.BaseAPIURL+"/heartbeat", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
// 		time.Sleep(time.Second * 30)
// 	}
// }
