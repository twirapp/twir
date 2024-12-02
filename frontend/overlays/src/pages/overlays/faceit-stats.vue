<script lang="ts" setup>
import { FaceitStatsWidget } from '@twir/frontend-faceit-stats'
import { useUrlSearchParams } from '@vueuse/core'

import type { Settings } from '@twir/frontend-faceit-stats'

const query = useUrlSearchParams()
const shouldDisplayKD = query.displayAvarageKdr === 'true'
const shouldDisplayRanking = query.displayAvarageKdr === 'true'
const shouldDisplayLastTwentyMatches = query.displayLastTwentyMatches === 'true'

const settings: Settings = {
	bgColor: query.bgColor as string,
	textColor: query.textColor as string,
	borderRadius: Number.parseInt(query.borderRadius as string),
	game: query.game as 'cs2',
	nickname: query.nickname as string,
	displayAvarageKdr: shouldDisplayKD,
	displayLastTwentyMatches: shouldDisplayLastTwentyMatches,
	displayWorldRanking: shouldDisplayRanking,
}
</script>

<template>
	<div class="wrapper">
		<FaceitStatsWidget :settings="settings" />
	</div>
</template>

<style scoped>
.wrapper {
	display: flex;
	width: max-content;
}
</style>
