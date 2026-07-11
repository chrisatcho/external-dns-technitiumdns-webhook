package configuration

import (
	"time"

	"github.com/caarlos0/env/v8"
	log "github.com/sirupsen/logrus"
)

// Config struct for configuration environmental variables
type Config struct {
	ServerHost           string        `env:"SERVER_HOST" envDefault:"localhost"`
	ServerPort           int           `env:"SERVER_PORT" envDefault:"8888"`
	MetricsPort          int           `env:"METRICS_PORT" envDefault:"8080"`
	MetricsServer        bool          `env:"METRICS_SERVER" envDefault:"false"`
	ServerReadTimeout    time.Duration `env:"SERVER_READ_TIMEOUT"`
	ServerWriteTimeout   time.Duration `env:"SERVER_WRITE_TIMEOUT"`
	DomainFilter         []string      `env:"DOMAIN_FILTER" envDefault:""`
	ExcludeDomains       []string      `env:"EXCLUDE_DOMAIN_FILTER" envDefault:""`
	RegexDomainFilter    string        `env:"REGEXP_DOMAIN_FILTER" envDefault:""`
	RegexDomainExclusion string        `env:"REGEXP_DOMAIN_FILTER_EXCLUSION" envDefault:""`
}

// Init sets up configuration by reading set environmental variables
func Init() Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("error reading configuration from environment: %v", err)
	}

	// Validate configuration
	if cfg.ServerPort < 1 || cfg.ServerPort > 65535 {
		log.Fatalf("invalid SERVER_PORT: %d (must be between 1 and 65535)", cfg.ServerPort)
	}
	if cfg.MetricsPort < 1 || cfg.MetricsPort > 65535 {
		log.Fatalf("invalid METRICS_PORT: %d (must be between 1 and 65535)", cfg.MetricsPort)
	}

	log.Infof("configuration loaded: host=%s, port=%d, metricsPort=%d",
		cfg.ServerHost, cfg.ServerPort, cfg.MetricsPort)

	return cfg
}
