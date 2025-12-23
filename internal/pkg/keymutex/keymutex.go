package keymutex

import (
	"sync"
)

// KeyMutex provides mutex locking mechanism based on a key.
// This is useful when you need to ensure that operations with the same key
// cannot run concurrently, while operations with different keys can.
type KeyMutex struct {
	mu   sync.Mutex
	maps sync.Map // map[any]*sync.Mutex
}

// New creates a new KeyMutex instance.
func New() *KeyMutex {
	return &KeyMutex{}
}

// Lock acquires the lock for the given key.
// If a lock for this key already exists, it will wait until it's released.
// The lock should be released by calling Unlock with the same key.
func (km *KeyMutex) Lock(key any) {
	mu, _ := km.maps.LoadOrStore(key, &sync.Mutex{})
	mu.(*sync.Mutex).Lock()
}

// Unlock releases the lock for the given key.
// The key must be the same as the one used for Lock.
func (km *KeyMutex) Unlock(key any) {
	mu, ok := km.maps.Load(key)
	if !ok {
		panic("KeyMutex.Unlock: unlock of unlocked key")
	}
	mu.(*sync.Mutex).Unlock()
}
