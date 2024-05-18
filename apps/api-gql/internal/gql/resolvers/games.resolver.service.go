package resolvers

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
)

func (r *queryResolver) gamesGetEightBall(ctx context.Context) (*gqlmodel.EightBallGame, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	entity := model.ChannelGames8Ball{
		ChannelId: dashboardId,
		Enabled:   false,
		Answers:   pq.StringArray{},
	}
	if err := r.gorm.
		WithContext(ctx).
		Where(`"channel_id" = ?`, dashboardId).
		FirstOrCreate(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("failed to get 8ball settings: %w", err)
	}

	return &gqlmodel.EightBallGame{
		Enabled: entity.Enabled,
		Answers: entity.Answers,
	}, nil
}

func (r *mutationResolver) gamesUpdateEightBall(
	ctx context.Context,
	opts gqlmodel.EightBallGameOpts,
) (*gqlmodel.EightBallGame, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	entity := model.ChannelGames8Ball{}
	if err := r.gorm.
		WithContext(ctx).
		Where(`"channel_id" = ?`, dashboardId).
		First(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("failed to get eight ball settings: %w", err)
	}

	if opts.Answers.IsSet() {
		if len(opts.Answers.Value()) > 25 {
			return nil, fmt.Errorf("max answers is 25")
		}

		for _, answer := range opts.Answers.Value() {
			if len(answer) > 500 {
				return nil, fmt.Errorf("max answer length is 500")
			}
		}

		entity.Answers = opts.Answers.Value()
	}

	if opts.Enabled.IsSet() {
		entity.Enabled = *opts.Enabled.Value()
	}

	if err := r.gorm.
		WithContext(ctx).
		Save(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("failed to save settings: %w", err)
	}

	return r.Query().GamesEightBall(ctx)
}

func (r *queryResolver) gamesGetDuel(ctx context.Context) (*gqlmodel.DuelGame, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	entity := model.ChannelGamesDuel{
		ChannelID:       dashboardId,
		StartMessage:    "@{target}, @{initiator} challenges you to a fight. Use {duelAcceptCommandName} for next {acceptSeconds} seconds to accept the challenge.",
		ResultMessage:   "Sadly, @{loser} couldn't find a way to dodge the bullet and falls apart into eternal slumber.",
		BothDieMessage:  "Unexpectedly @{initiator} and @{target} shoot each other. Only the time knows why this happened...",
		SecondsToAccept: 60,
		TimeoutSeconds:  600,
	}
	if err := r.gorm.
		WithContext(ctx).
		Where(`"channel_id" = ?`, dashboardId).
		FirstOrCreate(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("failed to get duel settings: %w", err)
	}

	return &gqlmodel.DuelGame{
		Enabled:         entity.Enabled,
		UserCooldown:    int(entity.UserCooldown),
		GlobalCooldown:  int(entity.GlobalCooldown),
		TimeoutSeconds:  int(entity.TimeoutSeconds),
		StartMessage:    entity.StartMessage,
		ResultMessage:   entity.ResultMessage,
		SecondsToAccept: int(entity.SecondsToAccept),
		PointsPerWin:    int(entity.PointsPerWin),
		PointsPerLose:   int(entity.PointsPerLose),
		BothDiePercent:  int(entity.BothDiePercent),
		BothDieMessage:  entity.BothDieMessage,
	}, nil
}

func (r *mutationResolver) gamesUpdateDuel(
	ctx context.Context,
	opts gqlmodel.DuelGameOpts,
) (*gqlmodel.DuelGame, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	entity := model.ChannelGamesDuel{}
	if err := r.gorm.
		WithContext(ctx).
		Where(`"channel_id" = ?`, dashboardId).
		First(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("failed to get duel settings: %w", err)
	}

	if opts.Enabled.IsSet() {
		entity.Enabled = *opts.Enabled.Value()
	}

	if opts.UserCooldown.IsSet() {
		entity.UserCooldown = int32(*opts.UserCooldown.Value())
	}

	if opts.GlobalCooldown.IsSet() {
		entity.GlobalCooldown = int32(*opts.GlobalCooldown.Value())
	}

	if opts.TimeoutSeconds.IsSet() {
		entity.TimeoutSeconds = int32(*opts.TimeoutSeconds.Value())
	}

	if opts.StartMessage.IsSet() {
		entity.StartMessage = *opts.StartMessage.Value()

		if len(entity.StartMessage) > 500 {
			return nil, fmt.Errorf("max start message length is 500")
		}
	}

	if opts.ResultMessage.IsSet() {
		entity.ResultMessage = *opts.ResultMessage.Value()

		if len(entity.ResultMessage) > 500 {
			return nil, fmt.Errorf("max result message length is 500")
		}
	}

	if opts.SecondsToAccept.IsSet() {
		entity.SecondsToAccept = int32(*opts.SecondsToAccept.Value())
	}

	if opts.PointsPerWin.IsSet() {
		entity.PointsPerWin = int32(*opts.PointsPerWin.Value())
	}

	if opts.PointsPerLose.IsSet() {
		entity.PointsPerLose = int32(*opts.PointsPerLose.Value())
	}

	if opts.BothDiePercent.IsSet() {
		entity.BothDiePercent = int32(*opts.BothDiePercent.Value())
	}

	if opts.BothDieMessage.IsSet() {
		entity.BothDieMessage = *opts.BothDieMessage.Value()

		if len(entity.BothDieMessage) > 500 {
			return nil, fmt.Errorf("max both die message length is 500")

		}
	}

	if err := r.gorm.
		WithContext(ctx).
		Save(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("failed to save settings: %w", err)
	}

	return r.Query().GamesDuel(ctx)
}

func (r *queryResolver) gamesGetRussianRoulette(ctx context.Context) (
	*gqlmodel.
		RussianRouletteGame, error,
) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	entity := model.ChannelGamesRussianRoulette{
		ChannelID:             dashboardId,
		InitMessage:           "{sender} has initiated a game of roulette. Is luck on their side?",
		SurviveMessage:        "{sender} survives the game of roulette! Luck smiles upon them.",
		DeathMessage:          "{sender} couldn't make it through the game of roulette. Unfortunately, luck wasn't on their side this time.",
		TimeoutSeconds:        60,
		TumberSize:            6,
		DecisionSeconds:       2,
		ChargedBullets:        1,
		CanBeUsedByModerators: false,
	}
	if err := r.gorm.
		WithContext(ctx).
		Where(`"channel_id" = ?`, dashboardId).
		FirstOrCreate(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("failed to get russian roulette settings: %w", err)
	}

	return &gqlmodel.RussianRouletteGame{
		Enabled:              entity.Enabled,
		CanBeUsedByModerator: entity.CanBeUsedByModerators,
		TimeoutSeconds:       entity.TimeoutSeconds,
		DecisionSeconds:      entity.DecisionSeconds,
		InitMessage:          entity.InitMessage,
		SurviveMessage:       entity.SurviveMessage,
		DeathMessage:         entity.DeathMessage,
		ChargedBullets:       entity.ChargedBullets,
		TumberSize:           entity.TumberSize,
	}, nil
}

func (r *mutationResolver) gamesUpdateRussianRoulette(
	ctx context.Context,
	opts gqlmodel.RussianRouletteGameOpts,
) (*gqlmodel.RussianRouletteGame, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	entity := model.ChannelGamesRussianRoulette{}
	if err := r.gorm.
		WithContext(ctx).
		Where(`"channel_id" = ?`, dashboardId).
		First(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("failed to get russian roulette settings: %w", err)
	}

	if opts.Enabled.IsSet() {
		entity.Enabled = *opts.Enabled.Value()
	}

	if opts.CanBeUsedByModerator.IsSet() {
		entity.CanBeUsedByModerators = *opts.CanBeUsedByModerator.Value()
	}

	if opts.TimeoutSeconds.IsSet() {
		entity.TimeoutSeconds = *opts.TimeoutSeconds.Value()
	}

	if opts.DecisionSeconds.IsSet() {
		entity.DecisionSeconds = *opts.DecisionSeconds.Value()
	}

	if opts.TumberSize.IsSet() {
		entity.TumberSize = *opts.TumberSize.Value()
	}

	if opts.ChargedBullets.IsSet() {
		entity.ChargedBullets = *opts.ChargedBullets.Value()
	}

	if opts.InitMessage.IsSet() {
		entity.InitMessage = *opts.InitMessage.Value()

		if len(entity.InitMessage) > 500 {
			return nil, fmt.Errorf("max init message length is 500")
		}
	}

	if opts.SurviveMessage.IsSet() {
		entity.SurviveMessage = *opts.SurviveMessage.Value()

		if len(entity.SurviveMessage) > 500 {
			return nil, fmt.Errorf("max survive message length is 500")
		}
	}

	if opts.DeathMessage.IsSet() {
		entity.DeathMessage = *opts.DeathMessage.Value()

		if len(entity.DeathMessage) > 500 {
			return nil, fmt.Errorf("max death message length is 500")
		}
	}

	if err := r.gorm.
		WithContext(ctx).
		Save(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("failed to save settings: %w", err)
	}

	return r.Query().GamesRussianRoulette(ctx)
}

func (r *queryResolver) gamesSeppuku(ctx context.Context) (*gqlmodel.SeppukuGame, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	entity := model.ChannelGamesSeppuku{
		ChannelID:         dashboardId,
		Enabled:           false,
		TimeoutSeconds:    60,
		TimeoutModerators: false,
		Message:           "{sender} said: my honor tarnished, I reclaim it through death. May my spirit find peace. Farewell.",
		MessageModerators: "{sender} drew his sword and ripped open his belly for the sad emperor.",
	}
	if err := r.gorm.
		WithContext(ctx).
		Where(`"channel_id" = ?`, dashboardId).
		FirstOrCreate(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("failed to get seppuku settings: %w", err)
	}

	return &gqlmodel.SeppukuGame{
		Enabled:           entity.Enabled,
		TimeoutSeconds:    int(entity.TimeoutSeconds),
		TimeoutModerators: entity.TimeoutModerators,
		Message:           entity.Message,
		MessageModerators: entity.MessageModerators,
	}, nil
}

func (r *mutationResolver) gamesUpdateSeppuku(
	ctx context.Context,
	opts gqlmodel.SeppukuGameOpts,
) (*gqlmodel.SeppukuGame, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	entity := model.ChannelGamesSeppuku{}
	if err := r.gorm.
		WithContext(ctx).
		Where(`"channel_id" = ?`, dashboardId).
		First(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("failed to get seppuku settings: %w", err)
	}

	if opts.Enabled.IsSet() {
		entity.Enabled = *opts.Enabled.Value()
	}

	if opts.TimeoutSeconds.IsSet() {
		entity.TimeoutSeconds = int8(*opts.TimeoutSeconds.Value())
	}

	if opts.TimeoutModerators.IsSet() {
		entity.TimeoutModerators = *opts.TimeoutModerators.Value()
	}

	if opts.Message.IsSet() {
		entity.Message = *opts.Message.Value()
	}

	if opts.MessageModerators.IsSet() {
		entity.MessageModerators = *opts.MessageModerators.Value()
	}

	if err := r.gorm.
		WithContext(ctx).
		Save(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("failed to save settings: %w", err)
	}

	return r.Query().GamesSeppuku(ctx)
}
