package services

import (
	"context"
	"go-scratch/generated"
	"go-scratch/internal/commons"
	"go-scratch/internal/repository"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type UserAuthRepository interface {
	InsertOne(ctx context.Context, userAuth repository.UserAuth) error
	FindOne(ctx context.Context, email string) (repository.UserAuth, error)
}

type UserRepository interface {
	InsertOne(ctx context.Context, userModel repository.User) error
	FindOne(ctx context.Context, email string) (repository.User, error)
}

type UserAuthService struct {
	uarp UserAuthRepository
	urp  UserRepository
}

func NewUserAuthService(uarp UserAuthRepository, urp UserRepository) *UserAuthService {
	return &UserAuthService{
		uarp: uarp,
		urp:  urp,
	}
}

func (s *UserAuthService) Register(ctx context.Context, user generated.RegisterUserRequest) error {
	hash, err := commons.HashPassword(user.Password)
	if err != nil {
		slog.Error("Error hashing password: %v", slog.Any("error", err.Error()))
		return err
	}

	accessToken, err := commons.GenerateAccessToken(user.Username, commons.STUDENT)
	if err != nil {
		slog.Error("Error generating access token: %v", slog.Any("error", err.Error()))
		return err
	}

	refreshToken, err := commons.GenerateRefreshToken(user.Username, commons.STUDENT)
	if err != nil {
		slog.Error("Error generating refresh token: %v", slog.Any("error", err.Error()))
		return err
	}

	accessToken, err = commons.Encrypt(accessToken)
	if err != nil {
		slog.Error("Error encrypting access token: %v", slog.Any("error", err.Error()))
		return err
	}

	refreshToken, err = commons.Encrypt(refreshToken)
	if err != nil {
		slog.Error("Error encrypting refresh token: %v", slog.Any("error", err.Error()))
		return err
	}

	userUid := uuid.NewString()
	currentTime := time.Now()

	userAuth := repository.UserAuth{
		UserID:       userUid,
		Username:     user.Username,
		Email:        user.Email,
		Password:     hash,
		LastLogin:    time.Now(),
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}

	userModel := repository.User{
		UserID:    userUid,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	err = s.urp.InsertOne(ctx, userModel)
	if err != nil {
		slog.Error("Error inserting user: %v", slog.Any("error", err.Error()))
		return err
	}

	err = s.uarp.InsertOne(ctx, userAuth)
	if err != nil {
		slog.Error("Error inserting user auth: %v", slog.Any("error", err.Error()))
		return err
	}

	slog.Info("User registered successfully: %v", slog.Any("user", user.Username))
	return nil

}

func (s *UserAuthService) Login(ctx context.Context, email, password string) (repository.UserAuth, error) {
	userAuth, err := s.uarp.FindOne(ctx, email)
	if err != nil {
		slog.Error("Error finding user auth: %v", slog.Any("error", err.Error()))
		return repository.UserAuth{}, err
	}

	if !commons.CheckPasswordHash(userAuth.Password, password) {
		slog.Error("Invalid password for user: %v", slog.Any("email", email))
		return repository.UserAuth{}, commons.ErrInvalidCredentials
	}

	accesstoken, err := commons.Decrypt(userAuth.AccessToken)
	if err != nil {
		slog.Error("Error decrypting access token: %v", slog.Any("error", err.Error()))
		return repository.UserAuth{}, err
	}

	refreshToken, err := commons.Decrypt(userAuth.RefreshToken)
	if err != nil {
		slog.Error("Error decrypting refresh token: %v", slog.Any("error", err.Error()))
		return repository.UserAuth{}, err
	}

	userAuth.AccessToken = accesstoken
	userAuth.RefreshToken = refreshToken

	return userAuth, nil
}
