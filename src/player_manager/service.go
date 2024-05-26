package player_manager

import (
	"fmt"

	"time"

	"github.com/altriusrs/netbeams/src/config"
	"github.com/altriusrs/netbeams/src/types"
)

// A Player Manager service instance
type PlayerManager struct {
	types.Service
	Players      map[int]*types.Player
	Reservations map[int]time.Time
}

// Create a new Player Manager service instance
func Service() *PlayerManager {
	pm := PlayerManager{
		Service:      types.SpinUp("Player Manager"),
		Players:      make(map[int]*types.Player),
		Reservations: make(map[int]time.Time),
	}

	pm.RegisterServiceHooks(pm.StartHook, pm.ShutdownHook, nil)

	return &pm
}

// Shutdown the Player Manager service
func (s *PlayerManager) ShutdownHook() (types.Status, error) {
	s.Info("Shutting down Player Manager service")

	// TODO: Close connections and remove players from the service

	return types.StatusShutdown, nil
}

// Start the Player Manager service
func (s *PlayerManager) StartHook() (types.Status, error) {
	s.Info("Starting Player Manager service")

	// start the reservation manager
	go s.Reservator()

	return types.StatusHealthy, nil
}

func (s *PlayerManager) Reservator() {
	for *s.Status == types.StatusHealthy || *s.Status == types.StatusStarting {
		time.Sleep(time.Second * 60)

		s.Info("Managing reservations")

		for id, t := range s.Reservations {
			if t.IsZero() {
				continue // Skip if the time is equivalent to EPOCH (zero time)
			}

			// Check if the time is in the past
			if t.Before(time.Now()) {
				delete(s.Reservations, id) // Remove the reservation to free the slot
				s.Players[id] = nil        // Set the player element to nil so that it can be avoided by other methods
				continue
			}

			// Check if the time is in the future
			if t.After(time.Now()) {

			}

			if time.Until(t) < time.Minute*10 {
				// types.App.GetService("TCP Server").(*tcp.Server)
			} else if time.Until(t) < time.Minute*5 {
				// types.App.GetService("TCP Server").(*tcp.Server)
			} else if time.Until(t) < time.Minute*1 {
				// types.App.GetService("TCP Server").(*tcp.Server)
			} else if time.Until(t) < time.Second*1 {
				// types.App.GetService("TCP Server").(*tcp.Server).KickPlayer(id)
				delete(s.Reservations, id) // Remove the reservation to free the slot
				s.Players[id] = nil        // Set the player element to nil so that it can be avoided by other methods
			} else {
				s.Debugf("Player %s (%d) has %f minutes left to play", s.Players[id].Account.Name, id, time.Until(t).Minutes())
			}
		}
	}
}

// Add a new player to the service
func (s *PlayerManager) AddPlayer(player *types.Player, id int) (*int, error) {
	s.Infof("Adding player %s (%s)", player.Account.Name, player.Account.Id)

	pid, err := s.ReserveSlotForPlay(id)

	if err != nil {
		return nil, err
	}

	s.Players[*pid] = player
	// Reserve the slot for 5 minutes
	// This is to allow the player to connect and load mods, without concern of disconnecting
	s.Reservations[*pid] = time.Now().Add(time.Minute * 5)

	return pid, nil
}

// Get the next available player ID
func (s *PlayerManager) GetNextID() (*int, error) {

	// Get the total number of players
	total := len(s.Reservations)

	// Check if the server is full
	if total >= config.Configuration.General.MaxPlayers {
		return nil, fmt.Errorf("server is full")
	}

	// Get the next available ID
	for id := range config.Configuration.General.MaxPlayers {

		if _, ok := s.Reservations[id]; !ok {
			return &id, nil
		}
	}

	// If we get here, the server is full
	return nil, fmt.Errorf("server is full")
}

// Reserve a slot for an incoming connection
// Slot reservations last for 60 seconds and are automatically released if the connection has not
// reached the "playing" state within that time.
// It may be reserved once again if the server has mods to synchronize, this reservation is valid for 5 minutes
func (s *PlayerManager) ReserveSlotForConnection(id *int) (*int, error) {
	if id != nil {
		if _, ok := s.Reservations[*id]; !ok {
			s.Reservations[*id] = time.Now().Add(time.Second * 60)
			return id, nil
		}
	}

	id, err := s.GetNextID()

	if err != nil {
		return nil, err
	}

	if *id == -1 {
		return nil, fmt.Errorf("server is full")
	}

	s.Reservations[*id] = time.Now().Add(time.Second * 60)

	return id, nil
}

// Reserve a slot for a loading period
func (s *PlayerManager) ReserveSlotForLoad(id int) (*int, error) {
	if _, ok := s.Reservations[id]; ok {

		if config.Configuration.Auth.Idle.Enable {
			s.Reservations[id] = time.Now().Add(time.Minute * config.Configuration.Auth.Idle.MaxTimeTime)
		} else {
			s.Reservations[id] = time.Time{}
		}

		s.Reservations[id] = time.Now().Add(time.Minute * config.Configuration.Auth.Idle.MaxTimeTime)
		return &id, nil
	}

	return nil, fmt.Errorf("cannot reserve slot for non-reserved player")
}

// Reserve a slot for play duration
func (s *PlayerManager) ReserveSlotForPlay(id int) (*int, error) {
	if _, ok := s.Reservations[id]; ok {
		if config.Configuration.Auth.Idle.Enable {
			s.Reservations[id] = time.Now().Add(time.Minute * config.Configuration.Auth.Online.QuotaTime)
		} else {
			s.Reservations[id] = time.Time{}
		}
		return &id, nil
	}

	return nil, fmt.Errorf("cannot reserve slot for non-reserved player")
}

// Get a player by their ID
func (s *PlayerManager) GetPlayer(id int) *types.Player {
	return s.Players[id]
}

// Get a player by their public key
func (s *PlayerManager) GetPlayerByPublicKey(key string) *types.Player {
	for _, p := range s.Players {
		if p.PublicKey == key {
			return p
		}
	}
	return nil
}
