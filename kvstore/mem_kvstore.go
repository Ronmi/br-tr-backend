package kvstore

import (
	"sync"
	"time"
)

type memoryStoreElement struct {
	value Data
	until time.Time
}
type memoryStore struct {
	data map[string]memoryStoreElement
	lock *sync.RWMutex
}

func (s *memoryStore) Get(k string) (ret Data) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.doGet(k)
}

func (s *memoryStore) Set(k string, v Data, ttl int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.doSet(k, v, ttl)
}

func (s *memoryStore) SetIf(k string, v, oldValue Data, ttl int) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	ret := false
	if s.doGet(k) == oldValue {
		s.doSet(k, v, ttl)
		ret = true
	}
	return ret
}

func (s *memoryStore) doGet(k string) (ret Data) {
	if data, ok := s.data[k]; ok && time.Now().Before(data.until) {
		ret = data.value
	}
	return
}

func (s *memoryStore) doSet(k string, v Data, ttl int) {
	if ttl <= 0 {
		delete(s.data, k)
		return
	}

	s.data[k] = memoryStoreElement{
		value: v,
		until: time.Now().Add(time.Duration(ttl) * time.Millisecond),
	}
}

func (s *memoryStore) clean() {
	s.lock.Lock()
	defer s.lock.Unlock()

	t := time.Now()
	for k, _ := range s.data {
		if t.Before(s.data[k].until) {
			delete(s.data, k)
		}
	}
}

// periodically cleanup data
func (s *memoryStore) run() {
	for {
		s.clean()
		time.Sleep(1 * time.Second)
	}
}

func NewMemStore() KVStore {
	return &memoryStore{
		data: map[string]memoryStoreElement{},
		lock: &sync.RWMutex{},
	}
}
