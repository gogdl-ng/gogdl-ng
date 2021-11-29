package models

type TaskStatus int64

const (
	New TaskStatus = iota
	Processing
	Done
)

func (ts TaskStatus) String() string {
	switch ts {
	case New:
		return "new"
	case Processing:
		return "processing"
	case Done:
		return "done"
	}

	return "unknown"
}
