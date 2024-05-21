package udp

import (
	"net"
	"strconv"
	"strings"

	"github.com/altriusrs/netbeams/src/config"
	"github.com/altriusrs/netbeams/src/types"
)

// A UDP server instance for the udp portion of the protocol
type Server struct {
	types.Service
	Addr     string       // The address to listen on
	Port     int          // The port to listen on
	Listener *net.UDPConn // The UDP listener instance
}

// Create a new UDP server instance
func Service() *Server {
	service := types.SpinUp("UDP Server")

	server := Server{
		Service:  service,
		Addr:     "0.0.0.0",
		Port:     config.Configuration.General.Port,
		Listener: nil,
	}

	service.RegisterServiceHooks(server.Start, server.Shutdown, nil)

	return &server
}

// Shutdown the UDP server
func (s *Server) Shutdown() (types.Status, error) {
	s.Info("Shutting down UDP server")
	err := s.Listener.Close()
	if err != nil {
		s.Error("Error closing UDP listener: " + err.Error())
		return types.StatusErrored, err
	}

	return types.StatusShutdown, nil
}

// Start the UDP server
func (s *Server) Start() (types.Status, error) {
	s.Info("Starting UDP server")
	s.SetStatus(types.StatusStarting)

	udpAddr, err := net.ResolveUDPAddr("udp", s.Addr+":"+strconv.Itoa(s.Port))

	if err != nil {
		s.Error("Error resolving UDP address: " + err.Error())
		return types.StatusErrored, err
	}

	listener, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		s.Error("Error starting UDP server: " + err.Error())
		return types.StatusErrored, err
	}

	s.Listener = listener

	go s.Listen()

	return types.StatusHealthy, nil
}

func (s *Server) Listen() {
	s.SetStatus(types.StatusHealthy)

	// While the server is healthy, listen for incoming UDP packets
	for *s.Status == types.StatusHealthy {
		packet, err := ReadPacketFromUDP(s.Listener)

		if err != nil {
			// If the error is a closed network connection, break out of the loop
			if strings.HasSuffix(err.Error(), "use of closed network connection") {
				s.Info("UDP listener closed")
				break
			}
			// Otherwise, continue attempting to listen to packets
			s.Error("Error reading UDP packet: " + err.Error())
			continue
		}

		s.Info("Received UDP packet from " + packet.Source.String())
		s.Infof("Packet: %s", packet.Data)

		// types.App.GetService("Player Manager").HandlePacket(packet)
	}
}
