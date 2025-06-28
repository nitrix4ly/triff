package storage

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/nitrix4ly/triff/core"
)

// MemoryEngine implements in-memory storage with optional persistence
type MemoryEngine struct {
	data            map[string]*core.TriffValue
	mu              sync.RWMutex
	persistencePath string
	autoSave        bool
	saveInterval    time.Duration
	stopChan        chan bool
}

// NewMemoryEngine creates a new memory storage engine
func NewMemoryEngine(persistencePath string, autoSave bool) *MemoryEngine {
	engine := &MemoryEngine{
		data:            make(map[string]*core.TriffValue),
		mu:              sync.RWMutex{},
		persistencePath: persistencePath,
		autoSave:        autoSave,
		saveInterval:    30 * time.Second, // Save every 30 seconds
		stopChan:        make(chan bool),
	}
	
	// Load existing data if available
	if persistencePath != "" {
		engine.loadFromDisk()
	}
	
	// Start auto-save routine if enabled
	if autoSave && persistencePath != "" {
		go engine.autoSaveRoutine()
	}
	
	return engine
}

// Get retrieves a value from memory
func (me *MemoryEngine) Get(key string) (*core.TriffValue, bool) {
	me.mu.RLock()
	defer me.mu.RUnlock()
	
	value, exists := me.data[key]
	if !exists {
		return nil, false
	}
	
	// Check if value has expired
	if value.TTL > 0 && time.Now().Unix() > value.TTL {
		delete(me.data, key)
		return nil, false
	}
	
	return value, true
}

// Set stores a value in memory
func (me *MemoryEngine) Set(key string, value *core.TriffValue) error {
	me.mu.Lock()
	defer me.mu.Unlock()
	
	now := time.Now()
	value.UpdatedAt = now
	
	if _, exists := me.data[key]; !exists {
		value.CreatedAt = now
	}
	
	me.data[key] = value
	return nil
}

// Delete removes a key from memory
func (me *MemoryEngine) Delete(key string) bool {
	me.mu.Lock()
	defer me.mu.Unlock()
	
	if _, exists := me.data[key]; exists {
		delete(me.data, key)
		return true
	}
	return false
}

// Exists checks if a key exists in memory
func (me *MemoryEngine) Exists(key string) bool {
	me.mu.RLock()
	defer me.mu.RUnlock()
	
	_, exists := me.data[key]
	return exists
}

// Keys returns all keys matching a pattern
func (me *MemoryEngine) Keys(pattern string) []string {
	me.mu.RLock()
	defer me.mu.RUnlock()
	
	keys := make([]string, 0)
	for key := range me.data {
		// Simple pattern matching - can be enhanced with regex
		if pattern == "*" || key == pattern {
			keys = append(keys, key)
		}
	}
	return keys
}

// FlushAll removes all data from memory
func (me *MemoryEngine) FlushAll() error {
	me.mu.Lock()
	defer me.mu.Unlock()
	
	me.data = make(map[string]*core.TriffValue)
	return nil
}

// Size returns the number of keys in memory
func (me *MemoryEngine) Size() int64 {
	me.mu.RLock()
	defer me.mu.RUnlock()
	
	return int64(len(me.data))
}

// CleanupExpired removes expired keys from memory
func (me *MemoryEngine) CleanupExpired() int {
	me.mu.Lock()
	defer me.mu.Unlock()
	
	now := time.Now().Unix()
	removed := 0
	
	for key, value := range me.data {
		if value.TTL > 0 && now > value.TTL {
			delete(me.data, key)
			removed++
		}
	}
	
	return removed
}

// SaveToDisk saves current data to disk
func (me *MemoryEngine) SaveToDisk() error {
	if me.persistencePath == "" {
		return nil
	}
	
	me.mu.RLock()
	defer me.mu.RUnlock()
	
	// Create a copy of data for serialization
	dataCopy := make(map[string]*core.TriffValue)
	for k, v := range me.data {
		dataCopy[k] = v
	}
	
	// Serialize to JSON
	jsonData, err := json.MarshalIndent(dataCopy, "", "  ")
	if err != nil {
		return err
	}
	
	// Write to file
	return os.WriteFile(me.persistencePath, jsonData, 0644)
}

// loadFromDisk loads data from disk if file exists
func (me *MemoryEngine) loadFromDisk() error {
	if me.persistencePath == "" {
		return nil
	}
	
	// Check if file exists
	if _, err := os.Stat(me.persistencePath); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to load
	}
	
	// Read file
	jsonData, err := os.ReadFile(me.persistencePath)
	if err != nil {
		return err
	}
	
	// Deserialize JSON
	var loadedData map[string]*core.TriffValue
	if err := json.Unmarshal(jsonData, &loadedData); err != nil {
		return err
	}
	
	me.mu.Lock()
	defer me.mu.Unlock()
	
	me.data = loadedData
	return nil
}

// autoSaveRoutine runs in background to periodically save data
func (me *MemoryEngine) autoSaveRoutine() {
	ticker := time.NewTicker(me.saveInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			if err := me.SaveToDisk(); err != nil {
				// Log error but don't stop the routine
				continue
			}
		case <-me.stopChan:
			// Final save before stopping
			me.SaveToDisk()
			return
		}
	}
}

// Stop stops the auto-save routine and saves data
func (me *MemoryEngine) Stop() error {
	if me.autoSave {
		me.stopChan <- true
	}
	return me.SaveToDisk()
}

// GetMemoryUsage returns approximate memory usage in bytes
func (me *MemoryEngine) GetMemoryUsage() int64 {
	me.mu.RLock()
	defer me.mu.RUnlock()
	
	// Simple estimation - can be enhanced with proper memory calculation
	usage := int64(0)
	for key, value := range me.data {
		usage += int64(len(key))
		switch v := value.Data.(type) {
		case string:
			usage += int64(len(v))
		case []interface{}:
			usage += int64(len(v) * 8) // Rough estimate
		case map[string]interface{}:
			usage += int64(len(v) * 16) // Rough estimate
		default:
			usage += 8 // Basic type estimate
		}
	}
	
	return usage
}

// GetStats returns storage statistics
func (me *MemoryEngine) GetStats() map[string]interface{} {
	me.mu.RLock()
	defer me.mu.RUnlock()
	
	stats := map[string]interface{}{
		"total_keys":     len(me.data),
		"memory_usage":   me.GetMemoryUsage(),
		"persistence":    me.persistencePath != "",
		"auto_save":      me.autoSave,
		"save_interval":  me.saveInterval.Seconds(),
	}
	
	// Count by data type
	typeCounts := make(map[string]int)
	for _, value := range me.data {
		switch value.Type {
		case core.STRING:
			typeCounts["string"]++
		case core.HASH:
			typeCounts["hash"]++
		case core.LIST:
			typeCounts["list"]++
		case core.SET:
			typeCounts["set"]++
		case core.ZSET:
			typeCounts["zset"]++
		}
	}
	
	stats["type_counts"] = typeCounts
	return stats
}
