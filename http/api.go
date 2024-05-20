package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"netbeams/environment"
	"netbeams/globals"
	"netbeams/logs"
	"time"
)

// API is a wrapper around http.Client tailored for the BeamMP API
type API struct {
	globals.Service
	Logger logs.Logger
	client *http.Client
}

func Service() *API {

	api := API{
		Service: globals.SpinUp("BeamMP API"),
		Logger:  logs.NetLogger("BeamMP API"),
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:       5,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: true,
			},
			Timeout: time.Second * 10,
		},
	}

	// Register the service hooks
	api.RegisterServiceHooks(api.Start, api.Stop, nil)

	return &api
}

func (a *API) Start() (globals.Status, error) {
	return globals.StatusHealthy, nil
}

func (a *API) Stop() (globals.Status, error) {
	return globals.StatusStopped, nil
}

// AuthenticatePlayer authenticates a player with the BeamMP API and returns a Player object
// if successful, or nil and an error if not
func (a *API) AuthenticatePlayer(key string) (*Player, error) {
	url := fmt.Sprintf("%s/pkToUser", globals.BaseAuthAPIURL)

	body := map[string]string{
		"key": key,
	}

	a.Debug("Authenticating player")

	payload, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("NetBeams/%s", environment.Context.Version))

	resp, err := a.client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var player Player

	err = json.NewDecoder(resp.Body).Decode(&player)

	if err != nil {
		return nil, err
	}

	return &player, nil
}
