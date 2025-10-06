package song

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	"github.com/twirapp/twir/libs/integrations/lastfm"
	"github.com/twirapp/twir/libs/integrations/spotify"
	"go.uber.org/zap"
)

type recentTrack struct {
	Title    string
	Artist   string
	PlayedAt time.Time
}

var History = &types.Variable{
	Name:                "songs.history",
	Description:         lo.ToPtr("Print combined history of played songs from spotify and lastfm"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}
		limit := 5

		integrations := parseCtx.Cacher.GetEnabledChannelIntegrations(ctx)
		lastFmIntegration, ok := lo.Find(
			integrations,
			func(integration *model.ChannelsIntegrations) bool {
				return integration.Integration.Service == "LASTFM"
			},
		)

		var lastfmService *lastfm.Lastfm
		if ok {
			i, err := lastfm.New(
				lastfm.Opts{
					Gorm:        parseCtx.Services.Gorm,
					Integration: lastFmIntegration,
				},
			)
			if err == nil {
				lastfmService = i
			}
		}

		var spotifyService *spotify.Spotify
		spotifyEntity, err := parseCtx.Services.SpotifyRepo.GetByChannelID(ctx, parseCtx.Channel.ID)
		if err != nil {
			parseCtx.Services.Logger.Error(i18n.GetCtx(ctx, locales.Translations.Variables.Song.Info.GetSpotifyEntity), zap.Error(err))
		} else if spotifyEntity.AccessToken != "" {
			spotifyIntegration := model.Integrations{}
			if err := parseCtx.Services.Gorm.
				Where("service = ?", "SPOTIFY").
				First(&spotifyIntegration).
				Error; err != nil {
				parseCtx.Services.Logger.Error(i18n.GetCtx(ctx, locales.Translations.Variables.Song.Info.FailedGetSpotifyIntegration), zap.Error(err))
			} else {
				spotifyService = spotify.New(
					spotifyIntegration,
					spotifyEntity,
					parseCtx.Services.SpotifyRepo,
				)
			}
		}

		if lastfmService == nil && spotifyService == nil {
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Song.Info.NoIntegrations)
			return result, nil
		}

		recentTracks := make(map[string]recentTrack)

		if lastfmService != nil {
			lfmTracks, err := lastfmService.GetRecentTracks(limit)
			if err != nil {
				result.Result = i18n.GetCtx(
					ctx,
					locales.Translations.Variables.Song.Errors.FetchTracksLastfm.
						SetVars(locales.KeysVariablesSongErrorsFetchTracksLastfmVars{Reason: err.Error()}),
				)
				return result, nil
			}

			for _, track := range lfmTracks {
				_, exists := recentTracks[track.Title]
				if exists {
					continue
				}

				if track.PlayedUTS == "" {
					continue
				}

				timestamp, err := strconv.ParseInt(track.PlayedUTS, 10, 64)
				if err != nil {
					continue
				}

				newTrack := recentTrack{
					Title:    track.Title,
					Artist:   track.Artist,
					PlayedAt: time.Unix(timestamp, 0),
				}

				recentTracks[track.Title] = newTrack
			}
		}

		if spotifyService != nil {
			spotifyTracks, err := spotifyService.GetRecentTracks(
				ctx,
				spotify.GetRecentTracksInput{Limit: limit},
			)
			if err != nil {
				if errors.Is(err, spotify.ErrNoNeededScope) {
					result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Song.Info.NoNeededScope)
					return result, nil
				}
				result.Result = i18n.GetCtx(
					ctx,
					locales.Translations.Variables.Song.Errors.FetchTracksSpotify.
						SetVars(locales.KeysVariablesSongErrorsFetchTracksSpotifyVars{Reason: err.Error()}),
				)
				return result, nil
			}

			for _, track := range spotifyTracks {
				_, exists := recentTracks[track.Title]
				if exists {
					continue
				}

				playedAt, err := time.Parse("2006-01-02T15:04:05Z", track.PlayedAt)
				if err != nil {
					parseCtx.Services.Logger.Error(i18n.GetCtx(ctx, locales.Translations.Variables.Song.Errors.ParsePlayedAt), zap.Error(err))
					continue
				}

				recentTracks[track.Title] = recentTrack{
					Title:    track.Title,
					Artist:   track.Artist,
					PlayedAt: playedAt,
				}
			}
		}

		recentTracksSlice := make([]recentTrack, 0, len(recentTracks))
		for _, t := range recentTracks {
			recentTracksSlice = append(recentTracksSlice, t)
		}

		slices.SortFunc(
			recentTracksSlice, func(a, b recentTrack) int {
				return b.PlayedAt.Compare(a.PlayedAt)
			},
		)

		mappedTracks := make([]string, len(recentTracksSlice))
		for i, track := range recentTracksSlice {
			ago := time.Now().UTC().Sub(track.PlayedAt)

			mappedTracks[i] = i18n.GetCtx(
				ctx,
				locales.Translations.Variables.Song.Info.History.
					SetVars(locales.KeysVariablesSongInfoHistoryVars{TrackTitle: track.Title, TrackArtist: track.Artist, Minutes: fmt.Sprintf("%.0f%%", ago.Minutes())}),
			)
		}

		result.Result = strings.Join(mappedTracks, " | ")

		return result, nil
	},
}
