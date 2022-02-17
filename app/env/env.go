package env

import (
	"os"
	"path/filepath"
)

const (
	config     = "config"
	downloads  = "downloads"
	completed  = "completed"
	incomplete = "incomplete"
)

func GetConfigurationFolder() (string, error) {
	wd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	dir := filepath.Join(wd, config)

	return dir, nil
}

func GetIncompleteFolder() (string, error) {
	return getDownloadsFolderPath(incomplete)
}

func GetCompletedFolder() (string, error) {
	return getDownloadsFolderPath(completed)
}

func GetDownloadFolder() (string, error) {
	return getDownloadsFolderPath("")
}

func getDownloadsFolderPath(lastPathSegment string) (string, error) {
	wd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	dir := filepath.Join(wd, downloads, lastPathSegment)

	return dir, nil
}
