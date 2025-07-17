package main

import (
	"context"
	"log"
	"os"
	"time"

	"task_manager/router"
	"task_manager/data"

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
	data.TaskCollection = client.Database("task_db").Collection("tasks")

	r := router.SetupRouter()
	r.Run(":8080")
}