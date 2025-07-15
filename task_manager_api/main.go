package main

import "fmt"
import "github.com/gin-gonic/gin"
import "time"
import "net/http"


func main() {
	router := gin.Default()

	router.GET("/ping", ping)
	router.GET("/tasks/:id", getTask)
	router.GET("/tasks", getTasks)
	router.PUT("/tasks/:id", updatedTask)
	router.DELETE("/tasks/:id", removeTask)
	router.POST("/tasks", addTask)

	fmt.Println("Task Manager API")
	router.Run()
}

func ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func getTasks(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func getTask(ctx *gin.Context){
	id := ctx.Param("id")

	for _, val := range tasks {
		if val.ID == id {
			ctx.JSON(http.StatusOK, val)
			return
		}
	}
}

func removeTask(ctx *gin.Context){
	id := ctx.Param("id")

	for i, val := range tasks {
		if val.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]... )
			ctx.JSON(http.StatusOK, gin.H{"message": "Task removed successfully"})
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

func updatedTask(ctx *gin.Context){
	id := ctx.Param("id")

	var updatedTask Task

	if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			if updatedTask.Title != ""{
				tasks[i].Title = updatedTask.Title
			}
			if updatedTask.Description != "" {
				tasks[i].Description = updatedTask.Description
			}
			if updatedTask.Status != "" {
				tasks[i].Status = updatedTask.Status
			}
			if !updatedTask.DueDate.IsZero() {
				tasks[i].DueDate = updatedTask.DueDate
			}
			ctx.JSON(http.StatusOK, gin.H{"message": "Task updated successfully", "task": tasks[i]})
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

func addTask(ctx *gin.Context) {
	var newTask Task

	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasks = append(tasks, newTask)
	ctx.JSON(http.StatusCreated, gin.H{"message": "Task created successfully"})
}

type Task struct {
	ID 		string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}

var tasks = []Task{
	{ID: "1", Title: "Task Manager Project", Description: "Add/View/Delete Tasks", DueDate: time.Now(), Status: "In Progress"},
	{ID: "2", Title: "Books Management Project", Description: "Add/View/Delete Books", DueDate: time.Now().AddDate(0, 0, -1), Status: "Completed"},
}