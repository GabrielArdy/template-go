package handler

import (
	"go-scratch/generated"
)

type (
	Handler struct {
	}
)

var _ generated.ServerInterface = (*Handler)(nil)

func NewHandler() *Handler {
	return &Handler{}
}
