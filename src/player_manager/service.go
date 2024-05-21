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
	s.Infof("Adding player %s", player.Name)

	id, err := s.ReserveSlot(player)

	if err != nil {
		return err
	}

	s.Players[id] = player
	s.Reservations[id] = player.PublicKey

	return nil
}

// Reserve a slot for a connection
func (s *PlayerManager) ReserveSlot(player *types.Player) (int, error) {
	if _, ok := s.Reservations[player.PublicKey]; ok {
		return nil, fmt.Errorf("player %s is already connected", player.Name)
	}

	id := s.GetNextID()
	if id == -1 {
		return nil, fmt.Errorf("server is full")
	}

	s.Reservations[id] = player.PublicKey

	return id, nil

}

// Get the next available player ID
func (s *PlayerManager) GetNextID() (int, error) {
	if len(s.ids) == 0 {
		return -1, fmt.Errorf("server is full")
	}

	for id := range config.Configuration.General.MaxPlayers {
		s.ids = s.ids[1:]
		if s.HasSlot(id) {
			continue
		}
		return id
	}
}
