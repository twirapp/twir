<script setup lang="ts">
import type { StreamStatsCounterItem } from '@/composables/stream-stats/use-stream-stats-counters.js'

import CounterValue from './counter-value.vue'
import { counterIcons } from './counter-meta.js'

defineProps<{
	items: StreamStatsCounterItem[]
}>()
</script>

<template>
	<div class="cards">
		<div v-for="item in items" :key="item.id" class="card">
			<component :is="counterIcons[item.key]" :size="16" :stroke-width="2.25" class="icon" />
			<CounterValue :value="item.value" :animated="item.key !== 'uptime'" class="value" />
			<span
				class="label"
				:class="{ 'label-platform': item.platformColor }"
				:style="item.platformColor ? { color: item.platformColor } : {}"
			>
				{{ item.label }}
			</span>
		</div>
	</div>
</template>

<style scoped>
.cards {
	display: inline-flex;
	align-items: stretch;
	gap: 10px;
	font-family: system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif;
}

.card {
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 4px;
	min-width: 76px;
	padding: 10px 16px;
	border: 1px solid rgba(255, 255, 255, 0.08);
	border-radius: 12px;
	background-color: rgba(10, 10, 14, 0.55);
	backdrop-filter: blur(8px);
	box-shadow: 0 4px 16px rgba(0, 0, 0, 0.35);
	color: #fff;
}

.icon {
	color: rgba(255, 255, 255, 0.7);
}

.value {
	font-size: 18px;
	font-weight: 700;
	line-height: 1.1;
}

.label {
	font-size: 9px;
	font-weight: 600;
	line-height: 1;
	letter-spacing: 0.12em;
	text-transform: uppercase;
	color: rgba(255, 255, 255, 0.5);
}

.label-platform {
	font-size: 10px;
}
</style>
