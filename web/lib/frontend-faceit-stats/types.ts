export interface Settings {
	nickname: string
	game: 'cs2'
	bgColor: string
	textColor: string
	borderRadius: number
	displayAvarageKdr: boolean
	displayWorldRanking: boolean
	displayLastTwentyMatches: boolean
}

export interface FaceitProfileResponse {
	player_id: string
	nickname: string
	country: string
	cover_image: string
	games: Record<string, {
		region: string
		skill_level: number
		faceit_elo: number
	}>
}

export interface FaceitMatchStats {
	'K/D Ratio': string
	'K/R Ratio': string
	'Kills': string
	'Headshots %': string
	'Headshots': string
}

export type FaceitMatchStatsKeys = keyof FaceitMatchStats

export interface FaceitLastGamesResponse {
	items: Array<{
		stats: FaceitMatchStats
	}>
}

export interface FaceitPlayerStatsResponse {
	lifetime: {
		'Average K/D Ratio': string
	}
}
