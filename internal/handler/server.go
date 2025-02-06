package handler

import (
	"go-scratch/generated"
	"go-scratch/internal/services"
)

type (
	Handler struct {
		svc     *services.CrawlerService
		authsvc *services.UserAuthService
	}
)

var _ generated.ServerInterface = (*Handler)(nil)

func NewHandler(svc *services.CrawlerService, auth *services.UserAuthService) *Handler {
	return &Handler{
		svc:     svc,
		authsvc: auth,
	}
}
