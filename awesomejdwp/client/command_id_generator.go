package client

import (
	"sync"
)

var (
	mu        sync.Mutex
	currentId uint32
)

func NextCommandId() uint32 {
	mu.Lock()
	defer mu.Unlock()

	currentId++
	return currentId
}
