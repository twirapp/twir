package pgx

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
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
		SELECT
			co.id AS overlay_id,
			co.channel_id,
			co.enable_spawn,
			co.excluded_emotes,
			co.enable_rave,
			co.created_at AS overlay_created_at,
			co.updated_at AS overlay_updated_at,
			ca.id AS animation_id,
			ca.overlay_id,
			ca.style,
			ca.count,
			ca.enabled,
			ca.created_at AS animation_created_at,
			ca.updated_at AS animation_updated_at,
			cap.id AS prefs_id,
			cap.animation_id AS prefs_animation_id,
			cap.size,
			cap.center,
			cap.speed,
			cap.faces,
			cap.message,
			cap.time,
			cap.created_at AS prefs_created_at,
			cap.updated_at AS prefs_updated_at
		FROM channels_overlays_kappagen co
		LEFT JOIN channels_overlays_kappagen_animations ca ON co.id = ca.overlay_id
		LEFT JOIN channels_overlays_kappagen_animations_prefs cap ON ca.id = cap.animation_id
		WHERE co.channel_id = $1
		ORDER BY ca.created_at, cap.created_at
	`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return kappagenmodel.KappagenOverlay{}, err
	}
	defer rows.Close()

	var overlay *kappagenmodel.KappagenOverlay
	animationsMap := make(map[string]*kappagenmodel.KappagenOverlayAnimationsSettings)

	for rows.Next() {
		var (
			overlayID         string
			selectedChannelId string
			enableSpawn       bool
			excludedEmotes    []string
			enableRave        bool
			overlayCreatedAt  time.Time
			overlayUpdatedAt  time.Time
			animationID       *string
			animOverlayID     *string
			animStyle         *string
			animCount         *int
			animEnabled       *bool
			animCreatedAt     *time.Time
			animUpdatedAt     *time.Time
			prefsID           *string
			prefsAnimationID  *string
			prefsSize         *float64
			prefsCenter       *bool
			prefsSpeed        *int
			prefsFaces        *bool
			prefsMessage      *[]string
			prefsTime         *int
			prefsCreatedAt    *time.Time
			prefsUpdatedAt    *time.Time
		)

		err := rows.Scan(
			&overlayID,
			&selectedChannelId,
			&enableSpawn,
			&excludedEmotes,
			&enableRave,
			&overlayCreatedAt,
			&overlayUpdatedAt,
			&animationID,
			&animOverlayID,
			&animStyle,
			&animCount,
			&animEnabled,
			&animCreatedAt,
			&animUpdatedAt,
			&prefsID,
			&prefsAnimationID,
			&prefsSize,
			&prefsCenter,
			&prefsSpeed,
			&prefsFaces,
			&prefsMessage,
			&prefsTime,
			&prefsCreatedAt,
			&prefsUpdatedAt,
		)
		if err != nil {
			return kappagenmodel.KappagenOverlay{}, err
		}

		parsedOverlayID, err := uuid.Parse(overlayID)
		if err != nil {
			return kappagenmodel.KappagenOverlay{}, err
		}

		if overlay == nil {
			overlay = &kappagenmodel.KappagenOverlay{
				ID:             parsedOverlayID,
				ChannelID:      channelID,
				EnableSpawn:    enableSpawn,
				ExcludedEmotes: excludedEmotes,
				EnableRave:     enableRave,
				CreatedAt:      overlayCreatedAt,
				UpdatedAt:      overlayUpdatedAt,
			}
		}

		if animationID != nil {
			parsedAnimationID, err := ulid.Parse(*animationID)
			if err != nil {
				return kappagenmodel.KappagenOverlay{}, err
			}

			if anim, exists := animationsMap[*animationID]; !exists {
				anim = &kappagenmodel.KappagenOverlayAnimationsSettings{
					ID:        parsedAnimationID,
					OverlayID: parsedOverlayID,
					Style:     *animStyle,
					Count:     *animCount,
					Enabled:   *animEnabled,
					CreatedAt: *animCreatedAt,
					UpdatedAt: *animUpdatedAt,
				}
				animationsMap[*animationID] = anim
			}

			if prefsID != nil {
				parsedPrefsID, err := ulid.Parse(*prefsID)
				if err != nil {
					return kappagenmodel.KappagenOverlay{}, err
				}

				prefs := kappagenmodel.KappagenOverlayAnimationsPrefsSettings{
					ID:          parsedPrefsID,
					AnimationID: parsedAnimationID,
					Size:        *prefsSize,
					Center:      *prefsCenter,
					Speed:       *prefsSpeed,
					Faces:       *prefsFaces,
					Message:     *prefsMessage,
					Time:        *prefsTime,
					CreatedAt:   *prefsCreatedAt,
					UpdatedAt:   *prefsUpdatedAt,
				}

				animationsMap[*animationID].Prefs = prefs
			}
		}
	}

	if rows.Err() != nil {
		return kappagenmodel.KappagenOverlay{}, rows.Err()
	}

	if overlay == nil {
		return kappagenmodel.KappagenOverlay{}, fmt.Errorf("not found")
	}

	for _, anim := range animationsMap {
		overlay.Animations = append(overlay.Animations, *anim)
	}

	slices.SortFunc(
		overlay.Animations,
		func(a, b kappagenmodel.KappagenOverlayAnimationsSettings) int {
			return a.ID.Compare(b.ID)
		},
	)

	return *overlay, nil
}

func (p *Pgx) Create(
	ctx context.Context,
	input overlays_kappagen.CreateInput,
) (kappagenmodel.KappagenOverlay, error) {
	panic("")
}

func (p *Pgx) Update(
	ctx context.Context,
	channelID string,
	input overlays_kappagen.UpdateInput,
) (kappagenmodel.KappagenOverlay, error) {
	panic("")
}
