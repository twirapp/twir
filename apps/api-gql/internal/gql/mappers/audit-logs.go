package mappers

import (
	model "github.com/satont/twir/libs/gomodels"
	pubsubauditlogs "github.com/satont/twir/libs/pubsub/audit-logs"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
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

func AuditTableNameToGqlSystem(t string) gqlmodel.AuditLogSystem {
	switch t {
	case "badges":
		return gqlmodel.AuditLogSystemBadge
	case "badge_users":
		return gqlmodel.AuditLogSystemBadgeUser
	case "channels_commands":
		return gqlmodel.AuditLogSystemChannelCommand
	case "channels_command_groups":
		return gqlmodel.AuditLogSystemChannelCommandGroup
	case "channels_customvars":
		return gqlmodel.AuditLogSystemChannelVariable
	case "channels_games_8ball":
		return gqlmodel.AuditLogSystemChannelGamesEightBall
	case "channels_games_duel":
		return gqlmodel.AuditLogSystemChannelGamesDuel
	case "channels_games_russian_roulette":
		return gqlmodel.AuditLogSystemChannelGamesRussianRoulette
	case "channels_games_seppuku":
		return gqlmodel.AuditLogSystemChannelGamesSeppuku
	case "channels_games_voteban":
		return gqlmodel.AuditLogSystemChannelGamesVoteban
	case "channels_greetings":
		return gqlmodel.AuditLogSystemChannelGreeting
	case "channels_keywords":
		return gqlmodel.AuditLogSystemChannelKeyword
	case "channels_moderation_settings":
		return gqlmodel.AuditLogSystemChannelModerationSetting
	case "channels_overlays_chat":
		return gqlmodel.AuditLogSystemChannelOverlayChat
	case "channels_overlays_dudes":
		return gqlmodel.AuditLogSystemChannelOverlayDudes
	case "channels_overlays_now_playing":
		return gqlmodel.AuditLogSystemChannelOverlayNowPlaying
	case "channels_roles":
		return gqlmodel.AuditLogSystemChannelRoles
	case "channels_timers":
		return gqlmodel.AuditLogSystemChannelTimers
	case "channels_song_requests_settings":
		return gqlmodel.AuditLogSystemChannelSongRequests
	case "channels_integrations":
		return gqlmodel.AuditLogSystemChannelIntegrations
	default:
		return ""
	}
}

func AuditSystemToTableName(s gqlmodel.AuditLogSystem) string {
	switch s {
	case gqlmodel.AuditLogSystemBadge:
		return "badges"
	case gqlmodel.AuditLogSystemBadgeUser:
		return "badge_users"
	case gqlmodel.AuditLogSystemChannelCommand:
		return "channels_commands"
	case gqlmodel.AuditLogSystemChannelCommandGroup:
		return "channels_command_groups"
	case gqlmodel.AuditLogSystemChannelVariable:
		return "channels_customvars"
	case gqlmodel.AuditLogSystemChannelGamesEightBall:
		return "channels_games_8ball"
	case gqlmodel.AuditLogSystemChannelGamesDuel:
		return "channels_games_duel"
	case gqlmodel.AuditLogSystemChannelGamesRussianRoulette:
		return "channels_games_russian_roulette"
	case gqlmodel.AuditLogSystemChannelGamesSeppuku:
		return "channels_games_seppuku"
	case gqlmodel.AuditLogSystemChannelGamesVoteban:
		return "channels_games_voteban"
	case gqlmodel.AuditLogSystemChannelGreeting:
		return "channels_greetings"
	case gqlmodel.AuditLogSystemChannelKeyword:
		return "channels_keywords"
	case gqlmodel.AuditLogSystemChannelModerationSetting:
		return "channels_moderation_settings"
	case gqlmodel.AuditLogSystemChannelOverlayChat:
		return "channels_overlays_chat"
	case gqlmodel.AuditLogSystemChannelOverlayDudes:
		return "channels_overlays_dudes"
	case gqlmodel.AuditLogSystemChannelOverlayNowPlaying:
		return "channels_overlays_now_playing"
	case gqlmodel.AuditLogSystemChannelRoles:
		return "channels_roles"
	case gqlmodel.AuditLogSystemChannelTimers:
		return "channels_timers"
	case gqlmodel.AuditLogSystemChannelSongRequests:
		return "channels_song_requests_settings"
	case gqlmodel.AuditLogSystemChannelIntegrations:
		return "channels_integrations"
	default:
		return ""
	}
}
