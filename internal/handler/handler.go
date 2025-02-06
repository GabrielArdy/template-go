package handler

import (
	"go-scratch/generated"
	"go-scratch/internal/commons"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) PostApiAuthRegister(ctx echo.Context) error {
	var req generated.RegisterUserRequest
	if err := ctx.Bind(&req); err != nil {
		slog.Error("Failed to bind request", slog.Any("error", err))
		return ctx.JSON(http.StatusBadRequest, commons.HttpErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
	}

	err := h.authsvc.Register(ctx.Request().Context(), req)
	if err != nil {
		slog.Error("Failed to register user", slog.Any("error", err))
		return ctx.JSON(http.StatusInternalServerError, commons.HttpErrorResponse(http.StatusInternalServerError, "Failed to register user", err.Error()))
	}

	return ctx.JSON(http.StatusOK, commons.HttpResponse(http.StatusOK, "User registered successfully", nil))
}

func (h *Handler) PostApiCrawl(ctx echo.Context) error {
	var req generated.CrawlRequest
	if err := ctx.Bind(&req); err != nil {
		slog.Error("Failed to bind request", slog.Any("error", err))
		return ctx.JSON(http.StatusBadRequest, commons.HttpErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
	}

	result, err := h.svc.Crawl(ctx.Request().Context(), req.Url, req.Depth)
	if err != nil {
		slog.Error("Failed to crawl", slog.Any("error", err))
		return ctx.JSON(http.StatusInternalServerError, commons.HttpErrorResponse(http.StatusInternalServerError, "Failed to crawl", err.Error()))
	}

	return ctx.JSON(http.StatusOK, result)

}
