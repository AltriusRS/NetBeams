package types

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
)

// TcpPacket is a wrapper around a byte array that contains a header and Data
type TcpPacket struct {
	Header int32
	Data   []byte
}

func NewEmptyTcpPacket() TcpPacket {
	return TcpPacket{
		Header: 0,
		Data:   []byte{},
	}
}

// NewTcpPacket creates a new TcpPacket from a given Data type
func NewTcpPacket(Data any) TcpPacket {
	var dataBytes []byte

	switch data := Data.(type) {
	case string:
		dataBytes = []byte(data)
	case []byte:
		dataBytes = data
	case int16:
		dataBytes = make([]byte, 2)
		binary.LittleEndian.PutUint16(dataBytes, uint16(data))
	case int32:
		dataBytes = make([]byte, 4)
		binary.LittleEndian.PutUint32(dataBytes, uint32(data))
	case int64:
		dataBytes = make([]byte, 8)
		binary.LittleEndian.PutUint64(dataBytes, uint64(data))
	default:
		panic("Unsupported Data type")
	}

	return TcpPacket{
		Header: int32(len(dataBytes)),
		Data:   dataBytes,
	}
}

func (p *TcpPacket) ReadHeader(c net.Conn) (int32, error) {
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

func (p *TcpPacket) ReadData(c net.Conn) ([]byte, error) {
	dataLength := p.Header
	Data := make([]byte, dataLength)
	bytesRead, err := c.Read(Data)
	if err != nil {
		return nil, err
	}

	if bytesRead != int(dataLength) {
		return nil, fmt.Errorf("invalid Data length")
	}

	p.Data = Data

	return Data, nil
}

func ReadTcpPacket(c net.Conn) (TcpPacket, error) {
	var TcpPacket TcpPacket
	_, err := TcpPacket.ReadHeader(c)

	if err != nil {
		return TcpPacket, err
	}

	_, err = TcpPacket.ReadData(c)

	if err != nil {
		return TcpPacket, err
	}

	return TcpPacket, nil
}

func (TcpPacket) FromString(s string) *TcpPacket {
	return &TcpPacket{
		Header: int32(len([]byte(s))),
		Data:   []byte(s),
	}
}

func (p *TcpPacket) ToString() string {
	return string(p.Data)
}

func (p *TcpPacket) WriteString(s string) {
	p.Data = append(p.Data, []byte(s)...)
	p.Header = int32(len(p.Data))
}

func (p *TcpPacket) WriteBytes(b []byte) {
	p.Data = append(p.Data, b...)
	p.Header = int32(len(p.Data))
}

func (p *TcpPacket) WriteInt8(i int8) {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(i))
	p.Data = append(p.Data, bytes[3])
	p.Header = int32(len(p.Data))
}

func (p *TcpPacket) WriteInt16(i int16) {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(i))
	p.Data = append(p.Data, bytes[1:3]...)
	p.Header = int32(len(p.Data))
}

func (p *TcpPacket) WriteInt32(i int32) {
	bin := make([]byte, 4)
	binary.LittleEndian.PutUint32(p.Data, uint32(i))
	p.Data = append(p.Data, bin...)
	p.Header = int32(len(p.Data))
}

func (p *TcpPacket) WriteInt64(i int64) {
	bin := make([]byte, 8)
	binary.LittleEndian.PutUint64(p.Data, uint64(i))
	p.Data = append(p.Data, bin...)
	p.Header = int32(len(p.Data))
}

func (p *TcpPacket) WriteFloat32(f float32) {
	bin := make([]byte, 4)
	binary.LittleEndian.PutUint32(p.Data, math.Float32bits(f))
	p.Data = append(p.Data, bin...)
	p.Header = int32(len(p.Data))
}

func (p *TcpPacket) WriteFloat64(f float64) {
	bin := make([]byte, 8)
	binary.LittleEndian.PutUint64(p.Data, math.Float64bits(f))
	p.Data = append(p.Data, bin...)
	p.Header = int32(len(p.Data))
}

func (p *TcpPacket) Serialize() []byte {
	payload := make([]byte, p.Header+4)
	binary.LittleEndian.PutUint32(payload, uint32(p.Header))
	payload = payload[:4]
	payload = append(payload, p.Data...)

	return payload
}

func (p *TcpPacket) Code(position int) rune {
	return rune(p.Data[position])
}

func (p *TcpPacket) Payload() []byte {
	return p.Data[1:]
}

func (p *TcpPacket) IsEmpty() bool {
	return p.Header == 0 && len(p.Data) == 0
}

func (p *TcpPacket) String() string {
	return string(p.Data)
}
