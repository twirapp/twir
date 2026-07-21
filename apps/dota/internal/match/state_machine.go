package match

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/dota/internal/gsi"
	busapi "github.com/twirapp/twir/libs/bus-core/api"
	busdota "github.com/twirapp/twir/libs/bus-core/dota"
	"github.com/twirapp/twir/libs/logger"
	dotarepository "github.com/twirapp/twir/libs/repositories/dota"
	"github.com/twirapp/twir/libs/repositories/dota/model"
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
	maxTransitionRetries = 3

	lifecycleActionsPerRevision = 2

	winProbabilityBasisPoints                = 10_000
	winProbabilityUpdateThresholdBasisPoints = 500

	heroNamePrefix = "npc_dota_hero_"

	eventRoshanKilled = "roshan_killed"
	eventAegisPicked  = "aegis_picked_up"

	eventClaimLease = 30 * time.Second
)

type EventEmitter interface {
	MatchStarted(ctx context.Context, msg busdota.MatchStartedMessage) error
	MatchEnded(ctx context.Context, msg busdota.MatchEndedMessage) error
	MatchAbandoned(ctx context.Context, msg busdota.MatchAbandonedMessage) error
	RoshanKilled(ctx context.Context, msg busdota.RoshanKilledMessage) error
	AegisPickup(ctx context.Context, msg busdota.AegisPickupMessage) error
	StateUpdate(ctx context.Context, msg busapi.DotaStateUpdateMessage) error
}

type Snapshot struct {
	ChannelID             uuid.UUID `json:"channelId"`
	Revision              uint64    `json:"revision"`
	State                 State     `json:"state"`
	InGame                bool      `json:"inGame"`
	MatchID               int64     `json:"matchId"`
	HeroName              string    `json:"heroName"`
	IsRadiant             bool      `json:"isRadiant"`
	TeamKnown             bool      `json:"teamKnown"`
	SteamAccountID        string    `json:"steamAccountId"`
	RadiantScore          int       `json:"radiantScore"`
	DireScore             int       `json:"direScore"`
	GameTime              int       `json:"gameTime"`
	LastProviderTimestamp int64     `json:"lastProviderTimestamp"`
	LastGameTime          int       `json:"lastGameTime"`
	Mmr                   int       `json:"mmr"`
	SessionWins           int       `json:"sessionWins"`
	SessionLosses         int       `json:"sessionLosses"`
	WinProbability        float64   `json:"winProbability"`

	SeenEvents  []string              `json:"seenEvents,omitempty"`
	EventClaims map[string]EventClaim `json:"eventClaims,omitempty"`
}

type EventClaim struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}

type StateMachine struct {
	repo    dotarepository.Repository
	emitter EventEmitter
	logger  *slog.Logger
}

var _ interface {
	Process(ctx context.Context, channelID uuid.UUID, payload gsi.Payload) error
} = (*StateMachine)(nil)

func New(
	repo dotarepository.Repository,
	emitter EventEmitter,
	logger *slog.Logger,
) *StateMachine {
	return &StateMachine{
		repo:    repo,
		emitter: emitter,
		logger:  logger,
	}
}

func mapGameState(gameState gsi.GameState) (State, bool) {
	switch gameState {
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

func (m *StateMachine) Process(ctx context.Context, channelID uuid.UUID, payload gsi.Payload) error {
	for attempt := 0; attempt < maxTransitionRetries; attempt++ {
		state, err := m.repo.GetMatchState(ctx, channelID)
		if err != nil {
			return fmt.Errorf("get dota match state: %w", err)
		}

		current, err := snapshotFromMatchState(state, channelID)
		if err != nil {
			return fmt.Errorf("decode dota match state: %w", err)
		}
		if isActivePayload(payload) && payload.Map.GameState != gsi.GameStatePostGame {
			m.loadSettings(ctx, &current)
		}

		next, actions, events, emitStateUpdate, changed := transition(current, payload)
		if !changed {
			m.retryPendingEvents(ctx, current, payload)
			return nil
		}

		next.Revision, err = nextRevision(state.Revision)
		if err != nil {
			return err
		}

		input, err := transitionInput(state, next, actions)
		if err != nil {
			return fmt.Errorf("encode dota match transition: %w", err)
		}
		committed, err := m.repo.ApplyMatchStateTransition(ctx, input)
		if err != nil {
			return fmt.Errorf("apply dota match transition: %w", err)
		}
		if !committed {
			continue
		}

		m.publishLifecycleNotifications(ctx, next, actions)
		m.publishEvents(ctx, next, events)
		if emitStateUpdate {
			m.emitStateUpdate(ctx, next)
		}

		return nil
	}

	return errors.New("dota match transition retry limit reached")
}

func (m *StateMachine) UpdateWinProbability(
	ctx context.Context,
	channelID uuid.UUID,
	expectedMatchID int64,
	probability float64,
) error {
	if math.IsNaN(probability) || math.IsInf(probability, 0) || probability < 0 || probability > 1 {
		return fmt.Errorf("invalid win probability %v: expected a value between 0 and 1", probability)
	}

	for attempt := 0; attempt < maxTransitionRetries; attempt++ {
		state, err := m.repo.GetMatchState(ctx, channelID)
		if err != nil {
			return fmt.Errorf("get dota match state: %w", err)
		}

		current, err := snapshotFromMatchState(state, channelID)
		if err != nil {
			return fmt.Errorf("decode dota match state: %w", err)
		}
		if !current.InGame || current.MatchID == 0 || current.MatchID != expectedMatchID {
			return nil
		}

		// Probabilities are compared at 0.01% precision; require a 5.00% (500 basis point) change.
		currentBasisPoints := int64(math.Round(current.WinProbability * winProbabilityBasisPoints))
		nextBasisPoints := int64(math.Round(probability * winProbabilityBasisPoints))
		if math.Abs(float64(currentBasisPoints-nextBasisPoints)) < winProbabilityUpdateThresholdBasisPoints {
			return nil
		}

		next := current
		next.WinProbability = probability
		next.Revision, err = nextRevision(state.Revision)
		if err != nil {
			return err
		}

		input, err := transitionInput(state, next, nil)
		if err != nil {
			return fmt.Errorf("encode dota win probability transition: %w", err)
		}
		committed, err := m.repo.ApplyMatchStateTransition(ctx, input)
		if err != nil {
			return fmt.Errorf("apply dota win probability transition: %w", err)
		}
		if !committed {
			continue
		}

		m.emitStateUpdate(ctx, next)
		return nil
	}

	return errors.New("dota win probability transition retry limit reached")
}

// UpdateStats copies settled settings into the latest state and publishes the resulting state update.
func (m *StateMachine) UpdateStats(
	ctx context.Context,
	channelID uuid.UUID,
	settings model.ChannelDotaSettings,
) error {
	for attempt := 0; attempt < maxTransitionRetries; attempt++ {
		state, err := m.repo.GetMatchState(ctx, channelID)
		if err != nil {
			return fmt.Errorf("get dota match state: %w", err)
		}

		current, err := snapshotFromMatchState(state, channelID)
		if err != nil {
			return fmt.Errorf("decode dota match state: %w", err)
		}
		if current.Mmr == settings.Mmr &&
			current.SessionWins == settings.SessionWins &&
			current.SessionLosses == settings.SessionLosses {
			if err := m.publishStateUpdate(ctx, current); err != nil {
				return fmt.Errorf("publish dota stats state update: %w", err)
			}
			return nil
		}

		next := current
		next.Mmr = settings.Mmr
		next.SessionWins = settings.SessionWins
		next.SessionLosses = settings.SessionLosses
		next.Revision, err = nextRevision(state.Revision)
		if err != nil {
			return err
		}

		input, err := transitionInput(state, next, nil)
		if err != nil {
			return fmt.Errorf("encode dota stats transition: %w", err)
		}
		committed, err := m.repo.ApplyMatchStateTransition(ctx, input)
		if err != nil {
			return fmt.Errorf("apply dota stats transition: %w", err)
		}
		if !committed {
			continue
		}

		if err := m.publishStateUpdate(ctx, next); err != nil {
			return fmt.Errorf("publish dota stats state update: %w", err)
		}
		return nil
	}

	return errors.New("dota stats transition retry limit reached")
}

func (m *StateMachine) GetSnapshot(ctx context.Context, channelID uuid.UUID) (Snapshot, error) {
	state, err := m.repo.GetMatchState(ctx, channelID)
	if err != nil {
		return Snapshot{}, fmt.Errorf("get dota match state: %w", err)
	}

	snapshot, err := snapshotFromMatchState(state, channelID)
	if err != nil {
		return Snapshot{}, fmt.Errorf("decode dota match state: %w", err)
	}

	return snapshot, nil
}

func snapshotFromMatchState(state model.MatchState, channelID uuid.UUID) (Snapshot, error) {
	if state.ChannelID == uuid.Nil {
		state.ChannelID = channelID
	}
	if state.ChannelID != channelID {
		return Snapshot{}, fmt.Errorf(
			"stored snapshot channel ID %s does not match requested channel %s",
			state.ChannelID,
			channelID,
		)
	}
	if state.Revision < 0 {
		return Snapshot{}, fmt.Errorf("stored revision %d must not be negative", state.Revision)
	}
	if state.ProviderTimestamp < 0 {
		return Snapshot{}, fmt.Errorf("stored provider timestamp %d must not be negative", state.ProviderTimestamp)
	}

	snapshot := Snapshot{
		ChannelID: channelID,
		State:     StateIdle,
	}
	if len(state.Snapshot) != 0 {
		if err := json.Unmarshal(state.Snapshot, &snapshot); err != nil {
			return Snapshot{}, err
		}
	}
	if snapshot.ChannelID == uuid.Nil {
		snapshot.ChannelID = channelID
	} else if snapshot.ChannelID != channelID {
		return Snapshot{}, fmt.Errorf(
			"snapshot channel ID %s does not match requested channel %s",
			snapshot.ChannelID,
			channelID,
		)
	}
	if snapshot.State == "" {
		snapshot.State = StateIdle
	}

	snapshot.Revision = uint64(state.Revision)
	snapshot.LastProviderTimestamp = state.ProviderTimestamp
	return snapshot, nil
}

func nextRevision(revision int64) (uint64, error) {
	if revision < 0 {
		return 0, fmt.Errorf("dota match revision %d must not be negative", revision)
	}
	if revision == int64(^uint64(0)>>1) {
		return 0, errors.New("dota match revision cannot be incremented")
	}

	return uint64(revision + 1), nil
}

func transitionInput(
	state model.MatchState,
	next Snapshot,
	actions []LifecycleAction,
) (dotarepository.ApplyMatchStateTransitionInput, error) {
	snapshot, err := json.Marshal(next)
	if err != nil {
		return dotarepository.ApplyMatchStateTransitionInput{}, err
	}

	outboxActions := make([]model.OutboxActionInput, 0, len(actions))
	for index, action := range actions {
		sequence, err := lifecycleActionSequence(next.Revision, index)
		if err != nil {
			return dotarepository.ApplyMatchStateTransitionInput{}, err
		}
		outboxAction, err := outboxActionInput(action, sequence)
		if err != nil {
			return dotarepository.ApplyMatchStateTransitionInput{}, err
		}
		outboxActions = append(outboxActions, outboxAction)
	}

	return dotarepository.ApplyMatchStateTransitionInput{
		ChannelID:         next.ChannelID,
		ExpectedRevision:  state.Revision,
		ProviderTimestamp: next.LastProviderTimestamp,
		Snapshot:          snapshot,
		Actions:           outboxActions,
	}, nil
}

func lifecycleActionSequence(revision uint64, index int) (int64, error) {
	if revision == 0 {
		return 0, errors.New("dota lifecycle action revision must be positive")
	}
	if index < 0 || index >= lifecycleActionsPerRevision {
		return 0, fmt.Errorf("dota lifecycle action index %d is out of range", index)
	}

	maxSequence := uint64(^uint64(0) >> 1)
	actionOffset := uint64(index + 1)
	if revision-1 > (maxSequence-actionOffset)/lifecycleActionsPerRevision {
		return 0, errors.New("dota lifecycle action sequence cannot be represented")
	}

	return int64((revision-1)*lifecycleActionsPerRevision + actionOffset), nil
}

func outboxActionInput(action LifecycleAction, sequence int64) (model.OutboxActionInput, error) {
	var actionKind model.OutboxAction
	switch action.Kind {
	case ActionCreate:
		actionKind = model.OutboxActionCreate
	case ActionResolve:
		actionKind = model.OutboxActionResolve
	case ActionCancel:
		actionKind = model.OutboxActionCancel
	default:
		return model.OutboxActionInput{}, fmt.Errorf("invalid lifecycle action kind %q", action.Kind)
	}

	payload, err := json.Marshal(action)
	if err != nil {
		return model.OutboxActionInput{}, err
	}

	return model.OutboxActionInput{
		ChannelID: action.ChannelID,
		MatchID:   action.MatchID,
		Action:    actionKind,
		Sequence:  sequence,
		Payload:   payload,
	}, nil
}

func isActivePayload(payload gsi.Payload) bool {
	return payload.Map != nil &&
		payload.Player != nil &&
		payload.Player.Activity == gsi.PlayerActivityPlaying
}

func transition(
	current Snapshot,
	payload gsi.Payload,
) (Snapshot, []LifecycleAction, []gsi.Event, bool, bool) {
	if payload.Map != nil && payload.Map.GameState == gsi.GameStatePostGame {
		if current.MatchID == 0 || payload.Map.MatchID != current.MatchID {
			return current, nil, nil, false, false
		}
		if !isNewerSource(current, payload.Provider.Timestamp, payload.Map.MatchID, payload.Map.GameTime) {
			return current, nil, nil, false, false
		}
	}

	if !isActivePayload(payload) {
		return transitionInactive(current, payload.Provider.Timestamp)
	}

	if !isNewerSource(current, payload.Provider.Timestamp, payload.Map.MatchID, payload.Map.GameTime) {
		return current, nil, nil, false, false
	}

	next := current
	next.LastProviderTimestamp = payload.Provider.Timestamp
	next.LastGameTime = payload.Map.GameTime
	next.RadiantScore = payload.Map.RadiantScore
	next.DireScore = payload.Map.DireScore
	next.GameTime = payload.Map.GameTime

	newState, known := mapGameState(payload.Map.GameState)
	if !known {
		newState = current.State
	}

	if newState == StatePostGame {
		return transitionPostGame(current, next, payload)
	}

	teamChanged := false
	if current.MatchID != 0 && current.MatchID == payload.Map.MatchID {
		teamChanged = updateTeam(&next, payload.Player.TeamName)
	}

	matchChanged := newState == StateInGame &&
		payload.Map.MatchID > 0 &&
		payload.Map.MatchID != current.MatchID
	if matchChanged {
		actions := make([]LifecycleAction, 0, 2)
		if current.MatchID > 0 {
			actions = append(actions, LifecycleAction{
				Kind:           ActionCancel,
				ChannelID:      current.ChannelID,
				MatchID:        current.MatchID,
				SteamAccountID: current.SteamAccountID,
			})
		}
		startMatch(&next, payload)
		actions = append(actions, LifecycleAction{
			Kind:           ActionCreate,
			ChannelID:      next.ChannelID,
			MatchID:        next.MatchID,
			SteamAccountID: next.SteamAccountID,
			HeroName:       next.HeroName,
			TeamKnown:      next.TeamKnown,
		})
		events := recordEvents(&next, payload.Events)
		emitStateUpdate := newState != current.State ||
			current.RadiantScore != next.RadiantScore ||
			current.DireScore != next.DireScore ||
			teamChanged
		return next, actions, events, emitStateUpdate, true
	}

	next.State = newState
	next.InGame = newState != StateIdle
	events := recordEvents(&next, payload.Events)
	emitStateUpdate := newState != current.State ||
		current.RadiantScore != next.RadiantScore ||
		current.DireScore != next.DireScore ||
		teamChanged
	return next, nil, events, emitStateUpdate, true
}

func transitionInactive(
	current Snapshot,
	providerTimestamp int64,
) (Snapshot, []LifecycleAction, []gsi.Event, bool, bool) {
	if providerTimestamp <= current.LastProviderTimestamp {
		return current, nil, nil, false, false
	}

	next := current
	next.LastProviderTimestamp = providerTimestamp
	if current.MatchID == 0 && current.State == StateIdle && !current.InGame {
		return next, nil, nil, false, true
	}

	var actions []LifecycleAction
	if current.MatchID > 0 {
		actions = append(actions, LifecycleAction{
			Kind:           ActionCancel,
			ChannelID:      current.ChannelID,
			MatchID:        current.MatchID,
			SteamAccountID: current.SteamAccountID,
		})
	}
	clearMatch(&next)
	return next, actions, nil, true, true
}

func transitionPostGame(
	current Snapshot,
	next Snapshot,
	payload gsi.Payload,
) (Snapshot, []LifecycleAction, []gsi.Event, bool, bool) {
	teamChanged := updateTeam(&next, payload.Player.TeamName)
	validResult := next.TeamKnown &&
		(payload.Map.WinTeam == gsi.WinTeamRadiant || payload.Map.WinTeam == gsi.WinTeamDire)
	if !validResult {
		next.State = StatePostGame
		next.InGame = true
		emitStateUpdate := next.State != current.State ||
			current.RadiantScore != next.RadiantScore ||
			current.DireScore != next.DireScore ||
			teamChanged
		return next, nil, nil, emitStateUpdate, true
	}

	won := (payload.Map.WinTeam == gsi.WinTeamRadiant) == next.IsRadiant
	action := LifecycleAction{
		Kind:           ActionResolve,
		ChannelID:      current.ChannelID,
		MatchID:        current.MatchID,
		SteamAccountID: current.SteamAccountID,
		Win:            won,
		HeroName:       current.HeroName,
		TeamKnown:      current.TeamKnown,
	}
	clearMatch(&next)
	return next, []LifecycleAction{action}, nil, true, true
}

func isNewerSource(
	current Snapshot,
	providerTimestamp int64,
	matchID int64,
	gameTime int,
) bool {
	if current.Revision == 0 {
		return true
	}
	if providerTimestamp > current.LastProviderTimestamp {
		return true
	}
	if providerTimestamp < current.LastProviderTimestamp {
		return false
	}

	return current.MatchID > 0 &&
		current.MatchID == matchID &&
		gameTime > current.LastGameTime
}

func startMatch(snapshot *Snapshot, payload gsi.Payload) {
	heroName := ""
	if payload.Hero != nil {
		heroName = stripHeroPrefix(payload.Hero.Name)
	}

	snapshot.MatchID = payload.Map.MatchID
	snapshot.WinProbability = 0
	snapshot.HeroName = heroName
	snapshot.TeamKnown = payload.Player.TeamName == "radiant" || payload.Player.TeamName == "dire"
	snapshot.IsRadiant = payload.Player.TeamName == "radiant"
	snapshot.SteamAccountID = strconv.FormatInt(payload.Player.AccountID, 10)
	snapshot.State = StateInGame
	snapshot.InGame = true
	snapshot.SeenEvents = nil
	snapshot.EventClaims = nil
}

func clearMatch(snapshot *Snapshot) {
	snapshot.State = StateIdle
	snapshot.InGame = false
	snapshot.MatchID = 0
	snapshot.HeroName = ""
	snapshot.IsRadiant = false
	snapshot.TeamKnown = false
	snapshot.SteamAccountID = ""
	snapshot.RadiantScore = 0
	snapshot.DireScore = 0
	snapshot.GameTime = 0
	snapshot.WinProbability = 0
	snapshot.SeenEvents = nil
	snapshot.EventClaims = nil
}

func updateTeam(snapshot *Snapshot, teamName string) bool {
	if teamName != "radiant" && teamName != "dire" {
		return false
	}

	isRadiant := teamName == "radiant"
	if snapshot.TeamKnown && snapshot.IsRadiant == isRadiant {
		return false
	}

	snapshot.TeamKnown = true
	snapshot.IsRadiant = isRadiant
	return true
}

func recordEvents(snapshot *Snapshot, events []gsi.Event) []gsi.Event {
	newEvents := make([]gsi.Event, 0, len(events))
	pendingEventKeys := make(map[string]struct{})
	for _, event := range events {
		key := eventKey(event)
		if seenEvent(snapshot.SeenEvents, key) {
			continue
		}

		switch event.EventType {
		case eventRoshanKilled, eventAegisPicked:
			if _, pending := pendingEventKeys[key]; pending {
				continue
			}
			pendingEventKeys[key] = struct{}{}
			newEvents = append(newEvents, event)
		default:
			snapshot.SeenEvents = append(snapshot.SeenEvents, key)
		}
	}

	return newEvents
}

func eventKey(event gsi.Event) string {
	return fmt.Sprintf("%s:%d", event.EventType, event.GameTime)
}

func (m *StateMachine) retryPendingEvents(ctx context.Context, snapshot Snapshot, payload gsi.Payload) {
	if !snapshot.InGame ||
		snapshot.MatchID == 0 ||
		!isActivePayload(payload) ||
		payload.Map.MatchID != snapshot.MatchID {
		return
	}

	m.publishEvents(ctx, snapshot, recordEvents(&snapshot, payload.Events))
}

func (m *StateMachine) claimEventDelivery(
	ctx context.Context,
	channelID uuid.UUID,
	expectedMatchID int64,
	key string,
) (string, bool, error) {
	claimID, err := uuid.NewRandom()
	if err != nil {
		return "", false, fmt.Errorf("generate dota event claim token: %w", err)
	}
	token := claimID.String()

	for attempt := 0; attempt < maxTransitionRetries; attempt++ {
		state, err := m.repo.GetMatchState(ctx, channelID)
		if err != nil {
			return "", false, fmt.Errorf("get dota match state for event claim: %w", err)
		}

		current, err := snapshotFromMatchState(state, channelID)
		if err != nil {
			return "", false, fmt.Errorf("decode dota match state for event claim: %w", err)
		}
		if !current.InGame || current.MatchID != expectedMatchID || seenEvent(current.SeenEvents, key) {
			return "", false, nil
		}

		now := time.Now()
		if claim, exists := current.EventClaims[key]; exists && claim.Token != "" && claim.ExpiresAt > now.UnixNano() {
			return "", false, nil
		}

		next := current
		next.EventClaims = cloneEventClaims(current.EventClaims)
		if next.EventClaims == nil {
			next.EventClaims = make(map[string]EventClaim)
		}
		next.EventClaims[key] = EventClaim{
			Token:     token,
			ExpiresAt: now.Add(eventClaimLease).UnixNano(),
		}
		next.Revision, err = nextRevision(state.Revision)
		if err != nil {
			return "", false, err
		}

		input, err := transitionInput(state, next, nil)
		if err != nil {
			return "", false, fmt.Errorf("encode dota event claim transition: %w", err)
		}
		committed, err := m.repo.ApplyMatchStateTransition(ctx, input)
		if err != nil {
			return "", false, fmt.Errorf("apply dota event claim transition: %w", err)
		}
		if committed {
			return token, true, nil
		}
	}

	return "", false, errors.New("dota event claim transition retry limit reached")
}

func (m *StateMachine) markEventDelivered(
	ctx context.Context,
	channelID uuid.UUID,
	expectedMatchID int64,
	key string,
	token string,
) error {
	for attempt := 0; attempt < maxTransitionRetries; attempt++ {
		state, err := m.repo.GetMatchState(ctx, channelID)
		if err != nil {
			return fmt.Errorf("get dota match state for event acknowledgement: %w", err)
		}

		current, err := snapshotFromMatchState(state, channelID)
		if err != nil {
			return fmt.Errorf("decode dota match state for event acknowledgement: %w", err)
		}
		if !current.InGame || current.MatchID != expectedMatchID || seenEvent(current.SeenEvents, key) {
			return nil
		}
		claim, exists := current.EventClaims[key]
		if !exists || claim.Token != token {
			return nil
		}

		next := current
		next.SeenEvents = append(append([]string(nil), current.SeenEvents...), key)
		next.EventClaims = cloneEventClaims(current.EventClaims)
		delete(next.EventClaims, key)
		next.Revision, err = nextRevision(state.Revision)
		if err != nil {
			return err
		}

		input, err := transitionInput(state, next, nil)
		if err != nil {
			return fmt.Errorf("encode dota event acknowledgement transition: %w", err)
		}
		committed, err := m.repo.ApplyMatchStateTransition(ctx, input)
		if err != nil {
			return fmt.Errorf("apply dota event acknowledgement transition: %w", err)
		}
		if committed {
			return nil
		}
	}

	return errors.New("dota event acknowledgement transition retry limit reached")
}

func (m *StateMachine) releaseEventClaim(
	ctx context.Context,
	channelID uuid.UUID,
	expectedMatchID int64,
	key string,
	token string,
) error {
	for attempt := 0; attempt < maxTransitionRetries; attempt++ {
		state, err := m.repo.GetMatchState(ctx, channelID)
		if err != nil {
			return fmt.Errorf("get dota match state for event claim release: %w", err)
		}

		current, err := snapshotFromMatchState(state, channelID)
		if err != nil {
			return fmt.Errorf("decode dota match state for event claim release: %w", err)
		}
		if !current.InGame || current.MatchID != expectedMatchID {
			return nil
		}
		claim, exists := current.EventClaims[key]
		if !exists || claim.Token != token {
			return nil
		}

		next := current
		next.EventClaims = cloneEventClaims(current.EventClaims)
		delete(next.EventClaims, key)
		next.Revision, err = nextRevision(state.Revision)
		if err != nil {
			return err
		}

		input, err := transitionInput(state, next, nil)
		if err != nil {
			return fmt.Errorf("encode dota event claim release transition: %w", err)
		}
		committed, err := m.repo.ApplyMatchStateTransition(ctx, input)
		if err != nil {
			return fmt.Errorf("apply dota event claim release: %w", err)
		}
		if committed {
			return nil
		}
	}

	return errors.New("dota event claim release retry limit reached")
}

func cloneEventClaims(claims map[string]EventClaim) map[string]EventClaim {
	if len(claims) == 0 {
		return nil
	}

	cloned := make(map[string]EventClaim, len(claims))
	for key, claim := range claims {
		cloned[key] = claim
	}

	return cloned
}

func seenEvent(events []string, key string) bool {
	for _, event := range events {
		if event == key {
			return true
		}
	}

	return false
}

func (m *StateMachine) loadSettings(ctx context.Context, snapshot *Snapshot) {
	settings, err := m.repo.GetByChannelID(ctx, snapshot.ChannelID)
	if err != nil {
		m.logger.WarnContext(ctx, "dota match: failed to load settings", logger.Error(err))
		return
	}

	snapshot.Mmr = settings.Mmr
	snapshot.SessionWins = settings.SessionWins
	snapshot.SessionLosses = settings.SessionLosses
}

func (m *StateMachine) publishLifecycleNotifications(
	ctx context.Context,
	snapshot Snapshot,
	actions []LifecycleAction,
) {
	for _, action := range actions {
		switch action.Kind {
		case ActionCreate:
			if err := m.emitter.MatchStarted(ctx, busdota.MatchStartedMessage{
				ChannelID:      snapshot.ChannelID.String(),
				SteamAccountID: snapshot.SteamAccountID,
				HeroName:       action.HeroName,
				MatchID:        action.MatchID,
				TeamKnown:      snapshot.TeamKnown,
			}); err != nil {
				m.logger.ErrorContext(ctx, "dota match: failed to emit match started", logger.Error(err))
			}
		case ActionCancel:
			if err := m.emitter.MatchAbandoned(ctx, busdota.MatchAbandonedMessage{
				ChannelID: action.ChannelID.String(),
				MatchID:   action.MatchID,
			}); err != nil {
				m.logger.ErrorContext(ctx, "dota match: failed to emit match abandoned", logger.Error(err))
			}
		}
	}
}

func (m *StateMachine) publishEvents(ctx context.Context, snapshot Snapshot, events []gsi.Event) {
	for _, event := range events {
		key := eventKey(event)
		token, claimed, err := m.claimEventDelivery(ctx, snapshot.ChannelID, snapshot.MatchID, key)
		if err != nil {
			m.logger.ErrorContext(ctx, "dota match: failed to claim event delivery", logger.Error(err))
			continue
		}
		if !claimed {
			continue
		}

		switch event.EventType {
		case eventRoshanKilled:
			if err := m.emitter.RoshanKilled(ctx, busdota.RoshanKilledMessage{
				ChannelID: snapshot.ChannelID.String(),
				Team:      event.KillerTeam,
				GameTime:  event.GameTime,
			}); err != nil {
				m.logger.ErrorContext(ctx, "dota match: failed to emit roshan killed", logger.Error(err))
				if releaseErr := m.releaseEventClaim(ctx, snapshot.ChannelID, snapshot.MatchID, key, token); releaseErr != nil {
					m.logger.ErrorContext(ctx, "dota match: failed to release event claim", logger.Error(releaseErr))
				}
				continue
			}
		case eventAegisPicked:
			if err := m.emitter.AegisPickup(ctx, busdota.AegisPickupMessage{
				ChannelID: snapshot.ChannelID.String(),
				PlayerID:  event.AegisPlayerID(),
				GameTime:  event.GameTime,
			}); err != nil {
				m.logger.ErrorContext(ctx, "dota match: failed to emit aegis pickup", logger.Error(err))
				if releaseErr := m.releaseEventClaim(ctx, snapshot.ChannelID, snapshot.MatchID, key, token); releaseErr != nil {
					m.logger.ErrorContext(ctx, "dota match: failed to release event claim", logger.Error(releaseErr))
				}
				continue
			}
		}

		if err := m.markEventDelivered(ctx, snapshot.ChannelID, snapshot.MatchID, key, token); err != nil {
			m.logger.ErrorContext(ctx, "dota match: failed to acknowledge event delivery", logger.Error(err))
		}
	}
}

func (m *StateMachine) publishStateUpdate(ctx context.Context, snapshot Snapshot) error {
	return m.emitter.StateUpdate(ctx, busapi.DotaStateUpdateMessage{
		ChannelID:      snapshot.ChannelID.String(),
		InGame:         snapshot.InGame,
		Mmr:            snapshot.Mmr,
		SessionWins:    snapshot.SessionWins,
		SessionLosses:  snapshot.SessionLosses,
		WinProbability: snapshot.WinProbability,
		HeroName:       snapshot.HeroName,
		MatchID:        snapshot.MatchID,
	})
}

func (m *StateMachine) emitStateUpdate(ctx context.Context, snapshot Snapshot) {
	if err := m.publishStateUpdate(ctx, snapshot); err != nil {
		m.logger.ErrorContext(ctx, "dota match: failed to emit state update", logger.Error(err))
	}
}
