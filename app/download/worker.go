package download

import (
	"io/fs"
	"path/filepath"
	"reflect"
	"time"

	"github.com/LegendaryB/gogdl-ng/app/env"
	"github.com/LegendaryB/gogdl-ng/app/gdrive"
	"github.com/LegendaryB/gogdl-ng/app/logging"
	"github.com/LegendaryB/gogdl-ng/app/utils"
	"google.golang.org/api/drive/v3"
)

var logger = logging.NewLogger()

var completedFolder, _ = env.GetCompletedFolder()
var incompleteFolder, _ = env.GetIncompleteFolder()

func Run() error {
	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		folders, err := utils.Subfolders(incompleteFolder)

		if err != nil {
			logger.Errorf("failed to retrieve subfolders. %w", err)
			return err
		}

		for _, folder := range folders {
			folderPath := filepath.Join(getFullPath(folder), folder.Name())
			driveId, err := readDriveIdFile(folderPath)

			if err != nil {
				logger.Errorf("failed to read drive id file. skipping. %w", err)
				continue
			}

			driveFiles, err := gdrive.GetFilesFromFolder(driveId)

			if err != nil {
				logger.Errorf("failed to retrieve files from google drive. skipping. %w", err)
				continue
			}

			if err := downloadFiles(folderPath, driveFiles); err != nil {
				logger.Errorf("failed to download files from google drive. skipping. %w", err)
				continue
			}

			if err = utils.Move(folderPath, filepath.Join(completedFolder, folder.Name())); err != nil {
				logger.Errorf("failed to move files to target path. skipping. %w", err)
				continue
			}
		}
	}

	return nil
}

func downloadFiles(targetPath string, files []*drive.File) error {
	for _, driveFile := range files {
		if err := gdrive.DownloadFile(targetPath, driveFile); err != nil {
			logger.Errorf("failed to download file (name: %s, id: %s). %w", driveFile.Name, driveFile.Id, err)
			return err
		}
	}

	if err := deleteDriveIdFile(targetPath); err != nil {
		logger.Errorf("failed to delete drive id file. %w", err)
		return err
	}

	return nil
}

func getFullPath(fi fs.FileInfo) string {
	fv := reflect.ValueOf(fi).Elem().FieldByName("path")

	return fv.String()
}
