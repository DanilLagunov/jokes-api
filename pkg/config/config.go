package config

import (
	"time"

	"github.com/caarlos0/env/v6"
)

// Config struct.
type Config struct {
	Port                   int           `env:"PORT" envDefault:"8000"`
	ReadHeaderTimeout      time.Duration `env:"READ_HEADER_TIMEOUT"`
	ReadTimeout            time.Duration `env:"READ_TIMEOUT"`
	WriteTimeout           time.Duration `env:"WRITE_TIMEOUT"`
	DbURI                  string        `env:"DB_URI"`
	DbName                 string        `env:"DB_NAME"`
	JokesCollection        string        `env:"JOKES_COLLECTION"`
	CacheDefaultExpiration time.Duration `env:"DEFAULT_EXPIRATION"`
	CacheCleanupInterval   time.Duration `env:"CLEANUP_INTERVAL"`
}

// NewConfig creating a new Config object.
func NewConfig() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
