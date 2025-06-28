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
	Client Client
	Server Server
}

type Client struct {
	Address          string
	UseTLS           bool // server-side TLS
	ServerName       string
	CertFile         string
	KeyFile          string
	CaFile           string
	KeepAlive        time.Duration
	KeepAliveTimeout time.Duration
	EnableRoundRobin bool
}

func (c Client) GetCertFile() string   { return c.CertFile }
func (c Client) GetKeyFile() string    { return c.KeyFile }
func (c Client) GetCaFile() string     { return c.CaFile }
func (c Client) GetServerName() string { return c.ServerName }

type Server struct {
	Address        string
	UseTLS         bool // server-side TLS
	ServerName     string
	CertFile       string
	KeyFile        string
	CaFile         string
	ClientCertAuth bool // require client cert for mTLS

	Metrics
}

func (s Server) GetCertFile() string   { return s.CertFile }
func (s Server) GetKeyFile() string    { return s.KeyFile }
func (s Server) GetCaFile() string     { return s.CaFile }
func (s Server) GetServerName() string { return s.ServerName }

type Metrics struct {
	Enabled bool   // enable metrics
	Address string // address for metrics endpoint
	Path    string // path for metrics endpoint
	Auth    struct {
		Username string // username for basic auth
		Password string // password for basic auth
	}
}

func Init(opts ...Option) (*Config, error) {

	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}

	if options.Prefix != "" {
		viper.SetEnvPrefix(options.Prefix)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	cfg := defaultConfig()

	cfg.Server.UseTLS = viper.GetBool("server.use_tls")
	cfg.Server.ClientCertAuth = viper.GetBool("server.client_cert_auth")
	if s := viper.GetString("server.address"); s != "" {
		cfg.Server.Address = s
	}
	if s := viper.GetString("server.server_name"); s != "" {
		cfg.Server.ServerName = s
	}
	if s := viper.GetString("server.cert_file"); s != "" {
		cfg.Server.CertFile = s
	}
	if s := viper.GetString("server.key_file"); s != "" {
		cfg.Server.KeyFile = s
	}
	if s := viper.GetString("server.ca_file"); s != "" {
		cfg.Server.CaFile = s
	}

	cfg.Server.Metrics.Enabled = viper.GetBool("server.metrics.enabled")
	if s := viper.GetString("server.metrics.address"); s != "" {
		cfg.Server.Metrics.Address = s
	}
	if s := viper.GetString("server.metrics.path"); s != "" {
		cfg.Server.Metrics.Path = s
	}
	if s := viper.GetString("server.metrics.auth.username"); s != "" {
		cfg.Server.Metrics.Auth.Username = s
	}
	if s := viper.GetString("server.metrics.auth.password"); s != "" {
		cfg.Server.Metrics.Auth.Password = s
	}

	cfg.Client.UseTLS = viper.GetBool("client.use_tls")
	if s := viper.GetString("client.address"); s != "" {
		cfg.Client.Address = s
	}

	if s := viper.GetString("client.server_name"); s != "" {
		cfg.Client.ServerName = s
	}

	if s := viper.GetString("client.cert_file"); s != "" {
		cfg.Client.CertFile = s
	}

	if s := viper.GetString("client.key_file"); s != "" {
		cfg.Client.KeyFile = s
	}

	if s := viper.GetString("client.ca_file"); s != "" {
		cfg.Client.CaFile = s
	}

	if d := viper.GetDuration("client.keep_alive"); d > 0 {
		cfg.Client.KeepAlive = d
	}

	if d := viper.GetDuration("client.keep_alive_timeout"); d > 0 {
		cfg.Client.KeepAliveTimeout = d
	}

	cfg.Client.EnableRoundRobin = viper.GetBool("client.enable_round_robin")

	return cfg, nil
}

func defaultConfig() *Config {
	return &Config{
		Client: Client{
			Address:          "localhost:9091",
			UseTLS:           false,
			ServerName:       "localhost",
			CertFile:         "",
			KeyFile:          "",
			CaFile:           "",
			KeepAlive:        10 * time.Second,
			KeepAliveTimeout: 5 * time.Second,
		},
		Server: Server{
			Address:        ":9091",
			UseTLS:         false,
			ServerName:     "localhost",
			CertFile:       "",
			KeyFile:        "",
			CaFile:         "",
			ClientCertAuth: false,
			Metrics: Metrics{
				Enabled: false,
				Address: ":9092",
				Path:    "/metrics",
			},
		},
	}
}
