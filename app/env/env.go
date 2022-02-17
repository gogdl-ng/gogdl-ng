package env

import (
	"os"
	"path/filepath"

	"github.com/LegendaryB/gogdl-ng/app/logging"
)

const (
	config     = "config"
	downloads  = "downloads"
	completed  = "completed"
	incomplete = "incomplete"
)

var logger = logging.NewLogger()

func GetConfigurationFolder() (string, error) {
	wd, err := os.Getwd()

	if err != nil {
		logger.Errorf("failed to get current directory. %w", err)
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
		logger.Errorf("failed to get current directory. %w", err)
		return "", err
	}

	dir := filepath.Join(wd, downloads, lastPathSegment)

	return dir, nil
}