package router

import (
	"github.com/gin-gonic/gin"
	"task_manager/controller"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	tasks:= router.Group("/tasks")

	{
		tasks.GET("/", controller.GetAllTasks)
		tasks.GET("/:id", controller.GetTask)
		tasks.POST("/", controller.AddTask)
		tasks.PUT("/:id", controller.UpdatedTask)
		tasks.DELETE("/:id", controller.RemoveTask)
	}
	return router
}
