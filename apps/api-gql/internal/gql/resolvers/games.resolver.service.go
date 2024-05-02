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
		Find(&entity).
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

	return &gqlmodel.EightBallGame{
		Enabled: entity.Enabled,
		Answers: entity.Answers,
	}, nil
}
