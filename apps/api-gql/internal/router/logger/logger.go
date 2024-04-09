package logger

import (
	"log/slog"
	"time"

	"github.com/go-chi/httplog/v2"
)

func New() *httplog.Logger {
	logger := httplog.NewLogger(
		"api",
		httplog.Options{
			// JSON:             true,
			LogLevel:         slog.LevelDebug,
			Concise:          true,
			RequestHeaders:   true,
			MessageFieldName: "message",
			QuietDownPeriod:  10 * time.Second,
		},
	)

	return logger
}
