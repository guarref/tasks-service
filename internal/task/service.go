package task

type TaskService struct {
	trepo TaskRepository
}

func NewTaskService(trepo TaskRepository) *TaskService {
	return &TaskService{trepo: trepo}
}

func (s *TaskService) CreateTask(task Task) (Task, error) {
	return s.trepo.CreateTask(task)
}

func (s *TaskService) GetAllTasks() ([]Task, error) {
	return s.trepo.GetAllTasks()
}

func (s *TaskService) GetTaskByID(id uint32) (Task, error) {
	return s.trepo.GetTaskByID(id)
}

func (s *TaskService) GetTasksByUserID(id uint32) ([]Task, error) {
	return s.trepo.GetTasksByUserID(id)
}

func (s *TaskService) UpdateTaskByID(id uint32, task Task) (Task, error) {
	return s.trepo.UpdateTaskByID(id, task)
}

func (s *TaskService) DeleteTaskByID(id uint32) error {
	return s.trepo.DeleteTaskByID(id)
}