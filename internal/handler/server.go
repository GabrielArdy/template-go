package handler

import (
	"go-scratch/generated"
	"go-scratch/internal/services"
)

type (
	Handler struct {
		uas *services.UserAuthService
	}
)

var _ generated.ServerInterface = (*Handler)(nil)

func NewHandler(uasSvc *services.UserAuthService) *Handler {
	return &Handler{
		uas: uasSvc,
	}
}
