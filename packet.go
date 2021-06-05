package SimplyP2P

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"time"
)

const (
	newPeerPacket     = 0x01
	changeStatePacket = 0x02
	ipv4Packet        = 0x04
	ipv6Packet        = 0x06
)

// AddPeer is a packet that adds a new peer to a node.
// Use an array instead of a slice to use it has map hashed key.
// +--------+--------------+---------+------+
// |  0x01  | Address Type | Address | Port |
// +--------+--------------+---------+------+
type AddPeer struct {
	Peer
}

// WriteTo encodes a Peer packet.
func (p AddPeer) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer

	// Packet ID
	buf.WriteByte(newPeerPacket)

	// Address type
	if p.addressLength == net.IPv4len {
		buf.WriteByte(ipv4Packet)
		buf.Write(p.address[:net.IPv4len])
	} else {
		buf.WriteByte(ipv6Packet)
		buf.Write(p.address[:net.IPv6len])
	}

	// Port
	port := make([]byte, 2)
	binary.LittleEndian.PutUint16(port, p.port)
	buf.Write(port)

	return buf.WriteTo(w)
}

// ReadFrom decodes a Peer packet.
func (p *AddPeer) ReadFrom(r io.Reader) (n int64, err error) {
	// Read address type
	addressType := make([]byte, 1)
	if _, err = r.Read(addressType); err != nil {
		return 0, err
	}

	// Read address
	if addressType[0] == ipv4Packet {
		_, err = r.Read(p.address[:net.IPv4len])
	} else if addressType[0] == ipv6Packet {
		_, err = r.Read(p.address[:net.IPv6len])
	} else {
		err = errors.New("invalid address type")
	}

	if err != nil {
		return 0, err
	}

	// Read port
	port := make([]byte, 2) // TCP port = 2 bytes
	if _, err = r.Read(port); err != nil {
		return 0, err
	}
	p.port = binary.LittleEndian.Uint16(port)

	return int64(len(addressType) + len(p.address) + len(port)), nil
}

// ChangeState is a packet that change global P2P state.
// +--------+-----------+-------+
// |  0x02  | Timestamp | State |
// +--------+-----------+-------+
type ChangeState struct {
	State bool
	time  time.Time
}

// WriteTo encodes a ChangeState packet.
func (s ChangeState) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer

	// Packet ID
	buf.WriteByte(changeStatePacket)

	// Send time
	timestamp := make([]byte, 8)
	binary.LittleEndian.PutUint64(timestamp, uint64(s.time.UnixNano()))
	buf.Write(timestamp)

	// Current state
	if s.State {
		buf.WriteByte(0x01)
	} else {
		buf.WriteByte(0x00)
	}

	return buf.WriteTo(w)
}

// ReadFrom decodes a ChangeState packet.
func (s *ChangeState) ReadFrom(r io.Reader) (n int64, err error) {
	// Read timestamp
	timestamp := make([]byte, 8)
	if _, err = r.Read(timestamp); err != nil {
		return 0, err
	}

	// Parse Time from unix timestamp
	s.time = time.Unix(0, int64(binary.LittleEndian.Uint64(timestamp)))

	// Read state
	state := make([]byte, 1)
	if _, err = r.Read(state); err != nil {
		return 0, err
	}
	s.State = state[0] == 0x01

	return int64(len(timestamp) + len(state)), nil
}

// GetTime returns current read packet time of dispatch.
func (s ChangeState) GetTime() time.Time {
	return s.time
}
