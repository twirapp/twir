package predictions

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	buscore "github.com/twirapp/twir/libs/bus-core"
	busdota "github.com/twirapp/twir/libs/bus-core/dota"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	dotarepository "github.com/twirapp/twir/libs/repositories/dota"
	dotamodel "github.com/twirapp/twir/libs/repositories/dota/model"
	"github.com/twirapp/twir/libs/twitch"
	"go.uber.org/fx"
)

const (
	predictionKeyPrefix          = "cache:twir:dota:prediction:"
	predictionTTL                = 12 * time.Hour
	predictionsSubscriptionGroup = "dota-predictions"
)

var (
	errPredictionNotFound = errors.New("dota prediction record not found")
	errPredictionPending  = errors.New("dota prediction creation pending")
)

type storedPrediction struct {
	PredictionID string `json:"predictionId"`
	YesOutcomeID string `json:"yesOutcomeId"`
	NoOutcomeID  string `json:"noOutcomeId"`
}

// Store keeps the Dota prediction that belongs to one channel and match.
type Store interface {
	Reserve(ctx context.Context, key string, token string, ttl time.Duration) (bool, error)
	Commit(ctx context.Context, key string, record storedPrediction, ttl time.Duration) error
	Get(ctx context.Context, key string) (storedPrediction, error)
	Release(ctx context.Context, key string, token string) error
	Delete(ctx context.Context, key string) error
}

type settingsRepository interface {
	GetByChannelID(ctx context.Context, channelID uuid.UUID) (dotamodel.ChannelDotaSettings, error)
}

type channelRepository interface {
	GetByID(ctx context.Context, channelID uuid.UUID) (channelsmodel.Channel, error)
}

type predictionClient interface {
	CreatePrediction(params *helix.CreatePredictionParams) (*helix.PredictionsResponse, error)
	GetPredictions(params *helix.PredictionsParams) (*helix.PredictionsResponse, error)
	EndPrediction(params *helix.EndPredictionParams) (*helix.PredictionsResponse, error)
}

type clientFactory interface {
	New(ctx context.Context, userID uuid.UUID) (predictionClient, error)
}

type twitchClientFactory struct {
	config cfg.Config
	bus    *buscore.Bus
}

func (f twitchClientFactory) New(ctx context.Context, userID uuid.UUID) (predictionClient, error) {
	return twitch.NewUserClientWithContext(ctx, userID, f.config, f.bus)
}

type subscription interface {
	Subscribe(group string) error
	Unsubscribe()
}

type subscriptionFunc struct {
	subscribe   func(group string) error
	unsubscribe func()
}

func (s subscriptionFunc) Subscribe(group string) error {
	return s.subscribe(group)
}

func (s subscriptionFunc) Unsubscribe() {
	s.unsubscribe()
}

type Opts struct {
	fx.In

	Lifecycle fx.Lifecycle
	Bus       *buscore.Bus
	Config    cfg.Config
	Logger    *slog.Logger

	SettingsRepository dotarepository.Repository
	ChannelsRepository channelsrepository.Repository
	Store              Store
}

type Predictions struct {
	settings settingsRepository
	channels channelRepository
	clients  clientFactory
	store    Store
	logger   *slog.Logger

	subscriptions []subscription

	handlersMu sync.Mutex
	handlers   sync.WaitGroup
	stopping   bool
}

func New(opts Opts) *Predictions {
	predictions := newPredictions(
		opts.SettingsRepository,
		opts.ChannelsRepository,
		twitchClientFactory{config: opts.Config, bus: opts.Bus},
		opts.Store,
		opts.Logger,
		nil,
		opts.Lifecycle,
	)
	predictions.subscriptions = newBusSubscriptions(opts.Bus, predictions)

	return predictions
}

func newPredictions(
	settings settingsRepository,
	channels channelRepository,
	clients clientFactory,
	store Store,
	logger *slog.Logger,
	subscriptions []subscription,
	lifecycle fx.Lifecycle,
) *Predictions {
	predictions := &Predictions{
		settings:      settings,
		channels:      channels,
		clients:       clients,
		store:         store,
		logger:        logger,
		subscriptions: subscriptions,
	}

	if lifecycle != nil {
		lifecycle.Append(fx.Hook{
			OnStart: predictions.Start,
			OnStop:  predictions.Stop,
		})
	}

	return predictions
}

func newBusSubscriptions(bus *buscore.Bus, predictions *Predictions) []subscription {
	return []subscription{
		subscriptionFunc{
			subscribe: func(group string) error {
				return bus.Dota.MatchStarted.SubscribeGroup(group, predictions.handleMatchStarted)
			},
			unsubscribe: bus.Dota.MatchStarted.Unsubscribe,
		},
		subscriptionFunc{
			subscribe: func(group string) error {
				return bus.Dota.MatchEnded.SubscribeGroup(group, predictions.handleMatchEnded)
			},
			unsubscribe: bus.Dota.MatchEnded.Unsubscribe,
		},
		subscriptionFunc{
			subscribe: func(group string) error {
				return bus.Dota.MatchAbandoned.SubscribeGroup(group, predictions.handleMatchAbandoned)
			},
			unsubscribe: bus.Dota.MatchAbandoned.Unsubscribe,
		},
	}
}

func (p *Predictions) Start(_ context.Context) error {
	started := make([]subscription, 0, len(p.subscriptions))
	for _, subscription := range p.subscriptions {
		if err := subscription.Subscribe(predictionsSubscriptionGroup); err != nil {
			for i := len(started) - 1; i >= 0; i-- {
				started[i].Unsubscribe()
			}
			return fmt.Errorf("subscribe to dota predictions: %w", err)
		}
		started = append(started, subscription)
	}

	return nil
}

func (p *Predictions) Stop(ctx context.Context) error {
	p.handlersMu.Lock()
	p.stopping = true
	for _, subscription := range p.subscriptions {
		subscription.Unsubscribe()
	}
	p.handlersMu.Unlock()

	handlersDone := make(chan struct{})
	go func() {
		p.handlers.Wait()
		close(handlersDone)
	}()

	select {
	case <-handlersDone:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *Predictions) handleMatchStarted(
	ctx context.Context,
	message busdota.MatchStartedMessage,
) (struct{}, error) {
	return struct{}{}, p.handleTracked(ctx, "match_started", func() error {
		return p.createPrediction(ctx, message)
	})
}

func (p *Predictions) handleMatchEnded(
	ctx context.Context,
	message busdota.MatchEndedMessage,
) (struct{}, error) {
	return struct{}{}, p.handleTracked(ctx, "match_ended", func() error {
		return p.finishPrediction(ctx, message.ChannelID, message.MatchID, message.Win, false)
	})
}

func (p *Predictions) handleMatchAbandoned(
	ctx context.Context,
	message busdota.MatchAbandonedMessage,
) (struct{}, error) {
	return struct{}{}, p.handleTracked(ctx, "match_abandoned", func() error {
		return p.finishPrediction(ctx, message.ChannelID, message.MatchID, false, true)
	})
}

func (p *Predictions) handleTracked(ctx context.Context, eventKind string, handler func() error) error {
	p.handlersMu.Lock()
	if p.stopping {
		p.handlersMu.Unlock()
		return context.Canceled
	}
	p.handlers.Add(1)
	p.handlersMu.Unlock()
	defer p.handlers.Done()

	if err := handler(); err != nil {
		p.logger.ErrorContext(
			ctx,
			"dota prediction handler failed",
			slog.String("event_kind", eventKind),
			logger.Error(err),
		)
		return err
	}

	return nil
}

func (p *Predictions) createPrediction(ctx context.Context, message busdota.MatchStartedMessage) error {
	if message.MatchID <= 0 || !message.TeamKnown {
		return nil
	}

	channelID, err := uuid.Parse(message.ChannelID)
	if err != nil {
		p.logger.WarnContext(ctx, "dota prediction skipped: invalid channel ID", logger.Error(err))
		return nil
	}

	settings, err := p.settings.GetByChannelID(ctx, channelID)
	if errors.Is(err, dotarepository.ErrNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("get dota settings: %w", err)
	}
	if !settings.Enabled || !settings.PredictionSettings.Enabled {
		return nil
	}

	title := strings.TrimSpace(settings.PredictionSettings.TitleTemplate)
	window := settings.PredictionSettings.WindowSeconds
	if title == "" || window < 1 || window > 1800 {
		p.logger.WarnContext(
			ctx,
			"dota prediction skipped: invalid settings",
			slog.String("channel_id", channelID.String()),
			slog.Int("window_seconds", window),
		)
		return nil
	}

	channel, err := p.channels.GetByID(ctx, channelID)
	if errors.Is(err, channelsrepository.ErrNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("get channel: %w", err)
	}
	if !channel.TwitchConnected() || channel.TwitchUserID == nil || channel.TwitchPlatformID == nil ||
		strings.TrimSpace(*channel.TwitchPlatformID) == "" {
		return nil
	}

	key := predictionKey(channelID, message.MatchID)
	token := uuid.NewString()
	reserved, err := p.store.Reserve(ctx, key, token, predictionTTL)
	if err != nil {
		return fmt.Errorf("reserve prediction: %w", err)
	}
	if !reserved {
		return nil
	}

	client, err := p.clients.New(ctx, *channel.TwitchUserID)
	if err != nil {
		return p.creationFailed(ctx, key, token, fmt.Errorf("create Twitch client: %w", err))
	}

	response, err := client.CreatePrediction(&helix.CreatePredictionParams{
		BroadcasterID:    strings.TrimSpace(*channel.TwitchPlatformID),
		Title:            title,
		PredictionWindow: window,
		Outcomes: []helix.PredictionChoiceParam{
			{Title: "Yes"},
			{Title: "No"},
		},
	})
	if err != nil {
		if isAlreadyActive(err) {
			if releaseErr := p.releaseReservation(ctx, key, token); releaseErr != nil {
				return releaseErr
			}
			return nil
		}
		return p.creationFailed(ctx, key, token, fmt.Errorf("create Twitch prediction: %w", err))
	}
	if err := predictionResponseError(response); err != nil {
		if isAlreadyActive(err) {
			if releaseErr := p.releaseReservation(ctx, key, token); releaseErr != nil {
				return releaseErr
			}
			return nil
		}
		return p.creationFailed(ctx, key, token, fmt.Errorf("create Twitch prediction: %w", err))
	}

	record, err := recordFromCreateResponse(response)
	if err != nil {
		return p.creationFailed(ctx, key, token, err)
	}
	if err := p.store.Commit(ctx, key, record, predictionTTL); err != nil {
		return p.creationFailed(ctx, key, token, fmt.Errorf("store created prediction: %w", err))
	}

	return nil
}

func (p *Predictions) creationFailed(ctx context.Context, key string, token string, cause error) error {
	if err := p.releaseReservation(ctx, key, token); err != nil {
		return fmt.Errorf("%w; release prediction reservation: %v", cause, err)
	}
	return cause
}

func (p *Predictions) releaseReservation(ctx context.Context, key string, token string) error {
	cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
	defer cancel()
	if err := p.store.Release(cleanupCtx, key, token); err != nil {
		p.logger.ErrorContext(cleanupCtx, "dota prediction reservation cleanup failed", logger.Error(err))
		return err
	}
	return nil
}

func (p *Predictions) finishPrediction(
	ctx context.Context,
	channelIDValue string,
	matchID int64,
	win bool,
	cancel bool,
) error {
	if matchID <= 0 {
		return nil
	}

	channelID, err := uuid.Parse(channelIDValue)
	if err != nil {
		p.logger.WarnContext(ctx, "dota prediction skipped: invalid channel ID", logger.Error(err))
		return nil
	}

	key := predictionKey(channelID, matchID)
	record, err := p.store.Get(ctx, key)
	if errors.Is(err, errPredictionNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("get stored prediction: %w", err)
	}

	channel, err := p.channels.GetByID(ctx, channelID)
	if err != nil {
		return fmt.Errorf("get channel for stored prediction: %w", err)
	}
	if !channel.TwitchConnected() || channel.TwitchUserID == nil || channel.TwitchPlatformID == nil ||
		strings.TrimSpace(*channel.TwitchPlatformID) == "" {
		return errors.New("stored prediction channel is not connected to Twitch")
	}

	client, err := p.clients.New(ctx, *channel.TwitchUserID)
	if err != nil {
		return fmt.Errorf("create Twitch client for stored prediction: %w", err)
	}
	broadcasterID := strings.TrimSpace(*channel.TwitchPlatformID)
	response, err := client.GetPredictions(&helix.PredictionsParams{
		BroadcasterID: broadcasterID,
		ID:            record.PredictionID,
	})
	if err != nil {
		return fmt.Errorf("get Twitch prediction: %w", err)
	}
	if err := predictionResponseError(response); err != nil {
		return fmt.Errorf("get Twitch prediction: %w", err)
	}

	prediction, active := activePrediction(response, record.PredictionID)
	if !active {
		if err := p.store.Delete(ctx, key); err != nil {
			return fmt.Errorf("delete inactive prediction: %w", err)
		}
		return nil
	}

	params := &helix.EndPredictionParams{
		BroadcasterID: broadcasterID,
		ID:            prediction.ID,
	}
	if cancel {
		params.Status = "CANCELED"
	} else {
		params.Status = "RESOLVED"
		if win {
			params.WinningOutcomeID = record.YesOutcomeID
		} else {
			params.WinningOutcomeID = record.NoOutcomeID
		}
	}

	endResponse, err := client.EndPrediction(params)
	if err != nil {
		return fmt.Errorf("end Twitch prediction: %w", err)
	}
	if err := predictionResponseError(endResponse); err != nil {
		return fmt.Errorf("end Twitch prediction: %w", err)
	}
	if err := p.store.Delete(ctx, key); err != nil {
		return fmt.Errorf("delete finished prediction: %w", err)
	}

	return nil
}

func predictionKey(channelID uuid.UUID, matchID int64) string {
	return predictionKeyPrefix + channelID.String() + ":" + strconv.FormatInt(matchID, 10)
}

func predictionResponseError(response *helix.PredictionsResponse) error {
	if response == nil {
		return errors.New("empty Twitch prediction response")
	}
	if message := strings.TrimSpace(response.ErrorMessage); message != "" {
		return errors.New(message)
	}
	return nil
}

func isAlreadyActive(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "already active")
}

func recordFromCreateResponse(response *helix.PredictionsResponse) (storedPrediction, error) {
	if len(response.Data.Predictions) == 0 {
		return storedPrediction{}, errors.New("create Twitch prediction returned no prediction")
	}

	prediction := response.Data.Predictions[0]
	if strings.TrimSpace(prediction.ID) == "" {
		return storedPrediction{}, errors.New("create Twitch prediction returned an empty prediction ID")
	}

	record := storedPrediction{PredictionID: prediction.ID}
	for _, outcome := range prediction.Outcomes {
		switch outcome.Title {
		case "Yes":
			record.YesOutcomeID = outcome.ID
		case "No":
			record.NoOutcomeID = outcome.ID
		}
	}
	if strings.TrimSpace(record.YesOutcomeID) == "" || strings.TrimSpace(record.NoOutcomeID) == "" {
		return storedPrediction{}, errors.New("create Twitch prediction returned incomplete outcomes")
	}

	return record, nil
}

func activePrediction(response *helix.PredictionsResponse, predictionID string) (helix.Prediction, bool) {
	for _, prediction := range response.Data.Predictions {
		if prediction.ID != predictionID {
			continue
		}
		return prediction, prediction.Status == "ACTIVE" || prediction.Status == "LOCKED"
	}

	return helix.Prediction{}, false
}
