package tcp

import (
	"fmt"
	"net"
	"netbeams/globals"
	"netbeams/logs"
	"strings"
	"time"
)

type Server struct {
	Addr           string
	Port           int
	Logger         logs.Logger
	Listener       *net.TCPListener
	StatusCallback func(globals.Status)
	Status         globals.Status
	Connections    map[string]TCPConnection
}

func NewServer(port int, l *logs.Logger, cb func(globals.Status)) Server {
	return Server{
		Addr:           "0.0.0.0", // Listen on all interfaces
		Port:           port,
		Logger:         l.Fork("TCP Server"),
		StatusCallback: cb,
		Connections:    make(map[string]TCPConnection),
	}
}

func (s *Server) SetStatus(status globals.Status) {
	if s.Status != status {
		s.Logger.Infof("Status changed from %s to %s", s.Status, status)
		s.Status = status
		s.StatusCallback(status)
	}
}

func (s *Server) Start() bool {
	s.SetStatus(globals.Starting)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", s.Addr, s.Port))

	if err != nil {
		s.Logger.Error("Error resolving TCP address - Additional output below")
		s.Logger.Fatal(err)
		s.SetStatus(globals.Errored)
		return false
	}

	s.Logger.Info("Starting TCP Server")

	listener, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		s.Logger.Error("Error starting TCP listener - Additional output below")
		s.Logger.Fatal(err)
		s.SetStatus(globals.Errored)
		return false
	}

	s.Logger.Info("TCP Server started")
	s.Listener = listener

	go s.Listen()
	return true
}

func (s *Server) Stop() {
	s.SetStatus(globals.Stopping)

	delay := time.Second

	for s.Status == globals.Stopping {
		time.Sleep(delay)
		s.Logger.Info("Waiting for listener to stop")
	}

	s.SetStatus(globals.Shutdown)
}

func (s *Server) Listen() {
	s.SetStatus(globals.Healthy)

	for {
		s.Listener.SetDeadline(time.Now().Add(time.Second))
		if s.Status != globals.Healthy {
			err := s.Listener.Close()
			if err != nil {
				s.Logger.Error("Error closing listener - Additional output below")
				s.Logger.Fatal(err)
			}
			s.SetStatus(globals.Stopped)
			break
		}

		conn, err := s.Listener.Accept()

		if err != nil {
			if strings.HasSuffix(err.Error(), "i/o timeout") {
				// s.Logger.Debug("Connection polling timed out")
				continue
			}
			s.Logger.Error("Error accepting connection - Additional output below")
			s.Logger.Fatal(err)
			return
		}

		addr := conn.RemoteAddr().String()

		s.Logger.Debugf("Incoming connection from %s", conn.RemoteAddr())
		connection := NewTCPConnection(conn, addr, s, &s.Logger)

		s.Connections[addr] = connection

		go connection.Listen()
	}

}
