package SimplyP2P

import (
	"encoding/binary"
	"errors"
	"net"
	"strconv"
)

type Connection struct {
	conn net.Conn
}

func (conn *Connection) Write(p []byte) (n int, err error) {
	return conn.conn.Write(p)
}

func (n *Node) Connect(address, port string) error {
	c := new(Connection)
	var err error

	dest := new(Peer)
	_ = dest.SetAddress(address)
	_ = dest.SetPort(port)

	// Skip current connection if there is already one
	if _, err = n.peers.Get(*dest); err == nil {
		return err
	}

	if c.conn, err = net.Dial("tcp", dest.GetAddressAndPort()); err != nil {
		return err
	}

	// Send P2P listening port
	portBytes := make([]byte, 2)
	portUint, err := strconv.ParseUint(n.listenPort, 10, 16)
	binary.LittleEndian.PutUint16(portBytes, uint16(portUint))
	if err := c.Send(portBytes); err != nil {
		return err
	}

	// Add peer to peers
	n.peers.Add(*dest, c)

	go n.handlePacket(c, *dest)
	return nil
}

func (conn *Connection) Send(data []byte) error {
	if conn.conn == nil {
		return errors.New("tcp connection is nil")
	}

	num, err := conn.conn.Write(data)
	if err != nil {
		return err
	}

	if num != len(data) {
		return errors.New("invalid number of sent bytes")
	}

	return nil
}

func (conn *Connection) Receive(length int) ([]byte, error) {
	if conn.conn == nil {
		return nil, errors.New("tcp connection is nil")
	}

	data := make([]byte, length)
	num, err := conn.conn.Read(data)
	if err != nil {
		return nil, err
	}

	if length > num {
		return nil, errors.New("some bytes are missing")
	}

	return data, nil
}

func (conn *Connection) Close() error {
	if conn.conn == nil {
		return errors.New("tcp connection is nil")
	}

	return conn.conn.Close()
}
