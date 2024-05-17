package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"netbeams/globals"
	"netbeams/logs"
	"time"
)

// API is a wrapper around http.Client tailored for the BeamMP API
type API struct {
	Logger logs.Logger
	client *http.Client
}

func NewAPI(l *logs.Logger) API {
	return API{
		Logger: l.Fork("API"),
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:       5,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: true,
			},
			Timeout: time.Second * 10,
		},
	}
}

func (a *API) AuthenticatePlayer(key string) (bool, error) {
	url := fmt.Sprintf("%s/pkToUser", globals.BaseAuthAPIURL)

	body := map[string]string{
		"key": key,
	}

	a.Logger.Debug("Authenticating player")

	payload, err := json.Marshal(body)

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))

	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", globals.UserAgent)

	resp, err := a.client.Do(req)

	if err != nil {
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var wholeBody []byte

	resp.Body = http.MaxBytesReader(nil, resp.Body, resp.ContentLength)

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header)
	fmt.Println(string(wholeBody))

	var respBody map[string]any

	err = json.NewDecoder(resp.Body).Decode(&respBody)

	if err != nil {
		return false, err
	}

	j, _ := json.MarshalIndent(respBody, "", "    ")
	fmt.Println(string(j))

	return true, nil
}
