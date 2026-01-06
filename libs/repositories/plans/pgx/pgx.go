package pgx

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/entities/plan"
	"github.com/twirapp/twir/libs/repositories/plans"
	"github.com/twirapp/twir/libs/repositories/plans/model"
)

type Opts struct {
	Db *pgxpool.Pool
}

func New(opts Opts) plans.Repository {
	return &repository{
		db: opts.Db,
	}
}

func NewFx(pool *pgxpool.Pool) plans.Repository {
	return New(Opts{Db: pool})
}

type repository struct {
	db *pgxpool.Pool
}

func (r *repository) GetByID(ctx context.Context, id string) (plan.Plan, error) {
	query, args, err := squirrel.Select(
		"id",
		"name",
		"max_commands",
		"max_timers",
		"max_variables",
		"max_alerts",
		"max_events",
		"max_chat_alerts_messages",
		"max_custom_overlays",
		"max_eightball_answers",
		"max_commands_responses",
		"max_moderation_rules",
		"max_keywords",
		"max_greetings",
		"created_at",
		"updated_at",
	).
		From("plans").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return plan.Nil, fmt.Errorf("failed to build query: %w", err)
	}

	var dbPlan model.Plan
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&dbPlan.ID,
		&dbPlan.Name,
		&dbPlan.MaxCommands,
		&dbPlan.MaxTimers,
		&dbPlan.MaxVariables,
		&dbPlan.MaxAlerts,
		&dbPlan.MaxEvents,
		&dbPlan.MaxChatAlertsMessages,
		&dbPlan.MaxCustomOverlays,
		&dbPlan.MaxEightballAnswers,
		&dbPlan.MaxCommandsResponses,
		&dbPlan.MaxModerationRules,
		&dbPlan.MaxKeywords,
		&dbPlan.MaxGreetings,
		&dbPlan.CreatedAt,
		&dbPlan.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return plan.Nil, nil
		}
		return plan.Nil, fmt.Errorf("failed to get plan: %w", err)
	}

	return r.dbToEntity(dbPlan), nil
}

func (r *repository) GetByNameID(ctx context.Context, nameID string) (plan.Plan, error) {
	query, args, err := squirrel.Select(
		"id",
		"name",
		"max_commands",
		"max_timers",
		"max_variables",
		"max_alerts",
		"max_events",
		"max_chat_alerts_messages",
		"max_custom_overlays",
		"max_eightball_answers",
		"max_commands_responses",
		"max_moderation_rules",
		"max_keywords",
		"max_greetings",
		"created_at",
		"updated_at",
	).
		From("plans").
		Where(squirrel.Eq{"name": nameID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return plan.Nil, fmt.Errorf("failed to build query: %w", err)
	}

	var dbPlan model.Plan
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&dbPlan.ID,
		&dbPlan.Name,
		&dbPlan.MaxCommands,
		&dbPlan.MaxTimers,
		&dbPlan.MaxVariables,
		&dbPlan.MaxAlerts,
		&dbPlan.MaxEvents,
		&dbPlan.MaxChatAlertsMessages,
		&dbPlan.MaxCustomOverlays,
		&dbPlan.MaxEightballAnswers,
		&dbPlan.MaxCommandsResponses,
		&dbPlan.MaxModerationRules,
		&dbPlan.MaxKeywords,
		&dbPlan.MaxGreetings,
		&dbPlan.CreatedAt,
		&dbPlan.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return plan.Nil, nil
		}
		return plan.Nil, fmt.Errorf("failed to get plan: %w", err)
	}

	return r.dbToEntity(dbPlan), nil
}

func (r *repository) GetByChannelID(ctx context.Context, channelID string) (plan.Plan, error) {
	query, args, err := squirrel.Select(
		"p.id",
		"p.name",
		"p.max_commands",
		"p.max_timers",
		"p.max_variables",
		"p.max_alerts",
		"p.max_events",
		"p.max_chat_alerts_messages",
		"p.max_custom_overlays",
		"p.max_eightball_answers",
		"p.max_commands_responses",
		"p.max_moderation_rules",
		"p.max_keywords",
		"p.max_greetings",
		"p.created_at",
		"p.updated_at",
	).
		From("plans p").
		Join("channels c ON c.plan_id = p.id").
		Where(squirrel.Eq{"c.id": channelID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return plan.Nil, fmt.Errorf("failed to build query: %w", err)
	}

	var dbPlan model.Plan
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&dbPlan.ID,
		&dbPlan.Name,
		&dbPlan.MaxCommands,
		&dbPlan.MaxTimers,
		&dbPlan.MaxVariables,
		&dbPlan.MaxAlerts,
		&dbPlan.MaxEvents,
		&dbPlan.MaxChatAlertsMessages,
		&dbPlan.MaxCustomOverlays,
		&dbPlan.MaxEightballAnswers,
		&dbPlan.MaxCommandsResponses,
		&dbPlan.MaxModerationRules,
		&dbPlan.MaxKeywords,
		&dbPlan.MaxGreetings,
		&dbPlan.CreatedAt,
		&dbPlan.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return plan.Nil, nil
		}
		return plan.Nil, fmt.Errorf("failed to get plan: %w", err)
	}

	return r.dbToEntity(dbPlan), nil
}

func (r *repository) dbToEntity(m model.Plan) plan.Plan {
	return plan.Plan{
		ID:                    m.ID,
		Name:                  m.Name,
		MaxCommands:           m.MaxCommands,
		MaxTimers:             m.MaxTimers,
		MaxVariables:          m.MaxVariables,
		MaxAlerts:             m.MaxAlerts,
		MaxEvents:             m.MaxEvents,
		MaxChatAlertsMessages: m.MaxChatAlertsMessages,
		MaxCustomOverlays:     m.MaxCustomOverlays,
		MaxEightballAnswers:   m.MaxEightballAnswers,
		MaxCommandsResponses:  m.MaxCommandsResponses,
		MaxModerationRules:    m.MaxModerationRules,
		MaxKeywords:           m.MaxKeywords,
		MaxGreetings:          m.MaxGreetings,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
	}
}
