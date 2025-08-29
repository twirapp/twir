<script setup lang="ts">
import { type Settings, ValorantStatsWidget } from '@twir/frontend-valorant-stats'

const route = useRoute()

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
		disabledBackground: query.disabledBackground === 'true',
		disabledBorder: query.disabledBorder === 'true',
		disabledGlowEffect: query.disabledGlowEffect === 'true',
		disabledWinLose: query.disabledWinLose === 'true',
		disabledProgress: query.disabledProgress === 'true',
		disabledTwentyLastMatches: query.disabledTwentyLastMatches === 'true',
	}
})

const apiKey = route.params.apiKey as string

useHead({
  link: [
    {
      rel: 'preconnect',
      href: 'https://fonts.googleapis.com'
    },
    {
      rel: 'preconnect',
      href: 'https://fonts.gstatic.com',
      crossorigin: ''
    },
    {
      rel: 'stylesheet',
      href: 'https://fonts.googleapis.com/css2?family=Inter:ital,opsz,wght@0,14..32,100..900;1,14..32,100..900&display=swap'
    }
  ]
})
</script>

<template>
	<div class="p-1">
		<ValorantStatsWidget :api-key="apiKey" :settings="settings" />
	</div>
</template>
