package song

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/i18n"
	"github.com/twirapp/twir/libs/integrations/lastfm"
)

var HistoryLastfm = &types.Variable{
	Name:                "songs.history.lastfm",
	Description:         lo.ToPtr("Print history of played songs from lastfm"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		limit := 10
		if variableData.Params != nil {
			newLimit, err := strconv.Atoi(*variableData.Params)
			if err == nil {
				limit = newLimit
			}
		}

		if limit > 50 {
			limit = 10
		}

		lastfmIntegration, err := parseCtx.Services.LastfmRepo.GetByChannelID(ctx, parseCtx.Channel.ID)
		if err != nil || lastfmIntegration.IsNil() || !lastfmIntegration.Enabled || lastfmIntegration.SessionKey == nil {
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Song.Info.LastfmIntegration)
			return result, nil
		}

		lastfmService, err := lastfm.New(
			lastfm.Opts{
				ApiKey:       parseCtx.Services.Config.LastFM.ApiKey,
				ClientSecret: parseCtx.Services.Config.LastFM.ClientSecret,
				SessionKey:   *lastfmIntegration.SessionKey,
			},
		)
		if err != nil {
			result.Result = i18n.GetCtx(
				ctx,
				locales.Translations.Variables.Song.Errors.CreateLastfmService.
					SetVars(locales.KeysVariablesSongErrorsCreateLastfmServiceVars{Reason: err.Error()}),
			)
			return result, nil
		}

		tracks, err := lastfmService.GetRecentTracks(limit)
		if err != nil {
			result.Result = i18n.GetCtx(
				ctx,
				locales.Translations.Variables.Song.Errors.FetchTracksLastfm.
					SetVars(locales.KeysVariablesSongErrorsFetchTracksLastfmVars{Reason: err.Error()}),
			)
			return result, nil
		}

		mappedTracks := make([]string, 0, len(tracks))
		for _, track := range tracks {
			if track.PlayedUTS == "" {
				continue
			}
			timestamp, err := strconv.ParseInt(track.PlayedUTS, 10, 64)
			if err != nil {
				continue
			}
			playedAt := time.Unix(timestamp, 0)

			ago := time.Now().UTC().Sub(playedAt)

			mappedTracks = append(
				mappedTracks,
				i18n.GetCtx(
					ctx,
					locales.Translations.Variables.Song.Info.History.
						SetVars(locales.KeysVariablesSongInfoHistoryVars{TrackTitle: track.Title, TrackArtist: track.Artist, Minutes: fmt.Sprintf("%.0f%%", ago.Minutes())}),
				),
			)
		}

		result.Result = strings.Join(mappedTracks, " | ")

		return result, nil
	},
}
