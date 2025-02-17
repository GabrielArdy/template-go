package services

import (
	"context"
	"go-scratch/internal/commons"
	"go-scratch/internal/repository"
	"log/slog"
)

type ActivityRepository interface {
	CreateActivityDocs(ctx context.Context, activity repository.ActivityDocs) error
	GetActivityDocs(ctx context.Context, teacherID string) (repository.ActivityDocs, error)
	GetAllActivityDocs(ctx context.Context) ([]repository.ActivityDocs, error)
	GetAllActivityDocsByCreatedAtRange(ctx context.Context, start, end string) ([]repository.ActivityDocs, error)
}
type LoggingService struct {
	activityRepo ActivityRepository
}

func NewLoggingService(ar ActivityRepository) *LoggingService {
	return &LoggingService{
		activityRepo: ar,
	}
}

func (ls *LoggingService) LogActivity(ctx context.Context, userId, activityType string) error {
	uid := commons.GenerateCustomUID()
	currentTime := commons.GetLocalTime()
	docs := repository.ActivityDocs{
		UID:       uid,
		UserId:    userId,
		Type:      activityType,
		CreatedAt: currentTime,
	}

	err := ls.activityRepo.CreateActivityDocs(ctx, docs)
	if err != nil {
		slog.Error("LoggingService :: LogActivity", slog.Any("error", err))
		return err
	}
	return nil
}
