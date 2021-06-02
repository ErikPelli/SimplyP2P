package SimplyP2P

import (
	"errors"
)

type Node struct {
	peers map[string]*Connection

}

// NewNode creates a new peer to peer node instance
// specifying local listening port and a slice of known peers.
func NewNode(listenPort string, peers []string) (n *Node, err error) {
	if len(peers) == 0 {
		return nil, errors.New("there is no peer")
	}

	n = new(Node)

	// Clean node if there is an error
	defer func() {
		if err != nil {
			_ = n.Close()
		}
	}()

	// Start listening
	if err = n.Listen(listenPort); err != nil {
		return nil, err
	}

	// Connect to a peer
	for _, peer := range peers {
		// if current peer is valid, else check the next
		if n.Connect(peer) != nil {
			break
		}
	}

	if len(n.peers) == 0 {
		return nil, errors.New("there is no peer")
	}

	return n, nil
}

// Close closes the connection with every peer.
func (n *Node) Close() error {
	for _, peer := range n.peers {
		_ = peer.Close()
	}
	return nil
}