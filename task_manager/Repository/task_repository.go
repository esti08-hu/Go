package repository

import (
	"context"
	"task_manager/Domain"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

type taskRepository struct {
	database *mongo.Database
	collection string
}

func NewTaskRepository(db *mongo.Database, collection string) domain.TaskRepository {
	return &taskRepository{
		database: db,
		collection: collection,
	}
}

func (tr *taskRepository) GetAllTasks(c context.Context, task *domain.Task) (*domain.Task, error) {
	collection := tr.database.Collection(tr.collection)

	_, err := collection.Find(c, bson.D{})
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (tr *taskRepository) GetTaskByID(c context.Context, id string) (*domain.Task, error) {
	collection := tr.database.Collection(tr.collection)

	filter := bson.M{"id": id}

	var task domain.Task
	err := collection.FindOne(c, filter).Decode(&task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (tr *taskRepository) CreateTask(c context.Context, task *domain.Task) error {
	collection := tr.database.Collection(tr.collection)

	_, err := collection.InsertOne(c, task)
	if err != nil {
		return err
	}

	return nil
}

func (tr *taskRepository) UpdateTask(c context.Context, id string, task *domain.Task) (*domain.Task, error) {
	collection := tr.database.Collection(tr.collection)

	filter := bson.M{"id": id}

	_, err := collection.UpdateOne(c, filter, bson.M{"$set": task})
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (tr *taskRepository) DeleteTask(c context.Context, id string) error {
	collection := tr.database.Collection(tr.collection)

	_, err := collection.DeleteOne(c, bson.M{"id": id})
	if err != nil {
		return err
	}

	return nil
}
