package udp

import (
	"net"
	"netbeams/globals"
	"netbeams/logs"
	"strconv"
)

// A UDP server instance for the udp portion of the protocol
type Server struct {
	Addr           string               // The address to listen on
	Port           int                  // The port to listen on
	Logger         logs.Logger          // The logger instance
	StatusCallback func(globals.Status) // The callback function to call when the status changes
	Status         globals.Status       // The current status of the server
	Listener       *net.UDPConn         // The UDP listener instance
}

// Create a new UDP server instance
func NewServer(addr string, port int, logger *logs.Logger, cb func(globals.Status)) Server {
	return Server{
		Addr:           addr,
		Port:           port,
		Logger:         logger.Fork("UDP Server"),
		Status:         globals.Starting,
		StatusCallback: cb,
	}
}

// Set the current status of the server
func (s *Server) SetStatus(status globals.Status) {
	if s.Status != status {
		s.Logger.Infof("Status changed from %s to %s", s.Status, status)
		s.Status = status
		s.StatusCallback(s.Status)
	}
}

// Shutdown the UDP server
func (s *Server) Shutdown() {
	s.Logger.Info("Shutting down UDP server")
	s.SetStatus(globals.Stopping)
	err := s.Listener.Close()
	if err != nil {
		s.Logger.Error("Error closing UDP listener: " + err.Error())
		s.SetStatus(globals.Errored)
	}
	s.SetStatus(globals.Shutdown)
}

// Start the UDP server
func (s *Server) Start() bool {
	s.Logger.Info("Starting UDP server")
	s.SetStatus(globals.Starting)

	udpAddr, err := net.ResolveUDPAddr("udp", s.Addr+":"+strconv.Itoa(s.Port))

	if err != nil {
		s.Logger.Error("Error resolving UDP address: " + err.Error())
		return false
	}

	listener, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		s.Logger.Error("Error starting UDP server: " + err.Error())
		return false
	}

	s.Listener = listener

	go s.Listen()

	return true
}

func (s *Server) Listen() {
	s.SetStatus(globals.Healthy)

	// While the server is healthy, listen for incoming UDP packets
	for s.Status == globals.Healthy {
		buf := make([]byte, 1024)
		n, addr, err := s.Listener.ReadFromUDP(buf)

		if err != nil {
			// If the error is a closed network connection, break out of the loop
			if err.Error() == "use of closed network connection" {
				s.Logger.Info("UDP listener closed")
				break
			}
			// Otherwise, continue attempting to listen to packets
			s.Logger.Error("Error reading UDP packet: " + err.Error())
			continue
		}

		s.Logger.Info("Received UDP packet from " + addr.String())
		s.Logger.Infof("Packet: %s", buf[:n])

	}
}
