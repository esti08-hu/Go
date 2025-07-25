package domain

import (
	"context"
	"errors"
	"time"
)

const (
	TaskCollection = "tasks"
	UserCollection = "users"
)

// MODELS
type Task struct {
	ID          string    
	UserID      string   
	Title       string 
	Description string
	DueDate     time.Time
	Status      string 
}

type User struct {
	ID       string
	Username string
	Email    string
	Password string
	Role   string
}

// REPOSITORIES
type TaskRepository interface {
	GetAllTasks(c context.Context, userId string) ([]*Task, error)
	GetTaskByID(c context.Context, taskId string) (*Task, error)
	CreateTask(c context.Context, task *Task) error
	UpdateTask(c context.Context, taskId string, task *Task) (*Task, error)
	DeleteTask(c context.Context, taskId string) error
}
type UserRepository interface {
	GetAllUsers(c context.Context, user *User) ([]*User, error)
	GetUserByID(c context.Context, userId string) (*User, error)
	GetUserByEmail(c context.Context, email string) (*User, error)
	GetUserByUsername(c context.Context, username string) (*User, error)
	CreateUser(c context.Context, user *User) (*User, error)
	PromoteUserToAdmin(c context.Context, userId string) error
	UserExists(c context.Context) (bool, error)
}

// USECASES
type TaskUsecases interface {
	GetAllTasks(ctx context.Context, userId string) ([]*Task, error)
	GetTaskByID(ctx context.Context, taskId string) (*Task, error)
	CreateTask(ctx context.Context, task *Task, userId string) error
	UpdateTask(ctx context.Context, taskId string, task *Task) (*Task, error)
	DeleteTask(ctx context.Context, taskId string) error
}
type UserUsecases interface {
	GetUserByID(ctx context.Context, userId string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
	PromoteUserToAdmin(ctx context.Context, userId string) error
	Login(ctx context.Context, email, password string) (string, error)
	GetCurrentUser(ctx context.Context) (*User, error)
}

type IPasswordService interface {
	HashPassword(passowrd string) (string, error)		
	VerifyPassword(user *User, password string) bool
}

type IJWTService interface {
	GenerateToken(user *User) (string, error)
}

var (
	ErrUserNotFound = errors.New("user not found")
	ErrTaskNotFound = errors.New("task not found")
	ErrTaskAlreadyExists = errors.New("task already exists")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized = errors.New("unauthorized")
)