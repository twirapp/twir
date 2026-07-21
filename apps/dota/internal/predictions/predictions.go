package predictions

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/apps/dota/internal/match"
	buscore "github.com/twirapp/twir/libs/bus-core"
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
	predictionKeyPrefix                = "cache:twir:dota:prediction:"
	matchEndedDeliveryKeyPrefix        = "cache:twir:dota:match-ended:"
	predictionTTL                      = 12 * time.Hour
	matchEndedDeliveryTTL              = 7 * 24 * time.Hour
	terminalClaimTTL                   = 2 * time.Minute
	terminalOperationTimeout           = 5 * time.Minute
	terminalRenewTimeout               = 2 * time.Second
	minPredictionWindow                = 30
	maxPredictionWindow                = 1_800
	maxPredictionTitleRunes            = 45
	predictionCorrelationBytes         = 8
	predictionCorrelationEncodedLength = 11
	predictionCorrelationMarkerPrefix  = " [d:"
	predictionCorrelationMarkerSuffix  = "]"
	predictionCorrelationMarkerRunes   = len(predictionCorrelationMarkerPrefix) +
		predictionCorrelationEncodedLength + len(predictionCorrelationMarkerSuffix)
	pendingIntentVersion               = 2
	pendingPredictionCorrelationWindow = 2 * time.Minute
)

var (
	errPredictionNotFound        = errors.New("dota prediction record not found")
	errPredictionPending         = errors.New("dota prediction creation pending")
	errPredictionIntentNotFound  = errors.New("dota prediction intent not found")
	errPredictionReservationLost = errors.New("dota prediction reservation ownership lost")
	errPredictionRecoveryUnsafe  = errors.New("dota prediction recovery is unsafe")

	ErrTerminalInProgress    = errors.New("dota prediction terminal operation in progress")
	ErrTerminalOwnershipLost = errors.New("dota prediction terminal ownership lost")
)

type storedPrediction struct {
	PredictionID string `json:"predictionId"`
	YesOutcomeID string `json:"yesOutcomeId"`
	NoOutcomeID  string `json:"noOutcomeId"`
}

type pendingPredictionIntent struct {
	Version         int       `json:"version"`
	Token           string    `json:"token"`
	Title           string    `json:"title"`
	Correlation     string    `json:"correlation"`
	YesOutcomeTitle string    `json:"yesOutcomeTitle"`
	NoOutcomeTitle  string    `json:"noOutcomeTitle"`
	ReservedAt      time.Time `json:"reservedAt"`
}

// Store keeps the Dota prediction that belongs to one channel and match.
type Store interface {
	Reserve(ctx context.Context, key string, intent pendingPredictionIntent, ttl time.Duration) (bool, error)
	Commit(ctx context.Context, key string, token string, record storedPrediction, ttl time.Duration) error
	Get(ctx context.Context, key string) (storedPrediction, error)
	GetPending(ctx context.Context, key string) (pendingPredictionIntent, error)
	Release(ctx context.Context, key string, token string) error
	ClaimTerminal(ctx context.Context, key string, token string, ttl time.Duration) (bool, error)
	RenewTerminal(ctx context.Context, key string, token string, ttl time.Duration) (bool, error)
	CompleteTerminal(ctx context.Context, key string, token string) (bool, error)
	ReleaseTerminal(ctx context.Context, key string, token string) error
	ClaimMatchEndedDelivery(
		ctx context.Context,
		key string,
		token string,
		ttl time.Duration,
	) (matchEndedDeliveryState, error)
	CompleteMatchEndedDelivery(ctx context.Context, key string, token string, ttl time.Duration) (bool, error)
	RenewMatchEndedDelivery(ctx context.Context, key string, token string, ttl time.Duration) (bool, error)
	ReleaseMatchEndedDelivery(ctx context.Context, key string, token string) error
}

type matchEndedDeliveryState uint8

const (
	matchEndedDeliveryAcquired matchEndedDeliveryState = iota + 1
	matchEndedDeliveryPending
	matchEndedDeliveryDelivered
)

type terminalLease struct {
	ttl           time.Duration
	renewInterval time.Duration
}

func defaultTerminalLease() terminalLease {
	return terminalLease{
		ttl:           terminalClaimTTL,
		renewInterval: terminalClaimTTL / 3,
	}
}

func (l terminalLease) normalized() terminalLease {
	if l.ttl <= 0 {
		l.ttl = terminalClaimTTL
	}
	if l.renewInterval <= 0 || l.renewInterval >= l.ttl {
		l.renewInterval = l.ttl / 3
	}
	if l.renewInterval <= 0 {
		l.renewInterval = time.Nanosecond
	}

	return l
}

type terminalHeartbeat struct {
	operationCtx    context.Context
	operationCancel context.CancelFunc
	stopCh          chan struct{}
	done            chan struct{}

	stopOnce sync.Once
	errMu    sync.Mutex
	err      error
}

func (h *terminalHeartbeat) stop() error {
	h.stopOnce.Do(func() { close(h.stopCh) })
	<-h.done

	h.errMu.Lock()
	defer h.errMu.Unlock()
	return h.err
}

func (h *terminalHeartbeat) setError(err error) {
	h.errMu.Lock()
	defer h.errMu.Unlock()
	if h.err == nil {
		h.err = err
	}
}

func (h *terminalHeartbeat) fail(err error) {
	h.setError(err)
	h.operationCancel()
}

type settingsRepository interface {
	GetByChannelID(ctx context.Context, channelID uuid.UUID) (dotamodel.ChannelDotaSettings, error)
	ApplyMatchResultOnce(
		ctx context.Context,
		input dotarepository.ApplyMatchResultInput,
	) (dotamodel.ChannelDotaSettings, error)
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

type Opts struct {
	fx.In

	Bus    *buscore.Bus
	Config cfg.Config
	Logger *slog.Logger

	SettingsRepository dotarepository.Repository
	ChannelsRepository channelsrepository.Repository
	Store              Store
}

type Predictions struct {
	settings      settingsRepository
	channels      channelRepository
	clients       clientFactory
	store         Store
	logger        *slog.Logger
	terminalLease terminalLease
}

func New(opts Opts) *Predictions {
	predictions := newPredictions(
		opts.SettingsRepository,
		opts.ChannelsRepository,
		twitchClientFactory{config: opts.Config, bus: opts.Bus},
		opts.Store,
		opts.Logger,
		defaultTerminalLease(),
	)

	return predictions
}

func newPredictions(
	settings settingsRepository,
	channels channelRepository,
	clients clientFactory,
	store Store,
	logger *slog.Logger,
	terminalLease terminalLease,
) *Predictions {
	return &Predictions{
		settings:      settings,
		channels:      channels,
		clients:       clients,
		store:         store,
		logger:        logger,
		terminalLease: terminalLease.normalized(),
	}
}

func (p *Predictions) Create(ctx context.Context, action match.LifecycleAction) error {
	if !action.TeamKnown {
		return nil
	}

	return p.createPrediction(ctx, action.ChannelID, action.MatchID)
}

func (p *Predictions) Resolve(
	ctx context.Context,
	action match.LifecycleAction,
) (dotamodel.ChannelDotaSettings, error) {
	if action.ChannelID == uuid.Nil || action.MatchID <= 0 {
		return dotamodel.Nil, nil
	}

	settings, err := p.settings.GetByChannelID(ctx, action.ChannelID)
	if err != nil {
		return dotamodel.Nil, fmt.Errorf("get dota settings: %w", err)
	}

	mmrDelta := settings.MmrDelta
	if !action.Win {
		mmrDelta = -mmrDelta
	}
	settings, err = p.settings.ApplyMatchResultOnce(ctx, dotarepository.ApplyMatchResultInput{
		ChannelID: action.ChannelID,
		MatchID:   action.MatchID,
		Won:       action.Win,
		MmrDelta:  mmrDelta,
	})
	if err != nil {
		return dotamodel.Nil, fmt.Errorf("apply dota match result: %w", err)
	}

	if err := p.finishPrediction(ctx, action.ChannelID, action.MatchID, action.Win, false); err != nil {
		return settings, err
	}

	return settings, nil
}

func (p *Predictions) Cancel(ctx context.Context, action match.LifecycleAction) error {
	return p.finishPrediction(ctx, action.ChannelID, action.MatchID, false, true)
}

func (p *Predictions) createPrediction(ctx context.Context, channelID uuid.UUID, matchID int64) error {
	if matchID <= 0 {
		return nil
	}

	if channelID == uuid.Nil {
		p.logger.WarnContext(ctx, "dota prediction skipped: invalid channel ID")
		return nil
	}
	key := predictionKey(channelID, matchID)
	if recovered, err := p.recoverPendingPrediction(ctx, key, channelID); recovered {
		return err
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

	template := settings.PredictionSettings.TitleTemplate
	window := settings.PredictionSettings.WindowSeconds
	if strings.TrimSpace(template) == "" ||
		utf8.RuneCountInString(template)+predictionCorrelationMarkerRunes > maxPredictionTitleRunes ||
		window < minPredictionWindow || window > maxPredictionWindow {
		p.logger.WarnContext(
			ctx,
			"dota prediction skipped: invalid settings",
			slog.String("channel_id", channelID.String()),
			slog.Int("window_seconds", window),
		)
		return nil
	}
	correlation, err := newPredictionCorrelation()
	if err != nil {
		return fmt.Errorf("generate prediction correlation: %w", err)
	}
	title := template + predictionCorrelationMarker(correlation)

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

	intent := pendingPredictionIntent{
		Version:         pendingIntentVersion,
		Token:           uuid.NewString(),
		Title:           title,
		Correlation:     correlation,
		YesOutcomeTitle: "Yes",
		NoOutcomeTitle:  "No",
		ReservedAt:      time.Now().UTC(),
	}
	reserved, err := p.store.Reserve(ctx, key, intent, predictionTTL)
	if err != nil {
		return fmt.Errorf("reserve prediction: %w", err)
	}
	if !reserved {
		return nil
	}

	client, err := p.clients.New(ctx, *channel.TwitchUserID)
	if err != nil {
		return p.creationFailed(ctx, key, intent.Token, fmt.Errorf("create Twitch client: %w", err))
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
			if releaseErr := p.releaseReservation(ctx, key, intent.Token); releaseErr != nil {
				return releaseErr
			}
			return nil
		}
		return p.creationFailed(ctx, key, intent.Token, fmt.Errorf("create Twitch prediction: %w", err))
	}
	if err := predictionResponseError(response); err != nil {
		if isAlreadyActive(err) {
			if releaseErr := p.releaseReservation(ctx, key, intent.Token); releaseErr != nil {
				return releaseErr
			}
			return nil
		}
		return p.creationFailed(ctx, key, intent.Token, fmt.Errorf("create Twitch prediction: %w", err))
	}

	record, err := recordFromCreateResponse(response)
	if err != nil {
		return p.creationFailed(ctx, key, intent.Token, err)
	}
	if err := p.store.Commit(ctx, key, intent.Token, record, predictionTTL); err != nil {
		return fmt.Errorf("store created prediction: %w", err)
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

func (p *Predictions) recoverPendingPrediction(
	ctx context.Context,
	key string,
	channelID uuid.UUID,
) (bool, error) {
	intent, err := p.store.GetPending(ctx, key)
	if errors.Is(err, errPredictionIntentNotFound) {
		return false, nil
	}
	if err != nil {
		return true, fmt.Errorf("get pending prediction intent: %w", err)
	}
	if err := intent.validate(); err != nil {
		return true, fmt.Errorf("%w: invalid pending intent: %v", errPredictionRecoveryUnsafe, err)
	}

	channel, err := p.channels.GetByID(ctx, channelID)
	if err != nil {
		return true, fmt.Errorf("get channel for pending prediction: %w", err)
	}
	if !channel.TwitchConnected() || channel.TwitchUserID == nil || channel.TwitchPlatformID == nil ||
		strings.TrimSpace(*channel.TwitchPlatformID) == "" {
		return true, errors.New("pending prediction channel is not connected to Twitch")
	}

	client, err := p.clients.New(ctx, *channel.TwitchUserID)
	if err != nil {
		return true, fmt.Errorf("create Twitch client for pending prediction: %w", err)
	}
	broadcasterID := strings.TrimSpace(*channel.TwitchPlatformID)
	response, err := client.GetPredictions(&helix.PredictionsParams{
		BroadcasterID: broadcasterID,
		First:         "100",
	})
	if err != nil {
		return true, fmt.Errorf("get Twitch predictions for pending prediction: %w", err)
	}
	if err := predictionResponseError(response); err != nil {
		return true, fmt.Errorf("get Twitch predictions for pending prediction: %w", err)
	}

	prediction, found := matchingPendingPrediction(response, intent)
	if !found {
		return true, fmt.Errorf(
			"%w: expected exactly one matching active or locked Twitch prediction",
			errPredictionRecoveryUnsafe,
		)
	}
	record, err := recordFromPrediction(prediction, intent.YesOutcomeTitle, intent.NoOutcomeTitle)
	if err != nil {
		return true, fmt.Errorf("%w: build recovered prediction record: %v", errPredictionRecoveryUnsafe, err)
	}
	if err := p.store.Commit(ctx, key, intent.Token, record, predictionTTL); err != nil {
		return true, fmt.Errorf("commit recovered prediction: %w", err)
	}

	return true, nil
}

func (p *Predictions) finishPrediction(
	ctx context.Context,
	channelID uuid.UUID,
	matchID int64,
	win bool,
	cancel bool,
) error {
	if matchID <= 0 {
		return nil
	}

	if channelID == uuid.Nil {
		p.logger.WarnContext(ctx, "dota prediction skipped: invalid channel ID")
		return nil
	}

	key := predictionKey(channelID, matchID)
	token := uuid.NewString()
	claimed, err := p.store.ClaimTerminal(ctx, key, token, p.terminalLease.ttl)
	if err != nil {
		return fmt.Errorf("claim prediction terminal: %w", err)
	}
	if !claimed {
		return ErrTerminalInProgress
	}
	heartbeat := p.startTerminalHeartbeat(ctx, key, token)
	operationCtx := heartbeat.operationCtx

	if recovered, err := p.recoverPendingPrediction(operationCtx, key, channelID); recovered {
		if err != nil {
			return p.terminalFailed(ctx, key, token, heartbeat, err)
		}
	}
	record, err := p.store.Get(operationCtx, key)
	if errors.Is(err, errPredictionNotFound) {
		return p.completeTerminal(ctx, key, token, heartbeat, "complete missing prediction")
	}
	if err != nil {
		return p.terminalFailed(ctx, key, token, heartbeat, fmt.Errorf("get stored prediction: %w", err))
	}

	channel, err := p.channels.GetByID(operationCtx, channelID)
	if err != nil {
		return p.terminalFailed(ctx, key, token, heartbeat, fmt.Errorf("get channel for stored prediction: %w", err))
	}
	if !channel.TwitchConnected() || channel.TwitchUserID == nil || channel.TwitchPlatformID == nil ||
		strings.TrimSpace(*channel.TwitchPlatformID) == "" {
		return p.terminalFailed(ctx, key, token, heartbeat, errors.New("stored prediction channel is not connected to Twitch"))
	}

	client, err := p.clients.New(operationCtx, *channel.TwitchUserID)
	if err != nil {
		return p.terminalFailed(ctx, key, token, heartbeat, fmt.Errorf("create Twitch client for stored prediction: %w", err))
	}
	broadcasterID := strings.TrimSpace(*channel.TwitchPlatformID)
	response, err := client.GetPredictions(&helix.PredictionsParams{
		BroadcasterID: broadcasterID,
		ID:            record.PredictionID,
	})
	if err != nil {
		return p.terminalFailed(ctx, key, token, heartbeat, fmt.Errorf("get Twitch prediction: %w", err))
	}
	if err := predictionResponseError(response); err != nil {
		return p.terminalFailed(ctx, key, token, heartbeat, fmt.Errorf("get Twitch prediction: %w", err))
	}

	prediction, active := activePrediction(response, record.PredictionID)
	if !active {
		return p.completeTerminal(ctx, key, token, heartbeat, "complete inactive prediction")
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
		return p.terminalFailed(ctx, key, token, heartbeat, fmt.Errorf("end Twitch prediction: %w", err))
	}
	if err := predictionResponseError(endResponse); err != nil {
		return p.terminalFailed(ctx, key, token, heartbeat, fmt.Errorf("end Twitch prediction: %w", err))
	}

	return p.completeTerminal(ctx, key, token, heartbeat, "complete finished prediction")
}

func (p *Predictions) startTerminalHeartbeat(ctx context.Context, key string, token string) *terminalHeartbeat {
	operationCtx, operationCancel := context.WithTimeout(ctx, terminalOperationTimeout)
	heartbeat := &terminalHeartbeat{
		operationCtx:    operationCtx,
		operationCancel: operationCancel,
		stopCh:          make(chan struct{}),
		done:            make(chan struct{}),
	}

	go func() {
		defer operationCancel()
		defer close(heartbeat.done)

		ticker := time.NewTicker(p.terminalLease.renewInterval)
		defer ticker.Stop()
		for {
			select {
			case <-heartbeat.stopCh:
				if err := operationCtx.Err(); err != nil {
					heartbeat.setError(err)
				}
				return
			case <-operationCtx.Done():
				heartbeat.setError(operationCtx.Err())
				return
			case <-ticker.C:
				renewCtx, renewCancel := context.WithTimeout(operationCtx, terminalRenewTimeout)
				renewed, err := p.store.RenewTerminal(
					renewCtx,
					key,
					token,
					p.terminalLease.ttl,
				)
				renewCancel()
				select {
				case <-heartbeat.stopCh:
					if err := operationCtx.Err(); err != nil {
						heartbeat.setError(err)
					}
					return
				default:
				}
				if err != nil {
					heartbeat.fail(fmt.Errorf("renew prediction terminal claim: %w", err))
					return
				}
				if !renewed {
					heartbeat.fail(ErrTerminalOwnershipLost)
					return
				}
			}
		}
	}()

	return heartbeat
}

func (p *Predictions) terminalFailed(
	ctx context.Context,
	key string,
	token string,
	heartbeat *terminalHeartbeat,
	cause error,
) error {
	if err := heartbeat.stop(); err != nil {
		cause = fmt.Errorf("%w; stop prediction terminal heartbeat: %v", cause, err)
	}
	return p.releaseTerminalFailure(ctx, key, token, cause)
}

func (p *Predictions) releaseTerminalFailure(ctx context.Context, key string, token string, cause error) error {
	if err := p.releaseTerminal(ctx, key, token); err != nil {
		return fmt.Errorf("%w; release prediction terminal claim: %v", cause, err)
	}
	return cause
}

func (p *Predictions) releaseTerminal(ctx context.Context, key string, token string) error {
	cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
	defer cancel()
	if err := p.store.ReleaseTerminal(cleanupCtx, key, token); err != nil {
		p.logger.ErrorContext(cleanupCtx, "dota prediction terminal claim cleanup failed", logger.Error(err))
		return err
	}
	return nil
}

func (p *Predictions) completeTerminal(
	ctx context.Context,
	key string,
	token string,
	heartbeat *terminalHeartbeat,
	operation string,
) error {
	if err := heartbeat.stop(); err != nil {
		return p.releaseTerminalFailure(
			ctx,
			key,
			token,
			fmt.Errorf("%s: stop prediction terminal heartbeat: %w", operation, err),
		)
	}

	cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
	defer cancel()
	completed, err := p.store.CompleteTerminal(cleanupCtx, key, token)
	if err != nil {
		return p.releaseTerminalFailure(ctx, key, token, fmt.Errorf("%s: %w", operation, err))
	}
	if !completed {
		return p.releaseTerminalFailure(ctx, key, token, fmt.Errorf("%s: %w", operation, ErrTerminalOwnershipLost))
	}
	return nil
}

func predictionKey(channelID uuid.UUID, matchID int64) string {
	return predictionKeyPrefix + channelID.String() + ":" + strconv.FormatInt(matchID, 10)
}

func matchEndedDeliveryKey(channelID uuid.UUID, matchID int64) string {
	return matchEndedDeliveryKeyPrefix + channelID.String() + ":" + strconv.FormatInt(matchID, 10)
}

func newPredictionCorrelation() (string, error) {
	bytes := make([]byte, predictionCorrelationBytes)
	if _, err := cryptorand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func predictionCorrelationMarker(correlation string) string {
	return predictionCorrelationMarkerPrefix + correlation + predictionCorrelationMarkerSuffix
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

	return recordFromPrediction(response.Data.Predictions[0], "Yes", "No")
}

func recordFromPrediction(
	prediction helix.Prediction,
	yesOutcomeTitle string,
	noOutcomeTitle string,
) (storedPrediction, error) {
	if strings.TrimSpace(prediction.ID) == "" {
		return storedPrediction{}, errors.New("create Twitch prediction returned an empty prediction ID")
	}
	if len(prediction.Outcomes) != 2 {
		return storedPrediction{}, errors.New("Twitch prediction returned unexpected outcomes")
	}

	record := storedPrediction{PredictionID: prediction.ID}
	for _, outcome := range prediction.Outcomes {
		switch outcome.Title {
		case yesOutcomeTitle:
			record.YesOutcomeID = outcome.ID
		case noOutcomeTitle:
			record.NoOutcomeID = outcome.ID
		default:
			return storedPrediction{}, errors.New("Twitch prediction returned unexpected outcomes")
		}
	}
	if strings.TrimSpace(record.YesOutcomeID) == "" || strings.TrimSpace(record.NoOutcomeID) == "" {
		return storedPrediction{}, errors.New("create Twitch prediction returned incomplete outcomes")
	}

	return record, nil
}

func (intent pendingPredictionIntent) validate() error {
	if intent.Version != pendingIntentVersion {
		return errors.New("unsupported pending intent version")
	}
	if strings.TrimSpace(intent.Token) == "" {
		return errors.New("missing reservation token")
	}
	if !validPredictionCorrelation(intent.Correlation) {
		return errors.New("invalid prediction correlation")
	}
	if strings.TrimSpace(intent.Title) == "" || utf8.RuneCountInString(intent.Title) < 1 ||
		utf8.RuneCountInString(intent.Title) > maxPredictionTitleRunes {
		return errors.New("invalid prediction title")
	}
	if !strings.HasSuffix(intent.Title, predictionCorrelationMarker(intent.Correlation)) {
		return errors.New("prediction title is missing correlation marker")
	}
	if strings.TrimSpace(intent.YesOutcomeTitle) == "" || strings.TrimSpace(intent.NoOutcomeTitle) == "" ||
		intent.YesOutcomeTitle == intent.NoOutcomeTitle {
		return errors.New("invalid prediction outcomes")
	}
	if intent.ReservedAt.IsZero() {
		return errors.New("missing reservation timestamp")
	}
	return nil
}

func validPredictionCorrelation(correlation string) bool {
	decoded, err := base64.RawURLEncoding.DecodeString(correlation)
	return err == nil && len(decoded) == predictionCorrelationBytes &&
		base64.RawURLEncoding.EncodeToString(decoded) == correlation
}

func matchingPendingPrediction(
	response *helix.PredictionsResponse,
	intent pendingPredictionIntent,
) (helix.Prediction, bool) {
	var match helix.Prediction
	found := false
	for _, prediction := range response.Data.Predictions {
		if !matchesPendingIntent(prediction, intent) {
			continue
		}
		if found {
			return helix.Prediction{}, false
		}
		match = prediction
		found = true
	}
	return match, found
}

func matchesPendingIntent(prediction helix.Prediction, intent pendingPredictionIntent) bool {
	if prediction.Status != "ACTIVE" && prediction.Status != "LOCKED" {
		return false
	}
	if prediction.Title != intent.Title ||
		!strings.HasSuffix(prediction.Title, predictionCorrelationMarker(intent.Correlation)) ||
		!matchesReservationTime(prediction.CreatedAt.Time, intent.ReservedAt) {
		return false
	}
	_, err := recordFromPrediction(prediction, intent.YesOutcomeTitle, intent.NoOutcomeTitle)
	return err == nil
}

func matchesReservationTime(createdAt time.Time, reservedAt time.Time) bool {
	if createdAt.IsZero() || reservedAt.IsZero() {
		return false
	}
	return !createdAt.Before(reservedAt.Add(-pendingPredictionCorrelationWindow)) &&
		!createdAt.After(reservedAt.Add(pendingPredictionCorrelationWindow))
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
