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
			fmt.Println(task.Status)
		}

	}

	errch <- nil
}

func getUnfinishedTasks() ([]task.Task, error) {
	var unfinishedTasks []task.Task

	return unfinishedTasks, nil
}
