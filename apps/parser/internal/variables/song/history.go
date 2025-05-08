package song

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/smithy-go/time"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/satont/twir/apps/parser/pkg/helpers"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/integrations/spotify"
	"go.uber.org/zap"
)

var History = &types.Variable{
	Name:                "songs.history",
	Description:         lo.ToPtr("Print history of played songs"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		spotifyEntity, err := parseCtx.Services.SpotifyRepo.GetByChannelID(ctx, parseCtx.Channel.ID)
		if err != nil {
			return nil, fmt.Errorf("cannot get spotify service for streamer: %w", err)
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
			result.Result = "cannot get spotify integration"
			return result, nil
		}

		spotifyService := spotify.New(spotifyIntegration, spotifyEntity, parseCtx.Services.SpotifyRepo)

		history, err := spotifyService.GetRecentTracks(ctx, spotify.GetRecentTracksInput{Limit: 5})
		if err != nil {
			if errors.Is(err, spotify.ErrNoNeededScope) {
				result.Result = `missed scope "user-read-recently-played", go to dashboard and reconnect spotify.`
				return result, nil
			}

			parseCtx.Services.Logger.Error("cannot get recent tracks", zap.Error(err))
			result.Result = "cannot get recent tracks"
			return result, nil
		}

		tracks := make([]string, 0, len(history))
		for _, track := range history {
			playedAt, err := time.ParseDateTime(track.PlayedAt)
			if err != nil {
				parseCtx.Services.Logger.Error("cannot parse played at", zap.Error(err))
				continue
			}

			playedAgo := helpers.Duration(
				playedAt,
				&helpers.DurationOpts{
					UseUtc: true,
					Hide: helpers.DurationOptsHide{
						Seconds: true,
					},
				},
			)

			tracks = append(tracks, fmt.Sprintf("%s — %s (%v ago)", track.Artist, track.Title, playedAgo))
		}

		result.Result = strings.Join(tracks, " | ")

		return result, nil
	},
}
