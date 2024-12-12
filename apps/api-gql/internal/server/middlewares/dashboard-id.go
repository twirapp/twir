package middlewares

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"github.com/twirapp/twir/libs/baseapp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (m *Middlewares) DashboardID(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())

	user, userErr := m.sessions.GetAuthenticatedUser(c.Request.Context())
	if userErr == nil {
		span.SetAttributes(
			attribute.String("user.id", user.ID),
		)
		sloggin.AddCustomAttributes(c, slog.String("userId", user.ID))
		c.Request = c.Request.WithContext(
			context.WithValue(c.Request.Context(), baseapp.RequesterUserIdContextKey, user.ID),
		)
	}

	selectedDashboardID, err := m.sessions.GetSelectedDashboard(c.Request.Context())
	var dashboardIdForSet string
	if err == nil {
		dashboardIdForSet = selectedDashboardID
	} else if userErr == nil {
		dashboardIdForSet = user.ID
	}

	if dashboardIdForSet != "" {
		span.SetAttributes(
			attribute.String("user.selectedDashboard", dashboardIdForSet),
		)

		c.Request = c.Request.WithContext(
			context.WithValue(
				c.Request.Context(),
				baseapp.SelectedDashboardContextKey,
				dashboardIdForSet,
			),
		)
	}
}
