package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/libs/entities/dashboard_widget"
)

func DashboardWidgetEntityToGQL(entity dashboard_widget.DashboardWidget) gqlmodel.DashboardWidgetLayout {
	return gqlmodel.DashboardWidgetLayout{
		ID:       entity.ID,
		WidgetID: entity.WidgetID,
		X:        entity.X,
		Y:        entity.Y,
		W:        entity.W,
		H:        entity.H,
		MinW:     entity.MinW,
		MinH:     entity.MinH,
		Visible:  entity.Visible,
	}
}

func DashboardWidgetGQLToEntity(input gqlmodel.DashboardWidgetLayoutInput) dashboard_widget.DashboardWidget {
	return dashboard_widget.DashboardWidget{
		WidgetID: input.WidgetID,
		X:        input.X,
		Y:        input.Y,
		W:        input.W,
		H:        input.H,
		MinW:     input.MinW,
		MinH:     input.MinH,
		Visible:  input.Visible,
	}
}
