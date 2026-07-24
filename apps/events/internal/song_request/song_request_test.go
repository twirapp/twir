package song_request

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/generic"
	"github.com/twirapp/twir/libs/bus-core/ytsr"
	cfg "github.com/twirapp/twir/libs/config"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type songRequestQueue[Req, Res any] struct {
	requestResponse *buscore.QueueResponse[Res]
	published       []Req
}

func (q *songRequestQueue[Req, Res]) Publish(_ context.Context, data Req) error {
	q.published = append(q.published, data)
	return nil
}

func (q *songRequestQueue[Req, Res]) Request(_ context.Context, _ Req) (*buscore.QueueResponse[Res], error) {
	return q.requestResponse, nil
}

func (q *songRequestQueue[Req, Res]) SubscribeGroup(
	string,
	buscore.QueueSubscribeCallback[Req, Res],
) error {
	return nil
}

func (q *songRequestQueue[Req, Res]) Subscribe(buscore.QueueSubscribeCallback[Req, Res]) error {
	return nil
}

func (q *songRequestQueue[Req, Res]) Unsubscribe() {}

type songRequestChannelsRepositoryFake struct {
	channelsrepository.Repository

	channel channelentity.Channel
}

func (r *songRequestChannelsRepositoryFake) GetByID(
	context.Context,
	uuid.UUID,
) (channelentity.Channel, error) {
	return r.channel, nil
}

func TestProcessFromDonationPublishesCanonicalTwitchCommands(t *testing.T) {
	channelID := uuid.New()
	bindingID := uuid.New()
	bindingUserID := uuid.New()
	platformChannelID := "twitch-channel-123"
	publishedMessages := &songRequestQueue[generic.ChatMessage, struct{}]{}
	searchQueue := &songRequestQueue[ytsr.SearchRequest, ytsr.SearchResponse]{
		requestResponse: &buscore.QueueResponse[ytsr.SearchResponse]{
			Data: ytsr.SearchResponse{
				Songs: []ytsr.Song{{Id: "first-song"}, {Id: "second-song"}},
			},
		},
	}
	bus := buscore.NewNatsBus(nil)
	bus.YTSRSearch = searchQueue
	bus.Parser.ProcessMessageAsCommand = publishedMessages

	channelService := channelservice.NewChannelService(
		&songRequestChannelsRepositoryFake{
			channel: channelentity.Channel{
				ID: channelID,
				Bindings: []channelplatformentity.ChannelPlatform{
					{
						ID:                bindingID,
						ChannelID:         channelID,
						Platform:          platform.PlatformTwitch,
						UserID:            bindingUserID,
						PlatformChannelID: platformChannelID,
					},
				},
			},
		},
		bus,
		cfg.Config{},
		nil,
		nil,
	)
	service := &SongRequest{
		gorm:           newSongRequestTestDB(t),
		twirBus:        bus,
		logger:         slog.Default(),
		channelService: channelService,
	}

	err := service.ProcessFromDonation(
		context.Background(),
		ProcessFromDonationInput{Text: "https://youtu.be/example", ChannelID: channelID.String()},
	)
	if err != nil {
		t.Fatalf("process donation song request: %v", err)
	}
	if len(publishedMessages.published) != 2 {
		t.Fatalf("published messages = %d, want 2", len(publishedMessages.published))
	}

	ids := make(map[string]struct{}, len(publishedMessages.published))
	for _, message := range publishedMessages.published {
		if message.ID == "" {
			t.Fatal("published message ID is empty")
		}
		if message.MessageID != message.ID {
			t.Fatalf("message ID = %q, want %q", message.MessageID, message.ID)
		}
		if _, exists := ids[message.ID]; exists {
			t.Fatalf("duplicate generated message ID %q", message.ID)
		}
		ids[message.ID] = struct{}{}

		if message.Platform != string(platform.PlatformTwitch) {
			t.Fatalf("platform = %q, want %q", message.Platform, platform.PlatformTwitch)
		}
		if message.ChannelID != channelID.String() {
			t.Fatalf("channel ID = %q, want %q", message.ChannelID, channelID)
		}
		if message.ChannelBindingID != bindingID.String() {
			t.Fatalf("channel binding ID = %q, want %q", message.ChannelBindingID, bindingID)
		}
		if message.UserID != bindingUserID.String() {
			t.Fatalf("user ID = %q, want binding user %q", message.UserID, bindingUserID)
		}

		for field, value := range map[string]string{
			"broadcaster user ID": message.BroadcasterUserId,
			"chatter user ID":     message.ChatterUserId,
			"platform channel ID": message.PlatformChannelID,
			"sender ID":           message.SenderID,
		} {
			if value != platformChannelID {
				t.Fatalf("%s = %q, want %q", field, value, platformChannelID)
			}
		}
	}
}

func newSongRequestTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(
		postgres.Open("host=127.0.0.1 user=twir dbname=twir sslmode=disable"),
		&gorm.Config{DisableAutomaticPing: true, DryRun: true},
	)
	if err != nil {
		t.Fatalf("open dry-run database: %v", err)
	}
	if err := db.Callback().Query().Before("gorm:query").Register(
		"song_request:seed_donation_command",
		func(tx *gorm.DB) {
			switch destination := tx.Statement.Dest.(type) {
			case *model.ChannelSongRequestsSettings:
				*destination = model.ChannelSongRequestsSettings{
					Enabled:                     true,
					TakeSongFromDonationMessage: true,
				}
			case *model.ChannelsCommands:
				*destination = model.ChannelsCommands{Enabled: true, Name: "sr"}
			}
			tx.RowsAffected = 1
		},
	); err != nil {
		t.Fatalf("register dry-run query callback: %v", err)
	}

	return db
}
