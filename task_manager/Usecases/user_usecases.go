package usecases

import (
	"context"
	"time"

	domain "task_manager/Domain"

	"github.com/google/uuid"
)

type userUsecases struct {
	userRepository domain.UserRepository
	passwordService domain.IPasswordService
	jwtService domain.IJWTService
	contextTimeout time.Duration
}

func NewUserUsecases(userRepository domain.UserRepository, ps domain.IPasswordService, js domain.IJWTService, contextTimeout time.Duration) domain.UserUsecases {
	return &userUsecases{
		userRepository: userRepository,
		passwordService: ps,
		jwtService: js,
		contextTimeout: contextTimeout,
	}
}

func (uu *userUsecases) Login(ctx context.Context, email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	user, err := uu.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if !uu.passwordService.VerifyPassword(user, password) {
		return "", err
	}

	// Generate JWT token
	token, err := uu.jwtService.GenerateToken(user)

	if err != nil {
		return "", err
	}
	return token, nil
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

	existingUser, err := uu.userRepository.GetUserByEmail(ctx, user.Email)
	if err != nil && err != domain.ErrUserNotFound {
		return nil, err
	}
	if existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}
	existingUser, err = uu.userRepository.GetUserByUsername(ctx, user.Username)
	if err != nil && err != domain.ErrUserNotFound {
		return nil, err
	}
	if existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	} 

	// Generate a new UUID for the user ID
	user.ID = uuid.New().String() 

	hashedPassword, err := uu.passwordService.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword
	
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

func (uu *userUsecases) GetCurrentUser(ctx context.Context) (*domain.User, error) {
	user, ok := ctx.Value("user").(*domain.User)
	if !ok || user == nil {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}
