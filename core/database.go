package core

import (
	"sync"
	"time"
)

// NewDatabase creates a new Triff database instance
func NewDatabase(config *Config) *Database {
	return &Database{
		Data:   make(map[string]*TriffValue),
		mu:     sync.RWMutex{},
		config: config,
	}
}

// Get retrieves a value from the database
func (db *Database) Get(key string) (*TriffValue, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	value, exists := db.Data[key]
	if !exists {
		return nil, false
	}
	
	// Check if value has expired
	if value.TTL > 0 && time.Now().Unix() > value.TTL {
		delete(db.Data, key)
		return nil, false
	}
	
	return value, true
}

// Set stores a value in the database
func (db *Database) Set(key string, value *TriffValue) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	now := time.Now()
	value.UpdatedAt = now
	
	if _, exists := db.Data[key]; !exists {
		value.CreatedAt = now
	}
	
	db.Data[key] = value
	return nil
}

// Delete removes a key from the database
func (db *Database) Delete(key string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	if _, exists := db.Data[key]; exists {
		delete(db.Data, key)
		return true
	}
	return false
}

// Exists checks if a key exists in the database
func (db *Database) Exists(key string) bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	_, exists := db.Data[key]
	return exists
}

// Keys returns all keys matching a pattern
func (db *Database) Keys(pattern string) []string {
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	keys := make([]string, 0)
	for key := range db.Data {
		// Simple pattern matching - can be enhanced
		if pattern == "*" || key == pattern {
			keys = append(keys, key)
		}
	}
	return keys
}

// FlushAll removes all data from the database
func (db *Database) FlushAll() error {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	db.Data = make(map[string]*TriffValue)
	return nil
}

// Size returns the number of keys in the database
func (db *Database) Size() int64 {
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	return int64(len(db.Data))
}

// SetTTL sets time to live for a key
func (db *Database) SetTTL(key string, seconds int64) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	if value, exists := db.Data[key]; exists {
		value.TTL = time.Now().Unix() + seconds
		return true
	}
	return false
}

// GetTTL returns time to live for a key
func (db *Database) GetTTL(key string) int64 {
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	if value, exists := db.Data[key]; exists {
		if value.TTL == 0 {
			return -1 // No expiration
		}
		remaining := value.TTL - time.Now().Unix()
		if remaining <= 0 {
			return -2 // Expired
		}
		return remaining
	}
	return -2 // Key doesn't exist
}

// CleanupExpired removes expired keys from the database
func (db *Database) CleanupExpired() {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	now := time.Now().Unix()
	for key, value := range db.Data {
		if value.TTL > 0 && now > value.TTL {
			delete(db.Data, key)
		}
	}
}

// Info returns database information
func (db *Database) Info() map[string]interface{} {
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	return map[string]interface{}{
		"version":    "1.0.0",
		"keys":       len(db.Data),
		"memory_mb":  db.getMemoryUsage(),
		"uptime":     time.Since(time.Now()).Seconds(),
		"tcp_port":   db.config.Port,
		"http_port":  db.config.HTTPPort,
	}
}

// getMemoryUsage calculates approximate memory usage
func (db *Database) getMemoryUsage() int64 {
	// Simple estimation - can be enhanced with proper memory calculation
	return int64(len(db.Data) * 100) // Rough estimate
}

// Ping returns pong - health check
func (db *Database) Ping() string {
	return "PONG"
}
