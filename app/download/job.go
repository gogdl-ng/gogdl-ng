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
	folder, err := createJobFolder(driveFolder.Name)

	if err != nil {
		return err
	}

	if err = createStateFile(folder, driveFolder.Id); err != nil {
		return err
	}

	return nil
}

func createJobFolder(name string) (string, error) {
	folderPath := filepath.Join(downloadFolder, name)

	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return "", err
	}

	return folderPath, nil
}
