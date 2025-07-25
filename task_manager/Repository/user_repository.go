package repository

import (
	"context"
	domain "task_manager/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	database *mongo.Database
	collection string
}

func NewUserRepository(db *mongo.Database, collection string) domain.UserRepository {
	return &userRepository{
		database: db,
		collection: collection,
	}
}

func (ur *userRepository) GetAllUsers(c context.Context, user *domain.User) ([]*domain.User, error) {
	collection := ur.database.Collection(ur.collection)

	cursor, err := collection.Find(c, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c)
	var users []*domain.User
	for cursor.Next(c) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (ur *userRepository) GetUserByID(c context.Context, id string) (*domain.User, error) {
	collection := ur.database.Collection(ur.collection)

	filter := bson.M{"id": id}

	var user domain.User
	err := collection.FindOne(c, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) GetUserByEmail(c context.Context, email string) (*domain.User, error) {
	collection := ur.database.Collection(ur.collection)
	filter := bson.M{"email": email}
	var user domain.User
	err := collection.FindOne(c, filter).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound // Translate mongodb to domain error
		}
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) GetUserByUsername(c context.Context, username string) (*domain.User, error) {
	collection := ur.database.Collection(ur.collection)
	filter := bson.M{"username": username}
	var user domain.User
	err := collection.FindOne(c, filter).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound // Translate mongodb to domain error
		}
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) CreateUser(c context.Context, user *domain.User) (*domain.User, error) {
	collection := ur.database.Collection(ur.collection)

	_, err := collection.InsertOne(c, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *userRepository) PromoteUserToAdmin(c context.Context, id string) error {
	collection := ur.database.Collection(ur.collection)

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"role": "admin"}}

	_, err := collection.UpdateOne(c, filter, update)
	return err
}

func (ur *userRepository) UserExists(c context.Context) (bool, error) {
	collection := ur.database.Collection(ur.collection)

	count, err := collection.CountDocuments(c, bson.D{})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}