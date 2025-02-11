package middleware

import (
	"go-scratch/internal/commons"
	"go-scratch/internal/repository"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	userAuthRepo *repository.UserAuthRepository
}

func NewAuthMiddleware(userAuthRepo *repository.UserAuthRepository) *AuthMiddleware {
	return &AuthMiddleware{
		userAuthRepo: userAuthRepo,
	}
}

// RequireAuth middleware checks auth from header or repository
func (m *AuthMiddleware) RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// First try Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" {
				// Validate Bearer token
				tokenString := strings.TrimPrefix(authHeader, "Bearer ")
				if tokenString == authHeader {
					return c.JSON(http.StatusUnauthorized, commons.HttpErrorResponse(
						http.StatusUnauthorized,
						"Invalid token format",
						"INVALID_TOKEN_FORMAT",
					))
				}

				claims, err := commons.ParseToken(tokenString)
				if err == nil {
					// Token is valid
					c.Set("claims", claims)
					return next(c)
				}
				// Token validation failed, continue to repository check
			}

			// Try repository-based authentication
			username := c.Request().Header.Get("X-Username")
			if username == "" {
				return c.JSON(http.StatusUnauthorized, commons.HttpErrorResponse(
					http.StatusUnauthorized,
					"Authentication required",
					"AUTH_REQUIRED",
				))
			}

			userAuth, err := m.userAuthRepo.FindOneByUsername(c.Request().Context(), username)
			if err != nil {
				slog.Error("Failed to find user", slog.Any("error", err))
				return c.JSON(http.StatusUnauthorized, commons.HttpErrorResponse(
					http.StatusUnauthorized,
					"Invalid credentials",
					"INVALID_CREDENTIALS",
				))
			}

			// Decrypt stored token
			decryptedToken, err := commons.Decrypt(userAuth.AccessToken)
			if err != nil {
				slog.Error("Failed to decrypt token", slog.Any("error", err))
				return c.JSON(http.StatusInternalServerError, commons.HttpErrorResponse(
					http.StatusInternalServerError,
					"Authentication error",
					"DECRYPT_ERROR",
				))
			}

			// Validate decrypted token
			claims, err := commons.ParseToken(decryptedToken)
			if err != nil {
				slog.Error("Invalid stored token", slog.Any("error", err))
				return c.JSON(http.StatusUnauthorized, commons.HttpErrorResponse(
					http.StatusUnauthorized,
					"Invalid stored token",
					"INVALID_STORED_TOKEN",
				))
			}

			c.Set("claims", claims)
			return next(c)
		}
	}
}

// RequireRole middleware remains the same
func (m *AuthMiddleware) RequireRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get("claims").(*commons.CustomClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, commons.HttpErrorResponse(
					http.StatusUnauthorized,
					"Invalid token claims",
					"INVALID_CLAIMS",
				))
			}

			hasRole := false
			for _, role := range roles {
				if claims.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				return c.JSON(http.StatusForbidden, commons.HttpErrorResponse(
					http.StatusForbidden,
					"Insufficient permissions",
					"FORBIDDEN",
				))
			}

			return next(c)
		}
	}
}
