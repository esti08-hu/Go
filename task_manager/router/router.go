package router

import (
	"task_manager/controller"
	"task_manager/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Public routes (no authentication required)
	router.POST("/register", controller.Register)
	router.POST("/login", controller.Login)

	// Admin-only routes
	router.POST("/promote", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controller.PromoteUser)

	// Task routes with authentication
	tasks := router.Group("/tasks")
	{
		// Routes accessible to all authenticated users
		tasks.GET("/", middleware.AuthMiddleware(), controller.GetAllTasks)
		tasks.GET("/:id", middleware.AuthMiddleware(), controller.GetTask)

		// Admin-only task routes
		tasks.POST("/", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controller.AddTask)
		tasks.PUT("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controller.UpdatedTask)
		tasks.DELETE("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(), controller.RemoveTask)
	}

	return router
}
