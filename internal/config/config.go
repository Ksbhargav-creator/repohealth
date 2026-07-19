package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Repos  []RepoRef       `yaml:"repos"`
	Orgs   []string        `yaml:"orgs"`
	Checks map[string]bool `yaml:"checks"`
}

type RepoRef struct {
	Owner string `yaml:"owner"`
	Name  string `yaml:"name"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
