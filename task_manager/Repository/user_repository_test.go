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

const testUserCollection = "test_users"

type userRepositoryTestSuite struct {
	suite.Suite
	db     *mongo.Database
	client *mongo.Client
	repo   domain.UserRepository
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *userRepositoryTestSuite) SetupSuite() {
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
	log.Println("Connected to MongoDB for testing:", client)

	// Assign the client to the struct field
	s.client = client

	s.db = client.Database("test_task_db")
	s.repo = repository.NewUserRepository(s.db, testUserCollection)
	s.ctx, s.cancel = context.WithTimeout(context.Background(), 10*time.Second)
}

func (s *userRepositoryTestSuite) TearDownSuite() {
	s.db.Collection(testUserCollection).Drop(s.ctx)
	s.cancel()
	s.client.Disconnect(s.ctx)
}

func (s *userRepositoryTestSuite) SetupTest() {
	// Clean collection before each test
	s.db.Collection(testUserCollection).DeleteMany(s.ctx, bson.M{})

}

func (s *userRepositoryTestSuite) TestCreateAndGetUser() {
	user := &domain.User{
		ID:       "user-123",
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "hashed-password",
		Role:     "user",
	}

	created, err := s.repo.CreateUser(s.ctx, user)
	assert := assert.New(s.T())
	assert.NoError(err)
	assert.Equal("johndoe", created.Username)

	// Get by ID
	found, err := s.repo.GetUserByID(s.ctx, "user-123")
	assert.NoError(err)
	assert.Equal("john@example.com", found.Email)

	// Get by Email
	foundByEmail, err := s.repo.GetUserByEmail(s.ctx, "john@example.com")
	assert.NoError(err)
	assert.Equal("johndoe", foundByEmail.Username)

	// Get by Username
	foundByUsername, err := s.repo.GetUserByUsername(s.ctx, "johndoe")
	assert.NoError(err)
	assert.Equal("john@example.com", foundByUsername.Email)
}

func (s *userRepositoryTestSuite) TestPromoteUserToAdmin() {
	user := &domain.User{
		ID:       "user-456",
		Username: "janedoe",
		Email:    "jane@example.com",
		Password: "hashed-password",
		Role:     "user",
	}
	_, err := s.repo.CreateUser(s.ctx, user)
	s.Require().NoError(err)

	err = s.repo.PromoteUserToAdmin(s.ctx, "user-456")
	s.Require().NoError(err)

	updated, err := s.repo.GetUserByID(s.ctx, "user-456")
	s.Require().NoError(err)
	assert.Equal(s.T(), "admin", updated.Role)
}

func (s *userRepositoryTestSuite) TestUserExists() {
	exists, err := s.repo.UserExists(s.ctx)
	assert.NoError(s.T(), err)
	assert.False(s.T(), exists)

	_, err = s.repo.CreateUser(s.ctx, &domain.User{
		ID:       "id-1",
		Username: "existtest",
		Email:    "exist@test.com",
		Password: "pw",
		Role:     "user",
	})
	s.Require().NoError(err)

	exists, err = s.repo.UserExists(s.ctx)
	assert.NoError(s.T(), err)
	assert.True(s.T(), exists)
}

func (s *userRepositoryTestSuite) TestGetAllUsers() {
	// Add 2 users
	s.repo.CreateUser(s.ctx, &domain.User{ID: "1", Username: "a", Email: "a@a.com", Password: "pw", Role: "user"})
	s.repo.CreateUser(s.ctx, &domain.User{ID: "2", Username: "b", Email: "b@b.com", Password: "pw", Role: "admin"})

	users, err := s.repo.GetAllUsers(s.ctx, nil)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), users, 2)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(userRepositoryTestSuite))
}
