package SimplyP2P

import (
	"errors"
	"io"
	"sync"
)

// Peers is a map that contains all connected peers
type Peers struct {
	peers      sync.Map
	counter    int
	counterMtx sync.Mutex
}

// Add adds a peer to the map
func (p *Peers) Add(key string, value *Connection) {
	_, old := p.peers.Load(key)
	if !old {
		// Add peer if not already present
		p.peers.Store(key, value)
		p.counterMtx.Lock()
		p.counter++
		p.counterMtx.Unlock()
	}
}

// Remove removes a peer from the map
func (p *Peers) Remove(key string) {
	if peer, ok := p.peers.LoadAndDelete(key); ok {
		currentPeer := peer.(*Connection)
		_ = currentPeer.Close()

		// Decrement counter
		p.counterMtx.Lock()
		p.counter--
		p.counterMtx.Unlock()
	}
}

// Get returns a peer from the map
func (p *Peers) Get(key string) (*Connection, error) {
	if peer, ok := p.peers.Load(key); ok {
		return peer.(*Connection), nil
	} else {
		return nil, errors.New("peer not found")
	}
}

// Len returns the number of peers connected.
func (p *Peers) Len() int {
	return p.counter
}

// Broadcast sends a packet to all peers connected.
func (p *Peers) Broadcast(to io.WriterTo) {
	p.peers.Range(func(key interface{}, value interface{}) bool {
		peer := value.(*Connection)

		// Remove current peer if there is a connection problem
		if _, err := to.WriteTo(peer); err != nil {
			p.Remove(key.(string))
		}

		return true
	})
}