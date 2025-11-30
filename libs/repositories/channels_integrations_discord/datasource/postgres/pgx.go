package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	channelsintegrationsdiscord "github.com/twirapp/twir/libs/repositories/channels_integrations_discord"
	"github.com/twirapp/twir/libs/repositories/channels_integrations_discord/model"
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

var _ channelsintegrationsdiscord.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

var selectColumns = []string{
	"id",
	"channel_id",
	"guild_id",
	"live_notification_enabled",
	"live_notification_channels_ids",
	"live_notification_show_title",
	"live_notification_show_category",
	"live_notification_show_viewers",
	"live_notification_message",
	"live_notification_show_preview",
	"live_notification_show_profile_image",
	"offline_notification_message",
	"should_delete_message_on_offline",
	"additional_users_ids_for_live_check",
}

var selectColumnsStr = strings.Join(selectColumns, ", ")

func (p *Pgx) GetByChannelID(
	ctx context.Context,
	channelID string,
) ([]model.ChannelIntegrationDiscord, error) {
	query, args, err := sq.
		Select(selectColumnsStr).
		From("channels_integrations_discord").
		Where(squirrel.Eq{"channel_id": channelID}).
		OrderBy("id").
		ToSql()
	if err != nil {
		return nil, err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[model.ChannelIntegrationDiscord],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (p *Pgx) GetByChannelIDAndGuildID(
	ctx context.Context,
	channelID, guildID string,
) (model.ChannelIntegrationDiscord, error) {
	query, args, err := sq.
		Select(selectColumnsStr).
		From("channels_integrations_discord").
		Where(
			squirrel.Eq{
				"channel_id": channelID,
				"guild_id":   guildID,
			},
		).
		Limit(1).
		ToSql()
	if err != nil {
		return model.Nil, err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.ChannelIntegrationDiscord],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, nil
		}
		return model.Nil, err
	}

	return result, nil
}

func (p *Pgx) GetByGuildID(
	ctx context.Context,
	guildID string,
) ([]model.ChannelIntegrationDiscord, error) {
	query, args, err := sq.
		Select(selectColumnsStr).
		From("channels_integrations_discord").
		Where(squirrel.Eq{"guild_id": guildID}).
		OrderBy("id").
		ToSql()
	if err != nil {
		return nil, err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[model.ChannelIntegrationDiscord],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (p *Pgx) GetByAdditionalUserID(
	ctx context.Context,
	userID string,
) ([]model.ChannelIntegrationDiscord, error) {
	query, args, err := sq.
		Select(selectColumnsStr).
		From("channels_integrations_discord").
		Where(squirrel.Expr("@user_id = ANY(additional_users_ids_for_live_check)")).
		OrderBy("id").
		ToSql()
	if err != nil {
		return nil, err
	}
	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[model.ChannelIntegrationDiscord],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (p *Pgx) Create(
	ctx context.Context,
	input channelsintegrationsdiscord.CreateInput,
) (model.ChannelIntegrationDiscord, error) {
	qb := sq.
		Insert("channels_integrations_discord").
		SetMap(
			map[string]any{
				"channel_id":                           input.ChannelID,
				"guild_id":                             input.GuildID,
				"live_notification_enabled":            input.LiveNotificationEnabled,
				"live_notification_channels_ids":       input.LiveNotificationChannelsIds,
				"live_notification_show_title":         input.LiveNotificationShowTitle,
				"live_notification_show_category":      input.LiveNotificationShowCategory,
				"live_notification_show_viewers":       input.LiveNotificationShowViewers,
				"live_notification_message":            input.LiveNotificationMessage,
				"live_notification_show_preview":       input.LiveNotificationShowPreview,
				"live_notification_show_profile_image": input.LiveNotificationShowProfileImage,
				"offline_notification_message":         input.OfflineNotificationMessage,
				"should_delete_message_on_offline":     input.ShouldDeleteMessageOnOffline,
				"additional_users_ids_for_live_check":  input.AdditionalUsersIdsForLiveCheck,
			},
		).
		Suffix("RETURNING " + selectColumnsStr)

	query, args, err := qb.ToSql()
	if err != nil {
		return model.Nil, err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.ChannelIntegrationDiscord],
	)
	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (p *Pgx) Update(
	ctx context.Context,
	id int,
	input channelsintegrationsdiscord.UpdateInput,
) error {
	builder := sq.Update("channels_integrations_discord").Where(squirrel.Eq{"id": id})

	if input.LiveNotificationEnabled != nil {
		builder = builder.Set("live_notification_enabled", *input.LiveNotificationEnabled)
	}

	if input.LiveNotificationChannelsIds != nil {
		builder = builder.Set("live_notification_channels_ids", *input.LiveNotificationChannelsIds)
	}

	if input.LiveNotificationShowTitle != nil {
		builder = builder.Set("live_notification_show_title", *input.LiveNotificationShowTitle)
	}

	if input.LiveNotificationShowCategory != nil {
		builder = builder.Set("live_notification_show_category", *input.LiveNotificationShowCategory)
	}

	if input.LiveNotificationShowViewers != nil {
		builder = builder.Set("live_notification_show_viewers", *input.LiveNotificationShowViewers)
	}

	if input.LiveNotificationMessage != nil {
		builder = builder.Set("live_notification_message", *input.LiveNotificationMessage)
	}

	if input.LiveNotificationShowPreview != nil {
		builder = builder.Set("live_notification_show_preview", *input.LiveNotificationShowPreview)
	}

	if input.LiveNotificationShowProfileImage != nil {
		builder = builder.Set(
			"live_notification_show_profile_image",
			*input.LiveNotificationShowProfileImage,
		)
	}

	if input.OfflineNotificationMessage != nil {
		builder = builder.Set("offline_notification_message", *input.OfflineNotificationMessage)
	}

	if input.ShouldDeleteMessageOnOffline != nil {
		builder = builder.Set("should_delete_message_on_offline", *input.ShouldDeleteMessageOnOffline)
	}

	if input.AdditionalUsersIdsForLiveCheck != nil {
		builder = builder.Set(
			"additional_users_ids_for_live_check",
			*input.AdditionalUsersIdsForLiveCheck,
		)
	}

	builder = builder.Set("updated_at", squirrel.Expr("now()"))

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	_, err = conn.Exec(ctx, query, args...)
	return err
}

func (p *Pgx) Delete(
	ctx context.Context,
	id int,
) error {
	query := `DELETE FROM channels_integrations_discord WHERE id = @id`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	_, err := conn.Exec(
		ctx,
		query,
		pgx.NamedArgs{"id": id},
	)
	return err
}

func (p *Pgx) DeleteByChannelIDAndGuildID(
	ctx context.Context,
	channelID, guildID string,
) error {
	query := `DELETE FROM channels_integrations_discord WHERE channel_id = @channel_id AND guild_id = @guild_id`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	_, err := conn.Exec(
		ctx,
		query,
		pgx.NamedArgs{
			"channel_id": channelID,
			"guild_id":   guildID,
		},
	)
	return err
}
