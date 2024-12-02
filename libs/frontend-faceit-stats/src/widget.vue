<script setup lang="ts">
import { useIntervalFn, watchDebounced } from '@vueuse/core'
import { GlobeIcon } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'

import { computeStats } from './computeStats'
import SkillLevelIcon from './skill-level-icon.vue'

import type { Stats } from './computeStats'
import type { Settings } from './types'

const props = defineProps<{
	settings: Settings
}>()

useIntervalFn(() => {
	fetchStats()
}, 10 * 1000)

watchDebounced(props.settings, () => {
	fetchStats()
}, {
	debounce: 500,
	immediate: true,
})

onMounted(() => {
	fetchStats()
})

const stats = ref<Stats | null>(null)

async function fetchStats() {
	if (!props.settings.nickname) return
	const data = await computeStats(props.settings.nickname, props.settings.game)

	stats.value = data
}

const borderRadius = computed(() => {
	return `${props.settings.borderRadius}px`
})
</script>

<template>
	<div class="wrapper">
		<div
			class="widget"
			:style="{
				backgroundColor: settings.bgColor,
			}"
		>
			<div>
				<SkillLevelIcon :lvl="stats?.lvl ?? 1" />
				<span>{{ stats?.elo ?? '--' }}</span>
			</div>

			<template v-if="settings.displayAvarageKdr && stats?.avgKdr">
				<div class="separator" style="height: 8px;"></div>
				{{ stats.avgKdr }} KDR
			</template>

			<template v-if="settings.displayWorldRanking && stats?.worldRanking">
				<div class="separator" style="height: 8px;"></div>
				<div style="display: flex; align-items: center; gap: 4px">
					<GlobeIcon style="width: 16px; height: 16px" />
					#{{ stats.worldRanking }}
				</div>
			</template>
		</div>

		<div
			v-if="settings.displayLastTwentyMatches && stats?.lastMatches"
			class="stats"
			:style="{
				backgroundColor: settings.bgColor,
			}"
		>
			<span
				class="text-xs uppercase text-muted"
				style="opacity: 0.7"
				:style="{
					color: settings.textColor,
				}"
			>
				Latest 20 matches
			</span>

			<div class="items">
				<div class="item">
					<span class="text-sm">{{ stats.lastMatches.winRate }}%</span>
					<span class="text-xs">Winrate</span>
				</div>

				<div class="separator"></div>

				<div class="item">
					<span class="text-sm">{{ stats.lastMatches.avgKills }}/{{ stats.lastMatches.headshots }}</span>
					<span class="text-xs">Avg. Kills / HS</span>
				</div>

				<div class="separator"></div>

				<div class="item">
					<span class="text-sm">{{ stats.lastMatches.avgKd }}/{{ stats.lastMatches.avgKr }}%</span>
					<span class="text-xs">Avg. K/D / K/R</span>
				</div>
			</div>
		</div>
	</div>
</template>

<style>
@import url('https://fonts.googleapis.com/css2?family=Play:wght@400;700&family=Roboto:ital,wght@0,100;0,300;0,400;0,500;0,700;0,900;1,100;1,300;1,400;1,500;1,700;1,900&display=swap');

.wrapper {
	font-weight: 400;
  font-size: 16px;
	display: flex;
  flex-direction: column;
  gap: .25rem;
	color: v-bind('settings.textColor')
}

.widget {
	padding: 4px;
	border-radius: v-bind(borderRadius);
	display: flex;
	gap: 4px;
	height: max-content;
	flex-wrap: wrap;
	padding: 2px 10px 2px 4px;
	align-items: center;
	font-family: "Play", sans-serif;
  font-weight: 400;
  font-style: normal;
}

.widget > div {
	display: flex;
	gap: 2px;
	align-items: center;
}

.stats {
	display: flex;
	flex-direction: column;
	flex-wrap: wrap;
	width: 100%;
	border-radius: 12px;
	padding-top: 2px;
	padding-bottom: 2px;
	padding-left: 4px;
	padding-right: 4px;
	justify-items: center;
	align-items: center;
	font-family: "Play", sans-serif;
  font-weight: 400;
  font-style: normal;
}

.stats .items {
	display: flex;
	gap: 8px;
	letter-spacing: .02em;
}

.stats .items .item {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
}

.stats .items .item span:first-child {

}

.text-sm {
	font-size: 14px;
}

.text-xs {
	font-size: 10px;
}

.uppercase {
	text-transform: uppercase;
}

.separator {
	border-left-color: rgb(255, 255, 255);
	border-left: 1px solid;
	height: 24px;
	opacity: .3;
	margin-right: 3px;
	margin-left: 3px;
	align-self: center;
}
</style>
