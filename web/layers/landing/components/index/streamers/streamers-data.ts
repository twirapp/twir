import { useLandingStats } from '~~/layers/landing/api/stats'

export async function getStreamers() {
	const { data, error } = await useLandingStats()
	if (error.value) {
		console.error(error.value)
		return []
	}

	const sortedStreamers = data.value?.twirStats.streamers.sort((a, b) => b.followersCount - a.followersCount)
	if (!sortedStreamers) return []

	const streamers: typeof sortedStreamers[number][][] = []

	if (import.meta.dev) {
		streamers.push(...chunk(Array.from({ length: 100 }).map(() => sortedStreamers.at(0)!), 3))
	} else {
		streamers.push(...chunk(sortedStreamers, 3))
	}

	return streamers
}

function chunk<T>(arr: T[], size: number): T[][] {
	const result: T[][] = []

	for (let i = 0; i < arr.length; i += size) {
		result.push(arr.slice(i, i + size))
	}

	return result
}
