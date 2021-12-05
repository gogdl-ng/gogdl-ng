package persistence

import (
	"database/sql"

	"github.com/LegendaryB/gogdl-ng/app/environment"

	_ "github.com/mattn/go-sqlite3"
)

var Database *sql.DB

func NewDatabase() error {
	dbFilePath, err := environment.GetDatabaseFilePath()

	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", dbFilePath)

	if err != nil {
		return err
	}

	Database = db

	if err = createTables(); err != nil {
		return err
	}

	return nil
}

func createTables() error {
	if err := createTableForTasks(); err != nil {
		return err
	}

	if err := createTableForDownloads(); err != nil {
		return err
	}

	return nil
}

func createTableForTasks() error {
	err := createTable(
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

func createTableForDownloads() error {
	err := createTable(
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

func createTable(statement string) error {
	stmt, err := Database.Prepare(statement)

	if err != nil {
		return err
	}

	_, err = stmt.Exec()

	if err != nil {
		return err
	}

	return nil
}
