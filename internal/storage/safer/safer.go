package safer

import (
	"encoding/json"
	"log/slog"
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

func (sfm *SafeMap) Put(inOrder string, log *slog.Logger) {
	sfm.mu.Lock() // wr lock
	defer sfm.mu.Unlock()

	var order storage.Order
	jsonOrder := []byte(inOrder)

	err := json.Unmarshal(jsonOrder, &order)
	if err != nil {
		log.Error("Error unmarshaling json data: ", err.Error())
	}

	sfm.orders[order.OrderUuid] = order
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
