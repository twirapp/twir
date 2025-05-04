package pgx

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/repositories/overlays_kappagen"
	kappagenmodel "github.com/twirapp/twir/libs/repositories/overlays_kappagen/model"
)

type Opts struct {
	PgxPool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool:   opts.PgxPool,
		getter: trmpgx.DefaultCtxGetter,
	}
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(Opts{PgxPool: pool})
}

var _ overlays_kappagen.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

const kappagenOverlayType = "kappagen_overlay"

func (p *Pgx) Get(ctx context.Context, channelID string) (kappagenmodel.KappagenOverlay, error) {
	query := `
SELECT id, "channelId", settings
FROM channels_modules_settings
WHERE "channelId" = $1 AND type = $2
LIMIT 1
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	row := conn.QueryRow(ctx, query, channelID, kappagenOverlayType)

	var id uuid.UUID
	var dbChannelID string
	var settingsBytes []byte

	if err := row.Scan(&id, &dbChannelID, &settingsBytes); err != nil {
		return kappagenmodel.Nil, overlays_kappagen.ErrNotFound
	}

	var settings model.KappagenOverlaySettings
	if err := json.Unmarshal(settingsBytes, &settings); err != nil {
		return kappagenmodel.Nil, fmt.Errorf("failed to unmarshal settings: %w", err)
	}

	return mapModelToEntity(id, dbChannelID, settings), nil
}

func (p *Pgx) Create(ctx context.Context, input overlays_kappagen.CreateInput) (kappagenmodel.KappagenOverlay, error) {
	id := uuid.New()
	settings := mapCreateInputToModel(input)

	settingsBytes, err := json.Marshal(settings)
	if err != nil {
		return kappagenmodel.Nil, fmt.Errorf("failed to marshal settings: %w", err)
	}

	query := `
INSERT INTO channels_modules_settings (id, "channelId", type, settings)
VALUES ($1, $2, $3, $4)
RETURNING id, "channelId"
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	row := conn.QueryRow(ctx, query, id, input.ChannelID, kappagenOverlayType, settingsBytes)

	var resultID uuid.UUID
	var channelID string
	if err := row.Scan(&resultID, &channelID); err != nil {
		return kappagenmodel.Nil, fmt.Errorf("failed to create kappagen overlay: %w", err)
	}

	return mapModelToEntity(resultID, channelID, settings), nil
}

func (p *Pgx) Update(ctx context.Context, channelID string, input overlays_kappagen.UpdateInput) (kappagenmodel.KappagenOverlay, error) {
	// First, get the current settings
	currentOverlay, err := p.Get(ctx, channelID)
	if err != nil {
		return kappagenmodel.Nil, err
	}

	// Update with new values
	settings := model.KappagenOverlaySettings{
		EnableSpawn:    input.EnableSpawn,
		ExcludedEmotes: input.ExcludedEmotes,
		EnableRave:     input.EnableRave,
		Animation: model.KappagenOverlaySettingsAnimation{
			FadeIn:  input.Animation.FadeIn,
			FadeOut: input.Animation.FadeOut,
			ZoomIn:  input.Animation.ZoomIn,
			ZoomOut: input.Animation.ZoomOut,
		},
		Animations: mapAnimationsToModel(input.Animations),
		// Preserve existing values for these fields
		Emotes: model.KappagenOverlaySettingsEmotes{
			Time:           currentOverlay.Emotes.Time,
			Max:            currentOverlay.Emotes.Max,
			Queue:          currentOverlay.Emotes.Queue,
			FfzEnabled:     currentOverlay.Emotes.FfzEnabled,
			BttvEnabled:    currentOverlay.Emotes.BttvEnabled,
			SevenTvEnabled: currentOverlay.Emotes.SevenTvEnabled,
			EmojiStyle:     model.KappagenOverlaySettingsEmotesEmojiStyle(currentOverlay.Emotes.EmojiStyle),
		},
		Size: model.KappagenOverlaySettingsSize{
			RatioNormal: currentOverlay.Size.RatioNormal,
			RatioSmall:  currentOverlay.Size.RatioSmall,
			Min:         currentOverlay.Size.Min,
			Max:         currentOverlay.Size.Max,
		},
		Cube: model.KappagenOverlaySettingsCube{
			Speed: currentOverlay.Cube.Speed,
		},
	}

	settingsBytes, err := json.Marshal(settings)
	if err != nil {
		return kappagenmodel.Nil, fmt.Errorf("failed to marshal settings: %w", err)
	}

	query := `
UPDATE channels_modules_settings
SET settings = $1
WHERE "channelId" = $2 AND type = $3
RETURNING id, "channelId"
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	row := conn.QueryRow(ctx, query, settingsBytes, channelID, kappagenOverlayType)

	var id uuid.UUID
	var dbChannelID string
	if err := row.Scan(&id, &dbChannelID); err != nil {
		return kappagenmodel.Nil, fmt.Errorf("failed to update kappagen overlay: %w", err)
	}

	return mapModelToEntity(id, dbChannelID, settings), nil
}

// Helper functions to map between entity and model
func mapModelToEntity(id uuid.UUID, channelID string, model model.KappagenOverlaySettings) kappagenmodel.KappagenOverlay {
	animations := make([]kappagenmodel.KappagenOverlayAnimationsSettings, 0, len(model.Animations))
	for _, a := range model.Animations {
		var prefs *kappagenmodel.KappagenOverlayAnimationsPrefsSettings
		if a.Prefs != nil {
			prefs = &kappagenmodel.KappagenOverlayAnimationsPrefsSettings{
				Size:    a.Prefs.Size,
				Center:  a.Prefs.Center,
				Speed:   a.Prefs.Speed,
				Faces:   a.Prefs.Faces,
				Message: a.Prefs.Message,
				Time:    a.Prefs.Time,
			}
		}
		
		animations = append(animations, kappagenmodel.KappagenOverlayAnimationsSettings{
			Style:   a.Style,
			Prefs:   prefs,
			Count:   a.Count,
			Enabled: a.Enabled,
		})
	}

	return kappagenmodel.KappagenOverlay{
		ID:             id,
		ChannelID:      channelID,
		EnableSpawn:    model.EnableSpawn,
		ExcludedEmotes: model.ExcludedEmotes,
		EnableRave:     model.EnableRave,
		Animation: kappagenmodel.KappagenOverlayAnimationSettings{
			FadeIn:  model.Animation.FadeIn,
			FadeOut: model.Animation.FadeOut,
			ZoomIn:  model.Animation.ZoomIn,
			ZoomOut: model.Animation.ZoomOut,
		},
		Animations: animations,
		Emotes: kappagenmodel.KappagenOverlayEmotesSettings{
			Time:           model.Emotes.Time,
			Max:            model.Emotes.Max,
			Queue:          model.Emotes.Queue,
			FfzEnabled:     model.Emotes.FfzEnabled,
			BttvEnabled:    model.Emotes.BttvEnabled,
			SevenTvEnabled: model.Emotes.SevenTvEnabled,
			EmojiStyle:     kappagenmodel.KappagenEmojiStyle(model.Emotes.EmojiStyle),
		},
		Size: kappagenmodel.KappagenOverlaySizeSettings{
			RatioNormal: model.Size.RatioNormal,
			RatioSmall:  model.Size.RatioSmall,
			Min:         model.Size.Min,
			Max:         model.Size.Max,
		},
		Cube: kappagenmodel.KappagenOverlayCubeSettings{
			Speed: model.Cube.Speed,
		},
	}
}

func mapCreateInputToModel(input overlays_kappagen.CreateInput) model.KappagenOverlaySettings {
	return model.KappagenOverlaySettings{
		EnableSpawn:    input.EnableSpawn,
		ExcludedEmotes: input.ExcludedEmotes,
		EnableRave:     input.EnableRave,
		Animation: model.KappagenOverlaySettingsAnimation{
			FadeIn:  input.Animation.FadeIn,
			FadeOut: input.Animation.FadeOut,
			ZoomIn:  input.Animation.ZoomIn,
			ZoomOut: input.Animation.ZoomOut,
		},
		Animations: mapAnimationsToModel(input.Animations),
		Emotes: model.KappagenOverlaySettingsEmotes{
			Time:           input.Emotes.Time,
			Max:            input.Emotes.Max,
			Queue:          input.Emotes.Queue,
			FfzEnabled:     input.Emotes.FfzEnabled,
			BttvEnabled:    input.Emotes.BttvEnabled,
			SevenTvEnabled: input.Emotes.SevenTvEnabled,
			EmojiStyle:     model.KappagenOverlaySettingsEmotesEmojiStyle(input.Emotes.EmojiStyle),
		},
		Size: model.KappagenOverlaySettingsSize{
			RatioNormal: input.Size.RatioNormal,
			RatioSmall:  input.Size.RatioSmall,
			Min:         input.Size.Min,
			Max:         input.Size.Max,
		},
		Cube: model.KappagenOverlaySettingsCube{
			Speed: input.Cube.Speed,
		},
	}
}

func mapAnimationsToModel(animations []kappagenmodel.KappagenOverlayAnimationsSettings) []model.KappagenOverlaySettingsAnimationSettings {
	result := make([]model.KappagenOverlaySettingsAnimationSettings, 0, len(animations))
	for _, a := range animations {
		var prefs *model.KappagenOverlaySettingsAnimationSettingsPrefs
		if a.Prefs != nil {
			prefs = &model.KappagenOverlaySettingsAnimationSettingsPrefs{
				Size:    a.Prefs.Size,
				Center:  a.Prefs.Center,
				Speed:   a.Prefs.Speed,
				Faces:   a.Prefs.Faces,
				Message: a.Prefs.Message,
				Time:    a.Prefs.Time,
			}
		}
		
		result = append(result, model.KappagenOverlaySettingsAnimationSettings{
			Style:   a.Style,
			Prefs:   prefs,
			Count:   a.Count,
			Enabled: a.Enabled,
		})
	}
	return result
}
