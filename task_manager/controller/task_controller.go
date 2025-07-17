package controller

import (
	"fmt"
	"net/http"
	"task_manager/data"
	"task_manager/models"

	"github.com/gin-gonic/gin"
)

func GetAllTasks(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"tasks": data.GetTasks()})
}

func GetTask(ctx *gin.Context) {
	id := ctx.Param("id")

	task := data.GetTaskById(id)
	if task == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"task": task})
}

func RemoveTask(ctx *gin.Context) {
	id := ctx.Param("id")

	if task := data.GetTaskById(id); task != nil {
		data.RemoveTask(id)
		ctx.JSON(http.StatusOK, gin.H{"message": "Task removed successfully"})
		return
	}
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

func UpdatedTask(ctx *gin.Context) {
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
	var newTask models.Task
	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newTask.ID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}
	// Check if the task already exists
	if existingTask := data.GetTaskById(newTask.ID); existingTask != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Task with this ID already exists"})
		return
	}
	// Add the new task
	newTask, err := data.AddTask(newTask)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add task"})
		return
	}
	fmt.Println("New Task:", newTask)
	ctx.JSON(http.StatusCreated, gin.H{"message": "Task added successfully", "task": newTask})
}
