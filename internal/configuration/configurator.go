package configuration

import (
	"sync"
)

type Configurator interface {
	GetConfig() *Config
}

type ConcreteConfigurator struct {
	config *Config
}

func (configurator *ConcreteConfigurator) GetConfig() *Config {
	return configurator.config
}

var configuratorInstance *ConcreteConfigurator

func GetInstance() (*ConcreteConfigurator, error) {
	once := &sync.Once{}
	var err error

	if configuratorInstance == nil {
		once.Do(
			func() {
				var config *Config
				config, err = loadConfig()

				if err != nil {
					return
				}

				configuratorInstance = &ConcreteConfigurator{
					config,
				}
			})
	}

	if err != nil {
		return nil, err
	}

	return configuratorInstance, nil
}
