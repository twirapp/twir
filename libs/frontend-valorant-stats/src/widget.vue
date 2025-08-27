<script setup lang="ts">
import {
	Api,
	HttpClient,
	type IntegrationsValorantStatsOutput,
	type Stream,
} from '@twir/api/openapi'
import { useIntervalFn, watchDebounced } from '@vueuse/core'
import {
	ArrowDown,
	ArrowUp,
	Globe,
	LoaderCircle,
	Minus,
	TrendingDown,
	TrendingUp,
} from 'lucide-vue-next'
import { computed, nextTick, onMounted, ref } from 'vue'

import type { Settings } from './types'

const props = withDefaults(
	defineProps<{
		settings?: Settings
		apiKey?: string
	}>(),
	{
		settings: () => ({
			backgroundColor: '#07090e',
			textColor: '#FFFFFF',
			primaryTextColor: '#FFFFFF',
			winColor: '#00FFE3',
			loseColor: '#FF7986',

			disabledPeakRR: false,
			disabledLeaderboardPlace: false,
			disabledPeakRankIcon: false,
			disabledBackground: false,
			disabledBorder: false,
			disabledWinLose: false,
			disabledProgress: false,
			disabledGlowEffect: false,
			overlayFont: 'Inter',
		}),
	},
)

const apiClient = computed(() => {
	const headers: HeadersInit = {}
	if (props.apiKey) {
		headers['api-key'] = props.apiKey
	}

	const apiClient = new Api(
		new HttpClient({
			baseUrl: `${window.location.origin}/api`,
			baseApiParams: {
				headers,
			},
		}),
	)

	return apiClient
})

const stats = ref<IntegrationsValorantStatsOutput | undefined>()
const streamData = ref<Stream | undefined>()

async function fetchStats() {
	try {
		const response = await apiClient.value.v1.integrationsValorantStats()

		if (!response.ok) {
			console.error('Error response from API:', response.error)
			return
		}
		stats.value = response.data.data
	} catch (error) {
		console.error('Error fetching stats:', error)
	}
}

async function fetchStream() {
	try {
		const response = await apiClient.value.v1.channelsStreamsCurrent()
		if (!response.ok) {
			console.error('Error response from API:', response.error)
			return
		}

		streamData.value = response.data.data
	} catch (error) {
		console.error('Error fetching stream:', error)
	}
}

onMounted(async () => {
	await nextTick()
	fetchStats()
	fetchStream()
})

watchDebounced(
	() => props.settings,
	async () => {
		await fetchStream()
		await fetchStats()
	},
	{
		debounce: 500,
		immediate: true,
	},
)

useIntervalFn(async () => {
	await fetchStream()
	await fetchStats()
}, 1 * 1000)

const matches = computed(() => {
	if (!stats.value?.matches) return []

	return stats.value.matches.filter((match) => {
		if (streamData.value?.StartedAt) {
			return (
				new Date(match.meta.started_at).getTime() > new Date(streamData.value.StartedAt).getTime()
			)
		}

		return true
	})
})

const wins = computed(() => {
	return matches.value.filter((match) => {
		const team = match.stats.team.toLowerCase()

		return (
			(team === 'red' && match.teams.red > match.teams.blue)
			|| (team === 'blue' && match.teams.blue > match.teams.red)
		)
	}).length
})

const matchesWinrate = computed(() => {
	if (!matches.value || matches.value.length === 0) return 0

	return Math.round((wins.value / matches.value.length) * 100)
})

const avgKills = computed(() => {
	if (!matches.value || matches.value.length === 0) return 0

	const totalKills = matches.value.reduce((sum, match) => sum + match.stats.kills, 0)
	return (totalKills / matches.value.length).toFixed(1)
})

const headShots = computed(() => {
	if (!matches.value || matches.value.length === 0) return 0

	const totalHeadshots = matches.value.reduce((sum, match) => sum + match.stats.shots.head, 0)
	const totalShots = matches.value.reduce(
		(sum, match) => sum + match.stats.shots.head + match.stats.shots.body + match.stats.shots.leg,
		0,
	)

	return ((totalHeadshots / totalShots) * 100).toFixed(0)
})

const avgKD = computed(() => {
	if (!matches.value || matches.value.length === 0) return 0

	const totalKills = matches.value.reduce((sum, match) => sum + match.stats.kills, 0)
	const totalDeaths = matches.value.reduce((sum, match) => sum + match.stats.deaths, 0)
	if (totalDeaths === 0) return totalKills.toFixed(2)
	return (totalKills / totalDeaths).toFixed(2)
})

const avgKR = computed(() => {
	if (!matches.value || matches.value.length === 0) return 0

	const totalKills = matches.value.reduce((sum, match) => sum + match.stats.kills, 0)
	const totalRounds = matches.value.reduce(
		(sum, match) => sum + (match.teams.red + match.teams.blue),
		0,
	)
	if (totalRounds === 0) return totalKills.toFixed(2)
	return (totalKills / totalRounds).toFixed(2)
})

const rankSrc = computed(() => {
	return new URL(`./assets/ranks/${stats.value?.mmr.current.tier.id}.webp`, import.meta.url).href
})
</script>

<template>
	<transition
		name="slide"
		enter-active-class="transition duration-300 ease-out"
		leave-active-class="transition duration-300 ease-in"
		enter-from-class="transform -translate-y-4 opacity-0"
		leave-to-class="transform -translate-y-4 opacity-0"
	>
		<div
			v-if="!stats"
			class="bg-black/90 min-h-54 rounded-md p-2"
			:style="{ backgroundColor: `${settings.backgroundColor}99` }"
		>
			<LoaderCircle
				class="size-20 text-white animate-spin text-center mx-auto"
				aria-hidden="true"
			/>
		</div>
		<div v-else class="flex flex-col gap-1">
			<div
				class="minimal-style flex min-h-[50px] w-full min-[310px]:flex-row flex-col items-center justify-between gap-3 rounded-lg bg-black/90 px-3 py-1"
				:style="{
					backgroundColor: settings.disabledBackground
						? 'transparent'
						: `${settings.backgroundColor}99`,
					fontFamily: settings.overlayFont,
				}"
				:class="{
					'overflow-hidden border-[2px] border-white/10': !settings.disabledBackground,
					'border-none': settings.disabledBorder,
					'h-fit': settings.disabledBackground,
				}"
			>
				<div class="inline-flex items-center justify-center gap-2">
					<div class="relative">
						<div class="relative flex">
							<img :src="rankSrc" class="z-10" alt="" height="40" width="40" fetchpriority="high" />
							<img
								v-if="!settings.disabledGlowEffect"
								class="absolute top-1/2 left-1/2 size-10 max-w-[unset] -translate-x-1/2 -translate-y-1/2 transform blur-[10px]"
								:src="rankSrc"
								alt=""
								fetchpriority="high"
							/>
						</div>
						<span
							class="absolute top-0 right-0 z-10 flex size-4 flex-col items-center justify-center rounded-full bg-white text-sm leading-none font-medium text-black"
						>
							{{ stats.mmr.current.tier.name.split(' ').at(-1) }}
						</span>
					</div>
					<div
						class="font-bold text-[var(--primary-text-color)] uppercase"
						:class="{
							'drop-shadow-[0px_0px_6px_var(--primary-text-color)]': !settings.disabledGlowEffect,
						}"
					>
						{{ stats.mmr.current.rr }} RR
					</div>
				</div>
				<div v-if="!settings.disabledWinLose" class="inline-flex items-center gap-2">
					<span
						v-if="
							stats.mmr.current.leaderboard_placement.rank && !settings.disabledLeaderboardPlace
						"
						class="inline-flex items-center gap-1 font-bold text-[var(--primary-text-color)] uppercase"
						:class="{
							'drop-shadow-[0px_0px_6px_var(--primary-text-color)]': !settings.disabledGlowEffect,
						}"
					>
						<Globe />
						#{{ stats.mmr.current.leaderboard_placement.rank }}
					</span>
					<span
						class="inline-flex items-center gap-1 font-bold text-[var(--win-color)]"
						:class="{ 'drop-shadow-[0px_0px_6px_var(--win-color)]': !settings.disabledGlowEffect }"
					>
						<ArrowUp />
						{{ wins }}
					</span>
					<span
						class="inline-flex items-center gap-1 font-bold text-[var(--lose-color)]"
						:class="{ 'drop-shadow-[0px_0px_6px_var(--lose-color)]': !settings.disabledGlowEffect }"
					>
						<ArrowDown />
						{{ matches.length - wins }}
					</span>
					<span
						class="flex flex-row items-center gap-1 font-bold"
						:class="{
							'drop-shadow-[0px_0px_6px_var(--win-color)]':
								!settings.disabledGlowEffect && stats.mmr.current.last_change > 0,
							'drop-shadow-[0px_0px_6px_var(--lose-color)]':
								!settings.disabledGlowEffect && stats.mmr.current.last_change < 0,
							'drop-shadow-[0px_0px_6px_var(--primary-text-color)]':
								!settings.disabledGlowEffect && stats.mmr.current.last_change === 0,
							'text-[var(--win-color)]': stats.mmr.current.last_change > 0,
							'text-[var(--lose-color)]': stats.mmr.current.last_change < 0,
							'text-[var(--primary-text-color)]': stats.mmr.current.last_change === 0,
						}"
					>
						<TrendingUp v-if="stats.mmr.current.last_change > 0" />
						<TrendingDown v-else-if="stats.mmr.current.last_change < 0" />
						<Minus v-else />
						{{ stats.mmr.current.last_change }}
					</span>
				</div>
			</div>

			<div
				v-if="!settings.disabledTwentyLastMatches"
				class="minimal-style flex min-h-[50px] w-full flex-col items-center justify-between rounded-lg bg-black/90 px-3 py-1 text-[var(--primary-text-color)] uppercase"
				:style="{
					backgroundColor: settings.disabledBackground
						? 'transparent'
						: `${settings.backgroundColor}99`,
					fontFamily: settings.overlayFont,
				}"
				:class="{
					'overflow-hidden border-[2px] border-white/10': !settings.disabledBackground,
					'border-none': settings.disabledBorder,
					'h-fit': settings.disabledBackground,
				}"
			>
				<h1 class="text-[11px] text-bold text-center mx-auto" style="color: rgb(255, 255, 255)">
					Latest {{ matches.length }} matches
				</h1>

				<div class="flex flex-row gap-2 items flex-wrap w-full items-center justify-center">
					<div class="item">
						<span class="text-sm">{{ matchesWinrate }}%</span>
						<span class="text-xs">Win Rate</span>
					</div>
					<div class="separator max-[300px]:hidden min-[350px]:block" />
					<div class="item">
						<span class="text-sm">{{ avgKills }} / {{ headShots }}%</span>
						<span class="text-xs">Avg kills / HS</span>
					</div>
					<div class="separator max-[300px]:hidden min-[350px]:block" />
					<div class="item">
						<span class="text-sm">{{ avgKD }} ({{ avgKR }})</span>
						<span class="text-xs">Avg KD (KR)</span>
					</div>
				</div>
			</div>
		</div>
	</transition>
</template>

<style scoped>
.minimal-style {
	--text-color: v-bind(settings.textColor);
	--primary-text-color: v-bind(settings.primaryTextColor);
	--win-color: v-bind(settings.winColor);
	--lose-color: v-bind(settings.loseColor);
}

.items {
	letter-spacing: 0.02em;
}

.item {
	@apply flex flex-col items-center justify-center;
}

.separator {
	border-left: 1px solid;
	height: 24px;
	opacity: 0.3;
	margin-right: 3px;
	margin-left: 3px;
	align-self: center;
}

.slide-enter-active {
	@apply transition duration-300 ease-out;
}
.slide-leave-active {
	@apply transition duration-300 ease-in;
}
.slide-enter-from,
.slide-leave-to {
	@apply transform -translate-y-4 opacity-0;
}
</style>
