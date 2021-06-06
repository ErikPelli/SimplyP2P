package SimplyP2P

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

// Counter is the number of connected peers.
type Counter int32

// Increment adds 1 to the value of counter.
func (c *Counter) Increment() {
	atomic.AddInt32((*int32)(c), 1)
}

// Decrement subtracts 1 from the value of the counter.
func (c *Counter) Decrement() {
	atomic.AddInt32((*int32)(c), ^int32(0))
}

// Value returns current value of the counter.
func (c *Counter) Value() int32 {
	return atomic.LoadInt32((*int32)(c))
}

// Peers is a map that contains all connected peers.
type Peers struct {
	peers       sync.Map
	peersNumber Counter
}

// Add adds a peer to the peers map.
func (p *Peers) Add(key Peer, value *Connection) {
	_, old := p.peers.Load(key)
	if !old {
		// Add peer if not already present
		p.peers.Store(key, value)
		p.peersNumber.Increment()
	}
}

// Remove removes a peer from the peers map.
func (p *Peers) Remove(key Peer) {
	if peer, ok := p.peers.LoadAndDelete(key); ok {
		_ = peer.(*Connection).Close()
		p.peersNumber.Decrement()
	}
}

// Get returns a peer from the peers map.
func (p *Peers) Get(key Peer) (*Connection, error) {
	if peer, ok := p.peers.Load(key); ok {
		return peer.(*Connection), nil
	} else {
		return nil, errors.New("peer not found")
	}
}

// Len returns the number of peers connected.
func (p *Peers) Len() int32 {
	return p.peersNumber.Value()
}

// Broadcast sends a packet to all peers connected.
func (p *Peers) Broadcast(to io.WriterTo) {
	fmt.Println("Sent new broadcast packet")

	p.peers.Range(func(key interface{}, value interface{}) bool {
		// Remove current peer if there is a connection problem
		if _, err := to.WriteTo(value.(*Connection)); err != nil {
			p.Remove(key.(Peer))
		}
		return true
	})
}

// Close closes all connections with other peers.
func (p *Peers) Close() error {
	p.peers.Range(func(key interface{}, value interface{}) bool {
		p.Remove(key.(Peer))
		return true
	})
	return nil
}
