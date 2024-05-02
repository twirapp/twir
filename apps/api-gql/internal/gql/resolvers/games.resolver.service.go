package resolvers

import (
	"context"
	"fmt"

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
	}

	if opts.ResultMessage.IsSet() {
		entity.ResultMessage = *opts.ResultMessage.Value()
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
	}

	if err := r.gorm.
		WithContext(ctx).
		Save(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("failed to save settings: %w", err)
	}

	return r.Query().GamesDuel(ctx)
}
