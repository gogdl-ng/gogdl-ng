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
	if err := createJobFolder(driveFolder.Name); err != nil {
		return err
	}

	if err := createStateFile(driveFolder.Name, driveFolder.Id); err != nil {
		return err
	}

	return nil
}

func createJobFolder(folderName string) error {
	path := filepath.Join(baseFolder, folderName)

	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	return nil
}
