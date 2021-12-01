package persistence

import (
	"database/sql"

	"github.com/LegendaryB/gogdl-ng/app/models/task"

	_ "github.com/mattn/go-sqlite3"
)

type DbContext struct {
	Tasks *task.Repository
}

func NewDbContext() (*DbContext, error) {
	db, err := sql.Open("sqlite3", "./config/gogdl-ng.db")

	if err != nil {
		return nil, err
	}

	if err = createTasksTable(db); err != nil {
		return nil, err
	}

	context := &DbContext{}
	context.Tasks = &task.Repository{DB: db}

	return context, nil
}

func createTasksTable(db *sql.DB) error {
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS
			tasks (
				ID	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				FolderId   TEXT,
				FolderName TEXT,
				Status	   TEXT CHECK( Status IN ('new','processing','done') )
			);
	`)

	if err != nil {
		return err
	}

	_, err = stmt.Exec()

	if err != nil {
		return err
	}

	return nil
}
