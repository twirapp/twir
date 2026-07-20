package dota2

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	buscore "github.com/twirapp/twir/libs/bus-core"
	busdota "github.com/twirapp/twir/libs/bus-core/dota"
	"github.com/twirapp/twir/libs/i18n"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestReadCommandsAreVisible(t *testing.T) {
	commands := map[string]*types.DefaultCommand{
		"mmr": Mmr,
		"wl":  Wl,
		"lg":  Lg,
		"gm":  Gm,
		"np":  Np,
		"wp":  Wp,
	}

	for name, command := range commands {
		if !command.ChannelsCommands.Visible {
			t.Errorf("%s command is not publicly visible", name)
		}
	}
}

func TestRequestDotaDataUsesFiveSecondDeadlineAndCancels(t *testing.T) {
	startedAt := time.Now()
	var requestCtx context.Context

	_, err := requestDotaData(
		context.Background(),
		func(ctx context.Context, _ busdota.GetDataRequest) (*buscore.QueueResponse[busdota.GetDataResponse], error) {
			requestCtx = ctx
			return &buscore.QueueResponse[busdota.GetDataResponse]{}, nil
		},
		busdota.GetDataRequest{},
	)
	if err != nil {
		t.Fatalf("requestDotaData() error = %v", err)
	}
	if requestCtx == nil {
		t.Fatal("requestDotaData() did not call the request function")
	}

	deadline, ok := requestCtx.Deadline()
	if !ok {
		t.Fatal("request context has no deadline")
	}
	if got := deadline.Sub(startedAt); got < 5*time.Second-100*time.Millisecond || got > 5*time.Second+100*time.Millisecond {
		t.Fatalf("request deadline offset = %s, want about 5s", got)
	}

	select {
	case <-requestCtx.Done():
		if !errors.Is(requestCtx.Err(), context.Canceled) {
			t.Fatalf("request context error = %v, want context canceled", requestCtx.Err())
		}
	case <-time.After(time.Second):
		t.Fatal("request context was not canceled after requestDotaData returned")
	}
}

func TestRequestDotaDataReturnsLocalizedCommandError(t *testing.T) {
	initializeDotaTestI18n(t)
	requestErr := errors.New("Dota bus backend unavailable")

	_, err := requestDotaData(
		context.Background(),
		func(context.Context, busdota.GetDataRequest) (*buscore.QueueResponse[busdota.GetDataResponse], error) {
			return nil, requestErr
		},
		busdota.GetDataRequest{},
	)

	var commandErr *types.CommandHandlerError
	if !errors.As(err, &commandErr) {
		t.Fatalf("requestDotaData() error = %T %v, want CommandHandlerError", err, err)
	}

	want := i18n.GetCtx(context.Background(), locales.Translations.Commands.Dota.Errors.GetData)
	if commandErr.Message != want {
		t.Errorf("CommandHandlerError message = %q, want %q", commandErr.Message, want)
	}
	if strings.Contains(commandErr.Message, requestErr.Error()) {
		t.Errorf("CommandHandlerError message exposes request error: %q", commandErr.Message)
	}
	if !errors.Is(commandErr.Err, requestErr) {
		t.Errorf("CommandHandlerError cause = %v, want wrapped %v", commandErr.Err, requestErr)
	}
}

func TestWinProbabilityOutputRequiresActiveAvailableData(t *testing.T) {
	tests := []struct {
		name          string
		data          *busdota.GetDataResponse
		wantOutput    string
		wantAvailable bool
	}{
		{
			name:          "inactive match",
			data:          &busdota.GetDataResponse{WinProbabilityAvailable: true, WinProbability: 0.625},
			wantAvailable: false,
		},
		{
			name:          "active match without probability",
			data:          &busdota.GetDataResponse{InGame: true},
			wantAvailable: false,
		},
		{
			name:          "active match with zero probability",
			data:          &busdota.GetDataResponse{InGame: true, WinProbabilityAvailable: true},
			wantOutput:    "0.0%",
			wantAvailable: true,
		},
		{
			name:          "active match with probability",
			data:          &busdota.GetDataResponse{InGame: true, WinProbabilityAvailable: true, WinProbability: 0.625},
			wantOutput:    "62.5%",
			wantAvailable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, available := winProbabilityOutput(tt.data)
			if output != tt.wantOutput {
				t.Errorf("winProbabilityOutput() output = %q, want %q", output, tt.wantOutput)
			}
			if available != tt.wantAvailable {
				t.Errorf("winProbabilityOutput() available = %t, want %t", available, tt.wantAvailable)
			}
		})
	}
}

func TestUpdateDotaMmrBuildsScopedTimestampedUpdate(t *testing.T) {
	db, err := gorm.Open(
		postgres.New(postgres.Config{DSN: "host=localhost user=postgres dbname=postgres sslmode=disable"}),
		&gorm.Config{DisableAutomaticPing: true, DryRun: true, SkipDefaultTransaction: true},
	)
	if err != nil {
		t.Fatalf("open GORM DryRun database: %v", err)
	}

	update := updateDotaMmr(context.Background(), db, "channel-id", 1234)
	if update.Error != nil {
		t.Fatalf("updateDotaMmr() error = %v", update.Error)
	}

	sql := update.Statement.SQL.String()
	for _, fragment := range []string{
		`UPDATE "channels_dota_settings"`,
		`"mmr"`,
		`"updated_at"`,
		"WHERE channel_id =",
	} {
		if !strings.Contains(sql, fragment) {
			t.Errorf("generated UPDATE = %q, missing %q", sql, fragment)
		}
	}
}

func initializeDotaTestI18n(t *testing.T) {
	t.Helper()

	if _, err := i18n.New(i18n.Opts{Store: locales.Store, DefaultLocale: "en"}); err != nil {
		t.Fatalf("initialize translations: %v", err)
	}
}
