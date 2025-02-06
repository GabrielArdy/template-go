package config

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

// SlogConfig represents the configuration for structured logging
type SlogConfig struct {
	Level      slog.Level
	TimeFormat string
}

// ContextHandler wraps the slog.Handler with context-aware capabilities
type ContextHandler struct {
	Handler slog.Handler
}

func (h ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

// Handle implements slog.Handler interface
func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	attrs := h.addTraceFromContext(ctx)
	r.AddAttrs(attrs...)
	return h.Handler.Handle(ctx, r)
}

// WithAttrs implements slog.Handler interface
func (h ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{Handler: h.Handler.WithAttrs(attrs)}
}

// WithGroup implements slog.Handler interface
func (h ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{Handler: h.Handler.WithGroup(name)}
}

// NewSlogHandler creates a new slog handler with custom formatting
func NewSlogHandler(cfg SlogConfig) *ContextHandler {
	opts := &slog.HandlerOptions{
		Level: cfg.Level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				return slog.String(slog.TimeKey, t.Format(cfg.TimeFormat))
			}
			return a
		},
	}

	return &ContextHandler{
		Handler: slog.NewJSONHandler(os.Stdout, opts),
	}
}

// addTraceFromContext extracts trace information from context
func (h ContextHandler) addTraceFromContext(ctx context.Context) []slog.Attr {
	attrs := make([]slog.Attr, 0)
	if reqID := ctx.Value("request_id"); reqID != nil {
		attrs = append(attrs, slog.String("request_id", reqID.(string)))
	}
	return attrs
}

// SlogMiddleware creates an Echo middleware for logging
func SlogMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			req := c.Request()
			res := c.Response()

			slog.LogAttrs(req.Context(),
				slog.LevelInfo,
				"HTTP Request",
				slog.String("method", req.Method),
				slog.String("uri", req.RequestURI),
				slog.Int("status", res.Status),
				slog.Duration("latency", time.Since(start)),
				slog.String("ip", c.RealIP()),
				slog.String("user_agent", req.UserAgent()),
			)

			return err
		}
	}
}

// InitLogger initializes the global logger
func InitLogger() {
	cfg := SlogConfig{
		Level:      slog.LevelInfo,
		TimeFormat: time.RFC3339,
	}

	handler := NewSlogHandler(cfg)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
