package SimplyP2P

import (
	"errors"
	"net"
	"strconv"
)

// Peer is another peer in the P2P network.
type Peer struct {
	address       [net.IPv6len]byte
	addressLength int
	port          uint16
}

// GetAddress returns the string representation of IP address of current Peer.
func (p Peer) GetAddress() string {
	return net.IP(p.address[:]).String()
}

// GetPort returns the string representation of listening port of current Peer.
func (p Peer) GetPort() string {
	return strconv.FormatUint(uint64(p.port), 10)
}

// GetAddressAndPort returns the string representation of IP address and port of current Peer.
func (p Peer) GetAddressAndPort() string {
	add := p.GetAddress()
	port := p.GetPort()

	if p.addressLength == net.IPv4len {
		return add + ":" + port
	} else if p.addressLength == net.IPv6len {
		return "[" + add + "]:" + port
	} else {
		return ""
	}
}

// SetAddress parses an address and save it to current Peer.
func (p *Peer) SetAddress(address string) error {
	// Local IP address if the function argument is empty
	if address == "" {
		address = "127.0.0.1"
	}

	// Parse IP address
	ip := net.ParseIP(address)
	if ip == nil {
		return errors.New("invalid IP address")
	}
	p.addressLength = copy(p.address[:], ip)

	return nil
}

// SetPort parses a port and save it to current Peer.
func (p *Peer) SetPort(port string) error {
	parsedPort, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return err
	}
	p.port = uint16(parsedPort)

	return nil
}
