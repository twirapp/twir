package song

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
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

		enabledIntegrations := parseCtx.Cacher.GetEnabledChannelIntegrations(ctx)
		lastfmIntegration, ok := lo.Find(
			enabledIntegrations,
			func(integration *model.ChannelsIntegrations) bool {
				return integration.Integration.Service == "LASTFM"
			},
		)
		if !ok {
			result.Result = "lastfm integration not enabled"
			return result, nil
		}

		lastfmService, err := lastfm.New(
			lastfm.Opts{
				Gorm:        parseCtx.Services.Gorm,
				Integration: lastfmIntegration,
			},
		)
		if err != nil {
			result.Result = fmt.Sprintf("cannot create lastfm service: %s", err)
			return result, nil
		}

		tracks, err := lastfmService.GetRecentTracks(limit)
		if err != nil {
			result.Result = fmt.Sprintf("cannot fetch tracks from lastfm: %s", err)
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
				fmt.Sprintf(
					"%s — %s (~%.0fm ago)",
					track.Title,
					track.Artist,
					ago.Minutes(),
				),
			)
		}

		result.Result = strings.Join(mappedTracks, " | ")

		return result, nil
	},
}
