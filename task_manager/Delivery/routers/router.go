package router

import (
	controller "task_manager/Delivery/controller"
	infrastructure "task_manager/Infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRouter(engine *gin.Engine, ctrl controller.Controller)  {
	public := engine.Group("")

	// Public routes (no authentication required)
	public.POST("/register", ctrl.Register)
	public.POST("/login", ctrl.Login)

	//Protected route
	protected := engine.Group("")
	// Attache the AuthMiddleware 
	protected.Use(infrastructure.AuthMiddleware())

	//  Admin-only routes
	admin := protected.Group("")
	admin.Use(infrastructure.AdminMiddleware())

	admin.POST("/promote", ctrl.PromoteUser)
	
	tasks := admin.Group("/tasks")
	{
		tasks.POST("/", ctrl.AddTask)
		tasks.PUT("/:id", ctrl.UpdatedTask)
		tasks.DELETE("/:id", ctrl.RemoveTask)
	}

	// Task routes accessible to all authenticated users
	userTasks := protected.Group("/tasks")	
	{
		userTasks.GET("/", ctrl.GetAllTasks)
		userTasks.GET("/:id", ctrl.GetTask)
	}
}
