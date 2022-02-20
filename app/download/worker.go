package download

import (
	"io/fs"
	"path/filepath"
	"reflect"
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
		logger.Fatalf("failed to initialize Google Drive service. %v", err)
		return nil, err
	}

	return &Downloader{logger: logger, conf: conf, drive: drive}, nil
}

func (service *Downloader) Run() error {
	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		folders, err := utils.Subfolders(env.IncompleteFolder)

		if err != nil {
			service.logger.Errorf("failed to retrieve subfolders. %v", err)
			return err
		}

		for _, folder := range folders {
			folderPath := filepath.Join(getFullPath(folder), folder.Name())
			driveId, err := service.readDriveIdFile(folderPath)

			if err != nil {
				service.logger.Errorf("failed to read drive id file. skipping. %v", err)
				continue
			}

			driveFiles, err := service.drive.GetFilesFromFolder(driveId)

			if err != nil {
				service.logger.Errorf("failed to retrieve files from google drive. skipping. %v", err)
				continue
			}

			if err := service.downloadFiles(folderPath, driveFiles); err != nil {
				service.logger.Errorf("failed to download files from google drive. skipping. %v", err)
				continue
			}

			if err = utils.Move(folderPath, filepath.Join(env.CompletedFolder, folder.Name())); err != nil {
				service.logger.Errorf("failed to move files to target path. skipping. %v", err)
				continue
			}
		}
	}

	return nil
}

func (service *Downloader) downloadFiles(targetPath string, files []*drive.File) error {
	for _, driveFile := range files {
		if err := service.drive.DownloadFile(targetPath, driveFile); err != nil {
			service.logger.Errorf("failed to download file (name: %s, id: %s). %v", driveFile.Name, driveFile.Id, err)
			return err
		}
	}

	if err := service.deleteDriveIdFile(targetPath); err != nil {
		service.logger.Errorf("failed to delete drive id file. %v", err)
		return err
	}

	return nil
}

func getFullPath(fi fs.FileInfo) string {
	fv := reflect.ValueOf(fi).Elem().FieldByName("path")

	return fv.String()
}
