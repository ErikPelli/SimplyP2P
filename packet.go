package SimplyP2P

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strconv"
	"time"
)

const (
	newPeerPacket     = 0x01
	changeStatePacket = 0x02
	ipv4Packet        = 0x04
	ipv6Packet        = 0x06
)

// Peer is a packet that adds a new peer to a node.
// Use an array instead of a slice to use it has map hashed key.
// +--------+--------------+---------+------+
// |  0x01  | Address Type | Address | Port |
// +--------+--------------+---------+------+
type Peer struct {
	address       [net.IPv6len]byte
	addressLength int
	port          uint16
}

// WriteTo encodes a Peer packet.
func (p Peer) WriteTo(w io.Writer) (n int64, err error) {
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
func (p *Peer) ReadFrom(r io.Reader) (n int64, err error) {
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

// GetAddress returns the string representation of IP address of Peer.
func (p Peer) GetAddress() string {
	add := net.IP(p.address[:]).String()
	port := strconv.FormatUint(uint64(p.port), 10)

	if p.addressLength == net.IPv4len {
		return add + ":" + port
	} else if p.addressLength == net.IPv6len {
		return "[" + add + "]:" + port
	} else {
		return ""
	}
}

// SetAddress parses an address and save it to current Peer packet.
func (p *Peer) SetAddress(address, port string) error {
	// Local ip if address wasn't set
	if address == "" {
		address = "127.0.0.1"
	}

	// Parse IP address
	ip := net.ParseIP(address)
	if ip == nil {
		return errors.New("invalid ip address")
	}
	p.addressLength = copy(p.address[:], ip)

	if parsedPort, err := strconv.ParseUint(port, 10, 16); err != nil {
		return err
	} else {
		p.port = uint16(parsedPort)
		return nil
	}
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
	binary.LittleEndian.PutUint64(timestamp, uint64(time.Now().UnixNano()))
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
