package SimplyP2P

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
)

// Node is current running peer in the network.
// It has a listener on a specific port and contains
// connections with other peers.
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
	var err error

	// Save listening port
	n.listenPort = listenPort

	// Add current node to wait group
	wg.Add(1)
	n.wg = wg

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

	// Connect to a peer if there are some peers in the arguments
	if len(peers) != 0 {
		for _, peer := range peers {
			// if current peer is valid, else check the next
			if err := n.connect(peer[0], peer[1]); err != nil {
				fmt.Println("Unable to connect to peer: " + err.Error())
				break
			}
		}

		if n.peers.Len() == 0 {
			return nil, errors.New("there is no peer")
		}
	}

	// Initialize GUI
	go n.newGui()

	return n, nil
}

// connect connects to a specified host (address and port).
func (n *Node) connect(address, port string) error {
	// Set destination peer parameters
	dest := new(Peer)
	_ = dest.SetAddress(address)
	_ = dest.SetPort(port)

	// Skip current connection if there is already one
	if _, err := n.peers.Get(*dest); err == nil {
		return errors.New("current connection already exists")
	}

	// Connect to peer
	c := new(Connection)
	if err := c.Connect(dest.GetAddressAndPort()); err != nil {
		return err
	}

	// Send current node listening port
	portBytes := make([]byte, 2)
	portUint, _ := strconv.ParseUint(n.listenPort, 10, 16)
	binary.LittleEndian.PutUint16(portBytes, uint16(portUint))
	if err := c.Send(portBytes); err != nil {
		return err
	}

	n.peers.Add(*dest, c)
	go n.packetsHandler(c, *dest)

	return nil
}

// listen listens for connections on a specified port.
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

// newConnectionHandler handles a new connection received in
// listening port and send to the client initialization data.
func (n *Node) newConnectionHandler(conn Connection) {
	// Get listen port from client
	portBytes, err := conn.Receive(2)
	if err != nil {
		return
	}

	// Send address of every connected peer to the client
	n.peers.peers.Range(func(key interface{}, value interface{}) bool {
		_, _ = (AddPeer{key.(Peer)}).WriteTo(&conn)
		return true
	})

	// Send current P2P state
	if _, err = (ChangeState{State: n.state.GetState(), Time: n.state.GetTime()}.WriteTo(&conn)); err != nil {
		fmt.Println(err)
		return
	}

	// Add current connection to peers
	currentPeer := Peer{
		port: binary.LittleEndian.Uint16(portBytes),
	}
	if err := currentPeer.SetAddress(conn.conn.RemoteAddr().(*net.TCPAddr).IP.String()); err != nil {
		return
	}
	n.peers.Add(currentPeer, &conn)

	// Start packets handler in current goroutine
	n.packetsHandler(&conn, currentPeer)
}

// packetsHandler handle received packets in conn connection
// and modify the node fields according to them.
func (n *Node) packetsHandler(conn *Connection, mapKey Peer) {
	defer func() {
		n.peers.Remove(mapKey)
	}()

	// Handler for this connection
	for !n.closed {
		packetID, err := conn.Receive(1)
		if err != nil {
			fmt.Println("packets handler: " + err.Error())
			break
		}

		switch packetID[0] {
		case newPeerPacket:
			peer := new(AddPeer)
			if _, err := peer.ReadFrom(conn.conn); err != nil {
				fmt.Println("packets handler: " + err.Error())
				break
			}
			_ = n.connect(peer.GetAddress(), peer.GetPort())

		case changeStatePacket:
			state := new(ChangeState)
			if _, err := state.ReadFrom(conn.conn); err != nil {
				fmt.Println("packets handler: " + err.Error())
				break
			}
			// If updated successfully, broadcast this state update
			if n.state.Update(state.State, state.Time) {
				n.peers.Broadcast(state)
			}
		}
	}
}

// Close closes the connection with every peer.
func (n *Node) Close() error {
	_ = n.peers.Close()
	n.wg.Done()
	return nil
}
