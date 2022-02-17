package download

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"time"

	"github.com/LegendaryB/gogdl-ng/app/environment"
)

var downloadFolder, _ = environment.GetDownloadFolder()

func Run() error {
	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		folders, err := getSubfolders(downloadFolder)

		if err != nil {
			log.Fatal(err)
		}

		for _, folder := range folders {
			state, err := readJobState(folder.Name())

			if err != nil {
				return err
			}

			fmt.Print(state.Finished)
		}
	}

	return nil
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
