package plan

import "time"

type Plan struct {
	ID                    string
	Name                  string
	MaxCommands           int
	MaxTimers             int
	MaxVariables          int
	MaxAlerts             int
	MaxEvents             int
	MaxChatAlertsMessages int
	MaxCustomOverlays     int
	MaxEightballAnswers   int
	MaxCommandsResponses  int
	MaxModerationRules    int
	MaxKeywords           int
	MaxGreetings          int
	CreatedAt             time.Time
	UpdatedAt             time.Time

	isNil bool
}

func (p Plan) IsNil() bool {
	return p.isNil
}

var Nil = Plan{
	isNil: true,
}
