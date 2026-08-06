package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/entities/song_request_mode"
	entity "github.com/twirapp/twir/libs/entities/song_requests_settings"
	songrequestssettings "github.com/twirapp/twir/libs/repositories/song_requests_settings"
)

type Opts struct {
	PgxPool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{pool: opts.PgxPool}
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(Opts{PgxPool: pool})
}

var _ songrequestssettings.Repository = (*Pgx)(nil)

type Pgx struct {
	pool *pgxpool.Pool
}

type scanModel struct {
	ID                                   uuid.UUID              `db:"id"`
	ChannelID                            *uuid.UUID             `db:"channel_id"`
	Mode                                 song_request_mode.Mode `db:"mode"`
	Enabled                              bool                   `db:"enabled"`
	AcceptOnlyWhenOnline                 bool                   `db:"accept_only_when_online"`
	PlayerNoCookieMode                   bool                   `db:"player_no_cookie_mode"`
	TakeSongFromDonationMessage          bool                   `db:"take_song_from_donation_message"`
	MaxRequests                          *int                   `db:"max_requests"`
	ChannelPointsRewardID                *string                `db:"channel_points_reward_id"`
	AnnouncePlay                         bool                   `db:"announce_play"`
	NeededVotesForSkip                   *float64               `db:"needed_votes_for_skip"`
	UserMaxRequests                      *int                   `db:"user_max_requests"`
	UserMinWatchTime                     *int64                 `db:"user_min_watch_time"`
	UserMinMessages                      *int                   `db:"user_min_messages"`
	UserMinFollowTime                    *int                   `db:"user_min_follow_time"`
	SongMinLength                        *int                   `db:"song_min_length"`
	SongMaxLength                        *int                   `db:"song_max_length"`
	SongMinViews                         *int                   `db:"song_min_views"`
	SongAcceptedCategories               []string               `db:"song_accepted_categories"`
	SongWordsDenyList                    []string               `db:"song_words_deny_list"`
	DenyListUsers                        []string               `db:"deny_list_users"`
	DenyListSongs                        []string               `db:"deny_list_songs"`
	DenyListChannels                     []string               `db:"deny_list_channels"`
	DenyListArtistsNames                 []string               `db:"deny_list_artists_names"`
	DenyListWords                        []string               `db:"deny_list_words"`
	TranslationsNowPlaying               *string                `db:"translations_now_playing"`
	TranslationsNotEnabled               *string                `db:"translations_not_enabled"`
	TranslationsNoText                   *string                `db:"translations_no_text"`
	TranslationsAcceptOnlineWhenOnline   *string                `db:"translations_accept_online_when_online"`
	TranslationsUserDenied               *string                `db:"translations_user_denied"`
	TranslationsUserMaxRequests          *string                `db:"translations_user_max_requests"`
	TranslationsUserMinMessages          *string                `db:"translations_user_min_messages"`
	TranslationsUserMinWatched           *string                `db:"translations_user_min_watched"`
	TranslationsUserMinFollow            *string                `db:"translations_user_min_follow"`
	TranslationsSongDenied               *string                `db:"translations_song_denied"`
	TranslationsSongNotFound             *string                `db:"translations_song_not_found"`
	TranslationsSongAlreadyInQueue       *string                `db:"translations_song_already_in_queue"`
	TranslationsSongAgeRestrictions      *string                `db:"translations_song_age_restrictions"`
	TranslationsSongCannotGetInformation *string                `db:"translations_song_cannot_get_information"`
	TranslationsSongLive                 *string                `db:"translations_song_live"`
	TranslationsSongMaxLength            *string                `db:"translations_song_max_length"`
	TranslationsSongMinLength            *string                `db:"translations_song_min_length"`
	TranslationsSongRequestedMessage     *string                `db:"translations_song_requested_message"`
	TranslationsSongMaximumOrdered       *string                `db:"translations_song_maximum_ordered"`
	TranslationsSongMinViews             *string                `db:"translations_song_min_views"`
	TranslationsChannelDenied            *string                `db:"translations_channel_denied"`
	Volume                               int                    `db:"volume"`
}

func (s scanModel) toEntity() entity.Settings {
	return entity.Settings{
		ID:                                   s.ID,
		ChannelID:                            lo.FromPtr(s.ChannelID),
		Mode:                                 s.Mode,
		Enabled:                              s.Enabled,
		AcceptOnlyWhenOnline:                 s.AcceptOnlyWhenOnline,
		PlayerNoCookieMode:                   s.PlayerNoCookieMode,
		TakeSongFromDonationMessage:          s.TakeSongFromDonationMessage,
		MaxRequests:                          lo.FromPtr(s.MaxRequests),
		ChannelPointsRewardID:                s.ChannelPointsRewardID,
		AnnouncePlay:                         s.AnnouncePlay,
		NeededVotesForSkip:                   lo.FromPtr(s.NeededVotesForSkip),
		UserMaxRequests:                      lo.FromPtr(s.UserMaxRequests),
		UserMinWatchTime:                     lo.FromPtr(s.UserMinWatchTime),
		UserMinMessages:                      lo.FromPtr(s.UserMinMessages),
		UserMinFollowTime:                    lo.FromPtr(s.UserMinFollowTime),
		SongMinLength:                        lo.FromPtr(s.SongMinLength),
		SongMaxLength:                        lo.FromPtr(s.SongMaxLength),
		SongMinViews:                         lo.FromPtr(s.SongMinViews),
		SongAcceptedCategories:               s.SongAcceptedCategories,
		SongWordsDenyList:                    s.SongWordsDenyList,
		DenyListUsers:                        s.DenyListUsers,
		DenyListSongs:                        s.DenyListSongs,
		DenyListChannels:                     s.DenyListChannels,
		DenyListArtistsNames:                 s.DenyListArtistsNames,
		DenyListWords:                        s.DenyListWords,
		TranslationsNowPlaying:               lo.FromPtr(s.TranslationsNowPlaying),
		TranslationsNotEnabled:               lo.FromPtr(s.TranslationsNotEnabled),
		TranslationsNoText:                   lo.FromPtr(s.TranslationsNoText),
		TranslationsAcceptOnlineWhenOnline:   lo.FromPtr(s.TranslationsAcceptOnlineWhenOnline),
		TranslationsUserDenied:               lo.FromPtr(s.TranslationsUserDenied),
		TranslationsUserMaxRequests:          lo.FromPtr(s.TranslationsUserMaxRequests),
		TranslationsUserMinMessages:          lo.FromPtr(s.TranslationsUserMinMessages),
		TranslationsUserMinWatched:           lo.FromPtr(s.TranslationsUserMinWatched),
		TranslationsUserMinFollow:            lo.FromPtr(s.TranslationsUserMinFollow),
		TranslationsSongDenied:               lo.FromPtr(s.TranslationsSongDenied),
		TranslationsSongNotFound:             lo.FromPtr(s.TranslationsSongNotFound),
		TranslationsSongAlreadyInQueue:       lo.FromPtr(s.TranslationsSongAlreadyInQueue),
		TranslationsSongAgeRestrictions:      lo.FromPtr(s.TranslationsSongAgeRestrictions),
		TranslationsSongCannotGetInformation: lo.FromPtr(s.TranslationsSongCannotGetInformation),
		TranslationsSongLive:                 lo.FromPtr(s.TranslationsSongLive),
		TranslationsSongMaxLength:            lo.FromPtr(s.TranslationsSongMaxLength),
		TranslationsSongMinLength:            lo.FromPtr(s.TranslationsSongMinLength),
		TranslationsSongRequestedMessage:     lo.FromPtr(s.TranslationsSongRequestedMessage),
		TranslationsSongMaximumOrdered:       lo.FromPtr(s.TranslationsSongMaximumOrdered),
		TranslationsSongMinViews:             lo.FromPtr(s.TranslationsSongMinViews),
		TranslationsChannelDenied:            lo.FromPtr(s.TranslationsChannelDenied),
		Volume:                               s.Volume,
	}
}

const selectFields = `id,
	channel_id,
	mode,
	enabled,
	accept_only_when_online,
	player_no_cookie_mode,
	take_song_from_donation_message,
	max_requests,
	channel_points_reward_id,
	announce_play,
	needed_votes_for_skip,
	user_max_requests,
	user_min_watch_time,
	user_min_messages,
	user_min_follow_time,
	song_min_length,
	song_max_length,
	song_min_views,
	song_accepted_categories,
	song_words_deny_list,
	deny_list_users,
	deny_list_songs,
	deny_list_channels,
	deny_list_artists_names,
	deny_list_words,
	translations_now_playing,
	translations_not_enabled,
	translations_no_text,
	translations_accept_online_when_online,
	translations_user_denied,
	translations_user_max_requests,
	translations_user_min_messages,
	translations_user_min_watched,
	translations_user_min_follow,
	translations_song_denied,
	translations_song_not_found,
	translations_song_already_in_queue,
	translations_song_age_restrictions,
	translations_song_cannot_get_information,
	translations_song_live,
	translations_song_max_length,
	translations_song_min_length,
	translations_song_requested_message,
	translations_song_maximum_ordered,
	translations_song_min_views,
	translations_channel_denied,
	volume`

func (c *Pgx) GetByChannelID(
	ctx context.Context,
	channelID string,
) (entity.Settings, error) {
	query := `SELECT ` + selectFields + `
FROM channels_song_requests_settings
WHERE channel_id = CAST(@channel_id AS uuid)
LIMIT 1`

	rows, err := c.pool.Query(
		ctx,
		query,
		pgx.NamedArgs{"channel_id": channelID},
	)
	if err != nil {
		return entity.Nil, fmt.Errorf("query song requests settings: %w", err)
	}

	settings, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Nil, songrequestssettings.ErrNotFound
		}

		return entity.Nil, fmt.Errorf("collect song requests settings: %w", err)
	}

	return settings.toEntity(), nil
}

func (c *Pgx) Upsert(
	ctx context.Context,
	settings entity.Settings,
) (entity.Settings, error) {
	query := `INSERT INTO channels_song_requests_settings (
	channel_id,
	mode,
	enabled,
	accept_only_when_online,
	player_no_cookie_mode,
	take_song_from_donation_message,
	max_requests,
	channel_points_reward_id,
	announce_play,
	needed_votes_for_skip,
	user_max_requests,
	user_min_watch_time,
	user_min_messages,
	user_min_follow_time,
	song_min_length,
	song_max_length,
	song_min_views,
	song_accepted_categories,
	song_words_deny_list,
	deny_list_users,
	deny_list_songs,
	deny_list_channels,
	deny_list_artists_names,
	deny_list_words,
	translations_now_playing,
	translations_not_enabled,
	translations_no_text,
	translations_accept_online_when_online,
	translations_user_denied,
	translations_user_max_requests,
	translations_user_min_messages,
	translations_user_min_watched,
	translations_user_min_follow,
	translations_song_denied,
	translations_song_not_found,
	translations_song_already_in_queue,
	translations_song_age_restrictions,
	translations_song_cannot_get_information,
	translations_song_live,
	translations_song_max_length,
	translations_song_min_length,
	translations_song_requested_message,
	translations_song_maximum_ordered,
	translations_song_min_views,
	translations_channel_denied,
	volume
) VALUES (
	@channel_id,
	@mode,
	@enabled,
	@accept_only_when_online,
	@player_no_cookie_mode,
	@take_song_from_donation_message,
	@max_requests,
	@channel_points_reward_id,
	@announce_play,
	@needed_votes_for_skip,
	@user_max_requests,
	@user_min_watch_time,
	@user_min_messages,
	@user_min_follow_time,
	@song_min_length,
	@song_max_length,
	@song_min_views,
	@song_accepted_categories,
	@song_words_deny_list,
	@deny_list_users,
	@deny_list_songs,
	@deny_list_channels,
	@deny_list_artists_names,
	@deny_list_words,
	@translations_now_playing,
	@translations_not_enabled,
	@translations_no_text,
	@translations_accept_online_when_online,
	@translations_user_denied,
	@translations_user_max_requests,
	@translations_user_min_messages,
	@translations_user_min_watched,
	@translations_user_min_follow,
	@translations_song_denied,
	@translations_song_not_found,
	@translations_song_already_in_queue,
	@translations_song_age_restrictions,
	@translations_song_cannot_get_information,
	@translations_song_live,
	@translations_song_max_length,
	@translations_song_min_length,
	@translations_song_requested_message,
	@translations_song_maximum_ordered,
	@translations_song_min_views,
	@translations_channel_denied,
	@volume
)
ON CONFLICT (channel_id) DO UPDATE SET
	enabled = EXCLUDED.enabled,
	mode = EXCLUDED.mode,
	accept_only_when_online = EXCLUDED.accept_only_when_online,
	player_no_cookie_mode = EXCLUDED.player_no_cookie_mode,
	take_song_from_donation_message = EXCLUDED.take_song_from_donation_message,
	max_requests = EXCLUDED.max_requests,
	channel_points_reward_id = EXCLUDED.channel_points_reward_id,
	announce_play = EXCLUDED.announce_play,
	needed_votes_for_skip = EXCLUDED.needed_votes_for_skip,
	user_max_requests = EXCLUDED.user_max_requests,
	user_min_watch_time = EXCLUDED.user_min_watch_time,
	user_min_messages = EXCLUDED.user_min_messages,
	user_min_follow_time = EXCLUDED.user_min_follow_time,
	song_min_length = EXCLUDED.song_min_length,
	song_max_length = EXCLUDED.song_max_length,
	song_min_views = EXCLUDED.song_min_views,
	song_accepted_categories = EXCLUDED.song_accepted_categories,
	song_words_deny_list = EXCLUDED.song_words_deny_list,
	deny_list_users = EXCLUDED.deny_list_users,
	deny_list_songs = EXCLUDED.deny_list_songs,
	deny_list_channels = EXCLUDED.deny_list_channels,
	deny_list_artists_names = EXCLUDED.deny_list_artists_names,
	deny_list_words = EXCLUDED.deny_list_words,
	translations_now_playing = EXCLUDED.translations_now_playing,
	translations_not_enabled = EXCLUDED.translations_not_enabled,
	translations_no_text = EXCLUDED.translations_no_text,
	translations_accept_online_when_online = EXCLUDED.translations_accept_online_when_online,
	translations_user_denied = EXCLUDED.translations_user_denied,
	translations_user_max_requests = EXCLUDED.translations_user_max_requests,
	translations_user_min_messages = EXCLUDED.translations_user_min_messages,
	translations_user_min_watched = EXCLUDED.translations_user_min_watched,
	translations_user_min_follow = EXCLUDED.translations_user_min_follow,
	translations_song_denied = EXCLUDED.translations_song_denied,
	translations_song_not_found = EXCLUDED.translations_song_not_found,
	translations_song_already_in_queue = EXCLUDED.translations_song_already_in_queue,
	translations_song_age_restrictions = EXCLUDED.translations_song_age_restrictions,
	translations_song_cannot_get_information = EXCLUDED.translations_song_cannot_get_information,
	translations_song_live = EXCLUDED.translations_song_live,
	translations_song_max_length = EXCLUDED.translations_song_max_length,
	translations_song_min_length = EXCLUDED.translations_song_min_length,
	translations_song_requested_message = EXCLUDED.translations_song_requested_message,
	translations_song_maximum_ordered = EXCLUDED.translations_song_maximum_ordered,
	translations_song_min_views = EXCLUDED.translations_song_min_views,
	translations_channel_denied = EXCLUDED.translations_channel_denied,
	volume = EXCLUDED.volume
RETURNING ` + selectFields

	rows, err := c.pool.Query(
		ctx,
		query,
		pgx.NamedArgs{
			"channel_id":                               settings.ChannelID,
			"mode":                                     settings.Mode,
			"enabled":                                  settings.Enabled,
			"accept_only_when_online":                  settings.AcceptOnlyWhenOnline,
			"player_no_cookie_mode":                    settings.PlayerNoCookieMode,
			"take_song_from_donation_message":          settings.TakeSongFromDonationMessage,
			"max_requests":                             settings.MaxRequests,
			"channel_points_reward_id":                 settings.ChannelPointsRewardID,
			"announce_play":                            settings.AnnouncePlay,
			"needed_votes_for_skip":                    settings.NeededVotesForSkip,
			"user_max_requests":                        settings.UserMaxRequests,
			"user_min_watch_time":                      settings.UserMinWatchTime,
			"user_min_messages":                        settings.UserMinMessages,
			"user_min_follow_time":                     settings.UserMinFollowTime,
			"song_min_length":                          settings.SongMinLength,
			"song_max_length":                          settings.SongMaxLength,
			"song_min_views":                           settings.SongMinViews,
			"song_accepted_categories":                 settings.SongAcceptedCategories,
			"song_words_deny_list":                     settings.SongWordsDenyList,
			"deny_list_users":                          settings.DenyListUsers,
			"deny_list_songs":                          settings.DenyListSongs,
			"deny_list_channels":                       settings.DenyListChannels,
			"deny_list_artists_names":                  settings.DenyListArtistsNames,
			"deny_list_words":                          settings.DenyListWords,
			"translations_now_playing":                 settings.TranslationsNowPlaying,
			"translations_not_enabled":                 settings.TranslationsNotEnabled,
			"translations_no_text":                     settings.TranslationsNoText,
			"translations_accept_online_when_online":   settings.TranslationsAcceptOnlineWhenOnline,
			"translations_user_denied":                 settings.TranslationsUserDenied,
			"translations_user_max_requests":           settings.TranslationsUserMaxRequests,
			"translations_user_min_messages":           settings.TranslationsUserMinMessages,
			"translations_user_min_watched":            settings.TranslationsUserMinWatched,
			"translations_user_min_follow":             settings.TranslationsUserMinFollow,
			"translations_song_denied":                 settings.TranslationsSongDenied,
			"translations_song_not_found":              settings.TranslationsSongNotFound,
			"translations_song_already_in_queue":       settings.TranslationsSongAlreadyInQueue,
			"translations_song_age_restrictions":       settings.TranslationsSongAgeRestrictions,
			"translations_song_cannot_get_information": settings.TranslationsSongCannotGetInformation,
			"translations_song_live":                   settings.TranslationsSongLive,
			"translations_song_max_length":             settings.TranslationsSongMaxLength,
			"translations_song_min_length":             settings.TranslationsSongMinLength,
			"translations_song_requested_message":      settings.TranslationsSongRequestedMessage,
			"translations_song_maximum_ordered":        settings.TranslationsSongMaximumOrdered,
			"translations_song_min_views":              settings.TranslationsSongMinViews,
			"translations_channel_denied":              settings.TranslationsChannelDenied,
			"volume":                                   settings.Volume,
		},
	)
	if err != nil {
		return entity.Nil, fmt.Errorf("upsert song requests settings: %w", err)
	}

	upserted, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		return entity.Nil, fmt.Errorf("collect upserted song requests settings: %w", err)
	}

	return upserted.toEntity(), nil
}

func (c *Pgx) SetVolume(ctx context.Context, channelID string, volume int) error {
	query := `INSERT INTO channels_song_requests_settings (channel_id, enabled, volume)
VALUES (CAST(@channel_id AS uuid), false, @volume)
ON CONFLICT (channel_id) DO UPDATE SET volume = EXCLUDED.volume`

	if _, err := c.pool.Exec(
		ctx,
		query,
		pgx.NamedArgs{"channel_id": channelID, "volume": volume},
	); err != nil {
		return fmt.Errorf("set song requests volume: %w", err)
	}

	return nil
}
