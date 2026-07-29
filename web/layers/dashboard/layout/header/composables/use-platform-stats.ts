import { tryOnScopeDispose, useIntervalFn } from '@vueuse/core'
import { intervalToDuration } from 'date-fns'
import { computed, ref, toValue } from 'vue'

import type { MaybeRefOrGetter } from 'vue'

import { padTo2Digits } from '~~/layers/dashboard/helpers/convertMillisToTime'
import { Platform } from '~/gql/graphql.js'

import type { DashboardStatsSubscription } from '~/gql/graphql.js'

export type DashboardStats = DashboardStatsSubscription['dashboardStats']
export type DashboardPlatformStats = DashboardStats['platforms'][number]

const platformSortOrder: Platform[] = [Platform.Twitch, Platform.Kick, Platform.VkVideoLive]

export function formatStreamUptime(startedAt: string | Date | null | undefined, now: Date): string {
	if (!startedAt) return '00:00:00'

	const duration = intervalToDuration({
		start: new Date(startedAt),
		end: now,
	})

	const mappedDuration = [duration.hours ?? 0, duration.minutes ?? 0, duration.seconds ?? 0]
	if (duration.days !== undefined && duration.days !== 0) mappedDuration.unshift(duration.days)

	return mappedDuration.map((value) => padTo2Digits(value)).join(':')
}

export function usePlatformStats(stats: MaybeRefOrGetter<DashboardStats | null | undefined>) {
	const platforms = computed<DashboardPlatformStats[]>(() => toValue(stats)?.platforms ?? [])

	const sortedPlatforms = computed<DashboardPlatformStats[]>(() => {
		return [...platforms.value].sort((a, b) => {
			if (a.isLive !== b.isLive) return a.isLive ? -1 : 1
			return platformSortOrder.indexOf(a.platform) - platformSortOrder.indexOf(b.platform)
		})
	})

	const livePlatformsCount = computed(() => platforms.value.filter((p) => p.isLive).length)

	const totalViewers = computed(() => {
		return platforms.value.reduce((sum, platform) => {
			if (!platform.isLive) return sum
			return sum + (platform.viewers ?? 0)
		}, 0)
	})

	const totalFollowers = computed(() => {
		return platforms.value.reduce((sum, platform) => sum + (platform.followers ?? 0), 0)
	})

	const primaryPlatform = computed<DashboardPlatformStats | null>(() => {
		return sortedPlatforms.value.find((p) => p.isLive) ?? sortedPlatforms.value[0] ?? null
	})

	const now = ref(new Date())
	const { pause: pauseUptimeTicker } = useIntervalFn(() => {
		now.value = new Date()
	}, 1000)
	tryOnScopeDispose(pauseUptimeTicker)

	const uptimes = computed<Partial<Record<Platform, string>>>(() => {
		const result: Partial<Record<Platform, string>> = {}
		for (const platform of platforms.value) {
			result[platform.platform] =
				platform.isLive && platform.startedAt
					? formatStreamUptime(platform.startedAt, now.value)
					: '—'
		}
		return result
	})

	const globalUptime = computed(() => {
		return formatStreamUptime(toValue(stats)?.startedAt, now.value)
	})

	return {
		platforms,
		sortedPlatforms,
		livePlatformsCount,
		totalViewers,
		totalFollowers,
		primaryPlatform,
		uptimes,
		globalUptime,
		now,
	}
}
