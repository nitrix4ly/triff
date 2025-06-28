package storage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sync"
)

// PersistentStore provides simple JSON-based persistence for map[string]string
type PersistentStore struct {
	filePath string
	data     map[string]string
	mu       sync.RWMutex
}

func NewPersistentStore(filePath string) (*PersistentStore, error) {
	ps := &PersistentStore{
		filePath: filePath,
		data:     make(map[string]string),
	}
	err := ps.Load()
	return ps, err
}

func (ps *PersistentStore) Load() error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	file, err := os.Open(ps.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // First time setup
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&ps.data)
}

func (ps *PersistentStore) Save() error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	bytes, err := json.MarshalIndent(ps.data, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(ps.filePath, bytes, 0644)
}

func (ps *PersistentStore) Set(key, value string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.data[key] = value
	ps.Save()
}

func (ps *PersistentStore) Get(key string) (string, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	val, ok := ps.data[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return val, nil
}
