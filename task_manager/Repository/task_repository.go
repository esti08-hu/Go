package repository

import (
	"context"
	domain "task_manager/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type taskRepository struct {
	database *mongo.Database
	collection string
}

func NewTaskRepository(db *mongo.Database, collection string) domain.TaskRepository {
	return &taskRepository{
		database:   db,
		collection: collection,
	}
}

func (tr *taskRepository) GetAllTasks(c context.Context, id string) ([]*domain.Task, error) {
	collection := tr.database.Collection(tr.collection)
	filter := bson.M{"userid": id}
	cursor, err := collection.Find(c, filter)
	
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c)

	var tasks []*domain.Task
	for cursor.Next(c) {
		var t domain.Task
		if err := cursor.Decode(&t); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, mongo.ErrNoDocuments
	}
	return tasks, nil
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
