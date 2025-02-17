package handler

import (
	"go-scratch/generated"
	"go-scratch/internal/services"
)

type (
	Handler struct {
		uas *services.UserAuthService
		as  *services.AttendanceService
	}
)

var _ generated.ServerInterface = (*Handler)(nil)

func NewHandler(uasSvc *services.UserAuthService, asSvc *services.AttendanceService) *Handler {
	return &Handler{
		uas: uasSvc,
		as:  asSvc,
	}
}
