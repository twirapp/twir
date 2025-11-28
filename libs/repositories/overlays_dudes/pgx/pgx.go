package pgx

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	"github.com/twirapp/twir/libs/repositories/overlays_dudes"
	"github.com/twirapp/twir/libs/repositories/overlays_dudes/model"
)

type Opts struct {
	PgxPool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool: opts.PgxPool,
	}
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(Opts{PgxPool: pool})
}

var _ overlays_dudes.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool *pgxpool.Pool
}

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (model.OverlaysDudes, error) {
	query := `
		SELECT
			id, channel_id, created_at,
			dude_color, dude_eyes_color, dude_cosmetics_color, dude_max_life_time,
			dude_gravity, dude_scale, dude_sounds_enabled, dude_sounds_volume,
			dude_visible_name, dude_grow_time, dude_grow_max_scale, dude_max_on_screen, dude_default_sprite,
			message_box_enabled, message_box_border_radius, message_box_box_color, message_box_font_family,
			message_box_font_size, message_box_padding, message_box_show_time, message_box_fill,
			name_box_font_family, name_box_font_size, name_box_fill, name_box_line_join,
			name_box_stroke_thickness, name_box_stroke, name_box_fill_gradient_stops, name_box_fill_gradient_type,
			name_box_font_style, name_box_font_variant, name_box_font_weight, name_box_drop_shadow,
			name_box_drop_shadow_alpha, name_box_drop_shadow_angle, name_box_drop_shadow_blur,
			name_box_drop_shadow_distance, name_box_drop_shadow_color,
			ignore_commands, ignore_users, ignored_users,
			spitter_emote_enabled
		FROM channels_overlays_dudes
		WHERE id = $1
	`

	rows, err := c.pool.Query(ctx, query, id)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(
		rows, func(row pgx.CollectableRow) (model.OverlaysDudes, error) {
			var entity model.OverlaysDudes
			var nameBoxFill []string
			var nameBoxFillGradientStops pq.Float32Array
			var ignoredUsers []string

			err := row.Scan(
				&entity.ID,
				&entity.ChannelID,
				&entity.CreatedAt,
				&entity.DudeColor,
				&entity.DudeEyesColor,
				&entity.DudeCosmeticsColor,
				&entity.DudeMaxLifeTime,
				&entity.DudeGravity,
				&entity.DudeScale,
				&entity.DudeSoundsEnabled,
				&entity.DudeSoundsVolume,
				&entity.DudeVisibleName,
				&entity.DudeGrowTime,
				&entity.DudeGrowMaxScale,
				&entity.DudeMaxOnScreen,
				&entity.DudeDefaultSprite,
				&entity.MessageBoxEnabled,
				&entity.MessageBoxBorderRadius,
				&entity.MessageBoxBoxColor,
				&entity.MessageBoxFontFamily,
				&entity.MessageBoxFontSize,
				&entity.MessageBoxPadding,
				&entity.MessageBoxShowTime,
				&entity.MessageBoxFill,
				&entity.NameBoxFontFamily,
				&entity.NameBoxFontSize,
				&nameBoxFill,
				&entity.NameBoxLineJoin,
				&entity.NameBoxStrokeThickness,
				&entity.NameBoxStroke,
				&nameBoxFillGradientStops,
				&entity.NameBoxFillGradientType,
				&entity.NameBoxFontStyle,
				&entity.NameBoxFontVariant,
				&entity.NameBoxFontWeight,
				&entity.NameBoxDropShadow,
				&entity.NameBoxDropShadowAlpha,
				&entity.NameBoxDropShadowAngle,
				&entity.NameBoxDropShadowBlur,
				&entity.NameBoxDropShadowDistance,
				&entity.NameBoxDropShadowColor,
				&entity.IgnoreCommands,
				&entity.IgnoreUsers,
				&ignoredUsers,
				&entity.SpitterEmoteEnabled,
			)
			if err != nil {
				return model.Nil, err
			}

			entity.NameBoxFill = []string(nameBoxFill)
			entity.NameBoxFillGradientStops = []float32(nameBoxFillGradientStops)
			entity.IgnoredUsers = []string(ignoredUsers)

			return entity, nil
		},
	)

	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) GetManyByChannelID(ctx context.Context, channelID string) (
	[]model.OverlaysDudes,
	error,
) {
	query := `
		SELECT
			id, channel_id, created_at,
			dude_color, dude_eyes_color, dude_cosmetics_color, dude_max_life_time,
			dude_gravity, dude_scale, dude_sounds_enabled, dude_sounds_volume,
			dude_visible_name, dude_grow_time, dude_grow_max_scale, dude_max_on_screen, dude_default_sprite,
			message_box_enabled, message_box_border_radius, message_box_box_color, message_box_font_family,
			message_box_font_size, message_box_padding, message_box_show_time, message_box_fill,
			name_box_font_family, name_box_font_size, name_box_fill, name_box_line_join,
			name_box_stroke_thickness, name_box_stroke, name_box_fill_gradient_stops, name_box_fill_gradient_type,
			name_box_font_style, name_box_font_variant, name_box_font_weight, name_box_drop_shadow,
			name_box_drop_shadow_alpha, name_box_drop_shadow_angle, name_box_drop_shadow_blur,
			name_box_drop_shadow_distance, name_box_drop_shadow_color,
			ignore_commands, ignore_users, ignored_users,
			spitter_emote_enabled
		FROM channels_overlays_dudes
		WHERE channel_id = $1
		ORDER BY created_at DESC
	`

	rows, err := c.pool.Query(ctx, query, channelID)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(
		rows, func(row pgx.CollectableRow) (model.OverlaysDudes, error) {
			var entity model.OverlaysDudes
			var nameBoxFill []string
			var nameBoxFillGradientStops pq.Float32Array
			var ignoredUsers []string

			err := row.Scan(
				&entity.ID,
				&entity.ChannelID,
				&entity.CreatedAt,
				&entity.DudeColor,
				&entity.DudeEyesColor,
				&entity.DudeCosmeticsColor,
				&entity.DudeMaxLifeTime,
				&entity.DudeGravity,
				&entity.DudeScale,
				&entity.DudeSoundsEnabled,
				&entity.DudeSoundsVolume,
				&entity.DudeVisibleName,
				&entity.DudeGrowTime,
				&entity.DudeGrowMaxScale,
				&entity.DudeMaxOnScreen,
				&entity.DudeDefaultSprite,
				&entity.MessageBoxEnabled,
				&entity.MessageBoxBorderRadius,
				&entity.MessageBoxBoxColor,
				&entity.MessageBoxFontFamily,
				&entity.MessageBoxFontSize,
				&entity.MessageBoxPadding,
				&entity.MessageBoxShowTime,
				&entity.MessageBoxFill,
				&entity.NameBoxFontFamily,
				&entity.NameBoxFontSize,
				&nameBoxFill,
				&entity.NameBoxLineJoin,
				&entity.NameBoxStrokeThickness,
				&entity.NameBoxStroke,
				&nameBoxFillGradientStops,
				&entity.NameBoxFillGradientType,
				&entity.NameBoxFontStyle,
				&entity.NameBoxFontVariant,
				&entity.NameBoxFontWeight,
				&entity.NameBoxDropShadow,
				&entity.NameBoxDropShadowAlpha,
				&entity.NameBoxDropShadowAngle,
				&entity.NameBoxDropShadowBlur,
				&entity.NameBoxDropShadowDistance,
				&entity.NameBoxDropShadowColor,
				&entity.IgnoreCommands,
				&entity.IgnoreUsers,
				&ignoredUsers,
				&entity.SpitterEmoteEnabled,
			)
			if err != nil {
				return model.Nil, err
			}

			entity.NameBoxFill = []string(nameBoxFill)
			entity.NameBoxFillGradientStops = []float32(nameBoxFillGradientStops)
			entity.IgnoredUsers = []string(ignoredUsers)

			return entity, nil
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Pgx) Create(ctx context.Context, input overlays_dudes.CreateInput) (
	model.OverlaysDudes,
	error,
) {
	id := uuid.New()

	query := `
		INSERT INTO channels_overlays_dudes (
			id, channel_id,
			dude_color, dude_eyes_color, dude_cosmetics_color, dude_max_life_time,
			dude_gravity, dude_scale, dude_sounds_enabled, dude_sounds_volume,
			dude_visible_name, dude_grow_time, dude_grow_max_scale, dude_max_on_screen, dude_default_sprite,
			message_box_enabled, message_box_border_radius, message_box_box_color, message_box_font_family,
			message_box_font_size, message_box_padding, message_box_show_time, message_box_fill,
			name_box_font_family, name_box_font_size, name_box_fill, name_box_line_join,
			name_box_stroke_thickness, name_box_stroke, name_box_fill_gradient_stops, name_box_fill_gradient_type,
			name_box_font_style, name_box_font_variant, name_box_font_weight, name_box_drop_shadow,
			name_box_drop_shadow_alpha, name_box_drop_shadow_angle, name_box_drop_shadow_blur,
			name_box_drop_shadow_distance, name_box_drop_shadow_color,
			ignore_commands, ignore_users, ignored_users,
			spitter_emote_enabled
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23,
			$24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44
		)
		RETURNING
			id, channel_id, created_at,
			dude_color, dude_eyes_color, dude_cosmetics_color, dude_max_life_time,
			dude_gravity, dude_scale, dude_sounds_enabled, dude_sounds_volume,
			dude_visible_name, dude_grow_time, dude_grow_max_scale, dude_max_on_screen, dude_default_sprite,
			message_box_enabled, message_box_border_radius, message_box_box_color, message_box_font_family,
			message_box_font_size, message_box_padding, message_box_show_time, message_box_fill,
			name_box_font_family, name_box_font_size, name_box_fill, name_box_line_join,
			name_box_stroke_thickness, name_box_stroke, name_box_fill_gradient_stops, name_box_fill_gradient_type,
			name_box_font_style, name_box_font_variant, name_box_font_weight, name_box_drop_shadow,
			name_box_drop_shadow_alpha, name_box_drop_shadow_angle, name_box_drop_shadow_blur,
			name_box_drop_shadow_distance, name_box_drop_shadow_color,
			ignore_commands, ignore_users, ignored_users,
			spitter_emote_enabled
	`

	rows, err := c.pool.Query(
		ctx,
		query,
		id,
		input.ChannelID,
		input.DudeColor,
		input.DudeEyesColor,
		input.DudeCosmeticsColor,
		input.DudeMaxLifeTime,
		input.DudeGravity,
		input.DudeScale,
		input.DudeSoundsEnabled,
		input.DudeSoundsVolume,
		input.DudeVisibleName,
		input.DudeGrowTime,
		input.DudeGrowMaxScale,
		input.DudeMaxOnScreen,
		input.DudeDefaultSprite,
		input.MessageBoxEnabled,
		input.MessageBoxBorderRadius,
		input.MessageBoxBoxColor,
		input.MessageBoxFontFamily,
		input.MessageBoxFontSize,
		input.MessageBoxPadding,
		input.MessageBoxShowTime,
		input.MessageBoxFill,
		input.NameBoxFontFamily,
		input.NameBoxFontSize,
		(input.NameBoxFill),
		input.NameBoxLineJoin,
		input.NameBoxStrokeThickness,
		input.NameBoxStroke,
		pq.Float32Array(input.NameBoxFillGradientStops),
		input.NameBoxFillGradientType,
		input.NameBoxFontStyle,
		input.NameBoxFontVariant,
		input.NameBoxFontWeight,
		input.NameBoxDropShadow,
		input.NameBoxDropShadowAlpha,
		input.NameBoxDropShadowAngle,
		input.NameBoxDropShadowBlur,
		input.NameBoxDropShadowDistance,
		input.NameBoxDropShadowColor,
		input.IgnoreCommands,
		input.IgnoreUsers,
		input.IgnoredUsers,
		input.SpitterEmoteEnabled,
	)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(
		rows, func(row pgx.CollectableRow) (model.OverlaysDudes, error) {
			var entity model.OverlaysDudes
			var nameBoxFill []string
			var nameBoxFillGradientStops pq.Float32Array
			var ignoredUsers []string

			err := row.Scan(
				&entity.ID,
				&entity.ChannelID,
				&entity.CreatedAt,
				&entity.DudeColor,
				&entity.DudeEyesColor,
				&entity.DudeCosmeticsColor,
				&entity.DudeMaxLifeTime,
				&entity.DudeGravity,
				&entity.DudeScale,
				&entity.DudeSoundsEnabled,
				&entity.DudeSoundsVolume,
				&entity.DudeVisibleName,
				&entity.DudeGrowTime,
				&entity.DudeGrowMaxScale,
				&entity.DudeMaxOnScreen,
				&entity.DudeDefaultSprite,
				&entity.MessageBoxEnabled,
				&entity.MessageBoxBorderRadius,
				&entity.MessageBoxBoxColor,
				&entity.MessageBoxFontFamily,
				&entity.MessageBoxFontSize,
				&entity.MessageBoxPadding,
				&entity.MessageBoxShowTime,
				&entity.MessageBoxFill,
				&entity.NameBoxFontFamily,
				&entity.NameBoxFontSize,
				&nameBoxFill,
				&entity.NameBoxLineJoin,
				&entity.NameBoxStrokeThickness,
				&entity.NameBoxStroke,
				&nameBoxFillGradientStops,
				&entity.NameBoxFillGradientType,
				&entity.NameBoxFontStyle,
				&entity.NameBoxFontVariant,
				&entity.NameBoxFontWeight,
				&entity.NameBoxDropShadow,
				&entity.NameBoxDropShadowAlpha,
				&entity.NameBoxDropShadowAngle,
				&entity.NameBoxDropShadowBlur,
				&entity.NameBoxDropShadowDistance,
				&entity.NameBoxDropShadowColor,
				&entity.IgnoreCommands,
				&entity.IgnoreUsers,
				&ignoredUsers,
				&entity.SpitterEmoteEnabled,
			)
			if err != nil {
				return model.Nil, err
			}

			entity.NameBoxFill = []string(nameBoxFill)
			entity.NameBoxFillGradientStops = []float32(nameBoxFillGradientStops)
			entity.IgnoredUsers = []string(ignoredUsers)

			return entity, nil
		},
	)

	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) Update(
	ctx context.Context,
	id uuid.UUID,
	input overlays_dudes.UpdateInput,
) (model.OverlaysDudes, error) {
	updateBuilder := sq.Update("channels_overlays_dudes").Where(squirrel.Eq{"id": id})

	if input.DudeColor != nil {
		updateBuilder = updateBuilder.Set("dude_color", *input.DudeColor)
	}
	if input.DudeEyesColor != nil {
		updateBuilder = updateBuilder.Set("dude_eyes_color", *input.DudeEyesColor)
	}
	if input.DudeCosmeticsColor != nil {
		updateBuilder = updateBuilder.Set("dude_cosmetics_color", *input.DudeCosmeticsColor)
	}
	if input.DudeMaxLifeTime != nil {
		updateBuilder = updateBuilder.Set("dude_max_life_time", *input.DudeMaxLifeTime)
	}
	if input.DudeGravity != nil {
		updateBuilder = updateBuilder.Set("dude_gravity", *input.DudeGravity)
	}
	if input.DudeScale != nil {
		updateBuilder = updateBuilder.Set("dude_scale", *input.DudeScale)
	}
	if input.DudeSoundsEnabled != nil {
		updateBuilder = updateBuilder.Set("dude_sounds_enabled", *input.DudeSoundsEnabled)
	}
	if input.DudeSoundsVolume != nil {
		updateBuilder = updateBuilder.Set("dude_sounds_volume", *input.DudeSoundsVolume)
	}
	if input.DudeVisibleName != nil {
		updateBuilder = updateBuilder.Set("dude_visible_name", *input.DudeVisibleName)
	}
	if input.DudeGrowTime != nil {
		updateBuilder = updateBuilder.Set("dude_grow_time", *input.DudeGrowTime)
	}
	if input.DudeGrowMaxScale != nil {
		updateBuilder = updateBuilder.Set("dude_grow_max_scale", *input.DudeGrowMaxScale)
	}
	if input.DudeMaxOnScreen != nil {
		updateBuilder = updateBuilder.Set("dude_max_on_screen", *input.DudeMaxOnScreen)
	}
	if input.DudeDefaultSprite != nil {
		updateBuilder = updateBuilder.Set("dude_default_sprite", *input.DudeDefaultSprite)
	}

	// Message box settings
	if input.MessageBoxEnabled != nil {
		updateBuilder = updateBuilder.Set("message_box_enabled", *input.MessageBoxEnabled)
	}
	if input.MessageBoxBorderRadius != nil {
		updateBuilder = updateBuilder.Set("message_box_border_radius", *input.MessageBoxBorderRadius)
	}
	if input.MessageBoxBoxColor != nil {
		updateBuilder = updateBuilder.Set("message_box_box_color", *input.MessageBoxBoxColor)
	}
	if input.MessageBoxFontFamily != nil {
		updateBuilder = updateBuilder.Set("message_box_font_family", *input.MessageBoxFontFamily)
	}
	if input.MessageBoxFontSize != nil {
		updateBuilder = updateBuilder.Set("message_box_font_size", *input.MessageBoxFontSize)
	}
	if input.MessageBoxPadding != nil {
		updateBuilder = updateBuilder.Set("message_box_padding", *input.MessageBoxPadding)
	}
	if input.MessageBoxShowTime != nil {
		updateBuilder = updateBuilder.Set("message_box_show_time", *input.MessageBoxShowTime)
	}
	if input.MessageBoxFill != nil {
		updateBuilder = updateBuilder.Set("message_box_fill", *input.MessageBoxFill)
	}

	// Name box settings
	if input.NameBoxFontFamily != nil {
		updateBuilder = updateBuilder.Set("name_box_font_family", *input.NameBoxFontFamily)
	}
	if input.NameBoxFontSize != nil {
		updateBuilder = updateBuilder.Set("name_box_font_size", *input.NameBoxFontSize)
	}
	if input.NameBoxFill != nil {
		updateBuilder = updateBuilder.Set("name_box_fill", *input.NameBoxFill)
	}
	if input.NameBoxLineJoin != nil {
		updateBuilder = updateBuilder.Set("name_box_line_join", *input.NameBoxLineJoin)
	}
	if input.NameBoxStrokeThickness != nil {
		updateBuilder = updateBuilder.Set("name_box_stroke_thickness", *input.NameBoxStrokeThickness)
	}
	if input.NameBoxStroke != nil {
		updateBuilder = updateBuilder.Set("name_box_stroke", *input.NameBoxStroke)
	}
	if input.NameBoxFillGradientStops != nil {
		updateBuilder = updateBuilder.Set(
			"name_box_fill_gradient_stops",
			pq.Float32Array(*input.NameBoxFillGradientStops),
		)
	}
	if input.NameBoxFillGradientType != nil {
		updateBuilder = updateBuilder.Set("name_box_fill_gradient_type", *input.NameBoxFillGradientType)
	}
	if input.NameBoxFontStyle != nil {
		updateBuilder = updateBuilder.Set("name_box_font_style", *input.NameBoxFontStyle)
	}
	if input.NameBoxFontVariant != nil {
		updateBuilder = updateBuilder.Set("name_box_font_variant", *input.NameBoxFontVariant)
	}
	if input.NameBoxFontWeight != nil {
		updateBuilder = updateBuilder.Set("name_box_font_weight", *input.NameBoxFontWeight)
	}
	if input.NameBoxDropShadow != nil {
		updateBuilder = updateBuilder.Set("name_box_drop_shadow", *input.NameBoxDropShadow)
	}
	if input.NameBoxDropShadowAlpha != nil {
		updateBuilder = updateBuilder.Set("name_box_drop_shadow_alpha", *input.NameBoxDropShadowAlpha)
	}
	if input.NameBoxDropShadowAngle != nil {
		updateBuilder = updateBuilder.Set("name_box_drop_shadow_angle", *input.NameBoxDropShadowAngle)
	}
	if input.NameBoxDropShadowBlur != nil {
		updateBuilder = updateBuilder.Set("name_box_drop_shadow_blur", *input.NameBoxDropShadowBlur)
	}
	if input.NameBoxDropShadowDistance != nil {
		updateBuilder = updateBuilder.Set(
			"name_box_drop_shadow_distance",
			*input.NameBoxDropShadowDistance,
		)
	}
	if input.NameBoxDropShadowColor != nil {
		updateBuilder = updateBuilder.Set("name_box_drop_shadow_color", *input.NameBoxDropShadowColor)
	}

	// Ignore settings
	if input.IgnoreCommands != nil {
		updateBuilder = updateBuilder.Set("ignore_commands", *input.IgnoreCommands)
	}
	if input.IgnoreUsers != nil {
		updateBuilder = updateBuilder.Set("ignore_users", *input.IgnoreUsers)
	}
	if input.IgnoredUsers != nil {
		updateBuilder = updateBuilder.Set("ignored_users", *input.IgnoredUsers)
	}

	// Spitter emote settings
	if input.SpitterEmoteEnabled != nil {
		updateBuilder = updateBuilder.Set("spitter_emote_enabled", *input.SpitterEmoteEnabled)
	}

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return model.Nil, err
	}

	_, err = c.pool.Exec(ctx, query, args...)
	if err != nil {
		return model.Nil, err
	}

	// Return the updated entity
	return c.GetByID(ctx, id)
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM channels_overlays_dudes WHERE id = $1"
	_, err := c.pool.Exec(ctx, query, id)
	return err
}
