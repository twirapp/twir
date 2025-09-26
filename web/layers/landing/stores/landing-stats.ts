import { useLandingStats } from '#layers/landing/api/stats'

export const LandingStatsStoreKey = 'landing-stats'
export const useLandingStatsStore = defineStore(LandingStatsStoreKey, () => {
	const { data, executeQuery } = useLandingStats()

	async function fetchLandingStats() {
		await executeQuery()
	}

	return {
		fetchLandingStats,
		stats: computed(() => data.value?.twirStats ?? null),
	}
})
