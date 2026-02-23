package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Environment string         `env:"APP_ENV"     envDefault:"development"`
	Server      ServerConfig   `envPrefix:"SERVER_"`
	Database    DatabaseConfig `envPrefix:"DATABASE_"`
	Auth        AuthConfig     `envPrefix:"AUTH_"`
	Logging     LoggingConfig  `envPrefix:"LOGGING_"`
	Defaults    DefaultsConfig `envPrefix:"DEFAULTS_"`
	CORS        CORSConfig     `envPrefix:"CORS_"`
}

type ServerConfig struct {
	Port         int           `env:"PORT"          envDefault:"8081"`
	Host         string        `env:"HOST"          envDefault:"0.0.0.0"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT"  envDefault:"15s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"15s"`
	SwaggerHost  string        `env:"SWAGGER_HOST"`
}

type DatabaseConfig struct {
	Postgres PostgresConfig `envPrefix:"POSTGRES_"`
}

type PostgresConfig struct {
	Host         string `env:"HOST"           envDefault:"localhost"`
	Port         int    `env:"PORT"           envDefault:"5432"`
	Database     string `env:"DATABASE"       envDefault:"edugo"`
	User         string `env:"USER"           envDefault:"edugo"`
	Password     string `env:"PASSWORD,required"`
	MaxOpenConns int    `env:"MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns int    `env:"MAX_IDLE_CONNS" envDefault:"10"`
	SSLMode      string `env:"SSL_MODE"       envDefault:"disable"`
}

type LoggingConfig struct {
	Level  string `env:"LEVEL"  envDefault:"info"`
	Format string `env:"FORMAT" envDefault:"json"`
}

type AuthConfig struct {
	JWT            JWTConfig            `envPrefix:"JWT_"`
	APIIamPlatform APIIamPlatformConfig `envPrefix:"API_IAM_PLATFORM_"`
}

type JWTConfig struct {
	Secret string `env:"SECRET,required"`
	Issuer string `env:"ISSUER" envDefault:"edugo-central"`
}

type APIIamPlatformConfig struct {
	BaseURL         string        `env:"BASE_URL"         envDefault:"http://localhost:8070/api"`
	Timeout         time.Duration `env:"TIMEOUT"          envDefault:"5s"`
	CacheTTL        time.Duration `env:"CACHE_TTL"        envDefault:"60s"`
	CacheEnabled    bool          `env:"CACHE_ENABLED"    envDefault:"true"`
	RemoteEnabled   bool          `env:"REMOTE_ENABLED"   envDefault:"false"`
	FallbackEnabled bool          `env:"FALLBACK_ENABLED" envDefault:"false"`
}

type DefaultsConfig struct {
	School SchoolDefaults `envPrefix:"SCHOOL_"`
}

type SchoolDefaults struct {
	Country          string `env:"COUNTRY"           envDefault:"CO"`
	SubscriptionTier string `env:"SUBSCRIPTION_TIER" envDefault:"free"`
	MaxTeachers      int    `env:"MAX_TEACHERS"      envDefault:"50"`
	MaxStudents      int    `env:"MAX_STUDENTS"      envDefault:"500"`
}

type CORSConfig struct {
	AllowedOrigins string `env:"ALLOWED_ORIGINS" envDefault:"*"`
	AllowedMethods string `env:"ALLOWED_METHODS" envDefault:"GET,POST,PUT,PATCH,DELETE,OPTIONS"`
	AllowedHeaders string `env:"ALLOWED_HEADERS" envDefault:"Origin,Content-Type,Accept,Authorization,X-Request-ID"`
}

// DSN returns the PostgreSQL connection string
func (c *PostgresConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode)
}

// Load parses configuration from environment variables
func Load() (*Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("error parsing config from environment: %w", err)
	}
	return &cfg, nil
}
