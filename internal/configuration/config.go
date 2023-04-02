package configuration

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Logging struct {
		LogLevel string `yaml:"log_level"`
	} `yaml:"logging"`
}

func loadConfig() (*Config, error) {
	config := &Config{}

	file, err := os.Open("config.yml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
