package repository_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	domain "task_manager/Domain"
	repository "task_manager/Repository"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const testTaskCollection = "test_tasks"

type taskRepositoryTestSuite struct {
	suite.Suite
	db       *mongo.Database
	taskRepo domain.TaskRepository
	ctx      context.Context
	cancel   context.CancelFunc
	client   *mongo.Client
}

func TestTaskRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(taskRepositoryTestSuite))
}

func (s *taskRepositoryTestSuite) SetupSuite() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	testMongoURL := os.Getenv("DATABASE_URL")
	if testMongoURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(testMongoURL))
	s.Require().NoError(err)

	s.client = client
	s.db = client.Database("test_task_db")
	s.taskRepo = repository.NewTaskRepository(s.db, testTaskCollection)
	s.ctx, s.cancel = context.WithTimeout(context.Background(), 10*time.Second)
}

func (s *taskRepositoryTestSuite) TearDownSuite() {
	s.db.Collection(testUserCollection).Drop(s.ctx)
	s.cancel()
	s.client.Disconnect(s.ctx)
}
func (s *taskRepositoryTestSuite) SetupTest() {
	// Clean collection before each test
	_, err := s.db.Collection(testTaskCollection).DeleteMany(s.ctx, bson.M{})
	s.Require().NoError(err)
}

func (s *taskRepositoryTestSuite) TestCreateAndGetTaskByID() {
	assert := assert.New(s.T())

	task := &domain.Task{
		ID:          "task-1",
		UserID:      "user-1",
		Title:       "Test Task",
		Description: "This is a test",
		Status:      "pending",
	}

	err := s.taskRepo.CreateTask(s.ctx, task)
	assert.NoError(err)

	found, err := s.taskRepo.GetTaskByID(s.ctx, "task-1")
	assert.NoError(err)
	assert.Equal(task.ID, found.ID)
	assert.Equal(task.Title, found.Title)
}

func (s *taskRepositoryTestSuite) TestGetAllTasks() {
	assert := assert.New(s.T())

	task1 := &domain.Task{ID: "1", UserID: "user-1", Title: "Task 1"}
	task2 := &domain.Task{ID: "2", UserID: "user-1", Title: "Task 2"}

	_ = s.taskRepo.CreateTask(s.ctx, task1)
	_ = s.taskRepo.CreateTask(s.ctx, task2)

	tasks, err := s.taskRepo.GetAllTasks(s.ctx, "user-1")
	assert.NoError(err)
	assert.Len(tasks, 2)
}

func (s *taskRepositoryTestSuite) TestUpdateTask() {
	assert := assert.New(s.T())

	task := &domain.Task{ID: "task-update", UserID: "user-x", Title: "Old"}
	_ = s.taskRepo.CreateTask(s.ctx, task)

	task.Title = "Updated"
	updated, err := s.taskRepo.UpdateTask(s.ctx, task.ID, task)
	assert.NoError(err)
	assert.Equal("Updated", updated.Title)

	found, _ := s.taskRepo.GetTaskByID(s.ctx, task.ID)
	assert.Equal("Updated", found.Title)
}

func (s *taskRepositoryTestSuite) TestDeleteTask() {
	assert := assert.New(s.T())

	task := &domain.Task{ID: "delete-me", UserID: "user-x", Title: "To Delete"}
	_ = s.taskRepo.CreateTask(s.ctx, task)

	err := s.taskRepo.DeleteTask(s.ctx, task.ID)
	assert.NoError(err)

	result, err := s.taskRepo.GetTaskByID(s.ctx, task.ID)
	assert.Nil(result)
	assert.Error(err)
}
