package utils

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
	"github.com/nitrix4ly/triff/core"
)

// LoadConfig loads configuration from YAML file
func LoadConfig(filepath string) (*core.Config, error) {
	// Default configuration
	config := &core.Config{
		Port:            6379,
		HTTPPort:        8080,
		MaxMemory:       1024 * 1024 * 1024, // 1GB
		PersistencePath: "./triff.db",
		LogLevel:        "info",
		EnableHTTP:      true,
		EnableTCP:       true,
	}
	
	// If no config file specified, return default
	if filepath == "" {
		return config, nil
	}
	
	// Check if config file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return config, nil // Return default config if file doesn't exist
	}
	
	// Read config file
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	
	// Parse YAML
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}
	
	return config, nil
}

// SaveConfig saves configuration to YAML file
func SaveConfig(config *core.Config, filepath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	
	return os.WriteFile(filepath, data, 0644)
}

// GetEnvConfig gets configuration from environment variables
func GetEnvConfig() *core.Config {
	config := &core.Config{
		Port:            6379,
		HTTPPort:        8080,
		MaxMemory:       1024 * 1024 * 1024, // 1GB
		PersistencePath: "./triff.db",
		LogLevel:        "info",
		EnableHTTP:      true,
		EnableTCP:       true,
	}

	// Override with environment variables if they exist
	if port := os.Getenv("TRIFF_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Port = p
		}
	}

	if httpPort := os.Getenv("TRIFF_HTTP_PORT"); httpPort != "" {
		if p, err := strconv.Atoi(httpPort); err == nil {
			config.HTTPPort = p
		}
	}

	if maxMem := os.Getenv("TRIFF_MAX_MEMORY"); maxMem != "" {
		if m, err := strconv.ParseInt(maxMem, 10, 64); err == nil {
			config.MaxMemory = m
		}
	}

	if persistPath := os.Getenv("TRIFF_PERSISTENCE_PATH"); persistPath != "" {
		config.PersistencePath = persistPath
	}

	if logLevel := os.Getenv("TRIFF_LOG_LEVEL"); logLevel != "" {
		config.LogLevel = logLevel
	}

	if enableHTTP := os.Getenv("TRIFF_ENABLE_HTTP"); enableHTTP != "" {
		if b, err := strconv.ParseBool(enableHTTP); err == nil {
			config.EnableHTTP = b
		}
	}

	if enableTCP := os.Getenv("TRIFF_ENABLE_TCP"); enableTCP != "" {
		if b, err := strconv.ParseBool(enableTCP); err == nil {
			config.EnableTCP = b
		}
	}

	return config
}

// MergeConfigs merges multiple config sources with priority: env > file > default
func MergeConfigs(filepath string) (*core.Config, error) {
	// Start with file config (which includes defaults)
	config, err := LoadConfig(filepath)
	if err != nil {
		return nil, err
	}

	// Override with environment variables
	envConfig := GetEnvConfig()
	
	// Only override non-default values from env
	if os.Getenv("TRIFF_PORT") != "" {
		config.Port = envConfig.Port
	}
	if os.Getenv("TRIFF_HTTP_PORT") != "" {
		config.HTTPPort = envConfig.HTTPPort
	}
	if os.Getenv("TRIFF_MAX_MEMORY") != "" {
		config.MaxMemory = envConfig.MaxMemory
	}
	if os.Getenv("TRIFF_PERSISTENCE_PATH") != "" {
		config.PersistencePath = envConfig.PersistencePath
	}
	if os.Getenv("TRIFF_LOG_LEVEL") != "" {
		config.LogLevel = envConfig.LogLevel
	}
	if os.Getenv("TRIFF_ENABLE_HTTP") != "" {
		config.EnableHTTP = envConfig.EnableHTTP
	}
	if os.Getenv("TRIFF_ENABLE_TCP") != "" {
		config.EnableTCP = envConfig.EnableTCP
	}

	return config, nil
}

// ValidateConfig validates the configuration values
func ValidateConfig(config *core.Config) error {
	if config.Port < 1 || config.Port > 65535 {
		return fmt.Errorf("invalid port: %d (must be between 1-65535)", config.Port)
	}
	
	if config.HTTPPort < 1 || config.HTTPPort > 65535 {
		return fmt.Errorf("invalid HTTP port: %d (must be between 1-65535)", config.HTTPPort)
	}
	
	if config.MaxMemory < 1024*1024 { // Minimum 1MB
		return fmt.Errorf("max memory too small: %d (minimum 1MB)", config.MaxMemory)
	}
	
	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLogLevels[config.LogLevel] {
		return fmt.Errorf("invalid log level: %s (must be debug, info, warn, or error)", config.LogLevel)
	}
	
	if !config.EnableHTTP && !config.EnableTCP {
		return fmt.Errorf("at least one protocol (HTTP or TCP) must be enabled")
	}
	
	return nil
}
