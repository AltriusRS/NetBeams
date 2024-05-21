package tcp

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/altriusrs/netbeams/src/config"
	"github.com/altriusrs/netbeams/src/types"
)

type Server struct {
	types.Service
	Addr        string
	Port        int
	Listener    *net.TCPListener
	Connections map[string]TCPConnection
}

func Service() *Server {

	server := Server{
		Service:     types.SpinUp("TCP Server"),
		Addr:        "0.0.0.0", // Listen on all interfaces
		Port:        config.Configuration.General.Port,
		Connections: make(map[string]TCPConnection),
	}

	server.RegisterServiceHooks(server.Start, server.Stop, nil)

	return &server
}

func (s *Server) Start() (types.Status, error) {
	s.SetStatus(types.StatusStarting)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", s.Addr, s.Port))

	if err != nil {
		s.Error("Error resolving TCP address - Additional output below")
		return types.StatusErrored, err
	}

	s.Info("Starting TCP Server")

	listener, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		s.Error("Error starting TCP listener - Additional output below")
		return types.StatusErrored, err
	}

	s.Info("TCP Server started")
	s.Listener = listener

	go s.Listen()
	return types.StatusHealthy, nil
}

func (s *Server) Stop() (types.Status, error) {
	s.SetStatus(types.StatusStopping)

	delay := time.Second

	for *s.Status == types.StatusStopped {
		time.Sleep(delay)
		s.Info("Waiting for listener to stop")
	}

	return types.StatusShutdown, nil
}

func (s *Server) Listen() {
	s.SetStatus(types.StatusHealthy)

	for {
		s.Listener.SetDeadline(time.Now().Add(time.Second))
		if *s.Status != types.StatusHealthy {
			err := s.Listener.Close()
			if err != nil {
				s.Error("Error closing listener - Additional output below")
				s.Fatal(err)
			}
			s.SetStatus(types.StatusStopped)
			break
		}

		conn, err := s.Listener.Accept()

		if err != nil {
			if strings.HasSuffix(err.Error(), "i/o timeout") {
				// s.Debug("Connection polling timed out")
				continue
			}
			s.Error("Error accepting connection - Additional output below")
			s.Fatal(err)
			return
		}

		addr := conn.RemoteAddr().String()

		s.Debugf("Incoming connection from %s", conn.RemoteAddr())
		connection := NewTCPConnection(conn, addr, s)

		s.Connections[addr] = connection

		go connection.Listen()
	}

}
