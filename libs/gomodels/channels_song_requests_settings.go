package model

import (
	"github.com/guregu/null"
	"github.com/lib/pq"
)

type ChannelSongRequestsSettings struct {
	ID                                   string         `gorm:"column:id;type:uuid;primary_key;default:gen_random_uuid()"`
	ChannelID                            string         `gorm:"column:channel_id;type:text;not null"`
	Enabled                              bool           `gorm:"column:enabled;type:boolean;not null"`
	AcceptOnlyWhenOnline                 bool           `gorm:"column:accept_only_when_online;type:boolean;not null"`
	PlayerNoCookieMode                   bool           `gorm:"column:player_no_cookie_mode;type:boolean;not null"`
	TakeSongFromDonationMessage          bool           `gorm:"column:take_song_from_donation_message;type:boolean;not null;default:false"`
	MaxRequests                          int            `gorm:"column:max_requests;type:int;check:max_requests <= 500"`
	ChannelPointsRewardID                null.String    `gorm:"column:channel_points_reward_id;type:varchar(100);check:length(channel_points_reward_id) <= 100"`
	AnnouncePlay                         bool           `gorm:"column:announce_play;type:boolean;not null"`
	NeededVotesForSkip                   float64        `gorm:"column:needed_votes_for_skip;type:float;check:needed_votes_for_skip <= 100 AND needed_votes_for_skip >=1"`
	UserMaxRequests                      int            `gorm:"column:user_max_requests;type:int"`
	UserMinWatchTime                     int64          `gorm:"column:user_min_watch_time;type:bigint"`
	UserMinMessages                      int            `gorm:"column:user_min_messages;type:int"`
	UserMinFollowTime                    int            `gorm:"column:user_min_follow_time;type:int"`
	SongMinLength                        int            `gorm:"column:song_min_length;type:int;check:song_min_length >= 0 AND song_min_length <= 86399"`
	SongMaxLength                        int            `gorm:"column:song_max_length;type:int;check:song_max_length <= 86400"`
	SongMinViews                         int            `gorm:"column:song_min_views;type:int;check:song_min_views <= 10000000000000"`
	SongAcceptedCategories               pq.StringArray `gorm:"column:song_accepted_categories;type:varchar[];check:array_length(song_accepted_categories, 1) <= 300"`
	SongWordsDenyList                    pq.StringArray `gorm:"column:song_words_deny_list;type:varchar[];check:array_length(song_words_deny_list, 1) <= 300"`
	DenyListUsers                        pq.StringArray `gorm:"column:deny_list_users;type:varchar[];check:array_length(deny_list_users, 1) >= 1"`
	DenyListSongs                        pq.StringArray `gorm:"column:deny_list_songs;type:varchar[];check:array_length(deny_list_songs, 1) >= 1"`
	DenyListChannels                     pq.StringArray `gorm:"column:deny_list_channels;type:varchar[];check:array_length(deny_list_channels, 1) >= 1"`
	DenyListArtistsNames                 pq.StringArray `gorm:"column:deny_list_artists_names;type:varchar[];check:array_length(deny_list_artists_names, 1) <= 300"`
	DenyListWords                        pq.StringArray `gorm:"column:deny_list_words;type:varchar[];check:array_length(deny_list_words, 1) <= 300"`
	TranslationsNowPlaying               string         `gorm:"column:translations_now_playing;type:varchar;check:length(translations_now_playing) <= 300"`
	TranslationsNotEnabled               string         `gorm:"column:translations_not_enabled;type:varchar;check:length(translations_not_enabled) <= 300"`
	TranslationsNoText                   string         `gorm:"column:translations_no_text;type:varchar;check:length(translations_no_text) <= 300"`
	TranslationsAcceptOnlineWhenOnline   string         `gorm:"column:translations_accept_online_when_online;type:varchar;check:length(translations_accept_online_when_online) <= 300"`
	TranslationsUserDenied               string         `gorm:"column:translations_user_denied;type:varchar;check:length(translations_user_denied) <= 300"`
	TranslationsUserMaxRequests          string         `gorm:"column:translations_user_max_requests;type:varchar;check:length(translations_user_max_requests) <= 300"`
	TranslationsUserMinMessages          string         `gorm:"column:translations_user_min_messages;type:varchar;check:length(translations_user_min_messages) <= 300"`
	TranslationsUserMinWatched           string         `gorm:"column:translations_user_min_watched;type:varchar;check:length(translations_user_min_watched) <= 300"`
	TranslationsUserMinFollow            string         `gorm:"column:translations_user_min_follow;type:varchar;check:length(translations_user_min_follow) <= 300"`
	TranslationsSongDenied               string         `gorm:"column:translations_song_denied;type:varchar;check:length(translations_song_denied) <= 300"`
	TranslationsSongNotFound             string         `gorm:"column:translations_song_not_found;type:varchar;check:length(translations_song_not_found) <= 300"`
	TranslationsSongAlreadyInQueue       string         `gorm:"column:translations_song_already_in_queue;type:varchar;check:length(translations_song_already_in_queue) <= 300"`
	TranslationsSongAgeRestrictions      string         `gorm:"column:translations_song_age_restrictions;type:varchar;check:length(translations_song_age_restrictions) <= 300"`
	TranslationsSongCannotGetInformation string         `gorm:"column:translations_song_cannot_get_information;type:varchar;check:length(translations_song_cannot_get_information) <= 300"`
	TranslationsSongLive                 string         `gorm:"column:translations_song_live;type:varchar;check:length(translations_song_live) <= 300"`
	TranslationsSongMaxLength            string         `gorm:"column:translations_song_max_length;type:varchar;check:length(translations_song_max_length) <= 300"`
	TranslationsSongMinLength            string         `gorm:"column:translations_song_min_length;type:varchar;check:length(translations_song_min_length) <= 300"`
	TranslationsSongRequestedMessage     string         `gorm:"column:translations_song_requested_message;type:varchar;check:length(translations_song_requested_message) <= 300"`
	TranslationsSongMaximumOrdered       string         `gorm:"column:translations_song_maximum_ordered;type:varchar;check:length(translations_song_maximum_ordered) <= 300"`
	TranslationsSongMinViews             string         `gorm:"column:translations_song_min_views;type:varchar;check:length(translations_song_min_views) <= 300"`
	TranslationsChannelDenied            string         `gorm:"column:translations_channel_denied;type:varchar;check:length(translations_channel_denied) <= 300"`
}

func (ChannelSongRequestsSettings) TableName() string {
	return "channels_song_requests_settings"
}
