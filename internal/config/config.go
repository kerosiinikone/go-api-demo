package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DBU string `yaml:"dbu"`
}

func Load(file string) (*Config, error) {
	cfg := Config{}

	bytes, err := os.ReadFile(file)

	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}