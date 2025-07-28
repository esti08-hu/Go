package controller_test

import (
	"bytes"
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"task_manager/Delivery/controller"
	domain "task_manager/Domain"

	"task_manager/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TaskControllerSuite struct {
	suite.Suite
	taskUsecase *mocks.TaskUsecases
	userUsecase *mocks.UserUsecases
	controller  *controller.Controller
	router *gin.Engine
}

func (s *TaskControllerSuite) SetupTest() {
	s.taskUsecase = new(mocks.TaskUsecases)
	s.userUsecase = new(mocks.UserUsecases)
	s.controller = controller.NewController(s.taskUsecase, s.userUsecase)
	s.router = gin.Default()

	s.router.GET("/tasks", s.controller.GetAllTasks)
	s.router.GET("/task/:id", s.controller.GetTask)
	s.router.POST("/task", s.controller.AddTask)
	s.router.PUT("/task/:id", s.controller.UpdatedTask)
	s.router.DELETE("/task/:id", s.controller.RemoveTask)
}

func (s *TaskControllerSuite) TestGetAllTasks_Success() {
	assert := assert.New(s.T())
	user := &domain.User{ID: "123"}
	tasks := []*domain.Task{{ID: "t1", Title: "Test Task", UserID: "123"}}

	s.userUsecase.On("GetCurrentUser", mock.Anything).Return(user, nil)
	s.taskUsecase.On("GetAllTasks", mock.Anything, "123").Return(tasks, nil)

	req, _ := http.NewRequest("GET", "/tasks", nil)
	res := httptest.NewRecorder()
	s.router.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)
	assert.Contains(res.Body.String(), "Test Task")
}

func (s *TaskControllerSuite) TestGetTask_Success() {
	assert := assert.New(s.T())
	task := &domain.Task{ID: "t1", Title: "Test Task"}
	s.taskUsecase.On("GetTaskByID", mock.Anything, "t1").Return(task, nil)

	req, _ := http.NewRequest("GET", "/task/t1", nil)
	res := httptest.NewRecorder()
	s.router.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)
	assert.Contains(res.Body.String(), "Test Task")
}

func (s *TaskControllerSuite) TestAddTask_Success() {
	assert := assert.New(s.T())
	user := &domain.User{ID: "123"}
	newTask := &domain.Task{Title: "New Task"}

	s.userUsecase.On("GetCurrentUser", mock.Anything).Return(user, nil)
	s.taskUsecase.On("CreateTask", mock.Anything, mock.Anything, "123").Return(nil)

	body, _ := json.Marshal(newTask)
	req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	s.router.ServeHTTP(res, req)

	assert.Equal(http.StatusCreated, res.Code)
	assert.Contains(res.Body.String(), "Task added successfully")
}

func (s *TaskControllerSuite) TestUpdatedTask_NotFound() {
	assert := assert.New(s.T())
	s.taskUsecase.On("GetTaskByID", mock.Anything, "t1").Return(nil, nil)

	body, _ := json.Marshal(&domain.Task{Title: "Updated"})
	req, _ := http.NewRequest("PUT", "/task/t1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	s.router.ServeHTTP(res, req)

	assert.Equal(http.StatusNotFound, res.Code)
	assert.Contains(res.Body.String(), "Task not found")
}

func (s *TaskControllerSuite) TestRemoveTask_Unauthorized() {
	assert := assert.New(s.T())
	s.userUsecase.On("GetCurrentUser", mock.Anything).Return(nil, nil)

	req, _ := http.NewRequest("DELETE", "/task/t1", nil)
	res := httptest.NewRecorder()
	s.router.ServeHTTP(res, req)

	assert.Equal(http.StatusUnauthorized, res.Code)
	assert.Contains(res.Body.String(), "User not authenticated")
}

func TestTaskControllerSuite(t *testing.T) {
	suite.Run(t, new(TaskControllerSuite))
}
