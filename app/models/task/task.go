package task

type Task struct {
	Id          int64  `json:"id"`
	DriveId     string `json:"driveId"`
	DriveName   string `json:"driveName"`
	LocalPath   string `json:"localPath"`
	IsCompleted bool   `json:"isCompleted"`
}
