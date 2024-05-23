package types

import "github.com/altriusrs/netbeams/src/logs"

type ServiceHook func() (Status, error)

type ServiceCompatible interface {
	GetName() string
	GetStatus() *Status
	SetStatus(status Status)
	StartService() error
	StopService() error
	RestartService()
}

// Service is a struct which represents a single service
type Service struct {
	logs.Logger
	Name        string
	Status      *Status
	StartHook   ServiceHook
	StopHook    ServiceHook
	CleanupHook ServiceHook
}

func SpinUp(name string) Service {
	return Service{
		Logger: logs.NetLogger(name),
		Name:   name,
		Status: nil,
	}
}

// Registers the service hooks used to start, stop, and cleanup the service
// These are called in order to maintain a healthy state
func (s *Service) RegisterServiceHooks(startHook ServiceHook, stopHook ServiceHook, cleanupHook ServiceHook) {
	s.StartHook = startHook
	s.StopHook = stopHook
	s.CleanupHook = cleanupHook
}

func (s *Service) StartService() error {

	if s.StartHook == nil {
		s.SetStatus(StatusHealthy)
		return nil
	}

	s.SetStatus(StatusStarting)
	state, err := s.StartHook()
	s.SetStatus(state)

	return err
}

func (s *Service) StopService() error {
	s.Debug("Stopping...")

	if s.StopHook == nil {
		s.SetStatus(StatusStopped)
		return nil
	}

	s.SetStatus(StatusStopping)
	state, err := s.StopHook()

	s.SetStatus(state)

	s.Debugf("Response from stop hook: %v - %v", state, err)

	return err
}

func (s *Service) RestartService() {
	if s.StartHook == nil {
		return
	}
	if s.StopHook == nil {
		return
	}

	if s.CleanupHook == nil {
		return
	}

	s.SetStatus(StatusRestarting)
	state, err := s.StopHook()

	if err != nil {
		s.SetStatus(state)
		return
	}

	state, err = s.CleanupHook()

	if err != nil {
		s.SetStatus(state)
		return
	}

	state, err = s.StartHook()

	if err != nil {
		s.SetStatus(state)
		return
	}
}

func (s *Service) SetStatus(status Status) {
	if s.Status == nil || *s.Status != status {
		s.Infof("Service %s status changed from %s to %s", s.Name, s.Status, status)
		s.Status = &status
	}
}

func (s *Service) GetStatus() *Status {
	return s.Status
}

func (s *Service) GetName() string {
	return s.Name
}
