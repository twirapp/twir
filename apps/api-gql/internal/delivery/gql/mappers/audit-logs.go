package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func AuditLogToGql(auditLog entity.AuditLog) gqlmodel.AuditLog {
	return gqlmodel.AuditLog{
		ID:            auditLog.ID,
		System:        AuditTableNameToGqlSystem(auditLog.TableName),
		OperationType: AuditTypeModelToGql(auditLog.OperationType),
		OldValue:      auditLog.OldValue,
		NewValue:      auditLog.NewValue,
		ObjectID:      auditLog.ObjectID,
		UserID:        auditLog.UserID,
		CreatedAt:     auditLog.CreatedAt,
	}
}

func AuditTypeModelToGql(t entity.AuditOperationType) gqlmodel.AuditOperationType {
	switch t {
	case entity.AuditOperationUpdate:
		return gqlmodel.AuditOperationTypeUpdate
	case entity.AuditOperationCreate:
		return gqlmodel.AuditOperationTypeCreate
	case entity.AuditOperationDelete:
		return gqlmodel.AuditOperationTypeDelete
	default:
		return ""
	}
}

func AuditTypeGqlToModel(t gqlmodel.AuditOperationType) entity.AuditOperationType {
	switch t {
	case gqlmodel.AuditOperationTypeUpdate:
		return entity.AuditOperationUpdate
	case gqlmodel.AuditOperationTypeCreate:
		return entity.AuditOperationCreate
	case gqlmodel.AuditOperationTypeDelete:
		return entity.AuditOperationDelete
	default:
		return ""
	}
}

var tableToGqlModel = map[string]gqlmodel.AuditLogSystem{
	"badges":                          gqlmodel.AuditLogSystemBadge,
	"badge_users":                     gqlmodel.AuditLogSystemBadgeUser,
	"channels_commands":               gqlmodel.AuditLogSystemChannelCommand,
	"channels_command_groups":         gqlmodel.AuditLogSystemChannelCommandGroup,
	"channels_customvars":             gqlmodel.AuditLogSystemChannelVariable,
	"channels_games_8ball":            gqlmodel.AuditLogSystemChannelGamesEightBall,
	"channels_games_duel":             gqlmodel.AuditLogSystemChannelGamesDuel,
	"channels_games_russian_roulette": gqlmodel.AuditLogSystemChannelGamesRussianRoulette,
	"channels_games_seppuku":          gqlmodel.AuditLogSystemChannelGamesSeppuku,
	"channels_games_voteban":          gqlmodel.AuditLogSystemChannelGamesVoteban,
	"channels_greetings":              gqlmodel.AuditLogSystemChannelGreeting,
	"channels_keywords":               gqlmodel.AuditLogSystemChannelKeyword,
	"channels_moderation_settings":    gqlmodel.AuditLogSystemChannelModerationSetting,
	"channels_overlays_chat":          gqlmodel.AuditLogSystemChannelOverlayChat,
	"channels_overlays_dudes":         gqlmodel.AuditLogSystemChannelOverlayDudes,
	"channels_overlays_now_playing":   gqlmodel.AuditLogSystemChannelOverlayNowPlaying,
	"channels_roles":                  gqlmodel.AuditLogSystemChannelRoles,
	"channels_timers":                 gqlmodel.AuditLogSystemChannelTimers,
	"channels_song_requests_settings": gqlmodel.AuditLogSystemChannelSongRequests,
	"channels_integrations":           gqlmodel.AuditLogSystemChannelIntegrations,
	"channels_alerts":                 gqlmodel.AuditLogSystemChannelsAlerts,
	"channels_chat_alerts":            gqlmodel.AuditLogSystemChannelsChatAlerts,
	"channels_chat_translation":       gqlmodel.AuditLogSystemChannelsChatTranslation,
}

func AuditTableNameToGqlSystem(t string) gqlmodel.AuditLogSystem {
	return tableToGqlModel[t]
}

func AuditSystemToTableName(s gqlmodel.AuditLogSystem) string {
	for k, v := range tableToGqlModel {
		if v == s {
			return k
		}
	}
	return ""
}
