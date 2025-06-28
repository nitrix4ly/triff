package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nitrix4ly/triff/core"
)

// StringCommands handles all string-related operations
type StringCommands struct {
	db *core.Database
}

// NewStringCommands creates a new string commands handler
func NewStringCommands(db *core.Database) *StringCommands {
	return &StringCommands{db: db}
}

// Set stores a string value
func (sc *StringCommands) Set(key, value string, ttl int64) *core.Response {
	triffValue := &core.TriffValue{
		Type: core.STRING,
		Data: value,
		TTL:  ttl,
	}
	
	if ttl > 0 {
		triffValue.TTL = time.Now().Unix() + ttl
	}
	
	err := sc.db.Set(key, triffValue)
	if err != nil {
		return &core.Response{
			Success: false,
			Error:   err.Error(),
			Type:    "string",
		}
	}
	
	return &core.Response{
		Success: true,
		Data:    "OK",
		Type:    "string",
	}
}

// Get retrieves a string value
func (sc *StringCommands) Get(key string) *core.Response {
	value, exists := sc.db.Get(key)
	if !exists {
		return &core.Response{
			Success: false,
			Data:    nil,
			Type:    "string",
		}
	}
	
	if value.Type != core.STRING {
		return &core.Response{
			Success: false,
			Error:   "value is not a string",
			Type:    "string",
		}
	}
	
	return &core.Response{
		Success: true,
		Data:    value.Data,
		Type:    "string",
	}
}

// Append appends a value to an existing string
func (sc *StringCommands) Append(key, value string) *core.Response {
	existing, exists := sc.db.Get(key)
	var newValue string
	
	if exists && existing.Type == core.STRING {
		newValue = existing.Data.(string) + value
	} else {
		newValue = value
	}
	
	triffValue := &core.TriffValue{
		Type: core.STRING,
		Data: newValue,
	}
	
	sc.db.Set(key, triffValue)
	
	return &core.Response{
		Success: true,
		Data:    len(newValue),
		Type:    "integer",
	}
}

// Strlen returns the length of a string
func (sc *StringCommands) Strlen(key string) *core.Response {
	value, exists := sc.db.Get(key)
	if !exists {
		return &core.Response{
			Success: true,
			Data:    0,
			Type:    "integer",
		}
	}
	
	if value.Type != core.STRING {
		return &core.Response{
			Success: false,
			Error:   "value is not a string",
			Type:    "string",
		}
	}
	
	length := len(value.Data.(string))
	return &core.Response{
		Success: true,
		Data:    length,
		Type:    "integer",
	}
}

// Incr increments a numeric string value
func (sc *StringCommands) Incr(key string) *core.Response {
	return sc.IncrBy(key, 1)
}

// IncrBy increments a numeric string value by a specific amount
func (sc *StringCommands) IncrBy(key string, increment int64) *core.Response {
	value, exists := sc.db.Get(key)
	var currentValue int64 = 0
	
	if exists {
		if value.Type != core.STRING {
			return &core.Response{
				Success: false,
				Error:   "value is not a string",
				Type:    "string",
			}
		}
		
		var err error
		currentValue, err = strconv.ParseInt(value.Data.(string), 10, 64)
		if err != nil {
			return &core.Response{
				Success: false,
				Error:   "value is not a valid integer",
				Type:    "string",
			}
		}
	}
	
	newValue := currentValue + increment
	triffValue := &core.TriffValue{
		Type: core.STRING,
		Data: fmt.Sprintf("%d", newValue),
	}
	
	sc.db.Set(key, triffValue)
	
	return &core.Response{
		Success: true,
		Data:    newValue,
		Type:    "integer",
	}
}

// Decr decrements a numeric string value
func (sc *StringCommands) Decr(key string) *core.Response {
	return sc.IncrBy(key, -1)
}

// MGet gets multiple string values
func (sc *StringCommands) MGet(keys []string) *core.Response {
	results := make([]interface{}, len(keys))
	
	for i, key := range keys {
		value, exists := sc.db.Get(key)
		if !exists || value.Type != core.STRING {
			results[i] = nil
		} else {
			results[i] = value.Data
		}
	}
	
	return &core.Response{
		Success: true,
		Data:    results,
		Type:    "array",
	}
}

// MSet sets multiple string values
func (sc *StringCommands) MSet(keyValues map[string]string) *core.Response {
	for key, value := range keyValues {
		triffValue := &core.TriffValue{
			Type: core.STRING,
			Data: value,
		}
		sc.db.Set(key, triffValue)
	}
	
	return &core.Response{
		Success: true,
		Data:    "OK",
		Type:    "string",
	}
}

// GetRange returns a substring of a string value
func (sc *StringCommands) GetRange(key string, start, end int) *core.Response {
	value, exists := sc.db.Get(key)
	if !exists {
		return &core.Response{
			Success: true,
			Data:    "",
			Type:    "string",
		}
	}
	
	if value.Type != core.STRING {
		return &core.Response{
			Success: false,
			Error:   "value is not a string",
			Type:    "string",
		}
	}
	
	str := value.Data.(string)
	length := len(str)
	
	// Handle negative indices
	if start < 0 {
		start = length + start
	}
	if end < 0 {
		end = length + end
	}
	
	// Bounds checking
	if start < 0 {
		start = 0
	}
	if end >= length {
		end = length - 1
	}
	if start > end {
		return &core.Response{
			Success: true,
			Data:    "",
			Type:    "string",
		}
	}
	
	result := str[start : end+1]
	return &core.Response{
		Success: true,
		Data:    result,
		Type:    "string",
	}
}
