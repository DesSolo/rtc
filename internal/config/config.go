package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config ...
type Config struct {
	Logging struct {
		Level int `yaml:"level"`
	} `yaml:"logging"`
	Server struct {
		Address           string        `yaml:"address"`
		ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
		Auth              struct {
			JWT struct {
				PrivateKey string        `yaml:"private_key"`
				PublicKey  string        `yaml:"public_key"`
				TTL        time.Duration `yaml:"ttl"`
			} `yaml:"jwt"`
		} `yaml:"auth"`
	} `yaml:"server"`
	Storage struct {
		DSN string `yaml:"dsn"`
	} `yaml:"storage"`
	ValuesStorage struct {
		Endpoints   []string      `yaml:"endpoints"`
		DialTimeout time.Duration `yaml:"dial_timeout"`
		Path        string        `yaml:"path"`
	} `yaml:"values_storage"`
}

// NewConfigFromFile ...
func NewConfigFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path) // nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("yaml.Unmarshal: %w", err)
	}

	return &config, nil
}
