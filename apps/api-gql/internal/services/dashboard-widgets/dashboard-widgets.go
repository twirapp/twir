package dashboard_widgets

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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

const DashboardWidgetsLayoutSubscriptionKey = "api.dashboardWidgetsLayout"

func DashboardWidgetsLayoutSubscriptionKeyCreate(channelID string) string {
	return DashboardWidgetsLayoutSubscriptionKey + "." + channelID
}

func (s *Service) Update(ctx context.Context, input UpdateInput) ([]dashboard_widget.DashboardWidget, error) {
	err := s.dashboardWidgetsRepository.UpsertMany(ctx, input.ChannelID, input.Layout)
	if err != nil {
		return nil, err
	}

	go func() {
		s.wsRouter.Publish(
			DashboardWidgetsLayoutSubscriptionKeyCreate(input.ChannelID),
			input.Layout,
		)
	}()

	return input.Layout, nil
}

type CreateCustomInput struct {
	ChannelID string
	Name      string
	URL       string
	X         int
	Y         int
	W         int
	H         int
}

func (s *Service) CreateCustom(ctx context.Context, input CreateCustomInput) (dashboard_widget.DashboardWidget, error) {
	widgetID := fmt.Sprintf("custom-%s", uuid.New().String())

	widget := dashboard_widget.DashboardWidget{
		WidgetID:   widgetID,
		ChannelID:  input.ChannelID,
		X:          input.X,
		Y:          input.Y,
		W:          input.W,
		H:          input.H,
		MinW:       2,
		MinH:       2,
		Visible:    true,
		StackId:    nil,
		StackOrder: 0,
		Type:       dashboard_widget.WidgetTypeCustom,
		CustomName: &input.Name,
		CustomUrl:  &input.URL,
	}

	err := s.dashboardWidgetsRepository.UpsertMany(
		ctx,
		input.ChannelID,
		[]dashboard_widget.DashboardWidget{widget},
	)
	if err != nil {
		return dashboard_widget.DashboardWidget{}, err
	}

	widgets, err := s.dashboardWidgetsRepository.GetByChannelID(ctx, input.ChannelID)
	if err != nil {
		return widget, err
	}

	go func() {
		s.wsRouter.Publish(
			DashboardWidgetsLayoutSubscriptionKeyCreate(input.ChannelID),
			widgets,
		)
	}()

	return widget, nil
}

type UpdateCustomInput struct {
	ChannelID string
	WidgetID  string
	Name      string
	URL       string
}

func (s *Service) UpdateCustom(ctx context.Context, input UpdateCustomInput) (dashboard_widget.DashboardWidget, error) {
	widgets, err := s.dashboardWidgetsRepository.GetByChannelID(ctx, input.ChannelID)
	if err != nil {
		return dashboard_widget.DashboardWidget{}, err
	}

	var targetWidget *dashboard_widget.DashboardWidget
	for i := range widgets {
		if widgets[i].WidgetID == input.WidgetID {
			targetWidget = &widgets[i]
			break
		}
	}

	if targetWidget == nil {
		return dashboard_widget.DashboardWidget{}, fmt.Errorf("widget not found: %s", input.WidgetID)
	}

	if targetWidget.Type != dashboard_widget.WidgetTypeCustom {
		return dashboard_widget.DashboardWidget{}, fmt.Errorf(
			"widget is not a custom widget: %s",
			input.WidgetID,
		)
	}

	targetWidget.CustomName = &input.Name
	targetWidget.CustomUrl = &input.URL

	err = s.dashboardWidgetsRepository.UpsertMany(
		ctx,
		input.ChannelID,
		[]dashboard_widget.DashboardWidget{*targetWidget},
	)
	if err != nil {
		return dashboard_widget.DashboardWidget{}, err
	}

	widgets, err = s.dashboardWidgetsRepository.GetByChannelID(ctx, input.ChannelID)
	if err != nil {
		return *targetWidget, err
	}

	go func() {
		s.wsRouter.Publish(
			DashboardWidgetsLayoutSubscriptionKeyCreate(input.ChannelID),
			widgets,
		)
	}()

	return *targetWidget, nil
}

type DeleteInput struct {
	ChannelID string
	WidgetID  string
}

func (s *Service) Delete(ctx context.Context, input DeleteInput) error {
	widgets, err := s.dashboardWidgetsRepository.GetByChannelID(ctx, input.ChannelID)
	if err != nil {
		return err
	}

	var updatedWidgets []dashboard_widget.DashboardWidget
	found := false
	for _, widget := range widgets {
		if widget.WidgetID != input.WidgetID {
			updatedWidgets = append(updatedWidgets, widget)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("widget not found: %s", input.WidgetID)
	}

	for i := range widgets {
		if widgets[i].WidgetID == input.WidgetID {
			widgets[i].Visible = false
			widgets[i].X = -1000
			widgets[i].Y = -1000
			err = s.dashboardWidgetsRepository.UpsertMany(
				ctx,
				input.ChannelID,
				[]dashboard_widget.DashboardWidget{widgets[i]},
			)
			if err != nil {
				return err
			}
			break
		}
	}

	widgets, err = s.dashboardWidgetsRepository.GetByChannelID(ctx, input.ChannelID)
	if err != nil {
		return err
	}

	go func() {
		s.wsRouter.Publish(
			DashboardWidgetsLayoutSubscriptionKeyCreate(input.ChannelID),
			widgets,
		)
	}()

	return nil
}
