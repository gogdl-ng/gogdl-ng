package download

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/LegendaryB/gogdl-ng/app/environment"
	"github.com/LegendaryB/gogdl-ng/app/gdrive"
)

const state_file_name = "state.json"

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

func readJobState(folderName string) (*JobState, error) {
	file, err := os.Open(filepath.Join(downloadFolder, folderName, state_file_name))

	if err != nil {
		return nil, err
	}

	defer file.Close()

	if err != nil {
		return nil, err
	}

	var state JobState
	if err = json.NewDecoder(file).Decode(&state); err != nil {
		return nil, err
	}

	return &state, nil
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

func createStateFile(path string, driveId string) error {
	state := JobState{
		DriveId: driveId,
	}

	json, err := json.MarshalIndent(state, "", " ")

	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(filepath.Join(path, state_file_name), json, 0755); err != nil {
		return err
	}

	return nil
}
