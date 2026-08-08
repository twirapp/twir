package song_requests_settings

import (
	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/song_request_mode"
)

type Settings struct {
	ID                                   uuid.UUID
	ChannelID                            uuid.UUID
	Mode                                 song_request_mode.Mode
	Enabled                              bool
	AcceptOnlyWhenOnline                 bool
	PlayerNoCookieMode                   bool
	TakeSongFromDonationMessage          bool
	MaxRequests                          int
	ChannelPointsRewardID                *string
	AnnouncePlay                         bool
	NeededVotesForSkip                   float64
	UserMaxRequests                      int
	UserMinWatchTime                     int64
	UserMinMessages                      int
	UserMinFollowTime                    int
	SongMinLength                        int
	SongMaxLength                        int
	SongMinViews                         int
	SongAcceptedCategories               []string
	SongWordsDenyList                    []string
	DenyListUsers                        []string
	DenyListSongs                        []string
	DenyListChannels                     []string
	DenyListArtistsNames                 []string
	DenyListWords                        []string
	TranslationsNowPlaying               string
	TranslationsNotEnabled               string
	TranslationsNoText                   string
	TranslationsAcceptOnlineWhenOnline   string
	TranslationsUserDenied               string
	TranslationsUserMaxRequests          string
	TranslationsUserMinMessages          string
	TranslationsUserMinWatched           string
	TranslationsUserMinFollow            string
	TranslationsSongDenied               string
	TranslationsSongNotFound             string
	TranslationsSongAlreadyInQueue       string
	TranslationsSongAgeRestrictions      string
	TranslationsSongCannotGetInformation string
	TranslationsSongLive                 string
	TranslationsSongMaxLength            string
	TranslationsSongMinLength            string
	TranslationsSongRequestedMessage     string
	TranslationsSongMaximumOrdered       string
	TranslationsSongMinViews             string
	TranslationsChannelDenied            string
	Volume                               int

	isNil bool
}

func (c Settings) IsNil() bool {
	return c.isNil
}

var Nil = Settings{isNil: true}
