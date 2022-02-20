package download

import (
	"os"
	"path/filepath"

	"github.com/LegendaryB/gogdl-ng/app/env"
)

type Job struct {
	DriveId string
}

func (service *Downloader) RegisterNewJob(driveId string) error {
	driveFolder, err := service.drive.GetFolderById(driveId)

	if err != nil {
		service.logger.Errorf("Failed to get google drive folder. %v", err)
		return err
	}

	path, err := service.createJobFolder(driveFolder.Name)

	if err != nil {
		service.logger.Errorf("Failed to create job folder. %v", err)
		return err
	}

	if err := service.createDriveIdFile(path, driveFolder.Id); err != nil {
		service.logger.Errorf("Failed to create drive id file. %v", err)
		return err
	}

	return nil
}

func (service *Downloader) createJobFolder(folderName string) (string, error) {
	path := filepath.Join(env.IncompleteFolder, folderName)

	if err := os.MkdirAll(path, 0755); err != nil {
		service.logger.Errorf("failed to create folder(s). %v", err)
		return "", err
	}

	return path, nil
}
