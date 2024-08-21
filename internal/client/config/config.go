package config

import (
	"os"

	"github.com/spf13/pflag"
)

type Config struct {
	ServerAddress string
	EnableTLS     bool
}

const (
	DefaultServerAddress = "localhost:4080"
	DefaultTLS           = false
)

func NewConfig(args []string) (*Config, error) {
	cfg := defaultConfig()

	var (
		flagAddress   string
		flagEnableTLS bool
	)

	flags := pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)

	flags.StringVarP(&flagAddress, "address", "a", DefaultServerAddress, "Server address")
	flags.BoolVarP(&flagEnableTLS, "tls", "t", DefaultTLS, "Use TLS, default is false")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	return cfg, nil
}

func defaultConfig() *Config {
	return &Config{
		ServerAddress: DefaultServerAddress,
		EnableTLS:     DefaultTLS,
	}
}
