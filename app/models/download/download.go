package download

type Download struct {
	Id        int64  `json:"id"`
	TaskId    int64  `json:"taskId"`
	DriveId   string `json:"driveId"`
	DriveName string `json:"driveName"`
	DriveHash string `json:"driveHash"`
	LocalPath string `json:"localPath"`
}
