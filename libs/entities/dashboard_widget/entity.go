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
	// StackId groups widgets into tabs - widgets with same StackId are displayed as tabs
	StackId *string
	// StackOrder determines the order of tabs within a stack (0, 1, 2, etc.)
	StackOrder int
	CreatedAt  time.Time
	UpdatedAt  time.Time

	isNil bool
}

func (d DashboardWidget) IsNil() bool {
	return d.isNil
}

var Nil = DashboardWidget{
	isNil: true,
}
