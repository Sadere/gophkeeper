package config

import (
	"errors"
	"os"

	"github.com/spf13/pflag"
)

var (
	DefaultAddress = "localhost:4080"
	DefaultTLS     = false
)

type Config struct {
	Address     string
	EnableTLS   bool
	PostgresDSN string
	LogLevel    string
	SecretKey   string // Secret key for jwt token signing
}

// Parses config from provided command line arguments or from environment variables and returns config
func NewConfig(args []string) (*Config, error) {
	cfg := defaultConfig()

	var (
		flagAddress     string
		flagEnableTLS   bool
		flagPostgresDSN string
	)

	flags := pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)

	flags.StringVarP(&cfg.LogLevel, "verbosity", "v", "fatal", "Log level, possible levels: debug, info, warn, error, dpanic, panic, fatal")
	flags.StringVarP(&flagAddress, "address", "a", DefaultAddress, "Address to listen on")
	flags.BoolVarP(&flagEnableTLS, "tls", "t", DefaultTLS, "Use TLS, default is false")
	flags.StringVarP(&flagPostgresDSN, "dsn", "d", "", "DSN string for PostgreSQL connection")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	if envAddr := os.Getenv("ADDRESS"); len(envAddr) > 0 {
		cfg.Address = envAddr
	} else {
		cfg.Address = flagAddress
	}

	if envTLS := os.Getenv("ENABLE_TLS"); envTLS == "true" {
		cfg.EnableTLS = true
	} else {
		cfg.EnableTLS = flagEnableTLS
	}

	if envDSN := os.Getenv("DATABASE_DSN"); len(envDSN) > 0 {
		cfg.PostgresDSN = envDSN
	} else {
		cfg.PostgresDSN = flagPostgresDSN
	}

	if envSecretKey := os.Getenv("SECRET_KEY"); len(envSecretKey) > 0 {
		cfg.SecretKey = envSecretKey
	} else {
		return nil, errors.New("SECRET_KEY env variable is not provided")
	}

	return cfg, nil
}

func defaultConfig() *Config {
	return &Config{
		Address:   DefaultAddress,
		EnableTLS: DefaultTLS,
	}
}
