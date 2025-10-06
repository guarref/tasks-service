package task

import (
	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(task Task) (Task, error)
	GetAllTasks() ([]Task, error)
	GetTaskByID(id uint32) (Task, error)
	GetTasksByUserID(id uint32) ([]Task, error)
	UpdateTaskByID(id uint32, newTask Task) (Task, error)
	DeleteTaskByID(id uint32) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *taskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(task Task) (Task, error) {
	result := r.db.Create(&task)
	if result.Error != nil {
		return Task{}, result.Error
	}

	return task, nil
}

func (r *taskRepository) GetAllTasks() ([]Task, error) {
	var tasks []Task
	err := r.db.Find(&tasks).Error

	return tasks, err
}

func (r *taskRepository) GetTaskByID(id uint32) (Task, error) {
	var curTask Task
	result := r.db.First(&curTask, id)
	if result.Error != nil {
		return Task{}, result.Error
	}

	return curTask, nil
}

func (r *taskRepository) GetTasksByUserID(id uint32) ([]Task, error) {
	var userTasks []Task
	if err := r.db.Where("user_id = ?", id).Find(&userTasks).Error; err != nil {
		return nil, err
	}
	return userTasks, nil
}

func (r *taskRepository) UpdateTaskByID(id uint32, newTask Task) (Task, error) {
	var task Task
	result := r.db.First(&task, id)
	if result.Error != nil {
		return Task{}, result.Error
	}

	newresult := r.db.Model(&task).Select("title").Updates(newTask)
	if newresult.Error != nil {
		return Task{}, newresult.Error
	}

	r.db.Model(&task).Find(&task, id)
	return task, nil
}

func (r *taskRepository) DeleteTaskByID(id uint32) error {
	var delTask Task
	result := r.db.First(&delTask, id)

	if result.Error != nil {
		return result.Error
	}

	r.db.Delete(&delTask)

	return nil
}
