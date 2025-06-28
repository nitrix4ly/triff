package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nitrix4ly/triff/commands"
	"github.com/nitrix4ly/triff/core"
	"github.com/nitrix4ly/triff/server"
	"github.com/nitrix4ly/triff/utils"
)

func main() {
	// Create database configuration
	config := &core.Config{
		Port:            6379,
		HTTPPort:        8080,
		MaxMemory:       1024 * 1024 * 1024, // 1GB
		PersistencePath: "./triff.db",
		LogLevel:        "info",
		EnableHTTP:      true,
		EnableTCP:       true,
	}

	// Initialize logger
	logger := utils.NewLogger(config.LogLevel)

	// Create new database instance
	db := core.NewDatabase(config)

	// Create string commands handler
	stringCmd := commands.NewStringCommands(db)

	// Basic usage examples
	fmt.Println("=== Triff Database Basic Usage Examples ===")

	// Set some values
	response := stringCmd.Set("user:1:name", "John Doe", 0)
	if response.Success {
		fmt.Println("✓ Set user:1:name = John Doe")
	}

	response = stringCmd.Set("user:1:age", "25", 0)
	if response.Success {
		fmt.Println("✓ Set user:1:age = 25")
	}

	// Set with TTL (expires in 60 seconds)
	response = stringCmd.Set("session:abc123", "active", 60)
	if response.Success {
		fmt.Println("✓ Set session:abc123 = active (expires in 60s)")
	}

	// Get values
	response = stringCmd.Get("user:1:name")
	if response.Success {
		fmt.Printf("✓ Get user:1:name = %s\n", response.Data)
	}

	response = stringCmd.Get("user:1:age")
	if response.Success {
		fmt.Printf("✓ Get user:1:age = %s\n", response.Data)
	}

	// Increment operations
	response = stringCmd.Incr("counter")
	if response.Success {
		fmt.Printf("✓ Increment counter = %d\n", response.Data)
	}

	response = stringCmd.Incr("counter")
	if response.Success {
		fmt.Printf("✓ Increment counter = %d\n", response.Data)
	}

	// String operations
	response = stringCmd.Append("user:1:name", " Smith")
	if response.Success {
		fmt.Printf("✓ Append to user:1:name, new length = %d\n", response.Data)
	}

	response = stringCmd.Get("user:1:name")
	if response.Success {
		fmt.Printf("✓ Get user:1:name = %s\n", response.Data)
	}

	// Multiple operations
	keyValues := map[string]string{
		"config:host": "localhost",
		"config:port": "3000",
		"config:env":  "production",
	}

	response = stringCmd.MSet(keyValues)
	if response.Success {
		fmt.Println("✓ Set multiple config values")
	}

	// Get multiple values
	keys := []string{"config:host", "config:port", "config:env"}
	response = stringCmd.MGet(keys)
	if response.Success {
		values := response.Data.([]interface{})
		for i, key := range keys {
			fmt.Printf("✓ %s = %v\n", key, values[i])
		}
	}

	// Database operations
	fmt.Printf("✓ Database size: %d keys\n", db.Size())
	fmt.Printf("✓ Database info: %v\n", db.Info())

	// Check if we want to start servers
	fmt.Println("\n=== Starting Servers ===")
	
	// Start TCP server in a goroutine
	if config.EnableTCP {
		tcpServer := server.NewTCPServer(db, config.Port, logger)
		go func() {
			if err := tcpServer.Start(); err != nil {
				log.Printf("TCP server error: %v", err)
			}
		}()
		fmt.Printf("✓ TCP server started on port %d\n", config.Port)
	}

	// Start HTTP server in a goroutine
	if config.EnableHTTP {
		httpServer := server.NewHTTPServer(db, config.HTTPPort, logger)
		go func() {
			if err := httpServer.Start(); err != nil {
				log.Printf("HTTP server error: %v", err)
			}
		}()
		fmt.Printf("✓ HTTP server started on port %d\n", config.HTTPPort)
	}

	// Keep the program running
	fmt.Println("\n✓ Triff database is running!")
	fmt.Println("Press Ctrl+C to stop...")
	
	// Cleanup expired keys every 10 seconds
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			db.CleanupExpired()
			logger.Debug("Cleaned up expired keys")
		}
	}
}
