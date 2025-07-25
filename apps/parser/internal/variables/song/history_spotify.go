package song

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
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
			result.Result = fmt.Sprintf("cannot get spotify integration: %s", err)
			return result, nil
		}
		if spotifyEntity.AccessToken == "" {
			result.Result = "spotify not connected"
			return result, nil
		}

		spotifyIntegration := model.Integrations{}
		if err := parseCtx.Services.Gorm.
			Where("service = ?", "SPOTIFY").
			First(&spotifyIntegration).
			Error; err != nil {
			parseCtx.Services.Logger.Error("failed to get spotify integration", zap.Error(err))
			result.Result = "internal error while getting spotify integration"
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
				result.Result = "no needed scope, reconnect spotify in dashboard"
				return result, nil
			}
			result.Result = fmt.Sprintf("cannot get recent tracks: %s", err)
			return result, nil
		}
		mappedTracks := make([]string, len(tracks))
		for i, track := range tracks {
			playedAt, err := time.Parse("2006-01-02T15:04:05Z", track.PlayedAt)
			if err != nil {
				parseCtx.Services.Logger.Error("cannot parse played at", zap.Error(err))
				continue
			}

			ago := time.Now().UTC().Sub(playedAt)

			mappedTracks[i] = fmt.Sprintf(
				"%s — %s (~%.0fm ago)",
				track.Title,
				track.Artist,
				ago.Minutes(),
			)
		}

		result.Result = strings.Join(mappedTracks, " | ")

		return result, nil
	},
}
