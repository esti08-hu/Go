package controller

import (
	"context"
	"fmt"
	"net/http"
	"task_manager/data"
	"task_manager/middleware"
	"task_manager/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// Register handles user registration
func Register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Validate required fields
	if user.Username == "" || user.Email == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username, email, and password are required"})
		return
	}

	// Check if user already exists by email
	existingUser, err := data.GetUserByEmail(user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if existingUser != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	// Check if user already exists by username
	existingUser, err = data.GetUserByUsername(user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if existingUser != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
		return
	}

	// Create the user (role will be set automatically in CreateUser)
	createdUser, err := data.CreateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Don't return the password in the response
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
func Login(ctx *gin.Context) {
	var loginRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	// Get user by email
	user, err := data.GetUserByEmail(loginRequest.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Verify password
	if !data.VerifyPassword(user, loginRequest.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := data.GenerateToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// PromoteUser allows admins to promote other users to admin role
func PromoteUser(ctx *gin.Context) {
	var promoteRequest struct {
		UserID string `json:"user_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&promoteRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}
	// Get the user to be promoted
	user, err := data.GetUserById(promoteRequest.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if user is already an admin
	if user.Role == "admin" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User is already an admin"})
		return
	}

	// Update user role to admin
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": promoteRequest.UserID}
	update := bson.M{"$set": bson.M{"role": "admin"}}

	_, err = data.UserCollection.UpdateOne(dbCtx, filter, update)
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
func GetAllTasks(ctx *gin.Context) {
	// Get user from context (set by AuthMiddleware)
	user, exists := middleware.GetUserFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	tasks := data.GetTasks()
	ctx.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

func GetTask(ctx *gin.Context) {
	// Get user from context (set by AuthMiddleware)
	user, exists := middleware.GetUserFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := ctx.Param("id")
	task := data.GetTaskById(id)
	if task == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"task": task,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

func RemoveTask(ctx *gin.Context) {
	// Get user from context (set by AuthMiddleware)
	_, exists := middleware.GetUserFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := ctx.Param("id")
	if task := data.GetTaskById(id); task != nil {
		data.RemoveTask(id)
		ctx.JSON(http.StatusOK, gin.H{"message": "Task removed successfully"})
		return
	}
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

func UpdatedTask(ctx *gin.Context) {
	// Get user from context (set by AuthMiddleware)
	user, exists := middleware.GetUserFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if user is admin
	if user.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin access required to update tasks"})
		return
	}

	id := ctx.Param("id")
	var updatedTask models.Task

	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if task := data.GetTaskById(id); task != nil {
		updatedTask, err := data.UpdatedTask(id, updatedTask)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Task updated successfully", "task": updatedTask})
		return
	}
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

func AddTask(ctx *gin.Context) {
	// Get user from context (set by AuthMiddleware)
	user, exists := middleware.GetUserFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if user is admin
	if user.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin access required to create tasks"})
		return
	}

	var newTask models.Task
	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate ID if not provided
	if newTask.ID == "" {
		newTask.ID = uuid.New().String()
	}

	// Check if the task already exists
	if existingTask := data.GetTaskById(newTask.ID); existingTask != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Task with this ID already exists"})
		return
	}

	// Add the new task
	createdTask, err := data.AddTask(newTask)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add task"})
		return
	}

	fmt.Println("New Task:", createdTask)
	ctx.JSON(http.StatusCreated, gin.H{"message": "Task added successfully", "task": createdTask})
}
