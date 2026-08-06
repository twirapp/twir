package sr_youtube

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/samber/lo"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/internal/types/services"
	"github.com/twirapp/twir/apps/parser/locales"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/api"
	song_request_mode "github.com/twirapp/twir/libs/entities/song_request_mode"
	songrequestssettingsentity "github.com/twirapp/twir/libs/entities/song_requests_settings"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type wrongQueue[Req, Res any] struct {
	requestResponse *buscore.QueueResponse[Res]
	published       []Req
}

func (q *wrongQueue[Req, Res]) Publish(_ context.Context, data Req) error {
	q.published = append(q.published, data)
	return nil
}

func (q *wrongQueue[Req, Res]) Request(context.Context, Req) (*buscore.QueueResponse[Res], error) {
	return q.requestResponse, nil
}

func (q *wrongQueue[Req, Res]) SubscribeGroup(
	string,
	buscore.QueueSubscribeCallback[Req, Res],
) error {
	return nil
}

func (q *wrongQueue[Req, Res]) Subscribe(buscore.QueueSubscribeCallback[Req, Res]) error {
	return nil
}

func (q *wrongQueue[Req, Res]) Unsubscribe() {}

type wrongSettingsRepositoryFake struct {
	settings songrequestssettingsentity.Settings
}

func (r *wrongSettingsRepositoryFake) GetByChannelID(
	context.Context,
	string,
) (songrequestssettingsentity.Settings, error) {
	return r.settings, nil
}

func (r *wrongSettingsRepositoryFake) Upsert(
	context.Context,
	songrequestssettingsentity.Settings,
) (songrequestssettingsentity.Settings, error) {
	return songrequestssettingsentity.Nil, nil
}

func (r *wrongSettingsRepositoryFake) SetVolume(context.Context, string, int) error {
	return nil
}

func newWrongTestDB(t *testing.T, songs []*model.RequestedSong) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(
		postgres.Open("host=127.0.0.1 user=twir dbname=twir sslmode=disable"),
		&gorm.Config{DisableAutomaticPing: true, DryRun: true, SkipDefaultTransaction: true},
	)
	if err != nil {
		t.Fatalf("open dry-run database: %v", err)
	}
	if err := db.Callback().Query().Before("gorm:query").Register(
		"sr_wrong:seed_songs",
		func(tx *gorm.DB) {
			if destination, ok := tx.Statement.Dest.(*[]*model.RequestedSong); ok {
				*destination = songs
			}
			tx.RowsAffected = 1
		},
	); err != nil {
		t.Fatalf("register dry-run query callback: %v", err)
	}

	return db
}

func buildWrongParseContext(
	t *testing.T,
	input string,
	songs []*model.RequestedSong,
	publishQueue *wrongQueue[api.SongRequestRemoveFromQueue, struct{}],
) *types.ParseContext {
	t.Helper()

	if _, err := i18n.New(i18n.Opts{Store: locales.Store, DefaultLocale: "en"}); err != nil {
		t.Fatalf("init i18n: %v", err)
	}

	argsParser, err := command_arguments.NewParser(
		command_arguments.Opts{
			Args: []command_arguments.Arg{
				command_arguments.Int{
					Name:     songSkipArgName,
					Optional: true,
					Min:      lo.ToPtr(1),
				},
			},
			Input: input,
		},
	)
	if err != nil {
		t.Fatalf("build args parser: %v", err)
	}

	bus := buscore.NewNatsBus(nil)
	bus.Api.SongRequestRemoveFromQueue = publishQueue

	return &types.ParseContext{
		Channel: &types.ParseContextChannel{DBChannelID: "channel-id"},
		Sender: &types.ParseContextSender{
			ID:   "sender-platform-id",
			Name: "viewer",
			DbUser: &model.Users{
				ID: "sender-db-id",
			},
		},
		Services: &services.Services{
			Gorm: newWrongTestDB(t, songs),
			Bus:  bus,
			SongRequestsSettingsRepo: &wrongSettingsRepositoryFake{
				settings: songrequestssettingsentity.Settings{
					Enabled: true,
					Mode:    song_request_mode.ModeYouTube,
				},
			},
		},
		ArgsParser: argsParser,
	}
}

func TestWrongCommand_outOfBoundsDoesNotPanicOrDelete(t *testing.T) {
	songs := []*model.RequestedSong{
		{ID: "song-1", Title: "First", VideoID: "video-1"},
		{ID: "song-2", Title: "Second", VideoID: "video-2"},
	}
	publishQueue := &wrongQueue[api.SongRequestRemoveFromQueue, struct{}]{}
	parseCtx := buildWrongParseContext(t, "5", songs, publishQueue)

	result, err := WrongCommand.Handler(context.Background(), parseCtx)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if len(result.Result) != 1 {
		t.Fatalf("result messages = %d, want 1", len(result.Result))
	}
	if !strings.Contains(result.Result[0], "2") {
		t.Fatalf("result message %q should mention available songs count", result.Result[0])
	}
	if len(publishQueue.published) != 0 {
		t.Fatalf("published messages = %d, want 0", len(publishQueue.published))
	}
}

func TestWrongCommand_defaultDeletesLatestSong(t *testing.T) {
	songs := []*model.RequestedSong{
		{ID: "song-latest", Title: "Latest", VideoID: "video-latest"},
		{ID: "song-older", Title: "Older", VideoID: "video-older"},
	}
	publishQueue := &wrongQueue[api.SongRequestRemoveFromQueue, struct{}]{}
	parseCtx := buildWrongParseContext(t, "", songs, publishQueue)

	result, err := WrongCommand.Handler(context.Background(), parseCtx)
	if err != nil {
		var handlerErr *types.CommandHandlerError
		if errors.As(err, &handlerErr) {
			t.Fatalf("handler error: %v (inner: %v)", handlerErr.Message, handlerErr.Err)
		}
		t.Fatalf("handler error: %v", err)
	}
	if len(result.Result) != 1 {
		t.Fatalf("result messages = %d, want 1", len(result.Result))
	}
	if !strings.Contains(result.Result[0], "Latest") {
		t.Fatalf("result message %q should mention deleted song title", result.Result[0])
	}
	if len(publishQueue.published) != 1 {
		t.Fatalf("published messages = %d, want 1", len(publishQueue.published))
	}
	if publishQueue.published[0].VideoID != "video-latest" {
		t.Fatalf("published video id = %q, want video-latest", publishQueue.published[0].VideoID)
	}
}

func TestWrongCommand_deletesSongByNumber(t *testing.T) {
	songs := []*model.RequestedSong{
		{ID: "song-latest", Title: "Latest", VideoID: "video-latest"},
		{ID: "song-older", Title: "Older", VideoID: "video-older"},
	}
	publishQueue := &wrongQueue[api.SongRequestRemoveFromQueue, struct{}]{}
	parseCtx := buildWrongParseContext(t, "2", songs, publishQueue)

	result, err := WrongCommand.Handler(context.Background(), parseCtx)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if len(result.Result) != 1 {
		t.Fatalf("result messages = %d, want 1", len(result.Result))
	}
	if !strings.Contains(result.Result[0], "Older") {
		t.Fatalf("result message %q should mention deleted song title", result.Result[0])
	}
	if len(publishQueue.published) != 1 || publishQueue.published[0].VideoID != "video-older" {
		t.Fatalf("published = %#v, want video-older removal", publishQueue.published)
	}
}

func TestWrongCommand_noSongsReturnsInfoMessage(t *testing.T) {
	publishQueue := &wrongQueue[api.SongRequestRemoveFromQueue, struct{}]{}
	parseCtx := buildWrongParseContext(t, "1", nil, publishQueue)

	result, err := WrongCommand.Handler(context.Background(), parseCtx)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if len(result.Result) != 1 {
		t.Fatalf("result messages = %d, want 1", len(result.Result))
	}
	if len(publishQueue.published) != 0 {
		t.Fatalf("published messages = %d, want 0", len(publishQueue.published))
	}
}
