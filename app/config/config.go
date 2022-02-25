package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type ApplicationConfiguration struct {
	ListenPort  int
	LogFilePath string
}

type TransferConfiguration struct {
	RetryThreeshold uint
}

type Configuration struct {
	Application ApplicationConfiguration
	Transfer    TransferConfiguration
}

const (
	configFolderName string = "config"
	configFileName   string = "config.toml"
)

func NewConfigurationFromFile() (*Configuration, error) {
	var conf Configuration
	path, err := getConfigurationPath()

	if err != nil {
		return nil, err
	}

	_, err = toml.DecodeFile(path, &conf)

	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func getConfigurationPath() (string, error) {
	wd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	path := filepath.Join(wd, configFolderName)

	if err := os.MkdirAll(path, 0644); err != nil {
		return "", err
	}

	path = filepath.Join(path, configFileName)

	return path, nil
}
