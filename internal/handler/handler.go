package handler

import (
	"go-scratch/generated"
	"go-scratch/internal/commons"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetApiActuatorHealth(ctx echo.Context) error {
	return nil
}

func (h *Handler) PostApiAuthRegister(ctx echo.Context) error {
	var req generated.UserSignup
	if err := ctx.Bind(&req); err != nil {
		slog.Error("User Handler ::: failed to bind request", slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, generated.Error{
			Message: "Invalid request",
			Fields:  "PostApiAuthRegister",
			Code:    http.StatusBadRequest,
		})
		return err
	}

	err := h.uas.RegisterUser(ctx.Request().Context(), req)
	if err != nil {
		slog.Error("User Handler ::: failed to register user", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, generated.Error{
			Message: "Internal server error",
			Fields:  "PostApiAuthRegister",
			Code:    http.StatusInternalServerError,
		})
		return err
	}

	return ctx.JSON(http.StatusCreated, "User registered successfully")
}

func (h *Handler) PostApiAuthLogin(ctx echo.Context) error {
	var req generated.UserLoginRequest
	if err := ctx.Bind(&req); err != nil {
		slog.Error("User Handler ::: failed to bind request", slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, generated.Error{
			Message: "Invalid request",
			Fields:  "PostApiAuthLogin",
			Code:    http.StatusBadRequest,
		})
		return err
	}

	res, err := h.uas.LoginUser(ctx.Request().Context(), req)
	if err != nil {
		slog.Error("User Handler ::: failed to login user", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, generated.Error{
			Message: "Internal server error",
			Fields:  "PostApiAuthLogin",
			Code:    http.StatusInternalServerError,
		})
		return err
	}
	return ctx.JSON(http.StatusOK, res)
}

func (h *Handler) GetApiAuthValidateId(ctx echo.Context, id string) error {
	if id == "" {
		slog.Error("User Handler ::: failed to validate user", slog.Any("error", "ID is required"))
		ctx.JSON(http.StatusBadRequest, generated.Error{
			Message: "Invalid request",
			Fields:  "GetApiAuthValidateId",
			Code:    http.StatusBadRequest,
		})
		return commons.ErrUserRequired
	}

	_, err := h.uas.ValidateUser(ctx.Request().Context(), id)
	if err != nil {
		slog.Error("User Handler ::: failed to validate user", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, generated.Error{
			Message: "Internal server error",
			Fields:  "GetApiAuthValidateId",
			Code:    http.StatusInternalServerError,
		})
		return err
	}
	return ctx.JSON(http.StatusOK, "User validated successfully")

}

func (h *Handler) PostApiAttendanceCheckin(ctx echo.Context) error {
	var req generated.CheckinRequest
	if err := ctx.Bind(&req); err != nil {
		slog.Error("Attendance Handler ::: failed to bind request", slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, generated.Error{
			Message: "Invalid request",
			Fields:  "PostApiAttendanceCheckin",
			Code:    http.StatusBadRequest,
		})
		return err
	}

	err := h.as.RecordCheckIn(ctx.Request().Context(), req)
	if err != nil {
		slog.Error("Attendance Handler ::: failed to record check-in", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, generated.Error{
			Message: "Internal server error",
			Fields:  "PostApiAttendanceCheckin",
			Code:    http.StatusInternalServerError,
		})
		return err
	}

	return ctx.JSON(http.StatusOK, "Check-in recorded successfully")
}

func (h *Handler) PostApiAttendanceCheckout(ctx echo.Context) error {
	var req generated.CheckoutRequest
	if err := ctx.Bind(&req); err != nil {
		slog.Error("Attendance Handler ::: failed to bind request", slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, generated.Error{
			Message: "Invalid request",
			Fields:  "PostApiAttendanceCheckout",
			Code:    http.StatusBadRequest,
		})
		return err
	}

	err := h.as.RecordCheckOut(ctx.Request().Context(), req)
	if err != nil {
		slog.Error("Attendance Handler ::: failed to record check-out", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, generated.Error{
			Message: "Internal server error",
			Fields:  "PostApiAttendanceCheckout",
			Code:    http.StatusInternalServerError,
		})
		return err
	}

	return ctx.JSON(http.StatusOK, "Check-out recorded successfully")
}
