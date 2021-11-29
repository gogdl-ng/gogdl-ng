package task

type Task struct {
	ID         int64  `json:"id"`
	FolderId   string `json:"folderId"`
	FolderName string `json:"folderName"`
	Status     string `json:"status"`
}

type TaskStatus int64

const (
	New TaskStatus = iota
	Processing
	Done
)

func (ts TaskStatus) String() string {
	switch ts {
	case New:
		return "new"
	case Processing:
		return "processing"
	case Done:
		return "done"
	}

	return "unknown"
}

func (s *SQLite) GetAll() ([]Task, error) {
	var tasks []Task

	var query = `SELECT ID, FolderId, FolderName, Status FROM tasks`

	rows, err := s.DB.Query(query)

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

func (s *SQLite) Get(id int64) (Task, error) {
	var task Task

	query := `SELECT FolderId, FolderName, Status FROM tasks WHERE ID=$1`

	row, err := s.DB.Query(query, id)

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

func (s *SQLite) Create(task Task) (*Task, error) {
	query := `INSERT INTO tasks(FolderId, FolderName, Status) VALUES($1, $2, $3);`

	result, err := s.DB.Exec(query, task.FolderId, task.FolderName, New.String())

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

func (s *SQLite) Update(task Task) error {
	query := `UPDATE tasks SET FolderId=$1, FolderName=$2, Status=$3 WHERE ID=$4;`

	_, err := s.DB.Exec(query, task.FolderId, task.FolderName, task.Status, task.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *SQLite) Delete(id int64) error {
	query := `DELETE FROM tasks WHERE ID=$1;`

	_, err := s.DB.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}
