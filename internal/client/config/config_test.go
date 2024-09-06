package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientConfig(t *testing.T) {
	var tests = []struct {
		name        string
		args        []string
		prepareWant func(c *Config)
	}{
		{
			name: "address",
			args: []string{"-a", "localhost:1111"},
			prepareWant: func(c *Config) {
				c.ServerAddress = "localhost:1111"
			},
		},
		{
			name: "enable tls",
			args: []string{"-t", "true"},
			prepareWant: func(c *Config) {
				c.EnableTLS = true
			},
		},
		{
			name: "download path",
			args: []string{"-d", "/test/"},
			prepareWant: func(c *Config) {
				c.DownloadPath = "/test/"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf, err := NewConfig(tt.args)

			assert.NoError(t, err)

			wantCfg := defaultConfig()
			tt.prepareWant(wantCfg)

			if !reflect.DeepEqual(conf, wantCfg) {
				t.Errorf("conf got %+v, want %+v", conf, wantCfg)
			}
		})
	}
}
