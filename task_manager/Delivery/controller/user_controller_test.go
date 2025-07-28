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

type ControllerSuite struct {
	suite.Suite
	userUsecase *mocks.UserUsecases
	controller  *controller.Controller
	router *gin.Engine
}

func (s *ControllerSuite) SetupTest() {
	s.userUsecase = new(mocks.UserUsecases)
	s.controller = controller.NewController(nil, s.userUsecase)
	s.router = gin.Default()
	s.router.POST("/register", s.controller.Register)
	s.router.POST("/login", s.controller.Login)
	s.router.POST("/promote", s.controller.PromoteUser)
}

func (s *ControllerSuite) TestRegister_ValidInput() {
	assert := assert.New(s.T())
	input := domain.User{Username: "john", Email: "john@example.com", Password: "secret"}
	expected := &domain.User{ID: "user-id", Username: "john", Email: "john@example.com", Role: "user"}
	s.userUsecase.On("CreateUser", mock.Anything, &input).Return(expected, nil)

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	s.router.ServeHTTP(res, req)

	assert.Equal(http.StatusCreated, res.Code)
	assert.Contains(res.Body.String(), "User registered successfully")
	s.userUsecase.AssertExpectations(s.T())
}

func (s *ControllerSuite) TestLogin_Success() {
	assert := assert.New(s.T())
	loginReq := map[string]string{"email": "john@example.com", "password": "secret"}
	s.userUsecase.On("Login", mock.Anything, "john@example.com", "secret").Return("mocked-token", nil)

	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	s.router.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)
	assert.Contains(res.Body.String(), "Login successful")
	assert.Contains(res.Body.String(), "mocked-token")
	s.userUsecase.AssertExpectations(s.T())
}

func (s *ControllerSuite) TestLogin_InvalidPayload() {
	assert := assert.New(s.T())
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	s.router.ServeHTTP(res, req)

	assert.Equal(http.StatusBadRequest, res.Code)
	assert.Contains(res.Body.String(), "Email and password are required")
}

func (s *ControllerSuite) TestPromoteUser_Success() {
	assert := assert.New(s.T())
	user := &domain.User{ID: "123", Username: "john", Email: "john@example.com", Role: "user"}
	promoteReq := map[string]string{"user_id": "123"}

	s.userUsecase.On("GetUserByID", mock.Anything, "123").Return(user, nil)
	s.userUsecase.On("PromoteUserToAdmin", mock.Anything, "123").Return(nil)

	body, _ := json.Marshal(promoteReq)
	req, _ := http.NewRequest("POST", "/promote", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	s.router.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)
	assert.Contains(res.Body.String(), "User promoted to admin successfully")
	s.userUsecase.AssertExpectations(s.T())
}

func (s *ControllerSuite) TestPromoteUser_AlreadyAdmin() {
	assert := assert.New(s.T())
	user := &domain.User{ID: "123", Role: "admin"}
	s.userUsecase.On("GetUserByID", mock.Anything, "123").Return(user, nil)

	body, _ := json.Marshal(map[string]string{"user_id": "123"})
	req, _ := http.NewRequest("POST", "/promote", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	s.router.ServeHTTP(res, req)

	assert.Equal(http.StatusBadRequest, res.Code)
	assert.Contains(res.Body.String(), "User is already an admin")
}

func (s *ControllerSuite) TestPromoteUser_NotFound() {
	assert := assert.New(s.T())
	s.userUsecase.On("GetUserByID", mock.Anything, "123").Return(nil, nil)

	body, _ := json.Marshal(map[string]string{"user_id": "123"})
	req, _ := http.NewRequest("POST", "/promote", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	s.router.ServeHTTP(res, req)

	assert.Equal(http.StatusNotFound, res.Code)
	assert.Contains(res.Body.String(), "User not found")
}

func TestControllerSuite(t *testing.T) {
	suite.Run(t, new(ControllerSuite))
}
