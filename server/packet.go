package server

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Packet struct {
	Header int32
	data   []byte
}

func NewPacket(data any) Packet {
	dataBytes := make([]byte, 0)

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

func (p *Packet) WriteString(s string) {
	p.data = append(p.data, []byte(s)...)
	p.Header = int32(len(p.data))
}

func (p *Packet) WriteBytes(b []byte) {
	p.data = append(p.data, b...)
	p.Header = int32(len(p.data))
}

func (p *Packet) WriteInt32(i int32) {
	unsigned := ConvertToSignedUint(i)

	data := make([]byte, 8)
	binary.LittleEndian.PutUint32(data, uint32(unsigned))
	p.data = append(p.data, data...)
	p.Header = int32(len(p.data))
}

func ReadPacket(c net.Conn) (Packet, error) {
	var packet Packet
	packet.data = make([]byte, 0)

	header := make([]byte, 4)
	_, err := c.Read(header)
	if err != nil {
		panic(err)
	}

	unsignedHeader := binary.LittleEndian.Uint32(header)

	println(fmt.Sprintf("Unsigned header: %032b", unsignedHeader))

	signedHeader := ConvertToSignedInt(unsignedHeader)

	println(fmt.Sprintf("Packet header: %x bytes", signedHeader))
	println(fmt.Sprintf("Packet header: %x bytes", signedHeader))

	packet.Header = int32(signedHeader)

	return packet, nil
}

func ConvertToSignedUint(i int32) uint32 {
	println(fmt.Sprintf("Start   : %b", i))

	unsigned := uint32(i)
	println(fmt.Sprintf("Unsigned: %b", i))

	if i < 0 {
		unsigned = 1 >> unsigned
	}

	println(fmt.Sprintf("Finished: %b", i))

	return unsigned
}

func ConvertToSignedInt(i uint32) int32 {
	println(fmt.Sprintf("Start   : %b", i))

	signed := int32(i)
	println(fmt.Sprintf("Signed  : %b", i))

	if i > 0x7FFFFFFF {
		signed = 1 << signed
	}

	println(fmt.Sprintf("Finished: %b", i))

	return signed
}
