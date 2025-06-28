package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DiskEngine struct {
	filePath string
	data     map[string]string
	mu       sync.RWMutex
}

func NewDiskEngine(path string) (*DiskEngine, error) {
	engine := &DiskEngine{
		filePath: path,
		data:     make(map[string]string),
	}
	if err := engine.load(); err != nil {
		return nil, err
	}
	return engine, nil
}

func (de *DiskEngine) load() error {
	de.mu.Lock()
	defer de.mu.Unlock()

	file, err := os.Open(de.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No file yet, skip
		}
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&de.data)
}

func (de *DiskEngine) persist() error {
	de.mu.RLock()
	defer de.mu.RUnlock()

	file, err := os.Create(de.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(de.data)
}

func (de *DiskEngine) Set(key, value string) error {
	de.mu.Lock()
	defer de.mu.Unlock()
	de.data[key] = value
	return de.persist()
}

func (de *DiskEngine) Get(key string) (string, error) {
	de.mu.RLock()
	defer de.mu.RUnlock()
	if val, exists := de.data[key]; exists {
		return val, nil
	}
	return "", errors.New("key not found")
}

func (de *DiskEngine) Delete(key string) error {
	de.mu.Lock()
	defer de.mu.Unlock()
	if _, exists := de.data[key]; exists {
		delete(de.data, key)
		return de.persist()
	}
	return errors.New("key not found")
}

func (de *DiskEngine) Exists(key string) bool {
	de.mu.RLock()
	defer de.mu.RUnlock()
	_, exists := de.data[key]
	return exists
}

func (de *DiskEngine) Keys() []string {
	de.mu.RLock()
	defer de.mu.RUnlock()
	keys := make([]string, 0, len(de.data))
	for key := range de.data {
		keys = append(keys, key)
	}
	return keys
}

func (de *DiskEngine) Flush() error {
	de.mu.Lock()
	defer de.mu.Unlock()
	de.data = make(map[string]string)
	return de.persist()
}
