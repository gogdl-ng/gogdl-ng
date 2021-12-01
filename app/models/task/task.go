package task

type Task struct {
	ID         int64  `json:"id"`
	FolderId   string `json:"folderId"`
	FolderName string `json:"folderName"`
	Status     string `json:"status"`
}

type TaskStatus int64

const (
	Created TaskStatus = iota
	Processing
	Done
)

func (ts TaskStatus) String() string {
	switch ts {
	case Created:
		return "new"
	case Processing:
		return "processing"
	case Done:
		return "done"
	}

	return "unknown"
}
