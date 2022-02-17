package download

import (
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"time"

	"github.com/LegendaryB/gogdl-ng/app/env"
	"github.com/LegendaryB/gogdl-ng/app/gdrive"
)

var completedFolder, _ = env.GetCompletedFolder()
var incompleteFolder, _ = env.GetIncompleteFolder()

func Run() error {
	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		folders, err := getSubfolders(incompleteFolder)

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
				gdrive.DownloadFile(folderPath, driveFile)
			}
		}
	}

	return nil
}

func getFullPath(fi fs.FileInfo) string {
	fv := reflect.ValueOf(fi).Elem().FieldByName("path")

	return fv.String()
}

func getSubfolders(path string) ([]fs.FileInfo, error) {
	items, err := ioutil.ReadDir(path)

	if err != nil {
		return nil, err
	}

	var subfolders []fs.FileInfo

	for _, item := range items {
		if !item.IsDir() {
			continue
		}

		subfolders = append(subfolders, item)
	}

	return subfolders, nil
}
