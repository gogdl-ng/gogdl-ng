package task

type Store interface {
	GetAll() ([]Task, error)
	Get(id int64) (Task, error)
	Create(task Task) (*Task, error)
	Update(task Task) error
	Delete(id int64) error
}
