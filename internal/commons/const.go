package commons

import (
	"errors"
	"time"
)

const (
	ADMIN      = "admin"
	INSTRUCTOR = "instructor"
	STUDENT    = "student"

	// default ttl for jwt token
	ACCESS_TTL  = time.Minute * 15
	REFRESH_TTL = time.Hour * 24 * 7
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInternalServer     = errors.New("internal server error")
)
