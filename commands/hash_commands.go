package commands

import (
	"errors"
	"sync"
)

type HashStore struct {
	data map[string]map[string]string
	mu   sync.RWMutex
}

func NewHashStore() *HashStore {
	return &HashStore{
		data: make(map[string]map[string]string),
	}
}

func (hs *HashStore) HSet(key, field, value string) {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	if _, exists := hs.data[key]; !exists {
		hs.data[key] = make(map[string]string)
	}
	hs.data[key][field] = value
}

func (hs *HashStore) HGet(key, field string) (string, error) {
	hs.mu.RLock()
	defer hs.mu.RUnlock()
	if fields, exists := hs.data[key]; exists {
		if val, ok := fields[field]; ok {
			return val, nil
		}
		return "", errors.New("field not found")
	}
	return "", errors.New("key not found")
}

func (hs *HashStore) HDel(key, field string) error {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	if fields, exists := hs.data[key]; exists {
		delete(fields, field)
		if len(fields) == 0 {
			delete(hs.data, key)
		}
		return nil
	}
	return errors.New("key not found")
}

func (hs *HashStore) HGetAll(key string) (map[string]string, error) {
	hs.mu.RLock()
	defer hs.mu.RUnlock()
	if fields, exists := hs.data[key]; exists {
		return fields, nil
	}
	return nil, errors.New("key not found")
}
