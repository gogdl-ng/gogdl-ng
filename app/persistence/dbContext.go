package persistence

import (
	"github.com/LegendaryB/gogdl-ng/app/environment"
	"github.com/LegendaryB/gostore"
)

const taskCollectionName = "tasks"

var Tasks *gostore.Collection

func NewDbContext() error {
	path, err := environment.GetDatabaseFilePath()

	if err != nil {
		return err
	}

	store, err := gostore.New(path, true)

	if err != nil {
		return err
	}

	tasks, err := store.CreateCollection(taskCollectionName)

	if err != nil {
		return err
	}

	Tasks = tasks

	return nil
}
