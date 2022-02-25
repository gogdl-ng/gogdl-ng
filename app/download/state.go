package download

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func (service *JobService) createDriveIdFile(path string, driveId string) error {
	path = filepath.Join(path, driveIdFileName)

	if err := ioutil.WriteFile(path, []byte(driveId), 0644); err != nil {
		service.logger.Errorf("Failed to write to drive id file. %v", err)
		return err
	}

	return nil
}

func (service *JobService) removeDriveIdFile(path string) error {
	path = filepath.Join(path, driveIdFileName)

	if err := os.Remove(path); err != nil {
		service.logger.Errorf("Failed to remove drive id file. %v", err)
		return err
	}

	return nil
}

func (service *JobService) readDriveIdFile(path string) (string, error) {
	path = filepath.Join(path, driveIdFileName)

	buf, err := ioutil.ReadFile(path)

	if err != nil {
		service.logger.Errorf("Failed to read from drive id file. %v", err)
		return "", err
	}

	driveId := string(buf)

	return driveId, nil
}
