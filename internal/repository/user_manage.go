package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database, collectionName string) *UserRepository {
	return &UserRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *UserRepository) InsertOne(ctx context.Context, user User) error {
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindOne(ctx context.Context, email string) (User, error) {
	var user User
	err := r.collection.FindOne(ctx, email).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
