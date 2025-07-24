package usecases

import (
	"context"
	"time"

	domain "task_manager/Domain"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type taskUsecases struct {
	taskRepository domain.TaskRepository
	contextTimeout time.Duration
}

func NewTaskUsecases(taskRepository domain.TaskRepository, contextTimeout time.Duration) domain.TaskUsecases {
	return &taskUsecases{
		taskRepository: taskRepository,
		contextTimeout: contextTimeout,
	}
}

func (tu *taskUsecases) GetAllTasks(ctx context.Context, id string) ([]*domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	tasks, err := tu.taskRepository.GetAllTasks(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, mongo.ErrNoDocuments
	}
	return tasks, nil
}

func (tu *taskUsecases) GetTaskByID(ctx context.Context, id string) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	return tu.taskRepository.GetTaskByID(ctx, id)
}

func (tu *taskUsecases) CreateTask(ctx context.Context, newTask *domain.Task, user_id string) error {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	newTask.ID = uuid.New().String()
	newTask.UserID = user_id
	
	return tu.taskRepository.CreateTask(ctx, newTask)
}

func (tu *taskUsecases) UpdateTask(ctx context.Context, id string, task *domain.Task) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	return tu.taskRepository.UpdateTask(ctx, id, task)
}

func (tu *taskUsecases) DeleteTask(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	return tu.taskRepository.DeleteTask(ctx, id)
}
