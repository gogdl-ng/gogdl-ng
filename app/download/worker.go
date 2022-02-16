package download

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/LegendaryB/gogdl-ng/app/environment"
	"github.com/LegendaryB/gogdl-ng/app/gdrive"
)

const state_file_name = "state.json"

func Run() {
	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {

	}
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
	folder, err := environment.GetDownloadDirectory()

	if err != nil {
		return "", err
	}

	folderPath := filepath.Join(folder, name)

	if err = os.MkdirAll(folderPath, 0755); err != nil {
		return "", err
	}

	return folderPath, nil
}

func createStateFilePath(path string) string {
	return filepath.Join(path, state_file_name)
}

func createStateFile(path string, driveId string) error {
	state := JobState{
		DriveId: driveId,
	}

	json, err := json.MarshalIndent(state, "", " ")

	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(createStateFilePath(path), json, 0755); err != nil {
		return err
	}

	return nil
}
