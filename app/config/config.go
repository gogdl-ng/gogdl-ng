package config

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/LegendaryB/gogdl-ng/app/env"
	"github.com/LegendaryB/gogdl-ng/app/logging"
)

type ApplicationConfiguration struct {
	ListenPort int
}

type GoogleDriveConfiguration struct {
	AcknowledgeAbuseFlag bool
}

type TransferConfiguration struct {
	RetryThreeshold        int
	MaxConcurrentDownloads int
}

type Configuration struct {
	Application ApplicationConfiguration
	GoogleDrive GoogleDriveConfiguration
	Transfer    TransferConfiguration
}

const CONFIGURATION_FILE = "config.toml"

var config *Configuration = nil

var logger = logging.NewLogger()

func GetConfiguration() (*Configuration, error) {
	if config != nil {
		return config, nil
	}

	var conf Configuration
	path := filepath.Join(env.ConfigurationFolder, CONFIGURATION_FILE)

	_, err := toml.DecodeFile(path, &conf)

	if err != nil {
		logger.Errorf("failed to read configuration file. %w", err)
		return nil, err
	}

	return &conf, nil
}
