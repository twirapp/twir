package mappers

import (
	model "github.com/satont/twir/libs/gomodels"
	pubsubauditlogs "github.com/satont/twir/libs/pubsub/audit-logs"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	buscoreauditlogs "github.com/twirapp/twir/libs/bus-core/audit-logs"
)

func AuditLogToGql(auditLog pubsubauditlogs.AuditLog) *gqlmodel.AuditLog {
	return &gqlmodel.AuditLog{
		ID:            auditLog.ID,
		System:        AuditTableNameToGqlSystem(auditLog.Table),
		OperationType: AuditLogOperationTypeToGql(auditLog.OperationType),
		OldValue:      auditLog.OldValue.Ptr(),
		NewValue:      auditLog.NewValue.Ptr(),
		ObjectID:      auditLog.ObjectID.Ptr(),
		UserID:        auditLog.UserID.Ptr(),
		CreatedAt:     auditLog.CreatedAt,
	}
}

func AuditLogOperationTypeToGql(t buscoreauditlogs.AuditOperationType) gqlmodel.AuditOperationType {
	switch t {
	case buscoreauditlogs.AuditOperationTypeUpdate:
		return gqlmodel.AuditOperationTypeUpdate
	case buscoreauditlogs.AuditOperationTypeCreate:
		return gqlmodel.AuditOperationTypeCreate
	case buscoreauditlogs.AuditOperationTypeDelete:
		return gqlmodel.AuditOperationTypeDelete
	default:
		return ""
	}
}

func AuditTypeModelToGql(t model.AuditOperationType) gqlmodel.AuditOperationType {
	switch t {
	case model.AuditOperationUpdate:
		return gqlmodel.AuditOperationTypeUpdate
	case model.AuditOperationCreate:
		return gqlmodel.AuditOperationTypeCreate
	case model.AuditOperationDelete:
		return gqlmodel.AuditOperationTypeDelete
	default:
		return ""
	}
}

func AuditTypeGqlToModel(t gqlmodel.AuditOperationType) model.AuditOperationType {
	switch t {
	case gqlmodel.AuditOperationTypeUpdate:
		return model.AuditOperationUpdate
	case gqlmodel.AuditOperationTypeCreate:
		return model.AuditOperationCreate
	case gqlmodel.AuditOperationTypeDelete:
		return model.AuditOperationDelete
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
