package migrations

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upSongsRequestToTable, downSongsRequestToTable)
}

type YouTubeUserSettings struct {
	MaxRequests   int   `json:"maxRequests"`
	MinWatchTime  int64 `json:"minWatchTime"`
	MinMessages   int   `json:"minMessages"`
	MinFollowTime int   `json:"minFollowTime"`
}

type YouTubeSongSettings struct {
	MinLength          int      `validate:"gte=0,lte=86399" json:"minLength"`
	MaxLength          int      `validate:"lte=86400"          json:"maxLength"`
	MinViews           int      `validate:"lte=10000000000000" json:"minViews"`
	AcceptedCategories []string `validate:"dive,max=300"       json:"acceptedCategories"`
	WordsDenyList      []string `validate:"dive,max=300"       json:"wordsDenyList"`
}

type YouTubeDenySettingsSongs struct {
	ID        string `validate:"required,min=1,max=300" json:"id"`
	Title     string `validate:"required,min=1,max=300" json:"title"`
	ThumbNail string `validate:"required,min=1,max=300" json:"thumbNail"`
}

type YouTubeDenySettingsChannels struct {
	ID        string `validate:"required,min=1"         json:"id"`
	Title     string `validate:"required,min=1,max=300" json:"title"`
	ThumbNail string `validate:"required,min=1,max=300" json:"thumbNail"`
}

type YouTubeDenyList struct {
	Users        []string `validate:"required,dive"         json:"users"`
	Songs        []string `validate:"required,dive"         json:"songs"`
	Channels     []string `validate:"required,dive"         json:"channels"`
	ArtistsNames []string `validate:"required,dive,max=300" json:"artistsNames"`
	Words        []string `validate:"required,dive,max=300" json:"words"`
}

type YouTubeUserTranslations struct {
	Denied      string `json:"denied"`
	MaxRequests string `json:"maxRequests"`
	MinMessages string `json:"minMessages"`
	MinWatched  string `json:"minWatched"`
	MinFollow   string `json:"minFollow"`
}

type YouTubeSongTranslations struct {
	Denied               string `json:"denied"`
	NotFound             string `json:"notFound"`
	AlreadyInQueue       string `json:"alreadyInQueue"`
	AgeRestrictions      string `json:"ageRestrictions"`
	CannotGetInformation string `json:"cannotGetInformation"`
	Live                 string `json:"live"`
	MaxLength            string `json:"maxLength"`
	MinLength            string `json:"minLength"`
	RequestedMessage     string `json:"requestedMessage"`
	MaximumOrdered       string `json:"maximumOrdered"`
	MinViews             string `json:"minViews"`
}

type YouTubeChannelTranslations struct {
	Denied string `json:"denied"`
}

type YouTubeTranslations struct {
	NowPlaying             string                     `json:"nowPlaying"`
	NotEnabled             string                     `json:"notEnabled"`
	NoText                 string                     `json:"noText"`
	AcceptOnlineWhenOnline string                     `json:"acceptOnlyWhenOnline"`
	User                   YouTubeUserTranslations    `json:"user"`
	Song                   YouTubeSongTranslations    `json:"song"`
	Channel                YouTubeChannelTranslations `json:"channel"`
}

type YouTubeSettings struct {
	Enabled                     *bool               `validate:"required" json:"enabled"`
	AcceptOnlyWhenOnline        *bool               `validate:"required" json:"acceptOnlyWhenOnline"`
	PlayerNoCookieMode          *bool               `validate:"required" json:"playerNoCookieMode"`
	TakeSongFromDonationMessage bool                `json:"takeSongFromDonationMessage"`
	MaxRequests                 int                 `validate:"lte=500"  json:"maxRequests"`
	ChannelPointsRewardId       string              `validate:"max=100"  json:"channelPointsRewardId"`
	AnnouncePlay                *bool               `validate:"required" json:"announcePlay"`
	NeededVotesVorSkip          float64             `validate:"max=100,min=1" json:"neededVotesVorSkip"`
	User                        YouTubeUserSettings `validate:"required" json:"user"`
	Song                        YouTubeSongSettings `validate:"required" json:"song"`
	DenyList                    YouTubeDenyList     `validate:"required" json:"denyList"`
	Translations                YouTubeTranslations `validate:"required" json:"translations"`
}

type SongRequestsToColumns struct {
	ID        string
	Settings  []byte
	ChannelId string `gorm:"column:channelId;type:text" json:"channelId"`
}

func upSongsRequestToTable(ctx context.Context, tx *sql.Tx) error {
	var entities []SongRequestsToColumns

	rows, err := tx.QueryContext(
		ctx,
		`SELECT id, settings, "channelId" FROM channels_modules_settings WHERE type = 'youtube_song_requests'`,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var entity SongRequestsToColumns
		if err := rows.Scan(&entity.ID, &entity.Settings, &entity.ChannelId); err != nil {
			return err
		}
		entities = append(entities, entity)
	}

	_, err = tx.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS channels_song_requests_settings (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			channel_id text references channels(id),
			enabled BOOLEAN NOT NULL,
			accept_only_when_online BOOLEAN NOT NULL DEFAULT false,
			player_no_cookie_mode BOOLEAN NOT NULL DEFAULT false,
			take_song_from_donation_message BOOLEAN NOT NULL DEFAULT false,
			max_requests INT CHECK (max_requests <= 500),
			channel_points_reward_id VARCHAR(100) CHECK (LENGTH(channel_points_reward_id) <= 100),
			announce_play BOOLEAN NOT NULL DEFAULT true,
			needed_votes_for_skip FLOAT CHECK (needed_votes_for_skip <= 100 AND needed_votes_for_skip >= 1),
			user_max_requests INT DEFAULT 20,
			user_min_watch_time BIGINT DEFAULT 0,
			user_min_messages INT DEFAULT 0,
			user_min_follow_time INT DEFAULT 0,
			song_min_length INT DEFAULT 0 CHECK (song_min_length >= 0 AND song_min_length <= 86399),
			song_max_length INT DEFAULT 10 CHECK (song_max_length <= 86400),
			song_min_views INT DEFAULT 50000 CHECK (song_min_views <= 10000000000000),
			song_accepted_categories VARCHAR[] DEFAULT '{}' CHECK (array_length(song_accepted_categories, 1) <= 300),
			song_words_deny_list VARCHAR[] DEFAULT '{}' CHECK (array_length(song_words_deny_list, 1) <= 300),
			deny_list_users VARCHAR[] DEFAULT '{}' CHECK (array_length(deny_list_users, 1) >= 1),
			deny_list_songs VARCHAR[]  DEFAULT '{}'CHECK (array_length(deny_list_songs, 1) >= 1),
			deny_list_channels VARCHAR[] DEFAULT '{}' CHECK (array_length(deny_list_channels, 1) >= 1),
			deny_list_artists_names VARCHAR[] DEFAULT '{}' CHECK (array_length(deny_list_artists_names, 1) <= 300),
			deny_list_words VARCHAR[] DEFAULT '{}' CHECK (array_length(deny_list_words, 1) <= 300),
			translations_now_playing VARCHAR CHECK (LENGTH(translations_now_playing) <= 300),
			translations_not_enabled VARCHAR CHECK (LENGTH(translations_not_enabled) <= 300),
			translations_no_text VARCHAR CHECK (LENGTH(translations_no_text) <= 300),
			translations_accept_online_when_online VARCHAR CHECK (LENGTH(translations_accept_online_when_online) <= 300),
			translations_user_denied VARCHAR CHECK (LENGTH(translations_user_denied) <= 300),
			translations_user_max_requests VARCHAR CHECK (LENGTH(translations_user_max_requests) <= 300),
			translations_user_min_messages VARCHAR CHECK (LENGTH(translations_user_min_messages) <= 300),
			translations_user_min_watched VARCHAR CHECK (LENGTH(translations_user_min_watched) <= 300),
			translations_user_min_follow VARCHAR CHECK (LENGTH(translations_user_min_follow) <= 300),
			translations_song_denied VARCHAR CHECK (LENGTH(translations_song_denied) <= 300),
			translations_song_not_found VARCHAR CHECK (LENGTH(translations_song_not_found) <= 300),
			translations_song_already_in_queue VARCHAR CHECK (LENGTH(translations_song_already_in_queue) <= 300),
			translations_song_age_restrictions VARCHAR CHECK (LENGTH(translations_song_age_restrictions) <= 300),
			translations_song_cannot_get_information VARCHAR CHECK (LENGTH(translations_song_cannot_get_information) <= 300),
			translations_song_live VARCHAR CHECK (LENGTH(translations_song_live) <= 300),
			translations_song_max_length VARCHAR CHECK (LENGTH(translations_song_max_length) <= 300),
			translations_song_min_length VARCHAR CHECK (LENGTH(translations_song_min_length) <= 300),
			translations_song_requested_message VARCHAR CHECK (LENGTH(translations_song_requested_message) <= 300),
			translations_song_maximum_ordered VARCHAR CHECK (LENGTH(translations_song_maximum_ordered) <= 300),
			translations_song_min_views VARCHAR CHECK (LENGTH(translations_song_min_views) <= 300),
			translations_channel_denied VARCHAR CHECK (LENGTH(translations_channel_denied) <= 300),
			UNIQUE (channel_id)
		)`,
	)
	if err != nil {
		return err
	}

	for _, entity := range entities {
		var settings YouTubeSettings
		if err := json.Unmarshal(entity.Settings, &settings); err != nil {
			return err
		}

		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO channels_song_requests_settings (
				id,
				channel_id,
				enabled,
				accept_only_when_online,
				player_no_cookie_mode,
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
				translations_channel_denied
			) VALUES (
				$1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7,
				$8,
				$9,
				$10,
				$11,
				$12,
				$13,
				$14,
				$15,
				$16,
				$17,
				$18,
				$19,
				$20,
				$21,
				$22,
				$23,
				$24,
				$25,
				$26,
				$27,
				$28,
				$29,
				$30,
				$31,
				$32,
				$33,
				$34,
				$35,
				$36,
				$37,
				$38,
				$39,
				$40,
				$41,
				$42,
				$43,
				$44
			)
			`,
			entity.ID,
			entity.ChannelId,
			settings.Enabled,
			settings.AcceptOnlyWhenOnline,
			settings.PlayerNoCookieMode,
			settings.MaxRequests,
			settings.ChannelPointsRewardId,
			settings.AnnouncePlay,
			settings.NeededVotesVorSkip,
			settings.User.MaxRequests,
			settings.User.MinWatchTime,
			settings.User.MinMessages,
			settings.User.MinFollowTime,
			settings.Song.MinLength,
			settings.Song.MaxLength,
			settings.Song.MinViews,
			append(pq.StringArray{}, settings.Song.AcceptedCategories...),
			append(pq.StringArray{}, settings.Song.WordsDenyList...),
			append(pq.StringArray{}, settings.DenyList.Users...),
			append(pq.StringArray{}, settings.DenyList.Songs...),
			append(pq.StringArray{}, settings.DenyList.Channels...),
			append(pq.StringArray{}, settings.DenyList.ArtistsNames...),
			append(pq.StringArray{}, settings.DenyList.Words...),
			settings.Translations.NowPlaying,
			settings.Translations.NotEnabled,
			settings.Translations.NoText,
			settings.Translations.AcceptOnlineWhenOnline,
			settings.Translations.User.Denied,
			settings.Translations.User.MaxRequests,
			settings.Translations.User.MinMessages,
			settings.Translations.User.MinWatched,
			settings.Translations.User.MinFollow,
			settings.Translations.Song.Denied,
			settings.Translations.Song.NotFound,
			settings.Translations.Song.AlreadyInQueue,
			settings.Translations.Song.AgeRestrictions,
			settings.Translations.Song.CannotGetInformation,
			settings.Translations.Song.Live,
			settings.Translations.Song.MaxLength,
			settings.Translations.Song.MinLength,
			settings.Translations.Song.RequestedMessage,
			settings.Translations.Song.MaximumOrdered,
			settings.Translations.Song.MinViews,
			settings.Translations.Channel.Denied,
		)

		if err != nil {
			return err
		}

		_, err = tx.ExecContext(
			ctx,
			`DELETE FROM channels_modules_settings WHERE id = $1`,
			entity.ID,
		)
	}

	// This code is executed when the migration is applied.
	return nil
}

func downSongsRequestToTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
