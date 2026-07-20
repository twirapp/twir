package processor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/twirapp/kv"
	kvoptions "github.com/twirapp/kv/options"
	"github.com/twirapp/twir/apps/dota/internal/gsi"
	"github.com/twirapp/twir/apps/dota/internal/match"
	busapi "github.com/twirapp/twir/libs/bus-core/api"
	busdota "github.com/twirapp/twir/libs/bus-core/dota"
	dotarepository "github.com/twirapp/twir/libs/repositories/dota"
	"github.com/twirapp/twir/libs/repositories/dota/model"
)

type fakeWinProbabilityProvider struct {
	probability float64
	err         error
	matchIDs    []int64
	onCall      func()
}

func (f *fakeWinProbabilityProvider) WinProbability(_ context.Context, matchID int64) (float64, error) {
	f.matchIDs = append(f.matchIDs, matchID)
	if f.onCall != nil {
		f.onCall()
	}

	return f.probability, f.err
}

type fakeRepo struct {
	settings model.ChannelDotaSettings
}

func (f *fakeRepo) GetByChannelID(
	_ context.Context,
	_ uuid.UUID,
) (model.ChannelDotaSettings, error) {
	return f.settings, nil
}

func (f *fakeRepo) GetByGsiToken(
	_ context.Context,
	_ string,
) (model.ChannelDotaSettings, error) {
	return model.Nil, dotarepository.ErrNotFound
}

func (f *fakeRepo) Create(
	_ context.Context,
	_ dotarepository.CreateInput,
) (model.ChannelDotaSettings, error) {
	return model.Nil, errors.New("not implemented")
}

func (f *fakeRepo) Update(
	_ context.Context,
	_ uuid.UUID,
	_ dotarepository.UpdateInput,
) (model.ChannelDotaSettings, error) {
	return model.Nil, errors.New("not implemented")
}

func (f *fakeRepo) UpdateMatchResult(
	_ context.Context,
	_ uuid.UUID,
	_ bool,
	_ int,
) (model.ChannelDotaSettings, error) {
	return f.settings, nil
}

func (f *fakeRepo) ResetSession(
	_ context.Context,
	_ uuid.UUID,
) (model.ChannelDotaSettings, error) {
	return model.Nil, errors.New("not implemented")
}

func (f *fakeRepo) RegenerateGsiToken(
	_ context.Context,
	_ uuid.UUID,
) (model.ChannelDotaSettings, error) {
	return model.Nil, errors.New("not implemented")
}

type fakeEmitter struct {
	roshanErr    error
	stateUpdates []busapi.DotaStateUpdateMessage
}

func (f *fakeEmitter) MatchStarted(_ context.Context, _ busdota.MatchStartedMessage) error {
	return nil
}

func (f *fakeEmitter) MatchEnded(_ context.Context, _ busdota.MatchEndedMessage) error {
	return nil
}

func (f *fakeEmitter) RoshanKilled(_ context.Context, _ busdota.RoshanKilledMessage) error {
	return f.roshanErr
}

func (f *fakeEmitter) AegisPickup(_ context.Context, _ busdota.AegisPickupMessage) error {
	return nil
}

func (f *fakeEmitter) StateUpdate(_ context.Context, msg busapi.DotaStateUpdateMessage) error {
	f.stateUpdates = append(f.stateUpdates, msg)
	return nil
}

type fakeValuer struct {
	bytes []byte
	err   error
}

func (v fakeValuer) Int() (int64, error)     { return 0, v.err }
func (v fakeValuer) String() (string, error) { return string(v.bytes), v.err }
func (v fakeValuer) Bytes() ([]byte, error)  { return v.bytes, v.err }
func (v fakeValuer) Bool() (bool, error)     { return false, v.err }
func (v fakeValuer) Float() (float64, error) { return 0, v.err }
func (v fakeValuer) Scan(_ any) error        { return v.err }
func (v fakeValuer) Err() error              { return v.err }

type fakeKV struct {
	store map[string][]byte
}

func newFakeKV() *fakeKV {
	return &fakeKV{store: make(map[string][]byte)}
}

func (f *fakeKV) Get(_ context.Context, key string) kv.Valuer {
	value, ok := f.store[key]
	if !ok {
		return fakeValuer{err: kv.ErrKeyNil}
	}

	return fakeValuer{bytes: value}
}

func (f *fakeKV) Set(_ context.Context, key string, value any, _ ...kvoptions.Option) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("unsupported value type %T", value)
	}

	f.store[key] = bytes
	return nil
}

func (f *fakeKV) SetMany(_ context.Context, values []kv.SetMany) error {
	for _, value := range values {
		if err := f.Set(context.Background(), value.Key, value.Value, value.Options...); err != nil {
			return err
		}
	}

	return nil
}

func (f *fakeKV) Delete(_ context.Context, key string) error {
	delete(f.store, key)
	return nil
}

func (f *fakeKV) DeleteMany(_ context.Context, keys []string) error {
	for _, key := range keys {
		delete(f.store, key)
	}

	return nil
}

func (f *fakeKV) Exists(_ context.Context, key string) (bool, error) {
	_, ok := f.store[key]
	return ok, nil
}

func (f *fakeKV) ExistsMany(_ context.Context, keys []string) ([]bool, error) {
	exists := make([]bool, len(keys))
	for i, key := range keys {
		exists[i], _ = f.Exists(context.Background(), key)
	}

	return exists, nil
}

func (f *fakeKV) GetKeysByPattern(_ context.Context, _ string) ([]string, error) {
	return nil, nil
}

func newStateMachine() (*match.StateMachine, *fakeEmitter, uuid.UUID) {
	channelID := uuid.New()
	repo := &fakeRepo{
		settings: model.ChannelDotaSettings{
			ChannelID: channelID,
			Mmr:       3000,
			MmrDelta:  25,
		},
	}
	emitter := &fakeEmitter{}

	return match.New(repo, emitter, newFakeKV(), slog.Default()), emitter, channelID
}

func inGamePayload(matchID int64) gsi.Payload {
	return gsi.Payload{
		Map: &gsi.Map{
			MatchID:   matchID,
			GameState: gsi.GameStateInProgress,
		},
		Player: &gsi.Player{
			Activity:  gsi.PlayerActivityPlaying,
			TeamName:  "radiant",
			AccountID: 12345,
		},
	}
}

func TestProcessRefreshesWinProbabilityAfterProcessingLiveMatch(t *testing.T) {
	sm, emitter, channelID := newStateMachine()
	provider := &fakeWinProbabilityProvider{probability: 0.625}
	processor := New(sm, provider, slog.Default())
	ctx := context.Background()

	provider.onCall = func() {
		snapshot, err := sm.GetSnapshot(ctx, channelID)
		require.NoError(t, err)
		require.True(t, snapshot.InGame)
		require.Equal(t, int64(2001), snapshot.MatchID)
	}

	require.NoError(t, processor.Process(ctx, channelID, inGamePayload(2001)))
	require.Equal(t, []int64{2001}, provider.matchIDs)

	snapshot, err := sm.GetSnapshot(ctx, channelID)
	require.NoError(t, err)
	require.Equal(t, 0.625, snapshot.WinProbability)
	require.Len(t, emitter.stateUpdates, 2)
	require.Equal(t, 0.625, emitter.stateUpdates[1].WinProbability)
}

func TestProcessKeepsGsiAvailableWhenWinProbabilityFails(t *testing.T) {
	sm, emitter, channelID := newStateMachine()
	provider := &fakeWinProbabilityProvider{err: errors.New("stratz unavailable")}
	processor := New(sm, provider, slog.Default())

	require.NoError(t, processor.Process(context.Background(), channelID, inGamePayload(2002)))
	require.Equal(t, []int64{2002}, provider.matchIDs)

	snapshot, err := sm.GetSnapshot(context.Background(), channelID)
	require.NoError(t, err)
	require.True(t, snapshot.InGame)
	require.Equal(t, int64(2002), snapshot.MatchID)
	require.Zero(t, snapshot.WinProbability)
	require.Len(t, emitter.stateUpdates, 1)
}

func TestProcessReturnsStateMachineErrorBeforeWinProbability(t *testing.T) {
	sm, emitter, channelID := newStateMachine()
	emitter.roshanErr = errors.New("publish failed")
	provider := &fakeWinProbabilityProvider{probability: 0.625}
	processor := New(sm, provider, slog.Default())
	payload := inGamePayload(2003)
	payload.Events = []gsi.Event{{EventType: "roshan_killed", GameTime: 300}}

	err := processor.Process(context.Background(), channelID, payload)
	require.ErrorIs(t, err, emitter.roshanErr)
	require.Empty(t, provider.matchIDs)
}
