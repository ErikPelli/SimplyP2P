package SimplyP2P

import (
	"errors"
	"fmt"
	"net"
)

type Connection struct {
	conn net.Conn
}

func (n *Node) Listen(port string) error {
	listener, err := net.Listen("tcp", ":" + port)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handler(conn)
	}
}

func (n *Node) Connect(address string) error {
	c := new(Connection)

	var err error
	if c.conn, err = net.Dial("tcp", address); err != nil {
		return err
	}

	// If already present, do nothing
	if _, ok := n.peers[address]; !ok {
		n.peers[address] = c
	}

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