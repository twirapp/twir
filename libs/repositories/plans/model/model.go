package model

import "time"

type Plan struct {
	ID                          string    `db:"id"`
	Name                        string    `db:"name"`
	MaxCommands                 int       `db:"max_commands"`
	MaxTimers                   int       `db:"max_timers"`
	MaxVariables                int       `db:"max_variables"`
	MaxAlerts                   int       `db:"max_alerts"`
	MaxEvents                   int       `db:"max_events"`
	MaxChatAlertsMessages       int       `db:"max_chat_alerts_messages"`
	MaxCustomOverlays           int       `db:"max_custom_overlays"`
	MaxEightballAnswers         int       `db:"max_eightball_answers"`
	MaxCommandsResponses        int       `db:"max_commands_responses"`
	MaxModerationRules          int       `db:"max_moderation_rules"`
	MaxKeywords                 int       `db:"max_keywords"`
	MaxGreetings                int       `db:"max_greetings"`
	LinksShortenerCustomDomains int       `db:"links_shortener_custom_domains"`
	CreatedAt                   time.Time `db:"created_at"`
	UpdatedAt                   time.Time `db:"updated_at"`

	isNil bool
}

func (p Plan) IsNil() bool {
	return p.isNil
}

var Nil = Plan{
	isNil: true,
}
