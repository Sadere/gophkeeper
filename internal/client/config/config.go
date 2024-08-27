package config

import (
	"os"

	"github.com/spf13/pflag"
)

type Config struct {
	ServerAddress string
	DownloadPath  string
	EnableTLS     bool
}

const (
	DefaultServerAddress = "localhost:8080"
	DefaultTLS           = false
	DefaultDownloadPath  = "./download/"
)

func NewConfig(args []string) (*Config, error) {
	cfg := defaultConfig()

	var (
		flagAddress      string
		flagEnableTLS    bool
		flagDownloadPath string
	)

	flags := pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)

	flags.StringVarP(&flagAddress, "address", "a", DefaultServerAddress, "Server address")
	flags.BoolVarP(&flagEnableTLS, "tls", "t", DefaultTLS, "Use TLS, default is false")
	flags.StringVarP(&flagDownloadPath, "dir", "d", DefaultDownloadPath, "Path for downloaded files")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	return cfg, nil
}

func defaultConfig() *Config {
	return &Config{
		ServerAddress: DefaultServerAddress,
		DownloadPath:  DefaultDownloadPath,
		EnableTLS:     DefaultTLS,
	}
}
