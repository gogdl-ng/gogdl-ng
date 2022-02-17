package download

import (
	"io/fs"
	"io/ioutil"
	"log"
	"time"

	"github.com/LegendaryB/gogdl-ng/app/environment"
)

var downloadFolder, _ = environment.GetDownloadFolder()

func Run() error {
	ticker := time.NewTicker(5 * time.Second)
	var finishedJobs []fs.FileInfo

	for range ticker.C {
		folders, err := getSubfolders(downloadFolder)

		if err != nil {
			log.Fatal(err)
		}

		for _, folder := range folders {
			if contains(finishedJobs, folder) {
				continue
			}

			state, err := readJobState(folder.Name())

			if err != nil {
				return err
			}

			if !state.Finished {
				// download

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
