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

const stateFileName = "state.json"

func createStateFile(path string, driveId string) error {
	state := JobState{
		DriveId: driveId,
	}

	return writeJobState(path, &state)
}

func writeJobState(folderName string, state *JobState) error {
	json, err := json.MarshalIndent(state, "", " ")

	if err != nil {
		return err
	}

	path := filepath.Join(baseFolder, folderName, stateFileName)

	if err = ioutil.WriteFile(path, json, 0755); err != nil {
		return err
	}

	return nil
}

func readJobState(folderName string) (*JobState, error) {
	path := filepath.Join(baseFolder, folderName, stateFileName)
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var state JobState
	if err = json.NewDecoder(file).Decode(&state); err != nil {
		return nil, err
	}

	return &state, nil
}
