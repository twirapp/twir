package song

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	"github.com/twirapp/twir/libs/integrations/spotify"
	"go.uber.org/zap"
)

var HistorySpotify = &types.Variable{
	Name:                "songs.history.spotify",
	Description:         lo.ToPtr("Print history of played songs from spotify"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		spotifyEntity, err := parseCtx.Services.SpotifyRepo.GetByChannelID(ctx, parseCtx.Channel.ID)
		if err != nil {
			result.Result = i18n.GetCtx(
				ctx,
				locales.Translations.Variables.Song.Info.GetSpotifyIntegration.
					SetVars(locales.KeysVariablesSongInfoGetSpotifyIntegrationVars{Reason: err.Error()}),
			)
			return result, nil
		}
		if spotifyEntity.AccessToken == "" {
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Song.Info.SpotifyNotConnected)
			return result, nil
		}

		spotifyIntegration := model.Integrations{}
		if err := parseCtx.Services.Gorm.
			Where("service = ?", "SPOTIFY").
			First(&spotifyIntegration).
			Error; err != nil {
			parseCtx.Services.Logger.Error(i18n.GetCtx(ctx, locales.Translations.Variables.Song.Info.FailedGetSpotifyIntegration), zap.Error(err))
			result.Result = i18n.GetCtx(ctx, locales.Translations.Errors.Generic.Internal)
			return result, nil
		}

		spotifyService := spotify.New(
			spotifyIntegration,
			spotifyEntity,
			parseCtx.Services.SpotifyRepo,
		)

		tracks, err := spotifyService.GetRecentTracks(ctx, spotify.GetRecentTracksInput{Limit: 10})
		if err != nil {
			if errors.Is(err, spotify.ErrNoNeededScope) {
				result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Song.Info.NoNeededScope)
				return result, nil
			}
			result.Result = i18n.GetCtx(
				ctx,
				locales.Translations.Variables.Song.Errors.GetRecentTracks.
					SetVars(locales.KeysVariablesSongErrorsGetRecentTracksVars{Reason: err.Error()}),
			)
			return result, nil
		}
		mappedTracks := make([]string, len(tracks))
		for i, track := range tracks {
			playedAt, err := time.Parse("2006-01-02T15:04:05Z", track.PlayedAt)
			if err != nil {
				parseCtx.Services.Logger.Error(i18n.GetCtx(ctx, locales.Translations.Variables.Song.Errors.ParsePlayedAt), zap.Error(err))
				continue
			}

			ago := time.Now().UTC().Sub(playedAt)

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
