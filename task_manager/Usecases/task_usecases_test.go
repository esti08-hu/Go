package usecases_test

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "task_manager/Domain"
	"task_manager/mocks"
	taskUsecases "task_manager/Usecases"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TaskUsecaseSuite struct {
	suite.Suite
	taskRepo *mocks.TaskRepository
	timeout  time.Duration
	taskUC   domain.TaskUsecases
}

func (s *TaskUsecaseSuite) SetupTest() {
	s.taskRepo = new(mocks.TaskRepository)
	s.timeout = time.Second * 2
	s.taskUC = taskUsecases.NewTaskUsecases(s.taskRepo, s.timeout)
}

func TestTaskUsecaseSuite(t *testing.T) {
	suite.Run(t, new(TaskUsecaseSuite))
}

func (s *TaskUsecaseSuite) TestGetAllTasks_Success() {
	assert := assert.New(s.T())
	tasks := []*domain.Task{{ID: "1", Title: "Task 1"}}

	s.taskRepo.On("GetAllTasks", mock.Anything, "user-id").Return(tasks, nil).Once()
	result, err := s.taskUC.GetAllTasks(context.Background(), "user-id")

	assert.NoError(err)
	assert.Len(result, 1)
	s.taskRepo.AssertExpectations(s.T())
}

func (s *TaskUsecaseSuite) TestGetAllTasks_NoTasks() {
	assert := assert.New(s.T())
	s.taskRepo.On("GetAllTasks", mock.Anything, "user-id").Return([]*domain.Task{}, nil).Once()
	result, err := s.taskUC.GetAllTasks(context.Background(), "user-id")

	assert.Error(err)
	assert.Nil(result)
	s.taskRepo.AssertExpectations(s.T())
}

func (s *TaskUsecaseSuite) TestGetTaskByID_Success() {
	assert := assert.New(s.T())
	task := &domain.Task{ID: "task-id", Title: "Task Title"}
	s.taskRepo.On("GetTaskByID", mock.Anything, "task-id").Return(task, nil).Once()

	result, err := s.taskUC.GetTaskByID(context.Background(), "task-id")

	assert.NoError(err)
	assert.Equal("task-id", result.ID)
	s.taskRepo.AssertExpectations(s.T())
}

func (s *TaskUsecaseSuite) TestCreateTask_Success() {
	assert := assert.New(s.T())
	task := &domain.Task{Title: "Create Me"}

	s.taskRepo.On("CreateTask", mock.Anything, mock.MatchedBy(func(t *domain.Task) bool {
		return t.ID != "" && t.UserID == "user-id"
	})).Return(nil).Once()

	err := s.taskUC.CreateTask(context.Background(), task, "user-id")

	assert.NoError(err)
	s.taskRepo.AssertExpectations(s.T())
}

func (s *TaskUsecaseSuite) TestUpdateTask_Success() {
	assert := assert.New(s.T())
	updated := &domain.Task{ID: "1", Title: "Updated Title"}

	s.taskRepo.On("UpdateTask", mock.Anything, "1", updated).Return(updated, nil).Once()

	result, err := s.taskUC.UpdateTask(context.Background(), "1", updated)

	assert.NoError(err)
	assert.Equal("Updated Title", result.Title)
	s.taskRepo.AssertExpectations(s.T())
}

func (s *TaskUsecaseSuite) TestDeleteTask_Success() {
	assert := assert.New(s.T())
	s.taskRepo.On("DeleteTask", mock.Anything, "1").Return(nil).Once()

	err := s.taskUC.DeleteTask(context.Background(), "1")

	assert.NoError(err)
	s.taskRepo.AssertExpectations(s.T())
}

func (s *TaskUsecaseSuite) TestDeleteTask_Error() {
	assert := assert.New(s.T())
	s.taskRepo.On("DeleteTask", mock.Anything, "1").Return(errors.New("delete failed")).Once()

	err := s.taskUC.DeleteTask(context.Background(), "1")

	assert.Error(err)
	assert.EqualError(err, "delete failed")
	s.taskRepo.AssertExpectations(s.T())
}
