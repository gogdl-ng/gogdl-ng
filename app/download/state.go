package download

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const driveIdFileName = "drive-id"

func (service *Downloader) createDriveIdFile(path string, driveId string) error {
	path = filepath.Join(path, driveIdFileName)

	if err := ioutil.WriteFile(path, []byte(driveId), 0755); err != nil {
		service.logger.Errorf("failed to write to drive id file. %v", err)
		return err
	}

	return nil
}

func (service *Downloader) deleteDriveIdFile(path string) error {
	path = filepath.Join(path, driveIdFileName)

	if err := os.Remove(path); err != nil {
		service.logger.Errorf("failed to delete drive id file. %v", err)
		return err
	}

	return nil
}

func (service *Downloader) readDriveIdFile(path string) (string, error) {
	path = filepath.Join(path, driveIdFileName)

	buf, err := ioutil.ReadFile(path)

	if err != nil {
		service.logger.Errorf("failed to read from drive id file. %v", err)
		return "", err
	}

	driveId := string(buf)

	return driveId, nil
}
