package SimplyP2P

import (
	"sync"
	"time"
)

// State is the current global P2P state.
type State struct {
	state bool
	t     time.Time
	mtx   sync.RWMutex
	event func(bool)
}

// SetEvent sets a function to call when the value change.
func (s *State) SetEvent(event func(bool)) {
	s.mtx.Lock()
	s.event = event
	s.mtx.Unlock()
}

// Update updates the state value if there is a more recent value change.
// Returns boolean update (true if has been updated).
func (s *State) Update(value bool, t time.Time) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if t.After(s.t) {
		s.state = value
		s.t = t

		// Execute event function
		if s.event != nil {
			s.event(value)
		}
		return true
	} else {
		return false
	}
}

// GetState returns current state value.
func (s *State) GetState() bool {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.state
}

// GetTime returns current time value.
func (s *State) GetTime() time.Time {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.t
}
