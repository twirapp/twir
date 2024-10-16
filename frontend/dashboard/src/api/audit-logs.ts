import { useQuery, useSubscription } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { ref, watch } from 'vue'

import type { UserLogFragmentFragment } from '@/gql/graphql'

import { graphql } from '@/gql'
import { AuditLogSystem, AuditOperationType } from '@/gql/graphql'

graphql(`
	fragment UserLogFragment on AuditLog {
      id
      createdAt
      objectId
      oldValue
      newValue
      operationType
      system
      user {
          id
          profileImageUrl
          displayName
          login
      }
	}
`)

export const useAuditLogs = createGlobalState(() => {
	const logs = ref<UserLogFragmentFragment[]>([])

	const { data: fetchedLogs, fetching } = useQuery({
		query: graphql(`
			query UserAuditLogs {
					auditLog {
							...UserLogFragment
          }
			}
		`),
		variables: {},
	})

	watch(fetchedLogs, (newLogs) => {
		if (!newLogs?.auditLog) return

		logs.value = newLogs.auditLog as UserLogFragmentFragment[]
	}, { immediate: true })

	const { data: realtimeLog } = useSubscription({
		query: graphql(`
			subscription UserAuditLogsSubscription {
					auditLog {
							...UserLogFragment
          }
			}
		`),
	})

	watch(realtimeLog, (newLog) => {
		if (!newLog?.auditLog) return

		logs.value.unshift(newLog.auditLog as UserLogFragmentFragment)
		console.log(logs.value)
	}, { immediate: true })

	return {
		logs,
		isLoading: fetching,
	}
})

export function mapSystemToTranslate(system: AuditLogSystem) {
	switch (system) {
		case AuditLogSystem.Badge:
			return 'dashboard.widgets.audit-logs.systems.badge'
		case AuditLogSystem.BadgeUser:
			return 'dashboard.widgets.audit-logs.systems.badge-user'
		case AuditLogSystem.ChannelCommand:
			return 'dashboard.widgets.audit-logs.systems.channel-command'
		case AuditLogSystem.ChannelCommandGroup:
			return 'dashboard.widgets.audit-logs.systems.channel-command-group'
		case AuditLogSystem.ChannelVariable:
			return 'dashboard.widgets.audit-logs.systems.channel-variable'
		case AuditLogSystem.ChannelGamesEightBall:
			return 'dashboard.widgets.audit-logs.systems.channel-games-eight-ball'
		case AuditLogSystem.ChannelGamesDuel:
			return 'dashboard.widgets.audit-logs.systems.channel-games-duel'
		case AuditLogSystem.ChannelGamesRussianRoulette:
			return 'dashboard.widgets.audit-logs.systems.channel-games-russian-roulette'
		case AuditLogSystem.ChannelGamesSeppuku:
			return 'dashboard.widgets.audit-logs.systems.channel-games-seppuku'
		case AuditLogSystem.ChannelGamesVoteban:
			return 'dashboard.widgets.audit-logs.systems.channel-games-voteban'
		case AuditLogSystem.ChannelGreeting:
			return 'dashboard.widgets.audit-logs.systems.channel-greeting'
		case AuditLogSystem.ChannelKeyword:
			return 'dashboard.widgets.audit-logs.systems.channel-keyword'
		case AuditLogSystem.ChannelModerationSetting:
			return 'dashboard.widgets.audit-logs.systems.channel-moderation-setting'
		case AuditLogSystem.ChannelOverlayChat:
			return 'dashboard.widgets.audit-logs.systems.channel-overlay-chat'
		case AuditLogSystem.ChannelOverlayDudes:
			return 'dashboard.widgets.audit-logs.systems.channel-overlay-dudes'
		case AuditLogSystem.ChannelOverlayNowPlaying:
			return 'dashboard.widgets.audit-logs.systems.channel-overlay-now-playing'
		case AuditLogSystem.ChannelRoles:
			return 'dashboard.widgets.audit-logs.systems.channel-roles'
		case AuditLogSystem.ChannelTimers:
			return 'dashboard.widgets.audit-logs.systems.channel-timers'
		case AuditLogSystem.ChannelSongRequests:
			return 'dashboard.widgets.audit-logs.systems.channel-song-requests'
		case AuditLogSystem.ChannelIntegrations:
			return 'dashboard.widgets.audit-logs.systems.channel-integrations'
	}
}

export function mapOperationTypeToTranslate(operationType: AuditOperationType) {
	switch (operationType) {
		case AuditOperationType.Create:
			return 'dashboard.widgets.audit-logs.operation-type.create'
		case AuditOperationType.Update:
			return 'dashboard.widgets.audit-logs.operation-type.update'
		case AuditOperationType.Delete:
			return 'dashboard.widgets.audit-logs.operation-type.delete'
	}
}
