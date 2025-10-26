import { useLandingStats } from '#layers/landing/api/stats'

export const LandingStatsStoreKey = 'landing-stats-key'
export const useLandingStatsStore = defineStore('landing-stats', () => {
	const { data, executeQuery } = useLandingStats()
	const stats = computed(() => data.value?.twirStats ?? null)

	async function fetchLandingStats() {
		const { data: newData } = await executeQuery()
		return newData.value
	}

	return {
		fetchLandingStats,
		stats,
	}
})
