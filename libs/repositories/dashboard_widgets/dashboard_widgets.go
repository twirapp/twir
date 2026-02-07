package dashboard_widgets

import (
	"context"

	"github.com/twirapp/twir/libs/entities/dashboard_widget"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) ([]dashboard_widget.DashboardWidget, error)
	UpsertMany(ctx context.Context, channelID string, widgets []dashboard_widget.DashboardWidget) error
}
