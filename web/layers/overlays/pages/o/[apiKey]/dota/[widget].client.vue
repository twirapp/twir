<script setup lang="ts">
import { computed, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'

import { useDotaState } from '../../../../composables/dota/use-dota-state.js'

const route = useRoute()

definePageMeta({
	layout: 'clean',
})

const widgets = ['medal', 'wl', 'wp'] as const
type Widget = (typeof widgets)[number]

const widget = computed<Widget>(() => {
	const param = route.params.widget
	return (widgets as readonly string[]).includes(param) ? (param as Widget) : 'medal'
})

const { state, connect, destroy } = useDotaState()
connect(route.params.apiKey as string)

onUnmounted(() => {
	destroy()
})

function medalForMmr(mmr: number): string {
	if (mmr < 770) return 'Herald'
	if (mmr < 1540) return 'Guardian'
	if (mmr < 2310) return 'Crusader'
	if (mmr < 3080) return 'Archon'
	if (mmr < 3850) return 'Legend'
	if (mmr < 4620) return 'Ancient'
	if (mmr < 5620) return 'Divine'
	return 'Immortal'
}

const medal = computed(() => medalForMmr(state.value?.mmr ?? 0))

const winRate = computed(() => {
	if (!state.value) return '0.0'
	const games = state.value.sessionWins + state.value.sessionLosses
	if (games === 0) return '0.0'
	return ((state.value.sessionWins / games) * 100).toFixed(1)
})

const winProbability = computed(() => {
	if (!state.value) return '0.0'
	const probability =
		state.value.teamKnown && !state.value.teamIsRadiant
			? 1 - state.value.winProbability
			: state.value.winProbability
	return (probability * 100).toFixed(1)
})

const showWinProbability = computed(() => state.value?.inGame ?? false)
</script>

<template>
	<div v-if="state" class="inline-block font-sans">
		<div
			v-if="widget === 'medal'"
			class="inline-flex items-center gap-2 rounded-xl bg-black/65 px-4 py-2 text-white"
		>
			<span class="text-xl font-semibold">{{ medal }}</span>
			<span class="text-2xl font-bold tabular-nums">{{ state.mmr }} MMR</span>
		</div>

		<div
			v-else-if="widget === 'wl'"
			class="inline-flex items-center gap-2 rounded-xl bg-black/65 px-4 py-2 text-white"
		>
			<span class="text-2xl font-bold tabular-nums text-green-400">{{ state.sessionWins }}W</span>
			<span class="text-xl opacity-50">·</span>
			<span class="text-2xl font-bold tabular-nums text-red-400">{{ state.sessionLosses }}L</span>
			<span class="text-xl opacity-50">·</span>
			<span class="text-2xl font-bold tabular-nums">{{ winRate }}%</span>
		</div>

		<div
			v-else-if="widget === 'wp' && showWinProbability"
			class="inline-flex items-center gap-2 rounded-xl bg-black/65 px-4 py-2 text-white"
		>
			<span class="text-2xl font-bold tabular-nums">{{ winProbability }}%</span>
		</div>
	</div>
</template>
