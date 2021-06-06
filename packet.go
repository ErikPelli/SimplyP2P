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
	stateFalse        = 0
	stateTrue         = 1
	uint16Bytes       = 2
)

// AddPeer is a packet that adds a new peer to a node.
// Use an array instead of a slice to use it has map hashed key.
type AddPeer struct {
	Peer
}

// WriteTo encodes a Peer packet.
func (p AddPeer) WriteTo(w io.Writer) (n int64, err error) {
	// +--------+--------------+---------+------+
	// |  0x01  | Address Type | Address | Port |
	// +--------+--------------+---------+------+

	var buf bytes.Buffer

	// Send Packet ID
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
	port := make([]byte, uint16Bytes)
	binary.LittleEndian.PutUint16(port, p.port)
	buf.Write(port)

	return buf.WriteTo(w)
}

// ReadFrom decodes a Peer packet.
func (p *AddPeer) ReadFrom(r io.Reader) (n int64, err error) {
	var addressLength int

	// Read address type
	addressType := make([]byte, 1)
	if _, err = r.Read(addressType); err != nil {
		return 0, err
	}

	// Read address
	if addressType[0] == ipv4Packet {
		addressLength = net.IPv4len
		_, err = r.Read(p.address[:addressLength])
	} else if addressType[0] == ipv6Packet {
		addressLength = net.IPv6len
		_, err = r.Read(p.address[:addressLength])
	} else {
		addressLength = 0
		err = errors.New("invalid address type")
	}

	// Check address read error
	if err != nil {
		return 0, err
	}

	// Read port
	port := make([]byte, uint16Bytes)
	if _, err = r.Read(port); err != nil {
		return 0, err
	}
	p.port = binary.LittleEndian.Uint16(port)

	return int64(len(addressType) + addressLength + len(port)), nil
}

// ChangeState is a packet that change global P2P state.
type ChangeState struct {
	State bool
	Time  time.Time
}

// WriteTo encodes a ChangeState packet.
// If the time hasn't been set, then it corresponds to the current time.
func (s ChangeState) WriteTo(w io.Writer) (n int64, err error) {
	// +--------+-----------+-------+
	// |  0x02  | Timestamp | State |
	// +--------+-----------+-------+

	var buf bytes.Buffer

	// Send Packet ID
	buf.WriteByte(changeStatePacket)

	// Set Time if it hasn't already been done
	if s.Time.IsZero() {
		s.Time = time.Now()
	}

	// Send Time
	timestamp := make([]byte, 8)
	binary.LittleEndian.PutUint64(timestamp, uint64(s.Time.UnixNano()))
	buf.Write(timestamp)

	// Current state
	if s.State {
		buf.WriteByte(stateTrue)
	} else {
		buf.WriteByte(stateFalse)
	}

	return buf.WriteTo(w)
}

// ReadFrom decodes a ChangeState packet.
func (s *ChangeState) ReadFrom(r io.Reader) (n int64, err error) {
	// Read Time
	timestamp := make([]byte, 8)
	if _, err = r.Read(timestamp); err != nil {
		return 0, err
	}

	// Parse Time from Unix timestamp
	s.Time = time.Unix(0, int64(binary.LittleEndian.Uint64(timestamp)))

	// Read state
	state := make([]byte, 1)
	if _, err = r.Read(state); err != nil {
		return 0, err
	}
	s.State = state[0] == stateTrue

	return int64(len(timestamp) + len(state)), nil
}
