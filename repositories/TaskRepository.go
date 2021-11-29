package repositories

import (
	"fmt"

	"github.com/LegendaryB/gogdl-ng/interfaces"
	"github.com/LegendaryB/gogdl-ng/models"
)

type TaskRepository struct {
	TaskRepository interfaces.ITaskRepository
}

func (s *SQLite) GetAll() ([]models.Task, error) {
	var tasks []models.Task

	rows, err := s.DB.Query("SELECT ID, FolderId, FolderName, Status FROM tasks")

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

		task := models.Task{
			ID:         id,
			FolderId:   folderId,
			FolderName: folderName,
			Status:     status,
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *SQLite) Get(id int64) (models.Task, error) {
	var task models.Task

	query := fmt.Sprintf("SELECT FolderId, FolderName, Status FROM tasks WHERE ID=%d", id)

	row, err := s.DB.Query(query)

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

		task = models.Task{
			ID:         id,
			FolderId:   folderId,
			FolderName: folderName,
			Status:     status,
		}
	}

	return task, nil
}

func (s *SQLite) Create(task models.Task) (*models.Task, error) {
	query := fmt.Sprintf("INSERT INTO tasks(FolderId, FolderName, Status) VALUES('%s', '%s', '%s');", task.FolderId, task.FolderName, models.New.String())

	result, err := s.DB.Exec(query)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return nil, err
	}

	created, err := s.Get(id)

	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (s *SQLite) Update(task models.Task) error {
	query := fmt.Sprintf("UPDATE tasks SET FolderId=%s, FolderName=%s, Status=%s WHERE ID=%d;", task.FolderId, task.FolderName, task.Status, task.ID)

	_, err := s.DB.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

func (s *SQLite) Delete(id int64) error {
	query := fmt.Sprintf("DELETE FROM tasks WHERE ID=%d;", id)

	_, err := s.DB.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}
