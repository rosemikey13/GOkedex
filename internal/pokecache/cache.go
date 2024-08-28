package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct{
	createdAt time.Time
	val []byte
}

type Cache struct{
	Data map[string]cacheEntry
	Mutex *sync.Mutex
}

func NewCache(interval int) Cache{
	newCache := Cache{
		map[string]cacheEntry{},
		&sync.Mutex{},
	}

	go newCache.readLoop(time.Duration(time.Duration(interval) * time.Second))

	return newCache
}

func (c Cache) Add(key string, val []byte) {
	defer c.Mutex.Unlock()
	c.Mutex.Lock()
	c.Data[key] = cacheEntry{time.Now(), val}
}

func (c Cache) Get(key string) ([]byte, bool) {
	defer c.Mutex.Unlock()
	c.Mutex.Lock()

	if entry, exists := c.Data[key]; exists{
		return entry.val, true
	}

	return nil, false
}

func (c Cache) readLoop(interval time.Duration){

	ticker := time.NewTicker(interval)

	for {
		currentTime :=  <-ticker.C
		for k, v := range c.Data{
			
			expirationTime := v.createdAt.Add(interval)
			
			if currentTime.After(expirationTime) {
				c.Mutex.Lock()
				delete(c.Data, k)
				c.Mutex.Unlock()	
			} 
		}
	}
} 