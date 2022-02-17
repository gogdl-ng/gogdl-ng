package download

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type JobState struct {
	Finished bool
	DriveId  string
}

const state_file_name = "state.json"

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
