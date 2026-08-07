import type { Component } from 'vue'

import ChatWidgetSettings from './components/widget-settings/ChatWidgetSettings.vue'
import DudesWidgetSettings from './components/widget-settings/DudesWidgetSettings.vue'
import FaceitStatsWidgetSettings from './components/widget-settings/FaceitStatsWidgetSettings.vue'
import NowPlayingWidgetSettings from './components/widget-settings/NowPlayingWidgetSettings.vue'
import ValorantStatsWidgetSettings from './components/widget-settings/ValorantStatsWidgetSettings.vue'

const faceitDefaults = {
	nickname: '',
	game: 'cs2',
	bgColor: '#1f1f22',
	textColor: '#ffffff',
	borderRadius: '24',
	displayAvarageKdr: 'false',
	displayWorldRanking: 'false',
	displayLastTwentyMatches: 'false',
} satisfies Record<string, string>

const valorantDefaults = {
	backgroundColor: '#07090e',
	textColor: '#f2f2f2',
	primaryTextColor: '#B9B4B4',
	winColor: '#00FFE3',
	loseColor: '#FF7986',
	disabledPeakRR: 'false',
	disabledLeaderboardPlace: 'false',
	disabledPeakRankIcon: 'false',
	disabledBorder: 'false',
	disabledWinLose: 'false',
	disabledProgress: 'false',
	disabledGlowEffect: 'false',
	disabledTwentyLastMatches: 'false',
} satisfies Record<string, string>

function withParams(base: string, params?: Record<string, string>) {
	if (!params || Object.keys(params).length === 0) return base

	const query = new URLSearchParams(params).toString()
	return `${base}?${query}`
}

export interface OverlayWidgetRegistryEntry {
	readonly key: string
	readonly nameKey: string
	readonly descriptionKey: string
	readonly icon: string
	readonly buildUrl: (ctx: {
		readonly origin: string
		readonly apiKey: string
		readonly params?: Record<string, string>
	}) => string
	readonly settingsComponent?: Component
}

export const overlayWidgetRegistry = [
	{
		key: 'chat',
		nameKey: 'overlayBuilder.widgets.chat.name',
		descriptionKey: 'overlayBuilder.widgets.chat.description',
		icon: 'lucide:messages-square',
		buildUrl: ({ origin, apiKey, params }) => {
			return withParams(`${origin}/overlays/${apiKey}/chat`, params)
		},
		settingsComponent: ChatWidgetSettings,
	},
	{
		key: 'now-playing',
		nameKey: 'overlayBuilder.widgets.nowPlaying.name',
		descriptionKey: 'overlayBuilder.widgets.nowPlaying.description',
		icon: 'lucide:audio-lines',
		buildUrl: ({ origin, apiKey, params }) => withParams(`${origin}/overlays/${apiKey}/now-playing`, params),
		settingsComponent: NowPlayingWidgetSettings,
	},
	{
		key: 'faceit-stats',
		nameKey: 'overlayBuilder.widgets.faceitStats.name',
		descriptionKey: 'overlayBuilder.widgets.faceitStats.description',
		icon: 'lucide:badge-check',
		buildUrl: ({ origin, params }) =>
			withParams(`${origin}/overlays/faceit-stats`, { ...faceitDefaults, ...params }),
		settingsComponent: FaceitStatsWidgetSettings,
	},
	{
		key: 'valorant-stats',
		nameKey: 'overlayBuilder.widgets.valorantStats.name',
		descriptionKey: 'overlayBuilder.widgets.valorantStats.description',
		icon: 'lucide:crosshair',
		buildUrl: ({ origin, apiKey, params }) =>
			withParams(`${origin}/o/${apiKey}/valorant-stats`, { ...valorantDefaults, ...params }),
		settingsComponent: ValorantStatsWidgetSettings,
	},
	{
		key: 'kappagen',
		nameKey: 'overlayBuilder.widgets.kappagen.name',
		descriptionKey: 'overlayBuilder.widgets.kappagen.description',
		icon: 'lucide:party-popper',
		buildUrl: ({ origin, apiKey }) => `${origin}/overlays/${apiKey}/kappagen`,
	},
	{
		key: 'dudes',
		nameKey: 'overlayBuilder.widgets.dudes.name',
		descriptionKey: 'overlayBuilder.widgets.dudes.description',
		icon: 'lucide:users',
		buildUrl: ({ origin, apiKey, params }) => withParams(`${origin}/overlays/${apiKey}/dudes`, params),
		settingsComponent: DudesWidgetSettings,
	},
	{
		key: 'afk',
		nameKey: 'overlayBuilder.widgets.afk.name',
		descriptionKey: 'overlayBuilder.widgets.afk.description',
		icon: 'lucide:coffee',
		buildUrl: ({ origin, apiKey }) => `${origin}/overlays/${apiKey}/brb`,
	},
] satisfies readonly OverlayWidgetRegistryEntry[]
