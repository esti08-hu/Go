package usecases_test

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "task_manager/Domain"
	"task_manager/mocks"
	userUsecases "task_manager/Usecases"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserUsecaseSuite struct {
	suite.Suite
	repo     *mocks.UserRepository
	ps       *mocks.IPasswordService
	jwt      *mocks.IJWTService
	uc       domain.UserUsecases
	timeout  time.Duration
}

func (s *UserUsecaseSuite) SetupTest() {
	s.repo = new(mocks.UserRepository)
	s.ps = new(mocks.IPasswordService)
	s.jwt = new(mocks.IJWTService)
	s.uc = userUsecases.NewUserUsecases(s.repo, s.ps, s.jwt, s.timeout)
}

func TestUserUsecaseSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseSuite))
}

// Login
func (s *UserUsecaseSuite) TestLogin_Success() {
	assert := assert.New(s.T())
	ctx := context.Background()

	user := &domain.User{ID: "u1", Email: "john@example.com", Password: "hashed"}

	s.repo.On("GetUserByEmail", mock.Anything, "john@example.com").Return(user, nil).Once()
	s.ps.On("VerifyPassword", user, "secret").Return(true).Once()
	s.jwt.On("GenerateToken", user).Return("jwt-token", nil).Once()

	token, err := s.uc.Login(ctx, "john@example.com", "secret")
	assert.NoError(err)
	assert.Equal("jwt-token", token)

	s.repo.AssertExpectations(s.T())
	s.ps.AssertExpectations(s.T())
	s.jwt.AssertExpectations(s.T())
}

func (s *UserUsecaseSuite) TestLogin_RepoError() {
	assert := assert.New(s.T())
	ctx := context.Background()

	s.repo.On("GetUserByEmail", mock.Anything, "john@example.com").Return(nil, errors.New("db error")).Once()

	token, err := s.uc.Login(ctx, "john@example.com", "secret")
	assert.Error(err)
	assert.Empty(token)

	s.repo.AssertExpectations(s.T())
}

func (s *UserUsecaseSuite) TestLogin_WrongPassword() {
	assert := assert.New(s.T())
	ctx := context.Background()

	user := &domain.User{ID: "u1", Email: "john@example.com", Password: "hashed"}
	s.repo.On("GetUserByEmail", mock.Anything, "john@example.com").Return(user, nil).Once()
	s.ps.On("VerifyPassword", user, "wrong").Return(false).Once()

	token, err := s.uc.Login(ctx, "john@example.com", "wrong")
	assert.NoError(err) // because code returns the existing err (nil)
	assert.Empty(token)

	s.repo.AssertExpectations(s.T())
	s.ps.AssertExpectations(s.T())
}

func (s *UserUsecaseSuite) TestLogin_JWTGenerationFails() {
	assert := assert.New(s.T())
	ctx := context.Background()

	user := &domain.User{ID: "u1", Email: "john@example.com", Password: "hashed"}
	s.repo.On("GetUserByEmail", mock.Anything, "john@example.com").Return(user, nil).Once()
	s.ps.On("VerifyPassword", user, "secret").Return(true).Once()
	s.jwt.On("GenerateToken", user).Return("", errors.New("jwt error")).Once()

	token, err := s.uc.Login(ctx, "john@example.com", "secret")
	assert.Error(err)
	assert.Empty(token)

	s.repo.AssertExpectations(s.T())
	s.ps.AssertExpectations(s.T())
	s.jwt.AssertExpectations(s.T())
}

// CreateUser
func (s *UserUsecaseSuite) TestCreateUser_FirstUserGetsAdminRole() {
	assert := assert.New(s.T())
	ctx := context.Background()

	in := &domain.User{Username: "john", Email: "john@example.com", Password: "plain"}

	// Email & username do not exist
	s.repo.On("GetUserByEmail", mock.Anything, in.Email).Return(nil, domain.ErrUserNotFound).Once()
	s.repo.On("GetUserByUsername", mock.Anything, in.Username).Return(nil, domain.ErrUserNotFound).Once()

	// Hash password
	s.ps.On("HashPassword", "plain").Return("hashed", nil).Once()

	// First user -> no users exist yet
	s.repo.On("UserExists", mock.Anything).Return(false, nil).Once()

	// Finally create
	s.repo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.Email == in.Email && u.Role == "admin" && u.Password == "hashed" && u.ID != ""
	})).Return(func(_ context.Context, u *domain.User) *domain.User { return u }, nil).Once()

	created, err := s.uc.CreateUser(ctx, in)
	assert.NoError(err)
	assert.Equal("admin", created.Role)
	assert.Equal("hashed", created.Password)
	assert.NotEmpty(created.ID)

	s.repo.AssertExpectations(s.T())
	s.ps.AssertExpectations(s.T())
}

func (s *UserUsecaseSuite) TestCreateUser_NextUsersGetUserRole() {
	assert := assert.New(s.T())
	ctx := context.Background()

	in := &domain.User{Username: "john", Email: "john@example.com", Password: "plain"}

	s.repo.On("GetUserByEmail", mock.Anything, in.Email).Return(nil, domain.ErrUserNotFound).Once()
	s.repo.On("GetUserByUsername", mock.Anything, in.Username).Return(nil, domain.ErrUserNotFound).Once()
	s.ps.On("HashPassword", "plain").Return("hashed", nil).Once()
	s.repo.On("UserExists", mock.Anything).Return(true, nil).Once()
	s.repo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.Role == "user"
	})).Return(func(_ context.Context, u *domain.User) *domain.User { return u }, nil).Once()

	created, err := s.uc.CreateUser(ctx, in)
	assert.NoError(err)
	assert.Equal("user", created.Role)

	s.repo.AssertExpectations(s.T())
	s.ps.AssertExpectations(s.T())
}

func (s *UserUsecaseSuite) TestCreateUser_EmailAlreadyExists() {
	assert := assert.New(s.T())
	ctx := context.Background()

	in := &domain.User{Username: "john", Email: "john@example.com", Password: "plain"}
	existing := &domain.User{ID: "u1"}

	s.repo.On("GetUserByEmail", mock.Anything, in.Email).Return(existing, nil).Once()

	created, err := s.uc.CreateUser(ctx, in)
	assert.Error(err)
	assert.Nil(created)
	assert.Equal(domain.ErrUserAlreadyExists, err)

	s.repo.AssertExpectations(s.T())
}

func (s *UserUsecaseSuite) TestCreateUser_UsernameAlreadyExists() {
	assert := assert.New(s.T())
	ctx := context.Background()

	in := &domain.User{Username: "john", Email: "john@example.com", Password: "plain"}

	s.repo.On("GetUserByEmail", mock.Anything, in.Email).Return(nil, domain.ErrUserNotFound).Once()
	s.repo.On("GetUserByUsername", mock.Anything, in.Username).Return(&domain.User{ID: "u2"}, nil).Once()

	created, err := s.uc.CreateUser(ctx, in)
	assert.Error(err)
	assert.Nil(created)
	assert.Equal(domain.ErrUserAlreadyExists, err)

	s.repo.AssertExpectations(s.T())
}

func (s *UserUsecaseSuite) TestCreateUser_HashPasswordFails() {
	assert := assert.New(s.T())
	ctx := context.Background()

	in := &domain.User{Username: "john", Email: "john@example.com", Password: "plain"}

	s.repo.On("GetUserByEmail", mock.Anything, in.Email).Return(nil, domain.ErrUserNotFound).Once()
	s.repo.On("GetUserByUsername", mock.Anything, in.Username).Return(nil, domain.ErrUserNotFound).Once()
	s.ps.On("HashPassword", "plain").Return("", errors.New("hash error")).Once()

	created, err := s.uc.CreateUser(ctx, in)
	assert.Error(err)
	assert.Nil(created)

	s.repo.AssertExpectations(s.T())
	s.ps.AssertExpectations(s.T())
}

func (s *UserUsecaseSuite) TestCreateUser_UserExistsCheckFails() {
	assert := assert.New(s.T())
	ctx := context.Background()

	in := &domain.User{Username: "john", Email: "john@example.com", Password: "plain"}

	s.repo.On("GetUserByEmail", mock.Anything, in.Email).Return(nil, domain.ErrUserNotFound).Once()
	s.repo.On("GetUserByUsername", mock.Anything, in.Username).Return(nil, domain.ErrUserNotFound).Once()
	s.ps.On("HashPassword", "plain").Return("hashed", nil).Once()
	s.repo.On("UserExists", mock.Anything).Return(false, errors.New("db error")).Once()

	created, err := s.uc.CreateUser(ctx, in)
	assert.Error(err)
	assert.Nil(created)

	s.repo.AssertExpectations(s.T())
	s.ps.AssertExpectations(s.T())
}

func (s *UserUsecaseSuite) TestCreateUser_CreateRepoFails() {
	assert := assert.New(s.T())
	ctx := context.Background()

	in := &domain.User{Username: "john", Email: "john@example.com", Password: "plain"}

	s.repo.On("GetUserByEmail", mock.Anything, in.Email).Return(nil, domain.ErrUserNotFound).Once()
	s.repo.On("GetUserByUsername", mock.Anything, in.Username).Return(nil, domain.ErrUserNotFound).Once()
	s.ps.On("HashPassword", "plain").Return("hashed", nil).Once()
	s.repo.On("UserExists", mock.Anything).Return(true, nil).Once()
	s.repo.On("CreateUser", mock.Anything, mock.Anything).Return((*domain.User)(nil), errors.New("insert err")).Once()

	created, err := s.uc.CreateUser(ctx, in)
	assert.Error(err)
	assert.Nil(created)

	s.repo.AssertExpectations(s.T())
	s.ps.AssertExpectations(s.T())
}

// Simple pass-throughs
func (s *UserUsecaseSuite) TestGetUserByID_Success() {
	assert := assert.New(s.T())
	ctx := context.Background()

	expected := &domain.User{ID: "u1"}
	s.repo.On("GetUserByID", mock.Anything, "u1").Return(expected, nil).Once()

	user, err := s.uc.GetUserByID(ctx, "u1")
	assert.NoError(err)
	assert.Equal(expected, user)
	s.repo.AssertExpectations(s.T())
}

func (s *UserUsecaseSuite) TestGetUserByEmail_Success() {
	assert := assert.New(s.T())
	ctx := context.Background()

	expected := &domain.User{ID: "u1", Email: "john@example.com"}
	s.repo.On("GetUserByEmail", mock.Anything, "john@example.com").Return(expected, nil).Once()

	user, err := s.uc.GetUserByEmail(ctx, "john@example.com")
	assert.NoError(err)
	assert.Equal(expected, user)
	s.repo.AssertExpectations(s.T())
}

func (s *UserUsecaseSuite) TestGetUserByUsername_Success() {
	assert := assert.New(s.T())
	ctx := context.Background()

	expected := &domain.User{ID: "u1", Username: "john"}
	s.repo.On("GetUserByUsername", mock.Anything, "john").Return(expected, nil).Once()

	user, err := s.uc.GetUserByUsername(ctx, "john")
	assert.NoError(err)
	assert.Equal(expected, user)
	s.repo.AssertExpectations(s.T())
}

func (s *UserUsecaseSuite) TestPromoteUserToAdmin_Success() {
	assert := assert.New(s.T())
	ctx := context.Background()

	s.repo.On("PromoteUserToAdmin", mock.Anything, "u1").Return(nil).Once()

	err := s.uc.PromoteUserToAdmin(ctx, "u1")
	assert.NoError(err)
	s.repo.AssertExpectations(s.T())
}

// GetCurrentUser
func (s *UserUsecaseSuite) TestGetCurrentUser_Success() {
	assert := assert.New(s.T())
	u := &domain.User{ID: "u1"}
	ctx := context.WithValue(context.Background(), "user", u)

	got, err := s.uc.GetCurrentUser(ctx)
	assert.NoError(err)
	assert.Equal(u, got)
}

func (s *UserUsecaseSuite) TestGetCurrentUser_NotFound() {
	assert := assert.New(s.T())
	ctx := context.Background()

	user, err := s.uc.GetCurrentUser(ctx)
	assert.Error(err)
	assert.Nil(user)
	assert.Equal(domain.ErrUserNotFound, err)
}
