package usecases

import (
	"context"
	domain "task_manager/Domain"
	"time"

	"github.com/google/uuid"
)

type userUsecases struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewUserUsecases(userRepository domain.UserRepository, contextTimeout time.Duration) domain.UserUsecases {
	return &userUsecases{
		userRepository: userRepository,
		contextTimeout: contextTimeout,
	}
}

func (uu *userUsecases) GetAllUsers(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	return uu.userRepository.GetAllUsers(ctx, user)
}

func (uu *userUsecases) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	return uu.userRepository.GetUserByID(ctx, id)
}


func (uu *userUsecases) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()
	return uu.userRepository.GetUserByEmail(ctx, email)
}

func (uu *userUsecases) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	return uu.userRepository.GetUserByUsername(ctx, username)
}

func (uu *userUsecases) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()
	user.ID = uuid.New().String() 
	exists, err := uu.userRepository.UserExists(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		user.Role = "admin"
	} else {
		user.Role = "user"
	}
	return uu.userRepository.CreateUser(ctx, user)
}

func (uu *userUsecases) PromoteUserToAdmin(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	return uu.userRepository.PromoteUserToAdmin(ctx, id)
}
