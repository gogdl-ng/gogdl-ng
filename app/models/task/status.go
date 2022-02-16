package task

import (
	"encoding/json"
	"errors"
)

type TaskStatus string

const (
	New        = "New"
	Processing = "Processing"
	Finished   = "Finished"
	Error      = "Error"
)

func (ts *TaskStatus) UnmarshalJSON(b []byte) error {
	type TS TaskStatus
	var r *TS = (*TS)(ts)
	err := json.Unmarshal(b, &r)
	if err != nil {
		return err
	}

	switch *ts {
	case New, Processing, Finished, Error:
		return nil
	}

	return errors.New("invalid task status")
}
