package core

import (
	"os"
	"strconv"

	"documents/internal/storage"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Storage storage.Config `yaml:"storage"`
	Metrics struct {
		Name string `yaml:"name"`
	} `yaml:"metrics"`
	Server struct {
		Host       string `yaml:"host"`
		Port       int    `yaml:"port"`
		CORSOrigin string `yaml:"cors_origin"`
		Callbacks  struct {
			BackRedirectURL string `yaml:"back_redirect_url"`
			Google          struct {
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

	if portStr, ok := os.LookupEnv("PORT"); ok {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}

		config.Server.Port = port
	}

	return &config, nil
}
