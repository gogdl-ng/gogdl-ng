package environment

import (
	"os"
	"path/filepath"
)

const (
	config     = "config"
	downloads  = "downloads"
	dbFileName = "gogdl-ng.db"
)

func GetDatabaseFilePath() (string, error) {
	configDir, err := GetConfigDir()

	if err != nil {
		return "", nil
	}

	return filepath.Join(configDir, dbFileName), nil
}

func GetConfigDir() (string, error) {
	wd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	dir := filepath.Join(wd, config)

	return dir, nil
}

func GetDownloadsDir() (string, error) {
	wd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	dir := filepath.Join(wd, downloads)

	return dir, nil
}
