/*
package download

import (
	"path/filepath"
	"time"

	"github.com/LegendaryB/gogdl-ng/app/config"
	"github.com/LegendaryB/gogdl-ng/app/env"
	"github.com/LegendaryB/gogdl-ng/app/gdrive"
	"github.com/LegendaryB/gogdl-ng/app/logging"
	"github.com/LegendaryB/gogdl-ng/app/utils"
	"google.golang.org/api/drive/v3"
)

type Downloader struct {
	logger logging.Logger
	conf   *config.TransferConfiguration

	drive *gdrive.DriveService
}

func NewDownloader(conf *config.TransferConfiguration, logger logging.Logger) (*Downloader, error) {
	drive, err := gdrive.NewDriveService(conf, logger)

	if err != nil {
		logger.Fatalf("Failed to initialize Google Drive service. %v", err)
		return nil, err
	}

	return &Downloader{conf: conf, logger: logger, drive: drive}, nil
}

func (service *Downloader) Run() error {
	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		folders, err := utils.Subfolders(env.IncompleteFolder)

		if err != nil {
			service.logger.Errorf("Failed to retrieve subfolders. %v", err)
			return err
		}

		for _, folder := range folders {
			path := filepath.Join(env.IncompleteFolder, folder.Name())
			driveId, err := service.readDriveIdFile(path)

			service.logger.Infof("Job: '%s' | id: '%s'", folder.Name(), driveId)

			if err != nil {
				service.logger.Errorf("Failed to read drive id file. Skipping.. %v", err)
				continue
			}

			driveFiles, err := service.drive.GetFilesFromFolder(driveId)

			if err != nil {
				service.logger.Errorf("Failed to retrieve files from google drive. Skipping.. %v", err)
				continue
			}

			if err := service.downloadFiles(path, driveFiles); err != nil {
				service.logger.Errorf("Failed to download files from google drive. Skipping.. %v", err)
				continue
			}

			targetPath := filepath.Join(env.CompletedFolder, folder.Name())

			if err = utils.Move(path, targetPath); err != nil {
				service.logger.Errorf("Failed to move files to target path. Skipping.. %v", err)
				continue
			}

			service.logger.Info("Job finished.")
		}
	}

	return nil
}

func (service *Downloader) downloadFiles(targetPath string, files []*drive.File) error {
	for _, driveFile := range files {
		if err := service.drive.DownloadFile(targetPath, driveFile); err != nil {
			service.logger.Errorf("Failed to download file (name: %s, id: %s). %v", driveFile.Name, driveFile.Id, err)
			return err
		}
	}

	if err := service.deleteDriveIdFile(targetPath); err != nil {
		service.logger.Errorf("Failed to delete drive id file. %v", err)
		return err
	}

	return nil
}
*/