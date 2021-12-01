package task

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	DB *sql.DB
}

func New(db *sql.DB) *Repository {
	repository := &Repository{db}
	return repository
}

func (repository *Repository) GetAll() ([]Task, error) {
	var tasks []Task

	rows, err := repository.DB.Query("SELECT ID, FolderId, FolderName, Status FROM tasks")

	if err != nil {
		return tasks, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int64
		var folderId, folderName, status string

		err := rows.Scan(&id, &folderId, &folderName, &status)

		if err != nil {
			return tasks, err
		}

		task := Task{
			ID:         id,
			FolderId:   folderId,
			FolderName: folderName,
			Status:     status,
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (repository *Repository) Get(id int64) (Task, error) {
	var task Task

	query := fmt.Sprintf("SELECT FolderId, FolderName, Status FROM tasks WHERE ID=%d", id)

	row, err := repository.DB.Query(query)

	if err != nil {
		return task, err
	}

	defer row.Close()

	if row.Next() {
		var folderId, folderName, status string

		err := row.Scan(&folderId, &folderName, &status)

		if err != nil {
			return task, err
		}

		task = Task{
			ID:         id,
			FolderId:   folderId,
			FolderName: folderName,
			Status:     status,
		}
	}

	return task, nil
}

func (repository *Repository) Create(task Task) (*Task, error) {
	query := fmt.Sprintf("INSERT INTO tasks(FolderId, FolderName, Status) VALUES('%s', '%s', '%s');", task.FolderId, task.FolderName, Created.String())

	result, err := repository.DB.Exec(query)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return nil, err
	}

	created, err := repository.Get(id)

	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (repository *Repository) Update(task Task) error {
	query := fmt.Sprintf("UPDATE tasks SET FolderId=%s, FolderName=%s, Status=%s WHERE ID=%d;", task.FolderId, task.FolderName, task.Status, task.ID)

	_, err := repository.DB.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

func (repository *Repository) Delete(id int64) error {
	query := fmt.Sprintf("DELETE FROM tasks WHERE ID=%d;", id)

	_, err := repository.DB.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}
