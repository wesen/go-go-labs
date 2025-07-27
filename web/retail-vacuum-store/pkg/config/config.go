package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Logging  LoggingConfig  `yaml:"logging"`
	Session  SessionConfig  `yaml:"session"`
	Storage  StorageConfig  `yaml:"storage"`
}

// ServerConfig contains server-related configuration
type ServerConfig struct {
	Port         int    `yaml:"port"`
	Host         string `yaml:"host"`
	LogLevel     string `yaml:"log_level"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
	IdleTimeout  int    `yaml:"idle_timeout"`
}

// DatabaseConfig contains database connection configuration
type DatabaseConfig struct {
	Host         string `mapstructure:"host" yaml:"host"`
	Port         int    `mapstructure:"port" yaml:"port"`
	User         string `mapstructure:"user" yaml:"user"`
	Password     string `mapstructure:"password" yaml:"password"`
	Name         string `mapstructure:"name" yaml:"name"`
	SSLMode      string `mapstructure:"ssl_mode" yaml:"ssl_mode"`
	MaxOpenConns int    `mapstructure:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`
	MaxLifetime  string `mapstructure:"max_lifetime" yaml:"max_lifetime"`
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

// SessionConfig contains session management configuration
type SessionConfig struct {
	Secret   string `yaml:"secret"`
	Duration int    `yaml:"duration"`
	Secure   bool   `yaml:"secure"`
}

// StorageConfig contains file storage configuration
type StorageConfig struct {
	Type       string `yaml:"type"`
	LocalPath  string `yaml:"local_path"`
	S3Bucket   string `yaml:"s3_bucket"`
	S3Region   string `yaml:"s3_region"`
	S3Endpoint string `yaml:"s3_endpoint"`
}

// Load loads configuration with default values
func Load() (*Config, error) {
	configPath := "config.yaml"
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		configPath = path
	}
	return LoadFromFile(configPath)
}

// LoadFromFile loads configuration from a specific file
func LoadFromFile(configPath string) (*Config, error) {
	// Set defaults
	config := &Config{
		Server: ServerConfig{
			Port:         8080,
			Host:         "0.0.0.0",
			LogLevel:     "info",
			ReadTimeout:  30,
			WriteTimeout: 30,
			IdleTimeout:  120,
		},
		Database: DatabaseConfig{
			Host:         "localhost",
			Port:         5432,
			User:         "postgres",
			Password:     "",
			Name:         "vacuum_store",
			SSLMode:      "disable",
			MaxOpenConns: 25,
			MaxIdleConns: 10,
			MaxLifetime:  "1h",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		},
		Session: SessionConfig{
			Secret:   "change-this-in-production",
			Duration: 86400 * 30, // 30 days
			Secure:   false,
		},
		Storage: StorageConfig{
			Type:      "local",
			LocalPath: "./uploads",
		},
	}

	// Load from file if it exists
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Override with environment variables
	if port := os.Getenv("PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Server.Port = p
		}
	}
	
	if host := os.Getenv("HOST"); host != "" {
		config.Server.Host = host
	}

	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.Database.Host = dbHost
	}

	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		if p, err := strconv.Atoi(dbPort); err == nil {
			config.Database.Port = p
		}
	}

	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		config.Database.User = dbUser
	}

	if dbPass := os.Getenv("DB_PASSWORD"); dbPass != "" {
		config.Database.Password = dbPass
	}

	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		config.Database.Name = dbName
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.Logging.Level = logLevel
	}

	if sessionSecret := os.Getenv("SESSION_SECRET"); sessionSecret != "" {
		config.Session.Secret = sessionSecret
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535")
	}

	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if c.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}

	if c.Session.Secret == "" || c.Session.Secret == "change-this-in-production" {
		return fmt.Errorf("session secret must be set and changed from default")
	}

	return nil
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Logging.Level == "debug" || os.Getenv("ENV") == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return os.Getenv("ENV") == "production"
}
