package mcp

import (
	"context"
	"fmt"
	"reflect"

	modelsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/twirapp/twir/apps/api-gql/internal/services/seventv_integration"
	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
	model "github.com/twirapp/twir/libs/gomodels"
)

type integrationResult struct {
	Connected bool   `json:"connected"`
	Data      any    `json:"data,omitempty"`
	Error     string `json:"error,omitempty"`
}

type toggleIntegrationInput struct {
	Integration string `json:"integration" jsonschema:"discord, spotify, lastfm, valorant, faceit, donationalerts, donatepay, donatestream, donatello, streamlabs, vk, or seventv"`
	Enabled     bool   `json:"enabled"`
	Credential  string `json:"credential,omitempty" jsonschema:"API key or confirmation code for credential-based integrations"`
	GuildID     string `json:"guildId,omitempty" jsonschema:"Discord guild ID when disconnecting"`
}

func integrationStatus(requestScope scope, data any, err error) integrationResult {
	if err != nil {
		return integrationResult{Error: err.Error()}
	}
	connected := data != nil && !reflect.ValueOf(data).IsZero()
	result := integrationResult{Connected: connected}
	if integrationDetailsAllowed(requestScope.AccessScopes) {
		result.Data = data
	}
	return result
}

func integrationDetailsAllowed(accessScopes toolAccessScopes) bool {
	scopes := make([]entity.Scope, 0, len(accessScopes))
	for scope := range accessScopes {
		scopes = append(scopes, scope)
	}
	return entity.HasScope(scopes, entity.ScopeGroupIntegrations, entity.ScopeActionEdit)
}

func (h *Handler) legacyIntegrationStatus(ctx context.Context, requestScope scope, service string) integrationResult {
	var item model.ChannelsIntegrations
	err := h.deps.Gorm.WithContext(ctx).Joins("JOIN integrations i ON i.id = channels_integrations.\"integrationId\"").Where(`channels_integrations."channelId" = ? AND i.service = ?`, requestScope.Channel.ID.String(), service).First(&item).Error
	if err != nil {
		return integrationResult{}
	}
	return legacyIntegrationResult(requestScope, item.Enabled)
}

func legacyIntegrationResult(requestScope scope, enabled bool) integrationResult {
	result := integrationResult{Connected: enabled}
	if integrationDetailsAllowed(requestScope.AccessScopes) {
		result.Data = map[string]any{"enabled": enabled}
	}
	return result
}

func (h *Handler) setLegacyIntegrationEnabled(ctx context.Context, channelID, service string, enabled bool) error {
	return h.deps.Gorm.WithContext(ctx).Model(&model.ChannelsIntegrations{}).
		Where(`"channelId" = ? AND "integrationId" IN (SELECT id FROM integrations WHERE service = ?)`, channelID, service).
		Update("enabled", enabled).Error
}

func (h *Handler) addIntegrationTools(s *modelsdk.Server, requestScope scope) {
	channelID := requestScope.Channel.ID.String()
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "get_integration_status", Description: "Get connection status and errors for all supported integrations. Detailed data requires write access."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
			discord, discordErr := h.deps.Discord.GetData(ctx, channelID)
			spotify, spotifyErr := h.deps.Spotify.GetSpotifyData(ctx, channelID)
			lastfm, lastfmErr := h.deps.LastFM.GetData(ctx, channelID)
			valorant, valorantErr := h.deps.Valorant.GetData(ctx, channelID)
			faceit, faceitErr := h.deps.Faceit.GetIntegrationData(ctx, channelID)
			donationAlerts, donationAlertsErr := h.deps.DonationAlerts.GetIntegrationData(ctx, channelID)
			donatePay, donatePayErr := h.deps.DonatePay.GetByChannelID(ctx, channelID)
			streamlabs, streamlabsErr := h.deps.Streamlabs.GetIntegrationData(ctx, channelID)
			vk, vkErr := h.deps.VK.GetIntegrationData(ctx, channelID)
			sevenTV, sevenTVErr := h.deps.SevenTV.GetSevenTvData(ctx, channelID)
			return nil, map[string]integrationResult{
				"discord": integrationStatus(requestScope, discord, discordErr), "spotify": integrationStatus(requestScope, spotify, spotifyErr),
				"lastfm": integrationStatus(requestScope, lastfm, lastfmErr), "valorant": integrationStatus(requestScope, valorant, valorantErr),
				"faceit": integrationStatus(requestScope, faceit, faceitErr), "donationalerts": integrationStatus(requestScope, donationAlerts, donationAlertsErr),
				"donatepay": integrationStatus(requestScope, donatePay, donatePayErr), "donatestream": h.legacyIntegrationStatus(ctx, requestScope, "DONATE_STREAM"),
				"donatello": h.legacyIntegrationStatus(ctx, requestScope, "DONATELLO"), "streamlabs": integrationStatus(requestScope, streamlabs, streamlabsErr),
				"vk": integrationStatus(requestScope, vk, vkErr), "seventv": integrationStatus(requestScope, sevenTV, sevenTVErr),
			}, nil
		})

	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "toggle_integration", Description: "Enable/connect or disable/disconnect an integration. OAuth integrations return an authorization URL when enabling."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, input toggleIntegrationInput) (*modelsdk.CallToolResult, any, error) {
			if input.Enabled {
				switch input.Integration {
				case "discord":
					link, err := h.deps.Discord.GetAuthLink(ctx)
					return nil, map[string]any{"authorizationUrl": link}, err
				case "spotify":
					link, err := h.deps.Spotify.GetAuthLink(ctx)
					return nil, map[string]any{"authorizationUrl": link}, err
				case "lastfm":
					link, err := h.deps.LastFM.GetAuthLink(ctx)
					return nil, map[string]any{"authorizationUrl": link}, err
				case "valorant":
					link, err := h.deps.Valorant.GetAuthLink(ctx)
					return nil, map[string]any{"authorizationUrl": link}, err
				case "faceit":
					link, err := h.deps.Faceit.GetAuthLink(ctx, channelID)
					return nil, map[string]any{"authorizationUrl": link}, err
				case "donationalerts":
					link, err := h.deps.DonationAlerts.GetAuthLink(ctx)
					return nil, map[string]any{"authorizationUrl": link}, err
				case "streamlabs":
					link, err := h.deps.Streamlabs.GetAuthLink(ctx)
					return nil, map[string]any{"authorizationUrl": link}, err
				case "vk":
					link, err := h.deps.VK.GetAuthLink(ctx)
					return nil, map[string]any{"authorizationUrl": link}, err
				case "donatepay":
					err := h.deps.DonatePay.CreateOrUpdate(ctx, channelID, input.Credential, true)
					return nil, map[string]bool{"enabled": err == nil}, err
				case "donatestream":
					err := h.deps.DonateStream.PostCode(ctx, channelID, input.Credential)
					return nil, map[string]bool{"enabled": err == nil}, err
				case "donatello":
					err := h.setLegacyIntegrationEnabled(ctx, channelID, "DONATELLO", true)
					return nil, map[string]bool{"enabled": err == nil}, err
				case "seventv":
					err := h.deps.SevenTV.CreateOrUpdateSevenTvData(ctx, seventv_integration.CreateInput{ChannelID: channelID})
					return nil, map[string]bool{"enabled": err == nil}, err
				default:
					return nil, nil, fmt.Errorf("unsupported integration %q", input.Integration)
				}
			}

			var err error
			switch input.Integration {
			case "discord":
				if input.GuildID == "" {
					err = fmt.Errorf("guildId is required")
				} else {
					err = h.deps.Discord.DisconnectGuild(ctx, channelID, input.GuildID)
				}
			case "spotify":
				err = h.deps.Spotify.Logout(ctx, channelID)
			case "lastfm":
				err = h.deps.LastFM.Logout(ctx, channelID)
			case "valorant":
				err = h.deps.Valorant.Logout(ctx, channelID)
			case "faceit":
				err = h.deps.Faceit.Logout(ctx, channelID)
			case "donationalerts":
				err = h.deps.DonationAlerts.Logout(ctx, channelID)
			case "donatepay":
				err = h.deps.DonatePay.CreateOrUpdate(ctx, channelID, input.Credential, false)
			case "donatestream":
				err = h.setLegacyIntegrationEnabled(ctx, channelID, "DONATE_STREAM", false)
			case "donatello":
				err = h.setLegacyIntegrationEnabled(ctx, channelID, "DONATELLO", false)
			case "streamlabs":
				err = h.deps.Streamlabs.Logout(ctx, channelID)
			case "vk":
				err = h.deps.VK.Logout(ctx, channelID)
			case "seventv":
				err = h.setLegacyIntegrationEnabled(ctx, channelID, "SEVENTV", false)
			default:
				err = fmt.Errorf("unsupported integration %q", input.Integration)
			}
			return nil, map[string]bool{"enabled": false}, err
		})
}
