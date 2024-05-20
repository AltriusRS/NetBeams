package tcp

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
)

// Packet is a wrapper around a byte array that contains a header and data
type Packet struct {
	Header int32
	data   []byte
}

func NewEmptyPacket() Packet {
	return Packet{
		Header: 0,
		data:   []byte{},
	}
}

// NewPacket creates a new packet from a given data type
func NewPacket(data any) Packet {
	var dataBytes []byte

	switch data.(type) {
	case string:
		dataBytes = []byte(data.(string))
	case []byte:
		dataBytes = data.([]byte)
	case int16:
		dataBytes = make([]byte, 2)
		binary.LittleEndian.PutUint16(dataBytes, uint16(data.(int16)))
	case int32:
		dataBytes = make([]byte, 4)
		binary.LittleEndian.PutUint32(dataBytes, uint32(data.(int32)))
	case int64:
		dataBytes = make([]byte, 8)
		binary.LittleEndian.PutUint64(dataBytes, uint64(data.(int64)))
	default:
		panic("Unsupported data type")
	}

	return Packet{
		Header: int32(len(dataBytes)),
		data:   dataBytes,
	}
}

func (p *Packet) ReadHeader(c net.Conn) (int32, error) {
	header := make([]byte, 4)
	bytesRead, err := c.Read(header)
	if err != nil {
		return 0, err
	}

	if bytesRead != 4 {
		return 0, fmt.Errorf("invalid header length")
	}

	FrameLength := int32(header[3])<<24 | int32(header[2])<<16 | int32(header[1])<<8 | int32(header[0])

	p.Header = FrameLength

	return FrameLength, nil
}

func (p *Packet) ReadData(c net.Conn) ([]byte, error) {
	dataLength := p.Header
	data := make([]byte, dataLength)
	bytesRead, err := c.Read(data)
	if err != nil {
		return nil, err
	}

	if bytesRead != int(dataLength) {
		return nil, fmt.Errorf("invalid data length")
	}

	p.data = data

	return data, nil
}

func ReadPacket(c net.Conn) (Packet, error) {
	var packet Packet
	_, err := packet.ReadHeader(c)

	if err != nil {
		return packet, err
	}

	_, err = packet.ReadData(c)

	if err != nil {
		return packet, err
	}

	return packet, nil
}

func (Packet) FromString(s string) *Packet {
	return &Packet{
		Header: int32(len([]byte(s))),
		data:   []byte(s),
	}
}

func (p *Packet) ToString() string {
	return string(p.data)
}

func (p *Packet) WriteString(s string) {
	p.data = append(p.data, []byte(s)...)
	p.Header = int32(len(p.data))
}

func (p *Packet) WriteBytes(b []byte) {
	p.data = append(p.data, b...)
	p.Header = int32(len(p.data))
}

func (p *Packet) WriteInt32(i int32) {
	bin := make([]byte, 4)
	binary.LittleEndian.PutUint32(p.data, uint32(i))
	p.data = append(p.data, bin...)
	p.Header = int32(len(p.data))
}

func (p *Packet) WriteInt64(i int64) {
	bin := make([]byte, 8)
	binary.LittleEndian.PutUint64(p.data, uint64(i))
	p.data = append(p.data, bin...)
	p.Header = int32(len(p.data))
}

func (p *Packet) WriteFloat32(f float32) {
	bin := make([]byte, 4)
	binary.LittleEndian.PutUint32(p.data, math.Float32bits(f))
	p.data = append(p.data, bin...)
	p.Header = int32(len(p.data))
}

func (p *Packet) WriteFloat64(f float64) {
	bin := make([]byte, 8)
	binary.LittleEndian.PutUint64(p.data, math.Float64bits(f))
	p.data = append(p.data, bin...)
	p.Header = int32(len(p.data))
}

func (p *Packet) Serialize() []byte {
	payload := make([]byte, p.Header+4)
	binary.LittleEndian.PutUint32(payload, uint32(p.Header))
	payload = append(payload, p.data...)
	return payload
}

func (p *Packet) Code() rune {
	return rune(p.data[0])
}

func (p *Packet) Data() []byte {
	return p.data[1:]
}

func (p *Packet) IsEmpty() bool {
	return p.Header == 0 && len(p.data) == 0
}
