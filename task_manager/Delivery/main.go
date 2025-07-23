package main

import (
	"context"
	"log"
	"os"
	"task_manager/Delivery/controller"
	router "task_manager/Delivery/routers"
	repository "task_manager/Repository"
	usecases "task_manager/Usecases"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("DATABASE_URL")
	if mongoURI == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}
	clientOptions := options.Client().ApplyURI(mongoURI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Collections
	db := client.Database("task_db")
	userRepo := repository.NewUserRepository(db, "users")
	taskRepo := repository.NewTaskRepository(db, "tasks")

	// Initialize usecases
	timeout := 10*time.Second
	userUsecase := usecases.NewUserUsecases(userRepo, timeout)
	taskUsecase := usecases.NewTaskUsecases(taskRepo, timeout)

	// Initialize controllers
	ctrl := controller.Controller{
		UserUsecases: userUsecase,
		TaskUsecases: taskUsecase,
	}

	// Setup router
	engine := gin.Default()
	router.SetupRouter(engine, ctrl)

	engine.Run(":8080")
}
