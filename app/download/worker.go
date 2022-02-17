package download

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"time"

	"github.com/LegendaryB/gogdl-ng/app/env"
	"github.com/LegendaryB/gogdl-ng/app/gdrive"
)

var baseFolder, _ = env.GetDownloadFolder()

func Run() error {
	ticker := time.NewTicker(5 * time.Second)
	var finishedJobs []fs.FileInfo

	for range ticker.C {
		folders, err := getSubfolders(baseFolder)

		if err != nil {
			log.Fatal(err)
		}

		for _, folder := range folders {
			if contains(finishedJobs, folder) {
				continue
			}

			state, err := readJobState(folder.Name())

			if err != nil {
				log.Printf("%s does not contain a state file.", folder.Name())
			}

			fmt.Print(state)

			if !state.Finished {
				driveFiles, err := gdrive.GetFilesFromFolder(state.DriveId)

				if err != nil {
					return err
				}

				for _, driveFile := range driveFiles {
					fmt.Print(driveFile.Name)
				}

				state.Finished = true

				if err = writeJobState(folder.Name(), state); err != nil {
					finishedJobs = append(finishedJobs, folder)
				}
			}
		}
	}

	return nil
}

func contains(fia []fs.FileInfo, fis fs.FileInfo) bool {
	for _, fi := range fia {
		if fi.Name() == fis.Name() {
			return true
		}
	}
	return false
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
