package services

import (
	"context"
	"go-scratch/internal/commons"
	"go-scratch/internal/repository"
	"log/slog"
	"time"
)

type JobService struct {
	arp AttendanceRepository
	urp UserRepository
	ls  *LoggingService
}

type QRObject struct {
	IssueDate time.Time
	IssuedBy  string
	QRToken   string
}

func NewJobService(ar AttendanceRepository, ur UserRepository, ls *LoggingService) *JobService {
	return &JobService{
		arp: ar,
		urp: ur,
		ls:  ls,
	}
}

func (js *JobService) GenerateAttendanceDocs(ctx context.Context) error {
	// Fetch all user docs
	userDocs, err := js.urp.GetAllUserDocs(ctx)
	if err != nil {
		return err
	}

	currentTime := commons.GetLocalTime()
	currentDate := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())

	// Create attendance docs for each user
	for _, user := range userDocs {
		attendanceDoc := repository.AttendanceDocs{
			UID:       commons.GenerateCustomUID(),
			Date:      currentDate,
			CheckIn:   time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()),
			CheckOut:  time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()),
			Status:    commons.ATTENDANCE_ABSENT,
			TeacherId: user.TeacherID,
		}

		err = js.arp.CreateAttendanceDocs(ctx, attendanceDoc)
		slog.Info("::: Successfully to create attendance for %s", user.TeacherID, slog.Any("service", "JobService"))
		if err != nil {
			return err
		}
	}
	return nil
}
