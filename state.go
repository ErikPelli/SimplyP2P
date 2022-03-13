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

	// Update state if new time is after the current time.
	// If instead the time is equal, give precedence to the higher value
	// (true in boolean), to have a consistent value in case of time collision.
	updateValue := t.After(s.t) || t.Equal(s.t) && s.state != value && value
	if updateValue {
		s.state = value
		s.t = t

		// Execute event function
		if s.event != nil {
			s.event(value)
		}
	}
	return updateValue
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
