import { computed, toValue } from 'vue'

import type { MaybeRefOrGetter } from 'vue'

import { counterColors } from '@/components/stream-stats/counter-meta.js'
import { Platform, StreamStatsOverlayViewersMode } from '@/gql/graphql'

import type {
	StreamStatsOverlayCountersSubscription,
	StreamStatsOverlaySettingsSubscription,
} from '@/gql/graphql'

export type StreamStatsSettings = StreamStatsOverlaySettingsSubscription['overlaysStreamStats']
export type StreamStatsCounters = StreamStatsOverlayCountersSubscription['overlaysStreamStatsCounters']

export type StreamStatsCounterKey = 'viewers' | 'messages' | 'uptime' | 'subscribers' | 'followers'

export type StreamStatsCounterItem = {
	id: string
	key: StreamStatsCounterKey
	label: string
	value: string
	color: string
	platform?: Platform
	platformColor?: string
}

const platformMeta: Partial<Record<Platform, { label: string, color: string }>> = {
	[Platform.Twitch]: { label: 'Twitch', color: '#9146FF' },
	[Platform.Kick]: { label: 'Kick', color: '#53FC18' },
	[Platform.VkVideoLive]: { label: 'VK', color: '#0077FF' },
	[Platform.Youtube]: { label: 'YouTube', color: '#FF0000' },
}

const numberFormatter = new Intl.NumberFormat('en-US')

export function formatStatNumber(value: number): string {
	return numberFormatter.format(value)
}

// 1:23:45 when hours present, 12:03 under an hour
export function formatUptime(totalSeconds: number): string {
	const hours = Math.floor(totalSeconds / 3600)
	const minutes = Math.floor((totalSeconds % 3600) / 60)
	const seconds = totalSeconds % 60
	const mm = minutes.toString().padStart(2, '0')
	const ss = seconds.toString().padStart(2, '0')
	return hours > 0 ? `${hours}:${mm}:${ss}` : `${mm}:${ss}`
}

export function useStreamStatsCounters(
	settings: MaybeRefOrGetter<StreamStatsSettings | null>,
	counters: MaybeRefOrGetter<StreamStatsCounters | null>,
	now: MaybeRefOrGetter<Date>,
) {
	const uptime = computed(() => {
		const currentCounters = toValue(counters)
		if (!currentCounters?.startedAt) return null

		const startedAt = new Date(currentCounters.startedAt).getTime()
		if (Number.isNaN(startedAt)) return null

		const totalSeconds = Math.max(0, Math.floor((toValue(now).getTime() - startedAt) / 1000))
		return formatUptime(totalSeconds)
	})

	const items = computed<StreamStatsCounterItem[]>(() => {
		const currentSettings = toValue(settings)
		const currentCounters = toValue(counters)
		if (!currentSettings || !currentCounters) return []

		const userColors: Record<StreamStatsCounterKey, string> = {
			viewers: currentSettings.viewersColor,
			messages: currentSettings.messagesColor,
			uptime: currentSettings.uptimeColor,
			subscribers: currentSettings.subscribersColor,
			followers: currentSettings.followersColor,
		}
		const resolveColor = (key: StreamStatsCounterKey): string =>
			userColors[key] || counterColors[key]

		const result: StreamStatsCounterItem[] = []

		if (currentSettings.viewersEnabled) {
			if (
				currentSettings.viewersMode === StreamStatsOverlayViewersMode.Separate &&
				currentCounters.platformViewers.length > 0
			) {
				for (const entry of currentCounters.platformViewers) {
					const meta = platformMeta[entry.platform]
					result.push({
						id: `viewers-${entry.platform}`,
						key: 'viewers',
						label: meta?.label ?? entry.platform,
						value: formatStatNumber(entry.viewers),
						color: resolveColor('viewers'),
						platform: entry.platform,
						...(meta ? { platformColor: meta.color } : {}),
					})
				}
			} else {
				result.push({
					id: 'viewers',
					key: 'viewers',
					label: 'Viewers',
					value: formatStatNumber(currentCounters.viewers),
					color: resolveColor('viewers'),
				})
			}
		}

		if (currentSettings.messagesEnabled) {
			result.push({
				id: 'messages',
				key: 'messages',
				label: 'Messages',
				value: formatStatNumber(currentCounters.messages),
				color: resolveColor('messages'),
			})
		}

		if (currentSettings.uptimeEnabled && uptime.value) {
			result.push({
				id: 'uptime',
				key: 'uptime',
				label: 'Uptime',
				value: uptime.value,
				color: resolveColor('uptime'),
			})
		}

		if (currentSettings.subscribersEnabled && currentCounters.subscribers != null) {
			result.push({
				id: 'subscribers',
				key: 'subscribers',
				label: 'Subscribers',
				value: formatStatNumber(currentCounters.subscribers),
				color: resolveColor('subscribers'),
			})
		}

		if (currentSettings.followersEnabled && currentCounters.followers != null) {
			result.push({
				id: 'followers',
				key: 'followers',
				label: 'Followers',
				value: formatStatNumber(currentCounters.followers),
				color: resolveColor('followers'),
			})
		}

		return result
	})

	const placeholders = computed<Record<string, string>>(() => {
		const currentCounters = toValue(counters)
		return {
			viewers: currentCounters ? formatStatNumber(currentCounters.viewers) : '0',
			messages: currentCounters ? formatStatNumber(currentCounters.messages) : '0',
			uptime: uptime.value ?? '0:00',
			subscribers:
				currentCounters?.subscribers != null ? formatStatNumber(currentCounters.subscribers) : '',
			followers:
				currentCounters?.followers != null ? formatStatNumber(currentCounters.followers) : '',
		}
	})

	return {
		items,
		placeholders,
		uptime,
	}
}
