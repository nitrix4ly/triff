package commands

import (
	"errors"
	"sync"
)

type ListStore struct {
	data map[string][]string
	mu   sync.RWMutex
}

func NewListStore() *ListStore {
	return &ListStore{
		data: make(map[string][]string),
	}
}

func (ls *ListStore) LPush(key, value string) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	ls.data[key] = append([]string{value}, ls.data[key]...)
}

func (ls *ListStore) RPush(key, value string) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	ls.data[key] = append(ls.data[key], value)
}

func (ls *ListStore) LPop(key string) (string, error) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	if list, exists := ls.data[key]; exists && len(list) > 0 {
		val := list[0]
		ls.data[key] = list[1:]
		return val, nil
	}
	return "", errors.New("list is empty or key not found")
}

func (ls *ListStore) RPop(key string) (string, error) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	if list, exists := ls.data[key]; exists && len(list) > 0 {
		val := list[len(list)-1]
		ls.data[key] = list[:len(list)-1]
		return val, nil
	}
	return "", errors.New("list is empty or key not found")
}

func (ls *ListStore) LRange(key string) ([]string, error) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()
	if list, exists := ls.data[key]; exists {
		return list, nil
	}
	return nil, errors.New("key not found")
}
