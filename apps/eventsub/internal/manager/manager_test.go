package manager

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

type fakeChannelsRepo struct {
	channel channelsmodel.Channel
	err     error
}

func (f *fakeChannelsRepo) GetMany(context.Context, channelsrepo.GetManyInput) ([]channelsmodel.Channel, error) {
	return nil, nil
}

func (f *fakeChannelsRepo) GetByID(_ context.Context, _ uuid.UUID) (channelsmodel.Channel, error) {
	if f.err != nil {
		return channelsmodel.Nil, f.err
	}

	return f.channel, nil
}

func (f *fakeChannelsRepo) GetByTwitchUserID(context.Context, uuid.UUID) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (f *fakeChannelsRepo) GetByTwitchPlatformID(context.Context, string) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (f *fakeChannelsRepo) GetByKickUserID(context.Context, uuid.UUID) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (f *fakeChannelsRepo) GetByKickPlatformID(context.Context, string) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (f *fakeChannelsRepo) GetCount(context.Context, channelsrepo.GetCountInput) (int, error) {
	return 0, nil
}

func (f *fakeChannelsRepo) Update(context.Context, uuid.UUID, channelsrepo.UpdateInput) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (f *fakeChannelsRepo) Create(context.Context, channelsrepo.CreateInput) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (f *fakeChannelsRepo) GetByApiKey(context.Context, string) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func TestResolveTwitchSubscriptionIdentities(t *testing.T) {
	twitchPlatformID := "tw-123"
	botID := "bot-456"
	m := &Manager{
		channelsRepo: &fakeChannelsRepo{
			channel: channelsmodel.Channel{TwitchPlatformID: &twitchPlatformID, BotID: botID},
		},
	}

	channelID := uuid.New().String()
	broadcasterID, resolvedBotID, err := m.resolveTwitchSubscriptionIdentities(context.Background(), channelID)
	if err != nil {
		t.Fatalf("resolveTwitchSubscriptionIdentities returned error: %v", err)
	}

	if broadcasterID != twitchPlatformID {
		t.Fatalf("expected broadcaster ID %q, got %q", twitchPlatformID, broadcasterID)
	}

	if resolvedBotID != botID {
		t.Fatalf("expected bot ID %q, got %q", botID, resolvedBotID)
	}
}

func TestResolveTwitchSubscriptionIdentitiesRequiresPlatformID(t *testing.T) {
	m := &Manager{
		channelsRepo: &fakeChannelsRepo{},
	}

	_, _, err := m.resolveTwitchSubscriptionIdentities(context.Background(), uuid.New().String())
	if err == nil {
		t.Fatal("expected error when Twitch platform ID is missing")
	}
}

func TestResolveTwitchSubscriptionIdentitiesReturnsRepositoryError(t *testing.T) {
	wantErr := errors.New("boom")
	m := &Manager{
		channelsRepo: &fakeChannelsRepo{err: wantErr},
	}

	_, _, err := m.resolveTwitchSubscriptionIdentities(context.Background(), uuid.New().String())
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected error %v, got %v", wantErr, err)
	}
}
