/* eslint-disable ts/ban-ts-comment */
import type { FaceitLastGamesResponse, FaceitPlayerStatsResponse, FaceitProfileResponse } from './types.js'

export interface Stats {
	lvl: number
	elo: number
	avgKdr?: string
	worldRanking?: string
	lastMatches?: {
		avgKills: string
		headshots: string
		winRate: string
		avgKd: string
		avgKr: string
	}
}

const requestOptions = {
	headers: {
		Authorization: `Bearer 5dbe323f-3fb4-4bd6-8b9f-b5688d63ebee`,
	},
}

export async function computeStats(nickname: string, game: string): Promise<Stats | null> {
	let result: Stats | null = null

	const profileResponse = await fetch(`https://open.faceit.com/data/v4/players?nickname=${nickname}&game=${game}`, requestOptions)
	if (!profileResponse.ok) {
		throw new Error(await profileResponse.text())
	}

	const profile = await profileResponse.json() as FaceitProfileResponse
	const profileGame = profile.games[game]
	if (!profileGame) return null

	const [lastGamesResponse, playerStatsResponse, positionResponse] = await Promise.all([
		fetch(`https://open.faceit.com/data/v4/players/${profile.player_id}/games/${game}/stats`, requestOptions),
		fetch(`https://open.faceit.com/data/v4/players/${profile.player_id}/stats/cs2`, requestOptions),
		getPlayerRanking(profile.player_id, profileGame.region, game),
	])

	result = {
		lvl: profileGame.skill_level,
		elo: profileGame.faceit_elo,
	}

	if (!lastGamesResponse.ok || !playerStatsResponse.ok) {
		return null
	}

	const lastGames = await lastGamesResponse.json() as FaceitLastGamesResponse
	const playerStats = await playerStatsResponse.json() as FaceitPlayerStatsResponse

	result.avgKdr = playerStats.lifetime['Average K/D Ratio']

	// @ts-expect-error
	const totals: Record<FaceitMatchStatsKeys, number> = {}
	// @ts-expect-error
	const counts: Record<FaceitMatchStatsKeys, number> = {}

	lastGames.items.forEach((item: any) => {
		const stats = item.stats

		for (const key in stats) {
			const value = Number.parseFloat(stats[key])

			if (!Number.isNaN(value)) {
				totals[key] = (totals[key] || 0) + value
				counts[key] = (counts[key] || 0) + 1
			}
		}

		// Count wins/losses dynamically
		const result = stats.Result
		counts[result] = (counts[result] || 0) + 1 // Increment count based on result
	})

	// @ts-expect-error
	const averages: Record<FaceitMatchStatsKeys, number> = new Map()
	for (const key in totals) {
		averages[key] = totals[key] / counts[key]
	}

	const winRate = (counts['1'] || 0) / lastGames.items.length // Assuming "1" indicates a wi

	result.lastMatches = {
		avgKills: averages.Kills.toFixed(),
		headshots: averages['Headshots %'].toFixed(),
		winRate: winRate.toString(),
		avgKd: averages['K/D Ratio'].toFixed(2),
		avgKr: averages['K/R Ratio'].toFixed(2),
	}

	result.worldRanking = positionResponse ?? ''

	return result
}

async function getPlayerRanking(player_id: string, region: string, game: string): Promise<string | null> {
	const response = await fetch(
		`https://open.faceit.com/data/v4/rankings/games/${game}/regions/${region}/players/${player_id}`,
		requestOptions,
	)
	if (!response.ok) {
		return null
	}

	const data = await response.json()

	return data.position
}
