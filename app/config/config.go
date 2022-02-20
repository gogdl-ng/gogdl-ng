package config

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/LegendaryB/gogdl-ng/app/env"
)

type ApplicationConfiguration struct {
	ListenPort  int
	LogFilePath string
}

type TransferConfiguration struct {
	AcknowledgeAbuseFlag bool
	RetryThreeshold      uint
}

type Configuration struct {
	Application ApplicationConfiguration
	Transfer    TransferConfiguration
}

func NewConfigurationFromFile() (*Configuration, error) {
	var conf Configuration
	path := filepath.Join(env.ConfigurationFolder, "config.toml")

	_, err := toml.DecodeFile(path, &conf)

	if err != nil {
		return nil, err
	}

	return &conf, nil
}
