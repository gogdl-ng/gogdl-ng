package repositories

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	DB *sql.DB
}

func NewTaskRepository(conn *sql.DB) *SQLite {
	stmt, _ := conn.Prepare(`
		CREATE TABLE IF NOT EXISTS
			tasks (
				ID	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				FolderId   TEXT,
				FolderName TEXT,
				Status	   TEXT CHECK( Status IN ('new','processing','done') )
			);
	`)

	stmt.Exec()

	return &SQLite{
		DB: conn,
	}
}
