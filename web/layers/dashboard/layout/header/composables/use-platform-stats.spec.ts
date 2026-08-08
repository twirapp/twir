import { ref } from 'vue'
import { describe, expect, it } from 'vitest'

import { Platform } from '~/gql/graphql.js'
import { usePlatformStats } from './use-platform-stats.js'

import type { DashboardPlatformStats, DashboardStats } from './use-platform-stats.js'

function makePlatform(overrides: Partial<DashboardPlatformStats>): DashboardPlatformStats {
	return {
		platform: Platform.Twitch,
		isLive: false,
		title: null,
		categoryId: null,
		categoryName: null,
		viewers: null,
		followers: null,
		startedAt: null,
		chatMessages: 0,
		usedEmotes: 0,
		canEditInfo: false,
		...overrides,
	}
}

function makeStats(platforms: DashboardPlatformStats[]): DashboardStats {
	return {
		categoryId: '',
		categoryName: '',
		viewers: null,
		startedAt: null,
		title: '',
		chatMessages: 0,
		followers: 0,
		usedEmotes: 0,
		requestedSongs: 0,
		subs: 0,
		platforms,
	}
}

describe('usePlatformStats', () => {
	it('sums viewers only across live platforms', () => {
		const stats = ref(
			makeStats([
				makePlatform({ platform: Platform.Twitch, isLive: true, viewers: 21 }),
				makePlatform({ platform: Platform.Kick, isLive: true, viewers: 12 }),
				makePlatform({ platform: Platform.VkVideoLive, isLive: false, viewers: 500 }),
			])
		)

		const { totalViewers } = usePlatformStats(stats)

		expect(totalViewers.value).toBe(33)
	})

	it('treats null viewers of a live platform as zero', () => {
		const stats = ref(
			makeStats([
				makePlatform({ platform: Platform.Twitch, isLive: true, viewers: null }),
				makePlatform({ platform: Platform.Kick, isLive: true, viewers: 12 }),
			])
		)

		const { totalViewers } = usePlatformStats(stats)

		expect(totalViewers.value).toBe(12)
	})

	it('sums followers skipping null values', () => {
		const stats = ref(
			makeStats([
				makePlatform({ platform: Platform.Twitch, isLive: true, followers: 126 }),
				makePlatform({ platform: Platform.Kick, isLive: true, followers: 48 }),
				makePlatform({ platform: Platform.VkVideoLive, isLive: false, followers: null }),
			])
		)

		const { totalFollowers } = usePlatformStats(stats)

		expect(totalFollowers.value).toBe(174)
	})

	it('sorts platforms live-first, then by platform order', () => {
		const stats = ref(
			makeStats([
				makePlatform({ platform: Platform.VkVideoLive, isLive: false }),
				makePlatform({ platform: Platform.Kick, isLive: true }),
				makePlatform({ platform: Platform.Twitch, isLive: false }),
				makePlatform({ platform: Platform.Twitch, isLive: true }),
			])
		)

		const { sortedPlatforms } = usePlatformStats(stats)

		expect(sortedPlatforms.value.map((p) => [p.platform, p.isLive])).toEqual([
			[Platform.Twitch, true],
			[Platform.Kick, true],
			[Platform.Twitch, false],
			[Platform.VkVideoLive, false],
		])
	})

	it('picks the first live platform as primary', () => {
		const stats = ref(
			makeStats([
				makePlatform({ platform: Platform.VkVideoLive, isLive: false }),
				makePlatform({ platform: Platform.Kick, isLive: true }),
				makePlatform({ platform: Platform.Twitch, isLive: true }),
			])
		)

		const { primaryPlatform } = usePlatformStats(stats)

		expect(primaryPlatform.value?.platform).toBe(Platform.Twitch)
	})

	it('falls back to the first platform when none are live', () => {
		const stats = ref(
			makeStats([
				makePlatform({ platform: Platform.Kick, isLive: false }),
				makePlatform({ platform: Platform.VkVideoLive, isLive: false }),
			])
		)

		const { primaryPlatform } = usePlatformStats(stats)

		expect(primaryPlatform.value?.platform).toBe(Platform.Kick)
	})

	it('returns null primary and zero totals for empty stats', () => {
		const stats = ref<DashboardStats | undefined>(undefined)

		const { primaryPlatform, totalViewers, totalFollowers, sortedPlatforms } =
			usePlatformStats(stats)

		expect(primaryPlatform.value).toBeNull()
		expect(totalViewers.value).toBe(0)
		expect(totalFollowers.value).toBe(0)
		expect(sortedPlatforms.value).toEqual([])
	})
})
