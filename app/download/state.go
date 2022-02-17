package download

import (
	"io/ioutil"
	"path/filepath"
)

const driveIdFileName = "drive-id"

func createDriveIdFile(path string, driveId string) error {
	path = filepath.Join(path, driveIdFileName)

	if err := ioutil.WriteFile(path, []byte(driveId), 0755); err != nil {
		return err
	}

	return nil
}

func readDriveIdFile(path string) (string, error) {
	path = filepath.Join(path, driveIdFileName)

	buf, err := ioutil.ReadFile(path)

	if err != nil {
		return "", err
	}

	driveId := string(buf)

	return driveId, nil
}
