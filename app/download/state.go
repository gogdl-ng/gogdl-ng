package download

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const driveIdFileName = "drive-id"

func createDriveIdFile(path string, driveId string) error {
	path = filepath.Join(path, driveIdFileName)

	if err := ioutil.WriteFile(path, []byte(driveId), 0755); err != nil {
		logger.Errorf("failed to write to drive id file. %w", err)
		return err
	}

	return nil
}

func deleteDriveIdFile(path string) error {
	path = filepath.Join(path, driveIdFileName)

	if err := os.Remove(path); err != nil {
		logger.Errorf("failed to delete drive id file. %w", err)
		return err
	}

	return nil
}

func readDriveIdFile(path string) (string, error) {
	path = filepath.Join(path, driveIdFileName)

	buf, err := ioutil.ReadFile(path)

	if err != nil {
		logger.Errorf("failed to read from drive id file. %w", err)
		return "", err
	}

	driveId := string(buf)

	return driveId, nil
}
