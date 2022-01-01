package download

import (
	"fmt"
	"time"

	"github.com/LegendaryB/gogdl-ng/app/models/task"
)

func Start(errch chan error) {
	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		tasks, err := getUnfinishedTasks()

		if err != nil {
			errch <- err
			break
		}

		for _, task := range tasks {
			fmt.Println(task.DriveId)
		}

	}

	errch <- nil
}

func getUnfinishedTasks() ([]task.Task, error) {
	var unfinishedTasks []task.Task
	tasks, err := task.Repository.GetAll()

	if err != nil {
		return unfinishedTasks, err
	}

	for _, v := range tasks {
		if !v.IsCompleted {
			unfinishedTasks = append(unfinishedTasks, v)
		}
	}

	return unfinishedTasks, nil
}
