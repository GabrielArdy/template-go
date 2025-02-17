package services

import (
	"context"
	"go-scratch/generated"
	"go-scratch/internal/commons"
	"go-scratch/internal/repository"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type AttendanceRepository interface {
	CreateAttendanceDocs(ctx context.Context, attendance repository.AttendanceDocs) error
	GetAttendanceDocs(ctx context.Context, teacherID string) (repository.AttendanceDocs, error)
	GetAllAttendanceDocs(ctx context.Context) ([]repository.AttendanceDocs, error)
	UpdateAttendanceDocs(ctx context.Context, teacherID string, attendance repository.AttendanceDocs) error
	GetAttendanceDocsByDateRange(ctx context.Context, teacherId, start, end string) ([]repository.AttendanceDocs, error)
}

type AttendanceService struct {
	c              redis.UniversalClient
	attendanceRepo AttendanceRepository
	log            *LoggingService
}

func NewAttendanceService(ar AttendanceRepository, ls *LoggingService, c redis.UniversalClient) *AttendanceService {
	return &AttendanceService{
		attendanceRepo: ar,
		log:            ls,
		c:              c,
	}
}

func (as *AttendanceService) RecordCheckIn(ctx context.Context, req generated.CheckinRequest) error {
	if req.TeacherId == "" {
		slog.Error("AttendanceService :: RecordCheckIn", slog.Any("error", commons.ErrInvalidTeacherID))
		return commons.ErrInvalidTeacherID
	}

	attendanceDoc, err := as.attendanceRepo.GetAttendanceDocs(ctx, req.TeacherId)
	if err != nil {
		slog.Error("AttendanceService :: RecordCheckIn", slog.Any("error", err))
		return err
	}

	ok, err := commons.VerifyQRCode(ctx, req.QRToken, as.c)
	if err != nil {
		slog.Error("AttendanceService :: RecordCheckIn", slog.Any("error", err))
		return err
	}

	if ok {
		attendanceDoc.CheckIn = commons.GetLocalTime()
		attendanceDoc.Status = commons.ATTENDANCE_ABSENT
		err = as.attendanceRepo.UpdateAttendanceDocs(ctx, req.TeacherId, attendanceDoc)
		if err != nil {
			slog.Error("AttendanceService :: RecordCheckIn", slog.Any("error", err))
			return err
		}
		as.log.LogActivity(ctx, req.TeacherId, commons.ACTIVITY_CHECK_IN)
		slog.Info("AttendanceService :: RecordCheckIn :::", slog.Any("teacherId", req.TeacherId))
	}

	return nil

}

func (as *AttendanceService) RecordCheckOut(ctx context.Context, req generated.CheckoutRequest) error {
	if req.TeacherId == "" {
		slog.Error("AttendanceService :: RecordCheckOut", slog.Any("error", commons.ErrInvalidTeacherID))
		return commons.ErrInvalidTeacherID
	}

	attendanceDoc, err := as.attendanceRepo.GetAttendanceDocs(ctx, req.TeacherId)
	if err != nil {
		slog.Error("AttendanceService :: RecordCheckOut", slog.Any("error", err))
		return err
	}

	ok, err := commons.VerifyQRCode(ctx, req.QRToken, as.c)
	if err != nil {
		slog.Error("AttendanceService :: RecordCheckOut", slog.Any("error", err))
		return err
	}

	if ok {
		attendanceDoc.CheckOut = commons.GetLocalTime()
		attendanceDoc.Status = commons.ATTENDANCE_PRESENT
		err = as.attendanceRepo.UpdateAttendanceDocs(ctx, req.TeacherId, attendanceDoc)
		if err != nil {
			slog.Error("AttendanceService :: RecordCheckOut", slog.Any("error", err))
			return err
		}
		as.log.LogActivity(ctx, req.TeacherId, commons.ACTIVITY_CHECK_OUT)
		slog.Info("AttendanceService :: RecordCheckOut :::", slog.Any("teacherId", req.TeacherId))
	}

	return nil
}
