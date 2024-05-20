package tcp

import (
	"fmt"
	"net"
	"netbeams/config"
	"netbeams/globals"
	"strings"
	"time"
)

type Server struct {
	globals.Service
	Addr        string
	Port        int
	Listener    *net.TCPListener
	Connections map[string]TCPConnection
}

func Service() *Server {

	server := Server{
		Service:     globals.SpinUp("TCP Server"),
		Addr:        "0.0.0.0", // Listen on all interfaces
		Port:        config.Configuration.General.Port,
		Connections: make(map[string]TCPConnection),
	}

	server.RegisterServiceHooks(server.Start, server.Stop, nil)

	// server.RegisterServiceHooks()

	return &server
}

func (s *Server) Start() (globals.Status, error) {
	s.SetStatus(globals.Starting)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", s.Addr, s.Port))

	if err != nil {
		s.Error("Error resolving TCP address - Additional output below")
		return globals.Errored, err
	}

	s.Info("Starting TCP Server")

	listener, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		s.Error("Error starting TCP listener - Additional output below")
		return globals.Errored, err
	}

	s.Info("TCP Server started")
	s.Listener = listener

	go s.Listen()
	return globals.Healthy, nil
}

func (s *Server) Stop() (globals.Status, error) {
	s.SetStatus(globals.Stopping)

	delay := time.Second

	for *s.Status == globals.Stopped {
		time.Sleep(delay)
		s.Info("Waiting for listener to stop")
	}

	return globals.Shutdown, nil
}

func (s *Server) Listen() {
	s.SetStatus(globals.Healthy)

	for {
		s.Listener.SetDeadline(time.Now().Add(time.Second))
		if *s.Status != globals.Healthy {
			err := s.Listener.Close()
			if err != nil {
				s.Error("Error closing listener - Additional output below")
				s.Fatal(err)
			}
			s.SetStatus(globals.Stopped)
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
		connection := NewTCPConnection(conn, addr, s, &s.Logger)

		s.Connections[addr] = connection

		go connection.Listen()
	}

}
