package env

import (
	"os"
	"path/filepath"
)

const (
	downloads  = "downloads"
	completed  = "completed"
	incomplete = "incomplete"
)

var CompletedFolder string
var IncompleteFolder string

func NewEnvironment() error {
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

func getIncompleteFolder() (string, error) {
	return getDownloadsFolderPath(incomplete)
}

func getCompletedFolder() (string, error) {
	return getDownloadsFolderPath(completed)
}

func getDownloadsFolderPath(lastPathSegment string) (string, error) {
	wd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	path := filepath.Join(wd, downloads, lastPathSegment)

	if err := os.MkdirAll(path, 0755); err != nil {
		return "", err
	}

	return path, nil
}
