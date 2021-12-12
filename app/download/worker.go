package download

import (
	"fmt"
	"time"

	"github.com/LegendaryB/gogdl-ng/app/models/task"
)

func Start(errch chan error) {
	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		fmt.Print("executed")

		// tasks, err := getUnfinishedTasks()

		// err := nil //errors.New("test")

		// if err != nil {
		// 	errch <- err
		// 	break
		// }

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
