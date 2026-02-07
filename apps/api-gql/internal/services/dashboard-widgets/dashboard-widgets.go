package dashboard_widgets

import (
	"context"

	"github.com/twirapp/twir/libs/entities/dashboard_widget"
	"github.com/twirapp/twir/libs/repositories/dashboard_widgets"
	"github.com/twirapp/twir/libs/wsrouter"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	DashboardWidgetsRepository dashboard_widgets.Repository
	WsRouter                   wsrouter.WsRouter
}

func New(opts Opts) *Service {
	return &Service{
		dashboardWidgetsRepository: opts.DashboardWidgetsRepository,
		wsRouter:                   opts.WsRouter,
	}
}

type Service struct {
	dashboardWidgetsRepository dashboard_widgets.Repository
	wsRouter                   wsrouter.WsRouter
}

func (s *Service) GetByChannelID(ctx context.Context, channelID string) ([]dashboard_widget.DashboardWidget, error) {
	return s.dashboardWidgetsRepository.GetByChannelID(ctx, channelID)
}

type UpdateInput struct {
	ChannelID string
	Layout    []dashboard_widget.DashboardWidget
}

func (s *Service) Update(ctx context.Context, input UpdateInput) ([]dashboard_widget.DashboardWidget, error) {
	err := s.dashboardWidgetsRepository.UpsertMany(ctx, input.ChannelID, input.Layout)
	if err != nil {
		return nil, err
	}

	return input.Layout, nil
}

func (s *Service) GetWsRouter() wsrouter.WsRouter {
	return s.wsRouter
}
