package environment

import (
	"os"
	"path/filepath"
)

const (
	config     = "config"
	downloads  = "downloads"
	dbFileName = "gogdl-ng.json"
)

func GetDatabaseFilePath() (string, error) {
	configDir, err := GetConfigurationDirectory()

	if err != nil {
		return "", nil
	}

	return filepath.Join(configDir, dbFileName), nil
}

func GetConfigurationDirectory() (string, error) {
	wd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	dir := filepath.Join(wd, config)

	return dir, nil
}

func GetDownloadDirectory() (string, error) {
	wd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	dir := filepath.Join(wd, downloads)

	return dir, nil
}

func CreateTaskDirectory(name string) (string, error) {
	downloadDirectory, err := GetDownloadDirectory()

	if err != nil {
		return "", err
	}

	taskDirectoryPath := filepath.Join(downloadDirectory, name)

	if err = os.MkdirAll(taskDirectoryPath, 0755); err != nil {
		return "", err
	}

	return taskDirectoryPath, nil
}
