package repository

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ActivityRepository struct {
	Collection *mongo.Collection
}

func NewActivityRepository(db *mongo.Database, collection string) *ActivityRepository {
	return &ActivityRepository{
		Collection: db.Collection(collection),
	}
}

func (ar *ActivityRepository) CreateActivityDocs(ctx context.Context, activity ActivityDocs) error {
	_, err := ar.Collection.InsertOne(ctx, activity)
	if err != nil {
		slog.Error("Activity Repository :: CreateActivityDocs", slog.Any("error", err))
		return err
	}
	return nil
}

func (ar *ActivityRepository) GetActivityDocs(ctx context.Context, teacherID string) (ActivityDocs, error) {
	var activity ActivityDocs
	err := ar.Collection.FindOne(ctx, map[string]string{"teacherID": teacherID}).Decode(&activity)
	if err != nil {
		slog.Error("Activity Repository :: GetActivityDocs", slog.Any("error", err))
		return activity, err
	}
	return activity, nil
}

func (ar *ActivityRepository) GetAllActivityDocs(ctx context.Context) ([]ActivityDocs, error) {
	var activities []ActivityDocs
	cursor, err := ar.Collection.Find(ctx, map[string]string{})
	if err != nil {
		slog.Error("Activity Repository :: GetAllActivityDocs", slog.Any("error", err))
		return activities, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var activity ActivityDocs
		if err := cursor.Decode(&activity); err != nil {
			slog.Error("Activity Repository :: GetAllActivityDocs", slog.Any("error", err))
			return activities, err
		}
		activities = append(activities, activity)
	}
	return activities, nil
}

// get all activity docs by createdat datetime range
func (ar *ActivityRepository) GetAllActivityDocsByCreatedAtRange(ctx context.Context, start, end string) ([]ActivityDocs, error) {
	var activities []ActivityDocs

	// Fix: Use bson.M for the entire query structure
	filter := bson.M{
		"createdAt": bson.M{
			"$gte": start,
			"$lt":  end,
		},
	}

	cursor, err := ar.Collection.Find(ctx, filter)
	if err != nil {
		slog.Error("Activity Repository :: GetAllActivityDocsByCreatedAtRange",
			slog.Any("error", err),
			slog.String("start", start),
			slog.String("end", end))
		return activities, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var activity ActivityDocs
		if err := cursor.Decode(&activity); err != nil {
			slog.Error("Activity Repository :: GetAllActivityDocsByCreatedAtRange",
				slog.Any("error", err))
			return activities, err
		}
		activities = append(activities, activity)
	}

	return activities, nil
}
