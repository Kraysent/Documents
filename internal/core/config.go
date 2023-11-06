package core

import (
	"os"

	"documents/internal/storage"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Storage storage.Config `yaml:"storage"`
	Server  struct {
		Host      string `yaml:"host"`
		Port      int    `yaml:"port"`
		Callbacks struct {
			Google struct {
				RedirectURL string `yaml:"redirect_url"`
			} `yaml:"google"`
		} `yaml:"callbacks"`
	} `yaml:"server"`
}

func ParseConfig(filename string) (*Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config

	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
