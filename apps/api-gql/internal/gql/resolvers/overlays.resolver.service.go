package resolvers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
)

func chatOverlayDbToGql(entity *model.ChatOverlaySettings) *gqlmodel.ChatOverlay {
	var animation *gqlmodel.ChatOverlayAnimation
	if entity.Animation != nil {
		animation = lo.ToPtr(gqlmodel.ChatOverlayAnimation(*entity.Animation))
	}

	return &gqlmodel.ChatOverlay{
		ID:                  entity.ID.String(),
		MessageHideTimeout:  int(entity.MessageHideTimeout),
		MessageShowDelay:    int(entity.MessageShowDelay),
		Preset:              entity.Preset,
		FontSize:            int(entity.FontSize),
		HideCommands:        entity.HideCommands,
		HideBots:            entity.HideBots,
		FontFamily:          entity.FontFamily,
		ShowBadges:          entity.ShowBadges,
		ShowAnnounceBadge:   entity.ShowAnnounceBadge,
		TextShadowColor:     entity.TextShadowColor,
		TextShadowSize:      int(entity.TextShadowSize),
		ChatBackgroundColor: entity.ChatBackgroundColor,
		Direction:           entity.Direction,
		FontWeight:          int(entity.FontWeight),
		FontStyle:           entity.FontStyle,
		PaddingContainer:    int(entity.PaddingContainer),
		Animation:           animation,
	}
}

func (r *queryResolver) chatOverlays(
	ctx context.Context,
) ([]gqlmodel.ChatOverlay, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	var entities []model.ChatOverlaySettings
	if err := r.gorm.
		WithContext(ctx).
		Where("channel_id = ?", dashboardId).
		Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("failed to get chat overlay settings: %w", err)
	}

	var result []gqlmodel.ChatOverlay
	for _, entity := range entities {
		result = append(result, *chatOverlayDbToGql(&entity))
	}

	return result, nil
}

func (r *Resolver) getChatOverlaySettings(
	ctx context.Context,
	id,
	channelId string,
) (*gqlmodel.ChatOverlay, error) {
	entity := model.ChatOverlaySettings{
		ID:        uuid.MustParse(id),
		ChannelID: channelId,
	}
	if err := r.gorm.
		WithContext(ctx).
		First(&entity).Error; err != nil {
		return nil, fmt.Errorf("failed to get chat overlay settings: %w", err)
	}

	return chatOverlayDbToGql(&entity), nil
}

func (r *mutationResolver) updateChatOverlay(
	ctx context.Context,
	opts gqlmodel.ChatOverlayUpdateOpts,
) (bool, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	entity := model.ChatOverlaySettings{}
	if err := r.gorm.
		WithContext(ctx).
		Where("channel_id = ? AND id = ?", dashboardId, opts.ID).
		First(&entity).Error; err != nil {
		return false, fmt.Errorf("failed to get chat overlay settings: %w", err)
	}

	if opts.MessageHideTimeout.IsSet() {
		entity.MessageHideTimeout = uint32(*opts.MessageHideTimeout.Value())
	}

	if opts.MessageShowDelay.IsSet() {
		entity.MessageShowDelay = uint32(*opts.MessageShowDelay.Value())
	}

	if opts.Preset.IsSet() {
		entity.Preset = *opts.Preset.Value()
	}

	if opts.FontSize.IsSet() {
		entity.FontSize = uint32(*opts.FontSize.Value())
	}

	if opts.HideCommands.IsSet() {
		entity.HideCommands = *opts.HideCommands.Value()
	}

	if opts.HideBots.IsSet() {
		entity.HideBots = *opts.HideBots.Value()
	}

	if opts.FontFamily.IsSet() {
		entity.FontFamily = *opts.FontFamily.Value()
	}

	if opts.ShowBadges.IsSet() {
		entity.ShowBadges = *opts.ShowBadges.Value()
	}

	if opts.ShowAnnounceBadge.IsSet() {
		entity.ShowAnnounceBadge = *opts.ShowAnnounceBadge.Value()
	}

	if opts.TextShadowColor.IsSet() {
		entity.TextShadowColor = *opts.TextShadowColor.Value()
	}

	if opts.TextShadowSize.IsSet() {
		entity.TextShadowSize = uint32(*opts.TextShadowSize.Value())
	}

	if opts.ChatBackgroundColor.IsSet() {
		entity.ChatBackgroundColor = *opts.ChatBackgroundColor.Value()
	}

	if opts.Direction.IsSet() {
		entity.Direction = *opts.Direction.Value()
	}

	if opts.FontWeight.IsSet() {
		entity.FontWeight = uint32(*opts.FontWeight.Value())
	}

	if opts.FontStyle.IsSet() {
		entity.FontStyle = *opts.FontStyle.Value()
	}

	if opts.PaddingContainer.IsSet() {
		entity.PaddingContainer = uint32(*opts.PaddingContainer.Value())
	}

	if opts.Animation.IsSet() {
		entity.Animation = lo.ToPtr(model.ChatOverlaySettingsAnimationType(*opts.Animation.Value()))
	}

	if err := r.gorm.
		WithContext(ctx).
		Save(&entity).Error; err != nil {
		return false, fmt.Errorf("failed to save chat overlay settings: %w", err)
	}

	if err := r.wsRouter.Publish(
		chatOverlaySubscriptionKeyCreate(entity.ID.String(), entity.ChannelID),
		chatOverlayDbToGql(&entity),
	); err != nil {
		r.logger.Error("failed to publish settings", slog.Any("err", err))
	}

	return true, nil
}
