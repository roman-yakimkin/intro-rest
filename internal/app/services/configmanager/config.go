package configmanager

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	BindAddr           string `yaml:"bind_addr"`
	MongoURI           string `yaml:"mongo_uri"`
	MaxEntities        int    `yaml:"max_entities"`
	CopyDataIntervalMS int    `yaml:"copy_data_interval_ms"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
	}
}

func (cm *Config) Init(configPath string) error {
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, cm)
	return err
}
