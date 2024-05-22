package player_manager

import (
	"fmt"

	"github.com/altriusrs/netbeams/src/config"
	"github.com/altriusrs/netbeams/src/types"
)

// A Player Manager service instance
type PlayerManager struct {
	types.Service
	Players      map[int]*types.Player
	Reservations map[int]string
}

// Create a new Player Manager service instance
func Service() *PlayerManager {
	pm := PlayerManager{

		Service:      types.SpinUp("Player Manager"),
		Players:      make(map[int]*types.Player),
		Reservations: make(map[int]string),
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

	return types.StatusHealthy, nil
}

// Add a new player to the service
func (s *PlayerManager) AddPlayer(player *types.Player) error {
	s.Infof("Adding player %s (%s)", player.Account.Name, player.Account.Id)

	id, err := s.ReserveSlot(player)

	if err != nil {
		return err
	}

	s.Players[*id] = player
	s.Reservations[*id] = player.PublicKey

	return nil
}

// Reserve a slot for a connection
func (s *PlayerManager) ReserveSlot(player *types.Player) (*int, error) {
	for _, p := range s.Reservations {
		if p == player.PublicKey {
			return nil, fmt.Errorf("player %s is already connected", player.Account.Name)
		}
	}

	id, err := s.GetNextID()

	if err != nil {
		return nil, err
	}

	if *id == -1 {
		return nil, fmt.Errorf("server is full")
	}

	s.Reservations[*id] = player.PublicKey

	return id, nil

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
