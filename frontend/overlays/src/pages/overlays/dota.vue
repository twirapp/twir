<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'

import { useDotaSocket } from '@/composables/dota/use-dota-socket.ts'

const route = useRoute()

const widget = computed(() => {
	const param = route.params.widget
	return typeof param === 'string' ? param : ''
})

const { state, connect, close } = useDotaSocket()

onMounted(() => {
	connect(route.params.apiKey as string)
})

onUnmounted(() => {
	close()
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

const streamerWinProbability = computed(() => {
	if (!state.value) return '0.0'
	const probability = state.value.teamKnown && !state.value.teamIsRadiant
		? 1 - state.value.winProbability
		: state.value.winProbability
	return (probability * 100).toFixed(1)
})

const showWinProbability = computed(() => state.value?.inGame ?? false)
</script>

<template>
	<div v-if="state" class="dota-overlay">
		<div v-if="widget === 'medal'" class="badge">
			<span class="badge-medal">{{ medal }}</span>
			<span class="badge-value">{{ state.mmr }} MMR</span>
		</div>

		<div v-else-if="widget === 'wl'" class="badge">
			<span class="badge-value wins">{{ state.sessionWins }}W</span>
			<span class="badge-separator">·</span>
			<span class="badge-value losses">{{ state.sessionLosses }}L</span>
			<span class="badge-separator">·</span>
			<span class="badge-value">{{ winRate }}%</span>
		</div>

		<div v-else-if="widget === 'wp' && showWinProbability" class="badge">
			<span class="badge-value">{{ streamerWinProbability }}%</span>
		</div>
	</div>
</template>

<style scoped>
.dota-overlay {
	font-family: 'Segoe UI', Roboto, sans-serif;
	display: inline-block;
}

.badge {
	display: inline-flex;
	align-items: center;
	gap: 8px;
	padding: 8px 16px;
	background: rgba(0, 0, 0, 0.65);
	border-radius: 12px;
	color: #fff;
}

.badge-medal {
	font-size: 20px;
	font-weight: 600;
}

.badge-value {
	font-size: 24px;
	font-weight: 700;
	font-variant-numeric: tabular-nums;
}

.badge-separator {
	font-size: 20px;
	opacity: 0.5;
}

.wins {
	color: #4ade80;
}

.losses {
	color: #f87171;
}
</style>
