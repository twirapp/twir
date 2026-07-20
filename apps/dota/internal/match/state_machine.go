package match

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/kv"
	kvoptions "github.com/twirapp/kv/options"
	"github.com/twirapp/twir/apps/dota/internal/gsi"
	busapi "github.com/twirapp/twir/libs/bus-core/api"
	busdota "github.com/twirapp/twir/libs/bus-core/dota"
	"github.com/twirapp/twir/libs/logger"
	dotarepository "github.com/twirapp/twir/libs/repositories/dota"
)

type State string

const (
	StateIdle          State = "idle"
	StateHeroSelection State = "hero_selection"
	StateStrategy      State = "strategy"
	StatePreGame       State = "pre_game"
	StateInGame        State = "in_game"
	StatePostGame      State = "post_game"
)

const (
	snapshotKeyPrefix = "cache:twir:dota:matchstate:"
	snapshotTTL       = 6 * time.Hour

	heroNamePrefix = "npc_dota_hero_"

	eventRoshanKilled = "roshan_killed"
	eventAegisPicked  = "aegis_picked_up"
)

type EventEmitter interface {
	MatchStarted(ctx context.Context, msg busdota.MatchStartedMessage) error
	MatchEnded(ctx context.Context, msg busdota.MatchEndedMessage) error
	RoshanKilled(ctx context.Context, msg busdota.RoshanKilledMessage) error
	AegisPickup(ctx context.Context, msg busdota.AegisPickupMessage) error
	StateUpdate(ctx context.Context, msg busapi.DotaStateUpdateMessage) error
}

type Snapshot struct {
	ChannelID      uuid.UUID `json:"channelId"`
	State          State     `json:"state"`
	InGame         bool      `json:"inGame"`
	MatchID        int64     `json:"matchId"`
	HeroName       string    `json:"heroName"`
	IsRadiant      bool      `json:"isRadiant"`
	SteamAccountID string    `json:"steamAccountId"`
	RadiantScore   int       `json:"radiantScore"`
	DireScore      int       `json:"direScore"`
	GameTime       int       `json:"gameTime"`
	Mmr            int       `json:"mmr"`
	SessionWins    int       `json:"sessionWins"`
	SessionLosses  int       `json:"sessionLosses"`
	WinProbability float64   `json:"winProbability"`

	SeenEvents []string `json:"seenEvents,omitempty"`
}

type channelState struct {
	mu          sync.Mutex
	loaded      bool
	snap        Snapshot
	seenEvents  map[string]struct{}
	mmrDelta    int
	settingsSet bool
}

type StateMachine struct {
	repo    dotarepository.Repository
	emitter EventEmitter
	kv      kv.KV
	logger  *slog.Logger

	mu       sync.Mutex
	channels map[uuid.UUID]*channelState
}

var _ interface {
	Process(ctx context.Context, channelID uuid.UUID, payload gsi.Payload) error
} = (*StateMachine)(nil)

func New(
	repo dotarepository.Repository,
	emitter EventEmitter,
	kvStore kv.KV,
	logger *slog.Logger,
) *StateMachine {
	return &StateMachine{
		repo:     repo,
		emitter:  emitter,
		kv:       kvStore,
		logger:   logger,
		channels: make(map[uuid.UUID]*channelState),
	}
}

func mapGameState(gs gsi.GameState) (State, bool) {
	switch gs {
	case gsi.GameStateHeroSelection:
		return StateHeroSelection, true
	case gsi.GameStateStrategyTime:
		return StateStrategy, true
	case gsi.GameStatePreGame:
		return StatePreGame, true
	case gsi.GameStateInProgress:
		return StateInGame, true
	case gsi.GameStatePostGame:
		return StatePostGame, true
	default:
		return "", false
	}
}

func stripHeroPrefix(name string) string {
	return strings.TrimPrefix(name, heroNamePrefix)
}

func snapshotKey(channelID uuid.UUID) string {
	return snapshotKeyPrefix + channelID.String()
}

func (m *StateMachine) channel(ctx context.Context, channelID uuid.UUID) *channelState {
	m.mu.Lock()
	cs, ok := m.channels[channelID]
	if !ok {
		cs = &channelState{
			seenEvents: make(map[string]struct{}),
			snap: Snapshot{
				ChannelID: channelID,
				State:     StateIdle,
			},
		}
		m.channels[channelID] = cs
	}
	m.mu.Unlock()

	cs.mu.Lock()
	defer cs.mu.Unlock()
	if !cs.loaded {
		cs.loaded = true
		m.restore(ctx, cs)
	}

	return cs
}

func (m *StateMachine) restore(ctx context.Context, cs *channelState) {
	val := m.kv.Get(ctx, snapshotKey(cs.snap.ChannelID))
	data, err := val.Bytes()
	if err != nil {
		if !errors.Is(err, kv.ErrKeyNil) {
			m.logger.WarnContext(
				ctx,
				"dota match: failed to read snapshot from kv",
				logger.Error(err),
			)
		}
		return
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		m.logger.WarnContext(
			ctx,
			"dota match: failed to unmarshal snapshot",
			logger.Error(err),
		)
		return
	}

	cs.snap = snap
	cs.seenEvents = make(map[string]struct{}, len(snap.SeenEvents))
	for _, e := range snap.SeenEvents {
		cs.seenEvents[e] = struct{}{}
	}
}

func (m *StateMachine) persist(ctx context.Context, cs *channelState) {
	snap := cs.snap
	snap.SeenEvents = make([]string, 0, len(cs.seenEvents))
	for e := range cs.seenEvents {
		snap.SeenEvents = append(snap.SeenEvents, e)
	}

	data, err := json.Marshal(snap)
	if err != nil {
		m.logger.WarnContext(ctx, "dota match: failed to marshal snapshot", logger.Error(err))
		return
	}

	if err := m.kv.Set(
		ctx,
		snapshotKey(cs.snap.ChannelID),
		data,
		kvoptions.WithExpire(snapshotTTL),
	); err != nil {
		m.logger.WarnContext(ctx, "dota match: failed to persist snapshot", logger.Error(err))
	}
}

func (m *StateMachine) Process(ctx context.Context, channelID uuid.UUID, payload gsi.Payload) error {
	cs := m.channel(ctx, channelID)

	cs.mu.Lock()
	defer cs.mu.Unlock()

	if payload.Map == nil || payload.Player == nil ||
		payload.Player.Activity != gsi.PlayerActivityPlaying {
		return m.goIdle(ctx, cs)
	}

	prevState := cs.snap.State
	scoreChanged := cs.snap.RadiantScore != payload.Map.RadiantScore ||
		cs.snap.DireScore != payload.Map.DireScore

	cs.snap.RadiantScore = payload.Map.RadiantScore
	cs.snap.DireScore = payload.Map.DireScore
	cs.snap.GameTime = payload.Map.GameTime

	newState, known := mapGameState(payload.Map.GameState)
	if !known {
		newState = prevState
	}

	if err := m.ensureSettings(ctx, cs); err != nil {
		m.logger.WarnContext(ctx, "dota match: failed to load settings", logger.Error(err))
	}

	if newState == StateInGame && payload.Map.MatchID != cs.snap.MatchID {
		m.startMatch(cs, payload)
		if err := m.emitter.MatchStarted(ctx, busdota.MatchStartedMessage{
			ChannelID:      cs.snap.ChannelID.String(),
			SteamAccountID: cs.snap.SteamAccountID,
			HeroName:       cs.snap.HeroName,
		}); err != nil {
			m.logger.ErrorContext(ctx, "dota match: failed to emit match started", logger.Error(err))
		}
	}

	if newState == StatePostGame {
		return m.finishMatch(ctx, cs, payload)
	}

	cs.snap.State = newState
	cs.snap.InGame = newState != StateIdle

	if err := m.processEvents(ctx, cs, payload.Events); err != nil {
		m.logger.WarnContext(ctx, "dota match: failed to emit events", logger.Error(err))
	}

	if newState != prevState || scoreChanged {
		m.persist(ctx, cs)
		m.emitStateUpdate(ctx, cs)
	}

	return nil
}

func (m *StateMachine) startMatch(cs *channelState, payload gsi.Payload) {
	heroName := ""
	if payload.Hero != nil {
		heroName = stripHeroPrefix(payload.Hero.Name)
	}

	cs.snap.MatchID = payload.Map.MatchID
	cs.snap.HeroName = heroName
	cs.snap.IsRadiant = payload.Player.TeamName == "radiant"
	cs.snap.SteamAccountID = strconv.FormatInt(payload.Player.AccountID, 10)
	cs.snap.State = StateInGame
	cs.snap.InGame = true
	cs.seenEvents = make(map[string]struct{})
}

func (m *StateMachine) finishMatch(
	ctx context.Context,
	cs *channelState,
	payload gsi.Payload,
) error {
	if cs.snap.MatchID == 0 ||
		(payload.Map.WinTeam != gsi.WinTeamRadiant && payload.Map.WinTeam != gsi.WinTeamDire) {
		cs.snap.State = StatePostGame
		return nil
	}

	won := (payload.Map.WinTeam == gsi.WinTeamRadiant) == cs.snap.IsRadiant

	delta := cs.mmrDelta
	settings, err := m.repo.GetByChannelID(ctx, cs.snap.ChannelID)
	if err != nil {
		m.logger.WarnContext(ctx, "dota match: failed to refresh settings", logger.Error(err))
	} else {
		delta = settings.MmrDelta
	}
	if !won {
		delta = -delta
	}

	updated, err := m.repo.UpdateMatchResult(ctx, cs.snap.ChannelID, won, delta)
	if err != nil {
		return fmt.Errorf("update match result: %w", err)
	}

	cs.snap.Mmr = updated.Mmr
	cs.snap.SessionWins = updated.SessionWins
	cs.snap.SessionLosses = updated.SessionLosses
	cs.mmrDelta = updated.MmrDelta
	cs.settingsSet = true

	if err := m.emitter.MatchEnded(ctx, busdota.MatchEndedMessage{
		ChannelID:      cs.snap.ChannelID.String(),
		SteamAccountID: cs.snap.SteamAccountID,
		Win:            won,
		HeroName:       cs.snap.HeroName,
		Mmr:            updated.Mmr,
		SessionWins:    updated.SessionWins,
		SessionLosses:  updated.SessionLosses,
	}); err != nil {
		m.logger.ErrorContext(ctx, "dota match: failed to emit match ended", logger.Error(err))
	}

	cs.snap.State = StateIdle
	cs.snap.InGame = false
	cs.snap.MatchID = 0
	cs.snap.HeroName = ""
	cs.snap.SteamAccountID = ""
	cs.snap.RadiantScore = 0
	cs.snap.DireScore = 0
	cs.snap.GameTime = 0
	cs.seenEvents = make(map[string]struct{})

	m.persist(ctx, cs)
	m.emitStateUpdate(ctx, cs)

	return nil
}

func (m *StateMachine) goIdle(ctx context.Context, cs *channelState) error {
	if cs.snap.State == StateIdle {
		return nil
	}

	cs.snap.State = StateIdle
	cs.snap.InGame = false
	cs.snap.MatchID = 0
	cs.snap.HeroName = ""
	cs.snap.SteamAccountID = ""
	cs.snap.RadiantScore = 0
	cs.snap.DireScore = 0
	cs.snap.GameTime = 0
	cs.seenEvents = make(map[string]struct{})

	m.persist(ctx, cs)
	m.emitStateUpdate(ctx, cs)

	return nil
}

func (m *StateMachine) processEvents(
	ctx context.Context,
	cs *channelState,
	events []gsi.Event,
) error {
	for _, event := range events {
		key := fmt.Sprintf("%s:%d", event.EventType, event.GameTime)
		if _, seen := cs.seenEvents[key]; seen {
			continue
		}
		cs.seenEvents[key] = struct{}{}

		switch event.EventType {
		case eventRoshanKilled:
			if err := m.emitter.RoshanKilled(ctx, busdota.RoshanKilledMessage{
				ChannelID: cs.snap.ChannelID.String(),
				Team:      event.KillerTeam,
				GameTime:  event.GameTime,
			}); err != nil {
				return fmt.Errorf("emit roshan killed: %w", err)
			}
		case eventAegisPicked:
			if err := m.emitter.AegisPickup(ctx, busdota.AegisPickupMessage{
				ChannelID: cs.snap.ChannelID.String(),
				GameTime:  event.GameTime,
			}); err != nil {
				return fmt.Errorf("emit aegis pickup: %w", err)
			}
		}
	}

	return nil
}

func (m *StateMachine) ensureSettings(ctx context.Context, cs *channelState) error {
	if cs.settingsSet {
		return nil
	}

	settings, err := m.repo.GetByChannelID(ctx, cs.snap.ChannelID)
	if err != nil {
		return fmt.Errorf("get settings: %w", err)
	}

	cs.snap.Mmr = settings.Mmr
	cs.snap.SessionWins = settings.SessionWins
	cs.snap.SessionLosses = settings.SessionLosses
	cs.mmrDelta = settings.MmrDelta
	cs.settingsSet = true

	return nil
}

func (m *StateMachine) emitStateUpdate(ctx context.Context, cs *channelState) {
	if err := m.emitter.StateUpdate(ctx, busapi.DotaStateUpdateMessage{
		ChannelID:      cs.snap.ChannelID.String(),
		InGame:         cs.snap.InGame,
		Mmr:            cs.snap.Mmr,
		SessionWins:    cs.snap.SessionWins,
		SessionLosses:  cs.snap.SessionLosses,
		WinProbability: 0,
		HeroName:       cs.snap.HeroName,
		MatchID:        cs.snap.MatchID,
	}); err != nil {
		m.logger.ErrorContext(ctx, "dota match: failed to emit state update", logger.Error(err))
	}
}

func (m *StateMachine) GetSnapshot(ctx context.Context, channelID uuid.UUID) (Snapshot, error) {
	cs := m.channel(ctx, channelID)

	cs.mu.Lock()
	defer cs.mu.Unlock()

	return cs.snap, nil
}
