package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config represents the full application configuration
type Config struct {
	Environment string         `mapstructure:"environment"`
	Server      ServerConfig   `mapstructure:"server"`
	Database    DatabaseConfig `mapstructure:"database"`
	Logging     LoggingConfig  `mapstructure:"logging"`
	Auth        AuthConfig     `mapstructure:"auth"`
	Defaults    DefaultsConfig `mapstructure:"defaults"`
	CORS        CORSConfig     `mapstructure:"cors"`
}

// ServerConfig configures the HTTP server
type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Host         string        `mapstructure:"host"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// DatabaseConfig configures the database connections
type DatabaseConfig struct {
	Postgres PostgresConfig `mapstructure:"postgres"`
}

// PostgresConfig configures PostgreSQL connection
type PostgresConfig struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	Database       string `mapstructure:"database"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	MaxConnections int    `mapstructure:"max_connections"`
	SSLMode        string `mapstructure:"ssl_mode"`
}

// LoggingConfig configures logging
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// AuthConfig configures authentication
type AuthConfig struct {
	JWT              JWTConfig              `mapstructure:"jwt"`
	InternalServices InternalServicesConfig `mapstructure:"internal_services"`
}

// JWTConfig configures JWT tokens
type JWTConfig struct {
	Secret               string        `mapstructure:"secret"`
	Issuer               string        `mapstructure:"issuer"`
	AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
}

// InternalServicesConfig configures internal service authentication
type InternalServicesConfig struct {
	APIKeys  string `mapstructure:"api_keys"`
	IPRanges string `mapstructure:"ip_ranges"`
}

// DefaultsConfig contains default values
type DefaultsConfig struct {
	School SchoolDefaults `mapstructure:"school"`
}

// SchoolDefaults contains default values for schools
type SchoolDefaults struct {
	Country          string `mapstructure:"country"`
	SubscriptionTier string `mapstructure:"subscription_tier"`
	MaxTeachers      int    `mapstructure:"max_teachers"`
	MaxStudents      int    `mapstructure:"max_students"`
}

// CORSConfig configures CORS
type CORSConfig struct {
	AllowedOrigins string `mapstructure:"allowed_origins"`
	AllowedMethods string `mapstructure:"allowed_methods"`
	AllowedHeaders string `mapstructure:"allowed_headers"`
}

// GetConnectionString returns the PostgreSQL connection string
func (c *PostgresConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode)
}

// Load reads configuration from file and environment variables
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// Environment variable overrides
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("environment", "local")
	viper.SetDefault("server.port", 8081)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.read_timeout", "15s")
	viper.SetDefault("server.write_timeout", "15s")
	viper.SetDefault("database.postgres.host", "localhost")
	viper.SetDefault("database.postgres.port", 5432)
	viper.SetDefault("database.postgres.database", "edugo")
	viper.SetDefault("database.postgres.user", "edugo")
	viper.SetDefault("database.postgres.ssl_mode", "disable")
	viper.SetDefault("database.postgres.max_connections", 25)
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("auth.jwt.issuer", "edugo-central")
	viper.SetDefault("auth.jwt.access_token_duration", "15m")
	viper.SetDefault("auth.jwt.refresh_token_duration", "168h")
	viper.SetDefault("defaults.school.country", "CO")
	viper.SetDefault("defaults.school.subscription_tier", "free")
	viper.SetDefault("defaults.school.max_teachers", 50)
	viper.SetDefault("defaults.school.max_students", 500)
	viper.SetDefault("cors.allowed_origins", "*")
	viper.SetDefault("cors.allowed_methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	viper.SetDefault("cors.allowed_headers", "Origin,Content-Type,Accept,Authorization,X-Request-ID")

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return &cfg, nil
}
