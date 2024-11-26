package dashboard_widget_events

import (
	"context"

	model "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard-widget-events/model"
)

type DashboardWidgetEventsService interface {
	GetDashboardWidgetsEvents(ctx context.Context, channelID string, limit int) ([]model.Event, error)
}
