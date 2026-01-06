package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/libs/entities/plan"
)

func PlanToGql(e plan.Plan) *gqlmodel.Plan {
	return &gqlmodel.Plan{
		ID:                    e.ID,
		Name:                  e.Name,
		MaxCommands:           e.MaxCommands,
		MaxTimers:             e.MaxTimers,
		MaxVariables:          e.MaxVariables,
		MaxAlerts:             e.MaxAlerts,
		MaxEvents:             e.MaxEvents,
		MaxChatAlertsMessages: e.MaxChatAlertsMessages,
		MaxCustomOverlays:     e.MaxCustomOverlays,
		MaxEightballAnswers:   e.MaxEightballAnswers,
		MaxCommandsResponses:  e.MaxCommandsResponses,
		MaxModerationRules:    e.MaxModerationRules,
		MaxKeywords:           e.MaxKeywords,
		MaxGreetings:          e.MaxGreetings,
	}
}
