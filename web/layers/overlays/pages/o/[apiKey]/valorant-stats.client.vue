<script setup lang="ts">
import { type Settings, ValorantStatsWidget } from '@twir/frontend-valorant-stats'
import { useRoute } from 'vue-router'

const route = useRoute<'o-apiKey-valorant-stats'>()

definePageMeta({
	layout: 'clean',
	colorMode: 'light',
})

const query = route.query

const settings = computed<Settings>(() => {
	return {
		...query,
		disabledPeakRR: query.disabledPeakRR === 'true',
		disabledLeaderboardPlace: query.disabledLeaderboardPlace === 'true',
		disabledPeakRankIcon: query.disabledPeakRankIcon === 'true',
		disabledBorder: query.disabledBorder === 'true',
		disabledGlowEffect: query.disabledGlowEffect === 'true',
		disabledWinLose: query.disabledWinLose === 'true',
		disabledProgress: query.disabledProgress === 'true',
		disabledTwentyLastMatches: query.disabledTwentyLastMatches === 'true',
	}
})

const apiKey = route.params.apiKey
</script>

<template>
	<div class="p-1">
		<ValorantStatsWidget :api-key="apiKey" :settings="settings" />
	</div>
</template>
