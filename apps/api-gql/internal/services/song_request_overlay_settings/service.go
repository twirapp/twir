package song_request_overlay_settings

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"

	entity "github.com/twirapp/twir/libs/entities/song_request_overlay_settings"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/song_request_overlay_settings"
	"github.com/twirapp/twir/libs/wsrouter"
	"go.uber.org/fx"
)

var colorPattern = regexp.MustCompile(`^#[0-9A-Fa-f]{6}([0-9A-Fa-f]{2})?$`)

type Opts struct {
	fx.In

	Repository song_request_overlay_settings.Repository
	WsRouter   wsrouter.WsRouter
	Logger     *slog.Logger
}

type Service struct {
	repository song_request_overlay_settings.Repository
	wsRouter   wsrouter.WsRouter
	logger     *slog.Logger
}

func New(opts Opts) *Service {
	return &Service{
		repository: opts.Repository,
		wsRouter:   opts.WsRouter,
		logger:     opts.Logger,
	}
}

func SubscriptionKey(channelID string) string {
	return "songrequests:overlay-settings:" + channelID
}

func (s *Service) Get(
	ctx context.Context,
	channelID string,
) (entity.SongRequestOverlaySettings, error) {
	settings, err := s.repository.GetByChannelID(ctx, channelID)
	if err != nil {
		if errors.Is(err, song_request_overlay_settings.ErrNotFound) {
			return entity.Default(channelID), nil
		}
		return entity.Nil, fmt.Errorf("get song request overlay settings: %w", err)
	}

	return settings, nil
}

type UpdateInput struct {
	ChannelID             string
	Style                 entity.Style
	AccentColor           string
	TickerBackgroundColor string
	TickerTextColor       string
	TickerSpeed           int
	HideOnPause           bool
}

func validate(input UpdateInput) error {
	if !input.Style.IsValid() {
		return fmt.Errorf("invalid overlay style %q", input.Style)
	}

	colors := map[string]string{
		"accentColor":           input.AccentColor,
		"tickerBackgroundColor": input.TickerBackgroundColor,
		"tickerTextColor":       input.TickerTextColor,
	}
	for name, color := range colors {
		if !colorPattern.MatchString(color) {
			return fmt.Errorf("%s must be #RRGGBB or #RRGGBBAA", name)
		}
	}

	if input.TickerSpeed < 10 || input.TickerSpeed > 100 {
		return fmt.Errorf("tickerSpeed must be between 10 and 100")
	}

	return nil
}

func (s *Service) Update(
	ctx context.Context,
	input UpdateInput,
) (entity.SongRequestOverlaySettings, error) {
	if err := validate(input); err != nil {
		return entity.Nil, err
	}

	settings, err := s.repository.Upsert(ctx, song_request_overlay_settings.UpsertInput{
		ChannelID:             input.ChannelID,
		Style:                 input.Style,
		AccentColor:           input.AccentColor,
		TickerBackgroundColor: input.TickerBackgroundColor,
		TickerTextColor:       input.TickerTextColor,
		TickerSpeed:           input.TickerSpeed,
		HideOnPause:           input.HideOnPause,
	})
	if err != nil {
		return entity.Nil, fmt.Errorf("update song request overlay settings: %w", err)
	}

	if err := s.wsRouter.Publish(SubscriptionKey(input.ChannelID), settings); err != nil {
		s.logger.ErrorContext(
			ctx,
			"failed to publish song request overlay settings",
			slog.String("channelId", input.ChannelID),
			logger.Error(err),
		)
	}

	return settings, nil
}
