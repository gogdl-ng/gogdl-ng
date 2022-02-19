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

var ConfigurationFolder string
var CompletedFolder string
var IncompleteFolder string

func InitializeEnvironment() error {
	configurationFolder, err := getConfigurationFolder()

	if err != nil {
		return err
	}

	ConfigurationFolder = configurationFolder

	completedFolder, err := getCompletedFolder()

	if err != nil {
		return err
	}

	CompletedFolder = completedFolder

	incompleteFolder, err := getIncompleteFolder()

	if err != nil {
		return err
	}

	IncompleteFolder = incompleteFolder

	return nil
}

func getConfigurationFolder() (string, error) {
	wd, err := os.Getwd()

	if err != nil {
		logger.Errorf("failed to get current directory. %w", err)
		return "", err
	}

	dir := filepath.Join(wd, config)

	return dir, nil
}

func getIncompleteFolder() (string, error) {
	return getDownloadsFolderPath(incomplete)
}

func getCompletedFolder() (string, error) {
	return getDownloadsFolderPath(completed)
}

func getDownloadFolder() (string, error) {
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
