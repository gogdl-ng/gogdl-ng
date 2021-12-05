package task

import (
	"database/sql"
	"fmt"
)

type DataSource struct {
	DB *sql.DB
}

var Repository *DataSource

func NewRepository(db *sql.DB) error {
	Repository = &DataSource{DB: db}

	return nil
}

func (dataSource *DataSource) GetAll() ([]Task, error) {
	var tasks []Task

	rows, err := dataSource.DB.Query("SELECT * FROM tasks")

	if err != nil {
		return tasks, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int64
		var driveId, driveName, localPath string
		var isCompleted bool

		err := rows.Scan(&id, &driveId, &driveName, &localPath, &isCompleted)

		if err != nil {
			return tasks, err
		}

		task := Task{
			Id:          id,
			DriveId:     driveId,
			DriveName:   driveName,
			LocalPath:   localPath,
			IsCompleted: isCompleted,
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (dataSource *DataSource) Get(id int64) (*Task, error) {
	var task Task

	query := fmt.Sprintf("SELECT * FROM tasks WHERE Id=%d", id)

	row, err := dataSource.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer row.Close()

	if row.Next() {
		var id int64
		var driveId, driveName, localPath string
		var isCompleted bool

		err := row.Scan(&id, &driveId, &driveName, &localPath, &isCompleted)

		if err != nil {
			return nil, err
		}

		task = Task{
			Id:          id,
			DriveId:     driveId,
			DriveName:   driveName,
			LocalPath:   localPath,
			IsCompleted: isCompleted,
		}

	}

	return &task, nil
}

func (dataSource *DataSource) Create(task Task) (*Task, error) {
	query := `INSERT INTO tasks(DriveId, DriveName, LocalPath, IsCompleted) VALUES($1, $2, $3, $4);`

	result, err := dataSource.DB.Exec(query, task.DriveId, task.DriveName, task.LocalPath, false)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return nil, err
	}

	created, err := dataSource.Get(id)

	if err != nil {
		return nil, err
	}

	return created, nil
}

func (dataSource *DataSource) Update(task Task) error {
	query := `UPDATE tasks SET DriveId=$1, DriveName=$2, LocalPath=$3, IsCompleted=$4 WHERE Id=$5`

	_, err := dataSource.DB.Exec(query, task.DriveId, task.DriveName, task.LocalPath, task.IsCompleted, task.Id)

	if err != nil {
		return err
	}

	return nil
}

func (dataSource *DataSource) Delete(id int64) error {
	query := `DELETE FROM tasks WHERE Id=$1`

	_, err := dataSource.DB.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}
