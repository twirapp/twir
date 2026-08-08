package pgx

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	jackcpgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/twirapp/twir/libs/repositories/dota"
	"github.com/twirapp/twir/libs/repositories/dota/model"
)

func TestApplyMatchResultOnce(t *testing.T) {
	channelID := uuid.New()
	settings := model.ChannelDotaSettings{
		ID:            uuid.New(),
		ChannelID:     channelID,
		Mmr:           1_525,
		SessionWins:   4,
		SessionLosses: 2,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	t.Run("applies a newly settled win once", func(t *testing.T) {
		executor := &matchResultExecutorFake{
			ledgerInserted: true,
			settings:       settings,
		}

		got, err := (&Pgx{}).applyMatchResultOnce(
			context.Background(),
			executor,
			dota.ApplyMatchResultInput{
				ChannelID: channelID,
				MatchID:   123,
				Won:       true,
				MmrDelta:  25,
			},
		)
		if err != nil {
			t.Fatalf("applyMatchResultOnce() error = %v", err)
		}
		if got != settings {
			t.Errorf("applyMatchResultOnce() = %#v, want %#v", got, settings)
		}
		if executor.ledgerInsertCalls != 1 {
			t.Errorf("ledger insert calls = %d, want 1", executor.ledgerInsertCalls)
		}
		if executor.settingsUpdateCalls != 1 {
			t.Errorf("settings update calls = %d, want 1", executor.settingsUpdateCalls)
		}
		if executor.settingsFetchCalls != 1 {
			t.Errorf("settings fetch calls = %d, want 1", executor.settingsFetchCalls)
		}
		if got := executor.settingsUpdateArguments; len(got) != 4 || got[0] != channelID || got[1] != 25 || got[2] != 1 || got[3] != 0 {
			t.Errorf("settings update arguments = %#v, want [%s 25 1 0]", got, channelID)
		}
	})

	t.Run("applies a newly settled loss once", func(t *testing.T) {
		executor := &matchResultExecutorFake{
			ledgerInserted: true,
			settings:       settings,
		}

		_, err := (&Pgx{}).applyMatchResultOnce(
			context.Background(),
			executor,
			dota.ApplyMatchResultInput{
				ChannelID: channelID,
				MatchID:   124,
				Won:       false,
				MmrDelta:  -25,
			},
		)
		if err != nil {
			t.Fatalf("applyMatchResultOnce() error = %v", err)
		}
		if got := executor.settingsUpdateArguments; len(got) != 4 || got[0] != channelID || got[1] != -25 || got[2] != 0 || got[3] != 1 {
			t.Errorf("settings update arguments = %#v, want [%s -25 0 1]", got, channelID)
		}
	})

	t.Run("returns current settings without replaying an existing settlement", func(t *testing.T) {
		executor := &matchResultExecutorFake{
			ledgerInserted: false,
			settings:       settings,
		}

		got, err := (&Pgx{}).applyMatchResultOnce(
			context.Background(),
			executor,
			dota.ApplyMatchResultInput{
				ChannelID: channelID,
				MatchID:   123,
				Won:       true,
				MmrDelta:  25,
			},
		)
		if err != nil {
			t.Fatalf("applyMatchResultOnce() error = %v", err)
		}
		if got != settings {
			t.Errorf("applyMatchResultOnce() = %#v, want %#v", got, settings)
		}
		if executor.ledgerInsertCalls != 1 {
			t.Errorf("ledger insert calls = %d, want 1", executor.ledgerInsertCalls)
		}
		if executor.settingsUpdateCalls != 0 {
			t.Errorf("settings update calls = %d, want 0", executor.settingsUpdateCalls)
		}
		if executor.settingsFetchCalls != 1 {
			t.Errorf("settings fetch calls = %d, want 1", executor.settingsFetchCalls)
		}
	})
}

type matchResultExecutorFake struct {
	ledgerInserted          bool
	ledgerInsertCalls       int
	settingsUpdateCalls     int
	settingsUpdateArguments []any
	settingsFetchCalls      int
	settings                model.ChannelDotaSettings
}

func (f *matchResultExecutorFake) Exec(
	_ context.Context,
	query string,
	arguments ...any,
) (pgconn.CommandTag, error) {
	switch {
	case strings.Contains(query, "INSERT INTO dota_match_settlements"):
		f.ledgerInsertCalls++
		if f.ledgerInserted {
			return pgconn.NewCommandTag("INSERT 0 1"), nil
		}

		return pgconn.NewCommandTag("INSERT 0 0"), nil
	case strings.Contains(query, "UPDATE channels_dota_settings"):
		f.settingsUpdateCalls++
		f.settingsUpdateArguments = arguments
		return pgconn.NewCommandTag("UPDATE 1"), nil
	default:
		return pgconn.CommandTag{}, fmt.Errorf("unexpected exec query: %s", query)
	}
}

func (f *matchResultExecutorFake) QueryRow(
	_ context.Context,
	query string,
	_ ...any,
) jackcpgx.Row {
	if !strings.Contains(query, "FROM channels_dota_settings") {
		return matchResultRow{err: fmt.Errorf("unexpected query row query: %s", query)}
	}

	f.settingsFetchCalls++
	return matchResultRow{settings: f.settings}
}

type matchResultRow struct {
	settings model.ChannelDotaSettings
	err      error
}

func (r matchResultRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	if len(dest) != 14 {
		return fmt.Errorf("scan destinations = %d, want 14", len(dest))
	}

	*dest[0].(*uuid.UUID) = r.settings.ID
	*dest[1].(*uuid.UUID) = r.settings.ChannelID
	*dest[2].(*bool) = r.settings.Enabled
	*dest[3].(**string) = r.settings.SteamAccountID
	*dest[4].(*string) = r.settings.GsiToken
	*dest[5].(*int) = r.settings.Mmr
	*dest[6].(*int) = r.settings.MmrDelta
	*dest[7].(*int) = r.settings.SessionWins
	*dest[8].(*int) = r.settings.SessionLosses
	*dest[9].(*model.PredictionSettings) = r.settings.PredictionSettings
	*dest[10].(*model.ChatEvents) = r.settings.ChatEvents
	*dest[11].(*model.CommandsSettings) = r.settings.CommandsSettings
	*dest[12].(*time.Time) = r.settings.CreatedAt
	*dest[13].(*time.Time) = r.settings.UpdatedAt

	return nil
}
