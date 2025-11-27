package middlewares

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

func (m *Middlewares) Logger() gin.HandlerFunc {
	return sloggin.NewWithConfig(
		m.logger,
		sloggin.Config{
			DefaultLevel:     slog.LevelInfo,
			ClientErrorLevel: slog.LevelWarn,
			ServerErrorLevel: slog.LevelError,
			Filters: []sloggin.Filter{
				sloggin.IgnoreStatus(200, 404),
			},
		},
	)
}
