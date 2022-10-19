package song

import (
	"fmt"
	"strings"
	model "tsuwari/models"
	lastfm "tsuwari/parser/internal/integrations/lastfm"
	spotify "tsuwari/parser/internal/integrations/spotify"
	vkIntegr "tsuwari/parser/internal/integrations/vk"
	"tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
)

const (
	SOUNDTRACK = "TWITCH_SOUNDTRACK"
	VK         = "VK"
	SPOTIFY    = "SPOTIFY"
	LASTFM     = "LASTFM"
)

var Variable = types.Variable{
	Name:        "currentsong",
	Description: lo.ToPtr("Current played song"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		integrations := *ctx.GetEnabledIntegrations()
		if integrations == nil {
			result.Result = "Haven't enabled integrations for fetching song"
			return result, nil
		}

		integrations = lo.Filter(
			integrations,
			func(integration model.ChannelInegrationWithRelation, _ int) bool {
				switch integration.Integration.Service {
				case SPOTIFY, VK, LASTFM:
					return integration.Enabled
				default:
					return false
				}
			},
		)

		lastFmIntegration, _ := lo.Find(
			integrations,
			func(integration model.ChannelInegrationWithRelation) bool {
				return integration.Integration.Service == "LASTFM"
			},
		)
		lfm := lastfm.New(&lastFmIntegration)

		spotifyIntegration, _ := lo.Find(
			integrations,
			func(integration model.ChannelInegrationWithRelation) bool {
				return integration.Integration.Service == "SPOTIFY"
			},
		)
		spoti := spotify.New(&spotifyIntegration, ctx.Services.Db)

		vkIntegration, _ := lo.Find(
			integrations,
			func(integration model.ChannelInegrationWithRelation) bool {
				return integration.Integration.Service == "VK"
			},
		)
		vk := vkIntegr.New(&vkIntegration)

		integrationsForFetch := lo.Map(
			integrations,
			func(integration model.ChannelInegrationWithRelation, _ int) string {
				return integration.Integration.Service
			},
		)

		integrationsForFetch = append(integrationsForFetch, SOUNDTRACK)

	checkServices:
		for _, integration := range integrationsForFetch {
			switch integration {
			case SPOTIFY:
				if spoti == nil {
					continue
				}
				track := spoti.GetTrack()
				if track != nil {
					result.Result = *track
					break checkServices
				}
			case LASTFM:
				if lfm == nil {
					continue
				}

				track := lfm.GetTrack()
				if track != nil {
					result.Result = *track
					break checkServices
				}
			case VK:
				if vk == nil {
					continue
				}
				track := vk.GetTrack()
				if track != nil {
					result.Result = *track
					break checkServices
				}
			case SOUNDTRACK:
				tracks, err := ctx.Services.Twitch.Client.GetSoundTrackCurrentTrack(&helix.SoundtrackCurrentTrackParams{
					BroadcasterID: ctx.ChannelId,
				})
				if err != nil {
					fmt.Println(err)
					continue
				}

				if len(tracks.Data.Tracks) == 0 {
					continue
				}

				track := tracks.Data.Tracks[0]
				artists := lo.Map(track.Track.Artists, func(artist helix.SoundtrackTrackArtist, _ int) string {
					return artist.Name
				})
				result.Result = fmt.Sprintf("%s — %s", strings.Join(artists, ", "), track.Track.Title)
				break checkServices
			}
		}

		return result, nil
	},
}
