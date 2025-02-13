package repository

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepository struct {
	Collection *mongo.Collection
}

func NewAuthRepository(db *mongo.Database, collection string) *AuthRepository {
	return &AuthRepository{
		Collection: db.Collection(collection),
	}
}

func (ar *AuthRepository) CreateAuthDocs(ctx context.Context, auth AuthDocs) error {
	_, err := ar.Collection.InsertOne(ctx, auth)
	if err != nil {
		slog.Error("Auth Repository :: CreateAuthDocs", slog.Any("error", err))
		return err
	}
	return nil
}

func (ar *AuthRepository) GetAuthDocs(ctx context.Context, teacherID string) (AuthDocs, error) {
	var auth AuthDocs
	err := ar.Collection.FindOne(ctx, map[string]string{"teacherID": teacherID}).Decode(&auth)
	if err != nil {
		slog.Error("Auth Repository :: GetAuthDocs", slog.Any("error", err))
		return auth, err
	}
	return auth, nil
}

func (ar *AuthRepository) GetAllAuthDocs(ctx context.Context) ([]AuthDocs, error) {
	var auths []AuthDocs
	cursor, err := ar.Collection.Find(ctx, map[string]string{})
	if err != nil {
		slog.Error("Auth Repository :: GetAllAuthDocs", slog.Any("error", err))
		return auths, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var auth AuthDocs
		if err := cursor.Decode(&auth); err != nil {
			slog.Error("Auth Repository :: GetAllAuthDocs", slog.Any("error", err))
			return auths, err
		}
		auths = append(auths, auth)
	}
	return auths, nil
}

func (ar *AuthRepository) UpdateAuthDocs(ctx context.Context, teacherID string, auth AuthDocs) error {
	_, err := ar.Collection.UpdateOne(ctx, map[string]string{"teacherID": teacherID}, auth)
	if err != nil {
		slog.Error("Auth Repository :: UpdateAuthDocs", slog.Any("error", err))
		return err
	}
	return nil
}

func (ar *AuthRepository) DeleteAuthDocs(ctx context.Context, teacherID string) error {
	_, err := ar.Collection.DeleteOne(ctx, map[string]string{"teacherID": teacherID})
	if err != nil {
		slog.Error("Auth Repository :: DeleteAuthDocs", slog.Any("error", err))
		return err
	}
	return nil
}
