package middlewares

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

func (m *Middlewares) Logging(c *gin.Context) {
	sloggin.NewWithConfig(
		m.logger.GetSlog(),
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
