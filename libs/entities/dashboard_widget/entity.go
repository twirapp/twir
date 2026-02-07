package dashboard_widget

import "time"

type DashboardWidget struct {
	ID        string
	ChannelID string
	WidgetID  string
	X         int
	Y         int
	W         int
	H         int
	MinW      int
	MinH      int
	Visible   bool
	CreatedAt time.Time
	UpdatedAt time.Time

	isNil bool
}

func (d DashboardWidget) IsNil() bool {
	return d.isNil
}

var Nil = DashboardWidget{
	isNil: true,
}
