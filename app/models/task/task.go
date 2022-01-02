package task

import "github.com/LegendaryB/gogdl-ng/app/gdrive"

type Task struct {
	Name   string
	Status TaskStatus
	Files  []*gdrive.DriveFile
}
