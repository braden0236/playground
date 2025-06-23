package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Options struct {
	Prefix string
}

type Option func(*Options)

func WithPrefix(prefix string) Option {
	return func(o *Options) {
		o.Prefix = prefix
	}
}

type Config struct {
	Server  ServerConfig
	Metrics MetricsConfig
}

type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type MetricsConfig struct {
	Username string
	Password string
}

const (
	DefaultServerPort          = 8080
	DefaultReadTimeoutSeconds  = 15
	DefaultWriteTimeoutSeconds = 20
	DefaultIdleTimeoutSeconds  = 60
)

var (
	DefaultReadTimeout  = time.Duration(DefaultReadTimeoutSeconds) * time.Second
	DefaultWriteTimeout = time.Duration(DefaultWriteTimeoutSeconds) * time.Second
	DefaultIdleTimeout  = time.Duration(DefaultIdleTimeoutSeconds) * time.Second
)

func Init(opts ...Option) (*Config, error) {

	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}

	if options.Prefix != "" {
		viper.SetEnvPrefix(options.Prefix)
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	}

	viper.AutomaticEnv()

	cfg := defaultConfig()
	if port := viper.GetInt("SERVER_PORT"); port >= 1 && port <= 65535 {
		cfg.Server.Port = port
	}
	if t := viper.GetDuration("SERVER_READ_TIMEOUT"); t > 0 {
		cfg.Server.ReadTimeout = t
	}
	if t := viper.GetDuration("SERVER_WRITE_TIMEOUT"); t > 0 {
		cfg.Server.WriteTimeout = t
	}
	if t := viper.GetDuration("SERVER_IDLE_TIMEOUT"); t > 0 {
		cfg.Server.IdleTimeout = t
	}

	if u := viper.GetString("METRICS_USERNAME"); u != "" {
		cfg.Metrics.Username = u
	}
	if p := viper.GetString("METRICS_PASSWORD"); p != "" {
		cfg.Metrics.Password = p
	}

	return cfg, nil
}

func defaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         DefaultServerPort,
			ReadTimeout:  DefaultReadTimeout,
			WriteTimeout: DefaultWriteTimeout,
			IdleTimeout:  DefaultIdleTimeout,
		},
		Metrics: MetricsConfig{
			Username: "",
			Password: "",
		},
	}
}
