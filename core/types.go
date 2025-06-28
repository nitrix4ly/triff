package core

import (
	"sync"
	"time"
)

// DataType represents different data types supported by Triff
type DataType int

const (
	STRING DataType = iota
	HASH
	LIST
	SET
	ZSET
)

// TriffValue represents a value stored in the database
type TriffValue struct {
	Type      DataType    `json:"type"`
	Data      interface{} `json:"data"`
	TTL       int64       `json:"ttl"`       // Time to live in seconds
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// Database represents the main database structure
type Database struct {
	Data      map[string]*TriffValue `json:"data"`
	mu        sync.RWMutex
	config    *Config
	persistence PersistenceEngine
}

// Config holds database configuration
type Config struct {
	Port            int    `yaml:"port"`
	HTTPPort        int    `yaml:"http_port"`
	MaxMemory       int64  `yaml:"max_memory"`
	PersistencePath string `yaml:"persistence_path"`
	LogLevel        string `yaml:"log_level"`
	EnableHTTP      bool   `yaml:"enable_http"`
	EnableTCP       bool   `yaml:"enable_tcp"`
}

// StorageEngine defines interface for storage implementations
type StorageEngine interface {
	Get(key string) (*TriffValue, bool)
	Set(key string, value *TriffValue) error
	Delete(key string) bool
	Exists(key string) bool
	Keys(pattern string) []string
	FlushAll() error
	Size() int64
}

// PersistenceEngine defines interface for data persistence
type PersistenceEngine interface {
	Save(data map[string]*TriffValue) error
	Load() (map[string]*TriffValue, error)
	SetPath(path string)
}

// Command represents a database command
type Command struct {
	Name string
	Args []string
}

// Response represents a command response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Type    string      `json:"type"`
}
