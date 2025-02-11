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

func (r *UserAuthRepository) FindOneByUsername(ctx context.Context, username string) (UserAuth, error) {
	var userAuth UserAuth
	err := r.collection.FindOne(nil, username).Decode(&userAuth)
	if err != nil {
		return UserAuth{}, err
	}
	return userAuth, nil
}

func (r *UserAuthRepository) UpdateRefreshToken(ctx context.Context, email, refreshToken string) error {
	_, err := r.collection.UpdateOne(nil, email, refreshToken)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserAuthRepository) UpdateAccessToken(ctx context.Context, email, accessToken string) error {
	_, err := r.collection.UpdateOne(nil, email, accessToken)
	if err != nil {
		return err
	}
	return nil
}
