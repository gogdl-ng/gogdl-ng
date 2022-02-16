package download

import (
	"fmt"
	"time"

	"github.com/LegendaryB/gogdl-ng/app/models/task"
)

func Run() {
	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		tasks, err := getUnfinishedTasks()

		if err != nil {
			// todo: log
			//( errch <- err
			break
		}

		for _, task := range tasks {
			fmt.Println(task.Status)
		}
	}
}

func getUnfinishedTasks() ([]task.Task, error) {
	var unfinishedTasks []task.Task

	return unfinishedTasks, nil
}
