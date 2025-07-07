package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds the application configuration
type Config struct {
	Server  ServerConfig
	Storage StorageConfig
	Engine  EngineConfig
	Logging LoggingConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// StorageConfig holds storage configuration
type StorageConfig struct {
	DataDir    string
	PluginsDir string
}

// EngineConfig holds flow engine configuration
type EngineConfig struct {
	MaxConcurrentFlows int
	DefaultTimeout     time.Duration
	DebugMode          bool
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string
	Format string
}

// Load loads configuration from environment variables with defaults
func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Address:      getEnv("SERVER_ADDRESS", ":8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Storage: StorageConfig{
			DataDir:    getEnv("DATA_DIR", "./data"),
			PluginsDir: getEnv("PLUGINS_DIR", "./data/plugins"),
		},
		Engine: EngineConfig{
			MaxConcurrentFlows: getIntEnv("MAX_CONCURRENT_FLOWS", 10),
			DefaultTimeout:     getDurationEnv("DEFAULT_TIMEOUT", 30*time.Second),
			DebugMode:          getBoolEnv("DEBUG_MODE", true),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}, nil
}

// Helper functions to get environment variables with defaults

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
