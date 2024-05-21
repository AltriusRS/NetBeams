package udp

import "net"

// A UDP packet
type Packet struct {
	Header int32
	Data   []byte
	Source net.Addr
}

func ReadPacketFromUDP(connection *net.UDPConn) (*Packet, error) {
	p := &Packet{}

	buf := make([]byte, 1024)
	_, addr, err := connection.ReadFromUDP(buf)

	if err != nil {
		return nil, err
	}

	p.Header = int32(buf[0])<<24 | int32(buf[1])<<16 | int32(buf[2])<<8 | int32(buf[3])
	p.Data = buf[4:]
	p.Source = addr

	return p, nil
}
