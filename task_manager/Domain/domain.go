package domain

import (
	"context"
	"time"
)

const (
	TaskCollection = "tasks"
	UserCollection = "users"
)

// MODELS
type Task struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string `json:"status"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role   string `json:"role"` // e.g., "admin", "user"
}

// REPOSITORIES
type TaskRepository interface {
	GetAllTasks(c context.Context, id string) ([]*Task, error)
	GetTaskByID(c context.Context, id string) (*Task, error)
	CreateTask(c context.Context, task *Task) error
	UpdateTask(c context.Context, id string, task *Task) (*Task, error)
	DeleteTask(c context.Context, id string) error
}
type UserRepository interface {
	GetAllUsers(c context.Context, user *User) error
	GetUserByID(c context.Context, id string) (*User, error)
	GetUserByEmail(c context.Context, email string) (*User, error)
	GetUserByUsername(c context.Context, username string) (*User, error)
	CreateUser(c context.Context, user *User) (*User, error)
	PromoteUserToAdmin(c context.Context, id string) error
	UserExists(c context.Context) (bool, error)
}

// USECASES
type TaskUsecases interface {
	GetAllTasks(ctx context.Context, id string) ([]*Task, error)
	GetTaskByID(ctx context.Context, id string) (*Task, error)
	CreateTask(ctx context.Context, task *Task, userId string) error
	UpdateTask(ctx context.Context, id string, task *Task) (*Task, error)
	DeleteTask(ctx context.Context, id string) error
}
type UserUsecases interface {
	GetAllUsers(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
	PromoteUserToAdmin(ctx context.Context, id string) error
}