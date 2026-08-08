package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	dotaentity "github.com/twirapp/twir/libs/entities/dota"
)

func MapDotaEntityToGQL(s dotaentity.ChannelDotaSettings) gqlmodel.DotaSettings {
	return gqlmodel.DotaSettings{
		ID:             s.ID,
		ChannelID:      s.ChannelID.String(),
		Enabled:        s.Enabled,
		SteamAccountID: s.SteamAccountID,
		GsiToken:       s.GsiToken,
		Mmr:            s.Mmr,
		MmrDelta:       s.MmrDelta,
		SessionWins:    s.SessionWins,
		SessionLosses:  s.SessionLosses,
		PredictionSettings: &gqlmodel.DotaPredictionSettings{
			Enabled:       s.PredictionSettings.Enabled,
			TitleTemplate: s.PredictionSettings.TitleTemplate,
			WindowSeconds: s.PredictionSettings.WindowSeconds,
		},
		ChatEvents: &gqlmodel.DotaChatEvents{
			MatchStarted: mapDotaChatEventToGQL(s.ChatEvents.MatchStarted),
			MatchEnded:   mapDotaChatEventToGQL(s.ChatEvents.MatchEnded),
			RoshanKilled: mapDotaChatEventToGQL(s.ChatEvents.RoshanKilled),
			AegisPickup:  mapDotaChatEventToGQL(s.ChatEvents.AegisPickup),
		},
		CommandsSettings: &gqlmodel.DotaCommandsSettings{
			Mmr: s.CommandsSettings.Mmr,
			Wl:  s.CommandsSettings.Wl,
			Lg:  s.CommandsSettings.Lg,
			Gm:  s.CommandsSettings.Gm,
			Np:  s.CommandsSettings.Np,
			Wp:  s.CommandsSettings.Wp,
		},
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func mapDotaChatEventToGQL(e dotaentity.ChatEventSettings) *gqlmodel.DotaChatEventSettings {
	return &gqlmodel.DotaChatEventSettings{
		Enabled:  e.Enabled,
		Template: e.Template,
		Cooldown: e.Cooldown,
	}
}

func MapDotaUpdateInputToEntity(input gqlmodel.DotaUpdateInput) dotaentity.ChannelDotaSettings {
	result := dotaentity.ChannelDotaSettings{
		Enabled:  input.Enabled,
		Mmr:      input.Mmr,
		MmrDelta: input.MmrDelta,
	}

	if input.PredictionSettings != nil {
		result.PredictionSettings = dotaentity.PredictionSettings{
			Enabled:       input.PredictionSettings.Enabled,
			TitleTemplate: input.PredictionSettings.TitleTemplate,
			WindowSeconds: input.PredictionSettings.WindowSeconds,
		}
	}

	if input.ChatEvents != nil {
		result.ChatEvents = dotaentity.ChatEvents{
			MatchStarted: mapDotaChatEventInputToEntity(input.ChatEvents.MatchStarted),
			MatchEnded:   mapDotaChatEventInputToEntity(input.ChatEvents.MatchEnded),
			RoshanKilled: mapDotaChatEventInputToEntity(input.ChatEvents.RoshanKilled),
			AegisPickup:  mapDotaChatEventInputToEntity(input.ChatEvents.AegisPickup),
		}
	}

	if input.CommandsSettings != nil {
		result.CommandsSettings = dotaentity.CommandsSettings{
			Mmr: input.CommandsSettings.Mmr,
			Wl:  input.CommandsSettings.Wl,
			Lg:  input.CommandsSettings.Lg,
			Gm:  input.CommandsSettings.Gm,
			Np:  input.CommandsSettings.Np,
			Wp:  input.CommandsSettings.Wp,
		}
	}

	return result
}

func mapDotaChatEventInputToEntity(e *gqlmodel.DotaChatEventSettingsInput) dotaentity.ChatEventSettings {
	if e == nil {
		return dotaentity.ChatEventSettings{}
	}

	return dotaentity.ChatEventSettings{
		Enabled:  e.Enabled,
		Template: e.Template,
		Cooldown: e.Cooldown,
	}
}
