package environment

import (
	"os"
	"path/filepath"
)

const (
	config    = "config"
	downloads = "downloads"
)

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
