package data

import (
	"context"
	"time"

	"task_manager/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var TaskCollection *mongo.Collection

func Ping() string {
	return "pong"
}

func GetTasks() []models.Task {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := TaskCollection.Find(ctx, bson.D{})
	if err != nil {
		return []models.Task{}
	}

	var tasks []models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return []models.Task{}
	}
	return tasks
}

func GetTaskById(id string) *models.Task {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var task models.Task
	err := TaskCollection.FindOne(ctx, bson.M{"id": id}).Decode(&task)
	if err != nil {
		return nil
	}
	return &task
}

func RemoveTask(id string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := TaskCollection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return
	}
}


func UpdatedTask(id string, updateTask models.Task) (models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{
		"$set": bson.M{
			"title":       updateTask.Title,
			"description": updateTask.Description,
			"status":      updateTask.Status,
			"due_date":    updateTask.DueDate,
		},
	}

	var updatedTask models.Task
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := TaskCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedTask)
	if err != nil {
		return models.Task{}, err
	}
	return updatedTask, nil
}

func AddTask(newTask models.Task) (models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := TaskCollection.InsertOne(ctx, newTask)
	if err != nil {
		return models.Task{}, err
	}
	return newTask, nil
}
