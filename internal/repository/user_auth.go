package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserAuthRepository struct {
	collection *mongo.Collection
}

func NewUserAuthRepository(db *mongo.Database, collectionName string) *UserAuthRepository {
	return &UserAuthRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *UserAuthRepository) InsertOne(ctx context.Context, userAuth UserAuth) error {
	_, err := r.collection.InsertOne(nil, userAuth)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserAuthRepository) FindOne(ctx context.Context, email string) (UserAuth, error) {
	var userAuth UserAuth
	err := r.collection.FindOne(nil, email).Decode(&userAuth)
	if err != nil {
		return UserAuth{}, err
	}
	return userAuth, nil
}
