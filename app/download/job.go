package download

import (
	"os"
	"path/filepath"

	"github.com/LegendaryB/gogdl-ng/app/gdrive"
)

type Job struct {
	DriveId string
}

func RegisterNewJob(driveFolder *gdrive.DriveFolder) error {
	path, err := createJobFolder(driveFolder.Name)

	if err != nil {
		return err
	}

	if err := createDriveIdFile(path, driveFolder.Id); err != nil {
		return err
	}

	return nil
}

func createJobFolder(folderName string) (string, error) {
	path := filepath.Join(incompleteFolder, folderName)

	if err := os.MkdirAll(path, 0755); err != nil {
		return "", err
	}

	return path, nil
}
