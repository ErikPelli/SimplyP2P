package SimplyP2P

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"sync"
)

type Node struct {
	listenPort string
	peers      Peers
	state      State
	closed     bool
	wg         *sync.WaitGroup
}

// NewNode creates a new peer to peer node instance
// specifying local listening port and a slice of known peers.
func NewNode(listenPort string, peers [][]string, wg *sync.WaitGroup) (*Node, error) {
	n := new(Node)
	n.listenPort = listenPort
	var err error

	// Add node to wait group
	n.wg = wg
	wg.Add(1)

	// Close node if there is an error
	defer func(n *Node) {
		if err != nil {
			_ = n.Close()
		}
	}(n)

	// Start listening
	if err = n.listen(listenPort); err != nil {
		return nil, err
	}

	// Connect to a peer if there are some specified
	if len(peers) != 0 {
		for _, peer := range peers {
			// if current peer is valid, else check the next
			if n.Connect(peer[0], peer[1]) != nil {
				break
			}
		}

		if n.peers.Len() == 0 {
			return nil, errors.New("there is no peer")
		}
	}

	return n, nil
}

// Listen listens for connections on a specified port.
func (n *Node) listen(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	go func() {
		defer listener.Close()
		for !n.closed {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println(err)
				return
			}
			go n.newConnectionHandler(Connection{conn: conn})
		}
	}()

	return nil
}

func (n *Node) newConnectionHandler(conn Connection) {
	// Get listen port
	portBytes, err := conn.Receive(2)
	if err != nil {
		return
	}

	// Send address of every connected peer to new peer
	n.peers.peers.Range(func(key interface{}, value interface{}) bool {
		peerAddress := key.(Peer)
		if _, err := peerAddress.WriteTo(&conn); err != nil {
			n.peers.Remove(key.(Peer))
		}
		return true
	})

	// Add current connection to peers
	remoteIP := new(Peer)
	if err := remoteIP.SetAddress(conn.conn.RemoteAddr().(*net.TCPAddr).IP.String(), "0"); err != nil {
		return
	}
	currentPeer := Peer{
		address:       remoteIP.address,
		addressLength: remoteIP.addressLength,
		port:          binary.LittleEndian.Uint16(portBytes),
	}
	n.peers.Add(currentPeer, &conn)

	// Send current P2P state
	if _, err = (ChangeState{State: n.state.GetState()}.WriteTo(&conn)); err != nil {
		fmt.Println(err)
		n.peers.Remove(currentPeer)
	}

	n.handlePacket(&conn, currentPeer)
}

func (n *Node) handlePacket(conn *Connection, mapKey Peer) {
	defer func() {
		n.peers.Remove(mapKey)
	}()

	// Handler for this connection
	for !n.closed {
		packetID, err := conn.Receive(1)
		if err != nil {
			break
		}

		switch packetID[0] {
		case newPeerPacket:
			peer := new(Peer)
			if _, err := peer.ReadFrom(conn.conn); err != nil {
				break
			}
			_ = n.Connect("tcp", peer.GetAddress())

		case changeStatePacket:
			state := new(ChangeState)
			if _, err := state.ReadFrom(conn.conn); err != nil {
				break
			}
			n.state.Update(state.State, state.time)
		}
	}
}

// Close closes the connection with every peer.
func (n *Node) Close() error {
	_ = n.peers.Close()
	n.wg.Done()
	return nil
}
