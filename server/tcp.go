package server

import (
	"fmt"
	"net"
	"netbeams/logs"
)

type TCPServer struct {
	Addr           string
	Port           int
	Logger         logs.Logger
	Listener       *net.TCPListener
	StatusCallback func(Status)
}

func NewTCPServer(port int, l *logs.Logger, cb func(Status)) TCPServer {
	return TCPServer{
		Addr:           "0.0.0.0",
		Port:           port,
		Logger:         l.Fork("TCP Server"),
		StatusCallback: cb,
	}
}

func (s *TCPServer) Start() bool {
	s.StatusCallback(Starting)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", s.Addr, s.Port))

	if err != nil {
		s.Logger.Error("Error resolving TCP address - Additional output below")
		s.Logger.Fatal(err)
		s.StatusCallback(Errored)
		return false
	}

	s.Logger.Info("Starting TCP Server")

	listener, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		s.Logger.Error("Error starting TCP listener - Additional output below")
		s.Logger.Fatal(err)
		s.StatusCallback(Errored)
		return false
	}

	s.Logger.Info("TCP Server started")
	s.Listener = listener
	go s.Listen()
	return true
}

func (s *TCPServer) Stop() {
	s.StatusCallback(Stopping)

	if s.Listener != nil {
		s.Listener.Close()
	}

	s.StatusCallback(Shutdown)
}

func (s *TCPServer) Listen() {
	s.StatusCallback(Healthy)
	for {
		conn, err := s.Listener.Accept()

		if err != nil {
			s.Logger.Error("Error accepting connection - Additional output below")
			s.Logger.Fatal(err)
			return
		}

		s.Logger.Debugf("Incoming connection from %s", conn.RemoteAddr())

		conn.Close()
	}
}

func (s *TCPServer) Accept() {

}

func (s *TCPServer) Handle() {

}
