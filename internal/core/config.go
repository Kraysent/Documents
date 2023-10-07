package core

import (
	"os"

	"documents/internal/storage"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Logging struct {
		StdoutPath string `yaml:"stdout_path"`
		StderrPath string `yaml:"stderr_path"`
	} `yaml:"logging"`
	Storage storage.Config `yaml:"storage"`
	Server  struct {
		Port int `yaml:"port"`
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
