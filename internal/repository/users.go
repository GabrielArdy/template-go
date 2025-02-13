package repository

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepostiory struct {
	Collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database, collection string) *UserRepostiory {
	return &UserRepostiory{
		Collection: db.Collection(collection),
	}
}

func (ur *UserRepostiory) CreateUserDocs(ctx context.Context, user UserDocs) error {
	_, err := ur.Collection.InsertOne(ctx, user)
	if err != nil {
		slog.Error("User Repository :: CreateUserDocs", slog.Any("error", err))
		return err
	}
	return nil
}

func (ur *UserRepostiory) GetUserDocs(ctx context.Context, teacherID string) (UserDocs, error) {
	var user UserDocs
	err := ur.Collection.FindOne(ctx, map[string]string{"teacherID": teacherID}).Decode(&user)
	if err != nil {
		slog.Error("User Repository :: GetUserDocs", slog.Any("error", err))
		return user, err
	}
	return user, nil
}

func (ur *UserRepostiory) GetAllUserDocs(ctx context.Context) ([]UserDocs, error) {
	var users []UserDocs
	cursor, err := ur.Collection.Find(ctx, map[string]string{})
	if err != nil {
		slog.Error("User Repository :: GetAllUserDocs", slog.Any("error", err))
		return users, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user UserDocs
		if err := cursor.Decode(&user); err != nil {
			slog.Error("User Repository :: GetAllUserDocs", slog.Any("error", err))
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (ur *UserRepostiory) UpdateUserDocs(ctx context.Context, teacherID string, user UserDocs) error {
	_, err := ur.Collection.UpdateOne(ctx, map[string]string{"teacherID": teacherID}, user)
	if err != nil {
		slog.Error("User Repository :: UpdateUserDocs", slog.Any("error", err))
		return err
	}
	return nil
}

func (ur *UserRepostiory) DeleteUserDocs(ctx context.Context, teacherID string) error {
	_, err := ur.Collection.DeleteOne(ctx, map[string]string{"teacherID": teacherID})
	if err != nil {
		slog.Error("User Repository :: DeleteUserDocs", slog.Any("error", err))
		return err
	}
	return nil
}
