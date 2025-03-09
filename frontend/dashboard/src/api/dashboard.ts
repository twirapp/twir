import { useSubscription } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { computed } from 'vue'

import type { DashboardEventsSubscription } from '@/gql/graphql'

import { useMutation } from '@/composables/use-mutation.ts'
import { graphql } from '@/gql'

export type DashboardWidgetEvent = DashboardEventsSubscription['dashboardWidgetsEvents']['events'][number]

export const useDashboardEvents = createGlobalState(() => {
	const { data, isPaused, fetching } = useSubscription({
		query: graphql(`
			subscription DashboardEvents {
				dashboardWidgetsEvents {
					events {
						userId
						data {
							donationAmount
							donationCurrency
							donationMessage
							donationUserName
							raidedViewersCount
							raidedFromUserName
							raidedFromDisplayName
							followUserName
							followUserDisplayName
							redemptionTitle
							redemptionInput
							redemptionUserName
							redemptionUserDisplayName
							redemptionCost
							subLevel
							subUserName
							subUserDisplayName
							reSubLevel
							reSubUserName
							reSubUserDisplayName
							reSubMonths
							reSubStreak
							subGiftLevel
							subGiftUserName
							subGiftUserDisplayName
							subGiftTargetUserName
							subGiftTargetUserDisplayName
							firstUserMessageUserName
							firstUserMessageUserDisplayName
							firstUserMessageMessage
							banReason
							banEndsInMinutes
							bannedUserName
							bannedUserLogin
							moderatorName
							moderatorDisplayName
							message
							userLogin
							userName
						}
						createdAt
						type
					}
				}
			}
		`),
	})

	const events = computed(() => {
		return data.value?.dashboardWidgetsEvents.events
	})

	return {
		events,
		isPaused,
		fetching,
	}
})

export function useRealtimeDashboardStats() {
	const { data, isPaused, fetching } = useSubscription({
		query: graphql(`
			subscription dashboardStats {
				dashboardStats {
					categoryId
					categoryName
					viewers
					startedAt
					title
					chatMessages
					followers
					usedEmotes
					requestedSongs
					subs
				}
			}
		`),
	})

	const stats = computed(() => {
		return data.value?.dashboardStats
	})

	return { stats, isPaused, fetching }
}

export function useBotStatus() {
	const { data, isPaused, fetching, executeSubscription } = useSubscription({
		query: graphql(`
			subscription botStatus {
				botStatus {
					isMod
					botId
					botName
					enabled
				}
			}
		`),
		context: {
			additionalTypenames: ['BotStatus'],
		},
	})

	const botStatus = computed(() => {
		return data.value?.botStatus
	})

	return { botStatus, isPaused, fetching, executeSubscription }
}

export function useBotJoinPart() {
	return useMutation(graphql(`
		mutation BotJoinPart($action: BotJoinLeaveAction!) {
			botJoinLeave(action: $action)
		}
	`), ['BotStatus'])
}
