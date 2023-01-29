package rwmutex

import (
	"sync"
)

// RWMutex is a sync.RWMutex with a switch for concurrent safe feature.
type RWMutex struct {
	// Underlying rwmutex.
	rwmutex *sync.RWMutex
}

// New creates and returns a new *RWMutex.
// The parameter `safe` is used to specify whether using this rwmutex in concurrent safety,
// which is false in default.
func New(safe ...bool) *RWMutex {
	mu := Create(safe...)
	return &mu
}

// Create creates and returns a new RWMutex object.
// The parameter `safe` is used to specify whether using this rwmutex in concurrent safety,
// which is false in default.
func Create(safe ...bool) RWMutex {
	if len(safe) > 0 && safe[0] {
		return RWMutex{
			rwmutex: new(sync.RWMutex),
		}
	}
	return RWMutex{}
}

// IsSafe checks and returns whether current rwmutex is in concurrent-safe usage.
func (mu *RWMutex) IsSafe() bool {
	return mu.rwmutex != nil
}

// Lock locks rwmutex for writing.
// It does nothing if it is not in concurrent-safe usage.
func (mu *RWMutex) Lock() {
	if mu.rwmutex != nil {
		mu.rwmutex.Lock()
	}
}

// Unlock unlocks rwmutex for writing.
// It does nothing if it is not in concurrent-safe usage.
func (mu *RWMutex) Unlock() {
	if mu.rwmutex != nil {
		mu.rwmutex.Unlock()
	}
}

// RLock locks rwmutex for reading.
// It does nothing if it is not in concurrent-safe usage.
func (mu *RWMutex) RLock() {
	if mu.rwmutex != nil {
		mu.rwmutex.RLock()
	}
}

// RUnlock unlocks rwmutex for reading.
// It does nothing if it is not in concurrent-safe usage.
func (mu *RWMutex) RUnlock() {
	if mu.rwmutex != nil {
		mu.rwmutex.RUnlock()
	}
}
