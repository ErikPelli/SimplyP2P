package SimplyP2P

import (
	"errors"
	"net"
)

// Connection is a TCP connection to another peer.
type Connection struct {
	conn net.Conn
}

// Connect connects to a remote peer with the destination IP and port in the Peer argument.
func (conn *Connection) Connect(destination Peer) error {
	var err error
	conn.conn, err = net.Dial("tcp", destination.GetAddressAndPort())
	return err
}

// Write implements the io.Writer interface for the Connection type.
func (conn *Connection) Write(p []byte) (n int, err error) {
	return conn.conn.Write(p)
}

// Send sends a byte slice to the other peer in this connection.
// Length of the data sent is the length of the data slice.
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

// Receive returns a byte slice with its length specified as argument,
// filled with bytes read from this connection.
func (conn *Connection) Receive(length int) ([]byte, error) {
	if conn.conn == nil {
		return nil, errors.New("tcp connection is nil")
	}

	data := make([]byte, length)
	num, err := conn.conn.Read(data)

	if err != nil {
		return nil, err
	}
	if num < length {
		return nil, errors.New("some bytes are missing")
	}

	return data, nil
}

// Close closes this connection.
// Implements the io.Closer interface.
func (conn *Connection) Close() error {
	if conn.conn == nil {
		return errors.New("tcp connection is nil")
	}

	return conn.conn.Close()
}
