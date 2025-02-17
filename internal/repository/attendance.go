package repository

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
)

type AttendanceRepository struct {
	Collection *mongo.Collection
}

func NewAttendanceRepository(db *mongo.Database, collection string) *AttendanceRepository {
	return &AttendanceRepository{
		Collection: db.Collection(collection),
	}
}

func (ar *AttendanceRepository) CreateAttendanceDocs(ctx context.Context, attendance AttendanceDocs) error {
	_, err := ar.Collection.InsertOne(ctx, attendance)
	if err != nil {
		slog.Error("Attendance Repository :: CreateAttendanceDocs", slog.Any("error", err))
		return err
	}
	return nil
}

func (ar *AttendanceRepository) GetAttendanceDocs(ctx context.Context, teacherID string) (AttendanceDocs, error) {
	var attendance AttendanceDocs
	err := ar.Collection.FindOne(ctx, map[string]string{"teacherID": teacherID}).Decode(&attendance)
	if err != nil {
		slog.Error("Attendance Repository :: GetAttendanceDocs", slog.Any("error", err))
		return attendance, err
	}
	return attendance, nil
}

func (ar *AttendanceRepository) GetAllAttendanceDocs(ctx context.Context) ([]AttendanceDocs, error) {
	var attendances []AttendanceDocs
	cursor, err := ar.Collection.Find(ctx, map[string]string{})
	if err != nil {
		slog.Error("Attendance Repository :: GetAllAttendanceDocs", slog.Any("error", err))
		return attendances, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var attendance AttendanceDocs
		if err := cursor.Decode(&attendance); err != nil {
			slog.Error("Attendance Repository :: GetAllAttendanceDocs", slog.Any("error", err))
			return attendances, err
		}
		attendances = append(attendances, attendance)
	}
	return attendances, nil
}

func (ar *AttendanceRepository) UpdateAttendanceDocs(ctx context.Context, teacherID string, attendance AttendanceDocs) error {
	_, err := ar.Collection.UpdateOne(ctx, map[string]string{"teacherID": teacherID}, attendance)
	if err != nil {
		slog.Error("Attendance Repository :: UpdateAttendanceDocs", slog.Any("error", err))
		return err
	}
	return nil
}

// GetAttendanceDocsByDateRange returns all attendance docs between the given date range and teacherID
func (ar *AttendanceRepository) GetAttendanceDocsByDateRange(ctx context.Context, teacherID, start, end string) ([]AttendanceDocs, error) {
	var attendances []AttendanceDocs
	cursor, err := ar.Collection.Find(ctx, map[string]interface{}{"teacherID": teacherID, "date": map[string]interface{}{"$gte": start, "$lte": end}})
	if err != nil {
		slog.Error("Attendance Repository :: GetAttendanceDocsByDateRange", slog.Any("error", err))
		return attendances, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var attendance AttendanceDocs
		if err := cursor.Decode(&attendance); err != nil {
			slog.Error("Attendance Repository :: GetAttendanceDocsByDateRange", slog.Any("error", err))
			return attendances, err
		}
		attendances = append(attendances, attendance)
	}
	return attendances, nil
}
