package handler

import (
	"go-scratch/generated"
	"go-scratch/internal/services"
)

type (
	Handler struct {
		uas *services.UserAuthService
		as  *services.AttendanceService
		js  *services.JobService
	}
)

var _ generated.ServerInterface = (*Handler)(nil)

func NewHandler(uasSvc *services.UserAuthService, asSvc *services.AttendanceService, jSvc *services.JobService) *Handler {
	return &Handler{
		uas: uasSvc,
		as:  asSvc,
		js:  jSvc,
	}
}
