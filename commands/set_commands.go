package commands

import (
	"errors"
	"sync"
)

type SetStore struct {
	data map[string]map[string]struct{}
	mu   sync.RWMutex
}

func NewSetStore() *SetStore {
	return &SetStore{
		data: make(map[string]map[string]struct{}),
	}
}

func (ss *SetStore) SAdd(key, value string) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if _, exists := ss.data[key]; !exists {
		ss.data[key] = make(map[string]struct{})
	}
	ss.data[key][value] = struct{}{}
}

func (ss *SetStore) SRem(key, value string) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if set, exists := ss.data[key]; exists {
		if _, found := set[value]; found {
			delete(set, value)
			return nil
		}
	}
	return errors.New("value not found in set")
}

func (ss *SetStore) SMembers(key string) ([]string, error) {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	if set, exists := ss.data[key]; exists {
		members := []string{}
		for val := range set {
			members = append(members, val)
		}
		return members, nil
	}
	return nil, errors.New("key not found")
}

func (ss *SetStore) SExists(key, value string) bool {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	if set, exists := ss.data[key]; exists {
		_, found := set[value]
		return found
	}
	return false
}
