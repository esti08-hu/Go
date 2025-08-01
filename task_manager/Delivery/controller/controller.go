package controller

import (
	"context"
	"net/http"
	domain "task_manager/Domain"
	"time"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	TaskUsecases domain.TaskUsecases
	UserUsecases domain.UserUsecases
}

func NewController(tu domain.TaskUsecases, uu domain.UserUsecases) *Controller {
	return &Controller{
		TaskUsecases: tu,
		UserUsecases: uu,
	}
}

// Register handles user registration
func (cr *Controller) Register(ctx *gin.Context) {
	var user domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}
	// Validate required fields
	if user.Username == "" || user.Email == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username, email, and password are required"})
		return
	}

	// Create the user (role will be set automatically in CreateUser)
	createdUser, err := cr.UserUsecases.CreateUser(ctx, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	response := gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":       createdUser.ID,
			"username": createdUser.Username,
			"email":    createdUser.Email,
			"role":     createdUser.Role,
		},
	}
	
	ctx.JSON(http.StatusCreated, response)
	
}
// Login handles user authentication
func (cr *Controller) Login(ctx *gin.Context) {
	var loginRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	token, _ := cr.UserUsecases.Login(ctx, loginRequest.Email, loginRequest.Password)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

// PromoteUser allows admins to promote other users to admin role
func (cr *Controller) PromoteUser(ctx *gin.Context) {
	var promoteRequest struct {
		UserID string `json:"user_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&promoteRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}
	// Get the user to be promoted
	user, err := cr.UserUsecases.GetUserByID(ctx, promoteRequest.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Role == "admin" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User is already an admin"})
		return	
	}
	
	// Update user role to admin
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = cr.UserUsecases.PromoteUserToAdmin(dbCtx, promoteRequest.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to promote user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User promoted to admin successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     "admin",
		},
	})
}

// Task Handlers (Updated with admin checks)
func (cr *Controller) GetAllTasks(ctx *gin.Context) {
	// Get user from context (set by AuthMiddlewarre)
	// user, _ := infrastructure.GetUserFromContext(ctx)
	user, _ := cr.UserUsecases.GetCurrentUser(ctx)

	tasks, err := cr.TaskUsecases.GetAllTasks(ctx, user.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
	})
}

func (cr *Controller) GetTask(ctx *gin.Context) {

	id := ctx.Param("id")
	task, err := cr.TaskUsecases.GetTaskByID(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		return
	}

	if task == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"task": task,
	})
}

func (cr *Controller) RemoveTask(ctx *gin.Context) {
	// Get user from context (set by AuthMiddleware)
	user, _ := cr.UserUsecases.GetCurrentUser(ctx)

	if user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := ctx.Param("id")
	 task, err := cr.TaskUsecases.GetTaskByID(ctx, id); 
	 if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		return
	 }
	 
	 if task != nil {
		cr.TaskUsecases.DeleteTask(ctx, id)
		ctx.JSON(http.StatusOK, gin.H{"message": "Task removed successfully"})
		return
	}
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

func (cr *Controller) UpdatedTask(ctx *gin.Context) {

	id := ctx.Param("id")
	var updatedTask *domain.Task

	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := cr.TaskUsecases.GetTaskByID(ctx, id); 
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		return
	}
	if task != nil {
		updatedTask, err := cr.TaskUsecases.UpdateTask(ctx, id, updatedTask)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Task updated successfully", "task": updatedTask})
		return
	}
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

func (cr *Controller) AddTask(ctx *gin.Context) {
	// Get user from context (set by AuthMiddleware)
	user, _ := cr.UserUsecases.GetCurrentUser(ctx)

	var newTask *domain.Task
	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add the new task
	err := cr.TaskUsecases.CreateTask(ctx, newTask, user.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add task"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Task added successfully"})
}
