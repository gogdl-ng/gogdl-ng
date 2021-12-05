package persistence

import (
	"database/sql"

	"github.com/LegendaryB/gogdl-ng/app/environment"

	_ "github.com/mattn/go-sqlite3"
)

func NewDatabase() (*sql.DB, error) {
	dbFilePath, err := environment.GetDatabaseFilePath()

	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbFilePath)

	if err != nil {
		return nil, err
	}

	if err = createTables(db); err != nil {
		return nil, err
	}

	return db, err
}

func createTables(db *sql.DB) error {
	if err := createTableForTasks(db); err != nil {
		return err
	}

	if err := createTableForDownloads(db); err != nil {
		return err
	}

	return nil
}

func createTableForTasks(db *sql.DB) error {
	err := createTable(db,
		`CREATE TABLE IF NOT EXISTS
			tasks (
				Id			INTEGER UNIQUE,
				DriveId		INTEGER NOT NULL UNIQUE,
				DriveName	TEXT NOT NULL,
				LocalPath	TEXT NOT NULL UNIQUE,
				IsCompleted	INTEGER NOT NULL,
				PRIMARY KEY(Id AUTOINCREMENT)
			)`)

	return err
}

func createTableForDownloads(db *sql.DB) error {
	err := createTable(db,
		`CREATE TABLE IF NOT EXISTS 
			downloads (
				Id			INTEGER UNIQUE,
				TaskId		INTEGER NOT NULL,
				DriveId		TEXT NOT NULL UNIQUE,
				DriveName	TEXT NOT NULL,
				DriveHash	TEXT NOT NULL,
				LocalPath	TEXT NOT NULL UNIQUE,
				PRIMARY KEY(Id AUTOINCREMENT),
				FOREIGN KEY(TaskId) REFERENCES tasks(id)
			)`)

	return err
}

func createTable(db *sql.DB, statement string) error {
	stmt, err := db.Prepare(statement)

	if err != nil {
		return err
	}

	_, err = stmt.Exec()

	if err != nil {
		return err
	}

	return nil
}
