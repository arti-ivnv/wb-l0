package safer

import (
	"ls-0/arti/order/internal/storage"
	"sync"
)

type SafeMap struct {
	mu     sync.RWMutex
	orders map[string]storage.Order
}

// initialize a new map
func NewSafeMap() *SafeMap {
	return &SafeMap{orders: make(map[string]storage.Order)}
}

func (sfm *SafeMap) Put(key string, order storage.Order) {
	sfm.mu.Lock() // wr lock
	defer sfm.mu.Unlock()

	sfm.orders[key] = order
}

func (sfm *SafeMap) Get(key string) (storage.Order, bool) {
	sfm.mu.RLock() // rd lock
	defer sfm.mu.RUnlock()

	order, ok := sfm.orders[key]

	return order, ok
}

func (sfm *SafeMap) Delete(key string) {
	sfm.mu.Lock()
	defer sfm.mu.Unlock()

	delete(sfm.orders, key)
}
