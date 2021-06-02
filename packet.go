package SimplyP2P

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strconv"
)

const (
	ipv4 = 0x04
	ipv6 = 0x06
)

// NewPeer is a packet that adds a new peer to a node.
// +--------+--------------+---------+------+
// |  0x01  | Address Type | Address | Port |
// +--------+--------------+---------+------+
type NewPeer struct {
	addressType byte
	address		[]byte
	port		uint16
}

// WriteTo encodes a NewPeer packet.
func (p NewPeer) WriteTo(w io.Writer) (n int64, err error) {
	var b bytes.Buffer

	b.WriteByte(0x01)
	b.WriteByte(p.addressType)
	b.Write(p.address)

	var port []byte
	binary.LittleEndian.PutUint16(port, p.port)
	b.Write(port)

	return b.WriteTo(w)
}

// ReadFrom decodes a NewPeer packet.
func (p *NewPeer) ReadFrom(r io.Reader) (n int64, err error) {
	// Read address type
	addressType := make([]byte, 1)
	if _, err = r.Read(addressType); err != nil {
		return 0, err
	}
	p.addressType = addressType[0]

	if p.addressType == ipv4 {
		p.address = make([]byte, net.IPv4len)
	} else if p.addressType == ipv6 {
		p.address = make([]byte, net.IPv6len)
	} else {
		return 0, errors.New("invalid address type")
	}

	// Read address
	if _, err = r.Read(p.address); err != nil {
		return 0, err
	}

	// Read port
	port := make([]byte, 2) // TCP port is 2 bytes
	if _, err = r.Read(port); err != nil {
		return 0, err
	}
	p.port = binary.LittleEndian.Uint16(port)

	return int64(len(addressType) + len(p.address) + len(port)), nil
}

// GetAddress returns the string representation of IP address.
func (p NewPeer) GetAddress() string {
	add := net.IP(p.address).String()
	port := strconv.FormatUint(uint64(p.port), 10)

	if p.addressType == ipv4 {
		return add+":"+port
	} else if p.addressType == ipv6 {
		return "["+add+"]:"+port
	} else {
		return ""
	}
}

// SetAddress parses an address and save it to current NewPeer packet.
func (p *NewPeer) SetAddress(address, port string) error {
	ip := net.ParseIP(address)
	if ip == nil {
		return errors.New("invalid ip address")
	}
	p.address = ip

	if parsedPort, err := strconv.ParseUint(port, 10, 16); err != nil {
		return err
	} else {
		p.port = uint16(parsedPort)
		return nil
	}
}