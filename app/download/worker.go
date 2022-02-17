package download

import (
	"io/fs"
	"log"
	"path/filepath"
	"reflect"
	"time"

	"github.com/LegendaryB/gogdl-ng/app/env"
	"github.com/LegendaryB/gogdl-ng/app/gdrive"
	"github.com/LegendaryB/gogdl-ng/app/utils"
)

var completedFolder, _ = env.GetCompletedFolder()
var incompleteFolder, _ = env.GetIncompleteFolder()

func Run() error {
	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		folders, err := utils.Subfolders(incompleteFolder)

		if err != nil {
			log.Fatal(err)
		}

		for _, folder := range folders {
			folderPath := filepath.Join(getFullPath(folder), folder.Name())
			driveId, err := readDriveIdFile(folderPath)

			if err != nil {
				log.Print("failed to read drive-id file")
			}

			driveFiles, err := gdrive.GetFilesFromFolder(driveId)

			if err != nil {
				log.Print("failed to retrieve files from google drive")
			}

			for _, driveFile := range driveFiles {
				if err = gdrive.DownloadFile(folderPath, driveFile); err != nil {
					log.Print("asdas")
				}

				if err := deleteDriveIdFile(folderPath); err != nil {
					log.Print("asdas")
				}

				err = utils.Move(folderPath, filepath.Join(completedFolder, folder.Name()))

				if err != nil {
					e := err.Error()
					log.Print(e)
				}
			}
		}
	}

	return nil
}

func getFullPath(fi fs.FileInfo) string {
	fv := reflect.ValueOf(fi).Elem().FieldByName("path")

	return fv.String()
}
