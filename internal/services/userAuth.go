package services

import (
	"context"
	"go-scratch/config"
	"go-scratch/generated"
	"go-scratch/internal/commons"
	"go-scratch/internal/repository"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUserDocs(context.Context, repository.UserDocs) error
	GetUserDocs(context.Context, string) (repository.UserDocs, error)
	GetAllUserDocs(context.Context) ([]repository.UserDocs, error)
	UpdateUserDocs(context.Context, string, repository.UserDocs) error
	DeleteUserDocs(context.Context, string) error
}

type AuthRepository interface {
	CreateAuthDocs(context.Context, repository.AuthDocs) error
	GetAuthDocs(context.Context, string) (repository.AuthDocs, error)
	GetAllAuthDocs(context.Context) ([]repository.AuthDocs, error)
	UpdateAuthDocs(context.Context, string, repository.AuthDocs) error
	DeleteAuthDocs(context.Context, string) error
}

type UserAuthService struct {
	UserRepository UserRepository
	AuthRepository AuthRepository
	Logging        *LoggingService
	Cache          redis.UniversalClient
}

func NewUserAuthService(ur UserRepository, ar AuthRepository, lr *LoggingService, c redis.UniversalClient) *UserAuthService {
	return &UserAuthService{
		UserRepository: ur,
		AuthRepository: ar,
		Logging:        lr,
		Cache:          c,
	}
}

func (ua *UserAuthService) RegisterUser(ctx context.Context, user generated.UserSignup) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("UserAuthService :: RegisterUser", slog.Any("error", err))
		return err
	}

	currentTime := commons.GetLocalTime()

	UID := commons.GenerateCustomUID()

	userDoc := repository.UserDocs{
		UID:       UID,
		Email:     string(user.Email),
		TeacherID: user.TeacherId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Grade:     user.Grade,
		Position:  user.Position,
		Status:    user.Status,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	authDoc := repository.AuthDocs{
		UID:        UID,
		Email:      string(user.Email),
		Password:   string(hashedPassword),
		TeacherID:  user.TeacherId,
		Role:       commons.ROLE_USER,
		CreatedAt:  currentTime,
		UpdatedAt:  currentTime,
		LastActive: currentTime,
	}

	err = ua.UserRepository.CreateUserDocs(ctx, userDoc)
	if err != nil {
		slog.Error("UserAuthService :: RegisterUser", slog.Any("error", err))
		return err
	}

	err = ua.AuthRepository.CreateAuthDocs(ctx, authDoc)
	if err != nil {
		slog.Error("UserAuthService :: RegisterUser", slog.Any("error", err))
		return err
	}

	return nil
}

func (ua *UserAuthService) LoginUser(ctx context.Context, req generated.UserLoginRequest) (generated.UserLoginResponse, error) {
	if req.User == "" {
		slog.Error("UserAuthService :: LoginUser", slog.Any("error", "Email is required"))
		return generated.UserLoginResponse{}, commons.ErrUserRequired
	}

	if req.Password == "" {
		slog.Error("UserAuthService :: LoginUser", slog.Any("error", "Password is required"))
		return generated.UserLoginResponse{}, commons.ErrPasswordRequired
	}

	authDoc, err := ua.AuthRepository.GetAuthDocs(ctx, req.User)
	if err != nil {
		slog.Error("UserAuthService :: LoginUser", slog.Any("error", err))
		return generated.UserLoginResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(authDoc.Password), []byte(req.Password))
	if err != nil {
		slog.Error("UserAuthService :: LoginUser", slog.Any("error", err))
		return generated.UserLoginResponse{}, commons.ErrInvalidCredentials
	}

	token, err := config.GenerateToken(authDoc.TeacherID, authDoc.Role)
	if err != nil {
		slog.Error("UserAuthService :: LoginUser", slog.Any("error", err))
		return generated.UserLoginResponse{}, err
	}

	err = ua.Cache.Set(ctx, commons.REDIS_KEY+authDoc.TeacherID, token, commons.REDIS_EXPIRY).Err()
	if err != nil {
		slog.Error("UserAuthService :: LoginUser", slog.Any("error", err))
		return generated.UserLoginResponse{}, err
	}

	err = ua.Logging.LogActivity(ctx, authDoc.TeacherID, commons.ACTIVITY_LOGIN)
	if err != nil {
		slog.Error("UserAuthService :: LoginUser", slog.Any("error", err))
	}

	return generated.UserLoginResponse{
		Token:     token,
		TeacherId: authDoc.TeacherID,
	}, nil
}

func (ua *UserAuthService) LogoutUser(ctx context.Context, teacherID string) error {
	err := ua.Cache.Del(ctx, commons.REDIS_KEY+teacherID).Err()
	if err != nil {
		slog.Error("UserAuthService :: LogoutUser", slog.Any("error", err))
		return err
	}

	return nil
}

func (ua *UserAuthService) ValidateUser(ctx context.Context, teacherId string) (bool, error) {
	token, err := ua.Cache.Get(ctx, commons.REDIS_KEY+teacherId).Result()
	if err != nil {
		slog.Error("UserAuthService :: ValidateUser", slog.Any("error", err))
		return false, err
	}

	return config.ValidateToken(token)
}
