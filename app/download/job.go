package download

import (
	"os"
	"path/filepath"

	"github.com/LegendaryB/gogdl-ng/app/env"
	"github.com/LegendaryB/gogdl-ng/app/gdrive"
)

type Job struct {
	DriveId string
}

func RegisterNewJob(driveFolder *gdrive.DriveFolder) error {
	path, err := createJobFolder(driveFolder.Name)

	if err != nil {
		logger.Errorf("failed to create job folder. %v", err)
		return err
	}

	if err := createDriveIdFile(path, driveFolder.Id); err != nil {
		logger.Errorf("failed to create drive id file. %v", err)
		return err
	}

	return nil
}

func createJobFolder(folderName string) (string, error) {
	path := filepath.Join(env.IncompleteFolder, folderName)

	if err := os.MkdirAll(path, 0755); err != nil {
		logger.Errorf("failed to create folder(s). %v", err)
		return "", err
	}

	return path, nil
}
