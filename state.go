package SimplyP2P

import (
	"sync"
	"time"
)

// State is the current global P2P state.
type State struct {
	state    bool
	stateMtx sync.Mutex
	event    func(bool)
	time.Time
}

// SetEvent sets a function to call when the value change.
func (s *State) SetEvent(event func(bool)) {
	s.event = event
}

// Update updates the state value if there is a more recent value change.
func (s *State) Update(value bool, t time.Time) {
	s.stateMtx.Lock()
	if t.After(s.Time) {
		s.state = value
		s.event(value)
	}
	s.stateMtx.Unlock()
}

// GetState returns current state value.
func (s *State) GetState() bool {
	return s.state
}
