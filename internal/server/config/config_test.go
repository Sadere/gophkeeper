package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerConfig(t *testing.T) {
	var tests = []struct {
		name        string
		args        []string
		env         map[string]string
		prepareWant func(c *Config)
		wantErr     bool
	}{
		{
			name: "address from arg",
			args: []string{"-a", "localhost:1111"},
			env: map[string]string{
				"SECRET_KEY": "test",
			},
			prepareWant: func(c *Config) {
				c.SecretKey = "test"
				c.Address = "localhost:1111"
			},
			wantErr: false,
		},
		{
			name: "address from env",
			args: []string{"-a", "localhost:1111"},
			env: map[string]string{
				"SECRET_KEY": "test",
				"ADDRESS":    "test:2222",
			},
			prepareWant: func(c *Config) {
				c.SecretKey = "test"
				c.Address = "test:2222"
			},
			wantErr: false,
		},
		{
			name: "tls from env",
			env: map[string]string{
				"SECRET_KEY": "test",
				"ENABLE_TLS": "true",
			},
			prepareWant: func(c *Config) {
				c.SecretKey = "test"
				c.EnableTLS = true
			},
			wantErr: false,
		},
		{
			name: "dsn from env",
			env: map[string]string{
				"SECRET_KEY":   "test",
				"DATABASE_DSN": "dsn://test",
			},
			prepareWant: func(c *Config) {
				c.SecretKey = "test"
				c.PostgresDSN = "dsn://test"
			},
			wantErr: false,
		},
		{
			name: "upload path from arg",
			args: []string{"--dir", "/test_dir/"},
			env: map[string]string{
				"SECRET_KEY": "test",
			},
			prepareWant: func(c *Config) {
				c.SecretKey = "test"
				c.UploadDir = "/test_dir/"
			},
			wantErr: false,
		},
		{
			name: "upload path from env",
			env: map[string]string{
				"SECRET_KEY": "test",
				"UPLOAD_DIR": "/env_test_dir/",
			},
			prepareWant: func(c *Config) {
				c.SecretKey = "test"
				c.UploadDir = "/env_test_dir/"
			},
			wantErr: false,
		},
		{
			name:        "secret key not provided",
			prepareWant: func(c *Config) {},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.env {
				t.Setenv(key, value)
			}

			conf, err := NewConfig(tt.args)

			if tt.wantErr {
				assert.Error(t, err)
				t.SkipNow()
			} else {
				assert.NoError(t, err)
			}

			wantCfg := defaultConfig()
			tt.prepareWant(wantCfg)

			if !reflect.DeepEqual(conf, wantCfg) {
				t.Errorf("conf got %+v, want %+v", conf, wantCfg)
			}

			// Remove envs for next tests
			for key := range tt.env {
				t.Setenv(key, "")
			}
		})
	}
}
