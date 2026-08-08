<script setup lang="ts">
import type { StreamStatsCounterItem } from '@/composables/stream-stats/use-stream-stats-counters.js'

import CounterValue from './counter-value.vue'
import { counterIcons } from './counter-meta.js'

defineProps<{
	items: StreamStatsCounterItem[]
}>()
</script>

<template>
	<div class="bar">
		<template v-for="(item, index) in items" :key="item.id">
			<div v-if="index > 0" class="divider" />
			<div class="item">
				<component :is="counterIcons[item.key]" :size="15" :stroke-width="2.25" class="icon" />
				<CounterValue :value="item.value" :animated="item.key !== 'uptime'" class="value" />
				<span
					v-if="item.platformColor"
					class="platform-label"
					:style="{ color: item.platformColor }"
				>
					{{ item.label }}
				</span>
			</div>
		</template>
	</div>
</template>

<style scoped>
.bar {
	display: inline-flex;
	align-items: center;
	gap: 14px;
	padding: 9px 18px;
	border: 1px solid rgba(255, 255, 255, 0.08);
	border-radius: 9999px;
	background-color: rgba(10, 10, 14, 0.55);
	backdrop-filter: blur(8px);
	box-shadow: 0 4px 16px rgba(0, 0, 0, 0.35);
	color: #fff;
	font-family: system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif;
}

.divider {
	width: 1px;
	height: 16px;
	background-color: rgba(255, 255, 255, 0.14);
}

.item {
	display: flex;
	align-items: center;
	gap: 7px;
}

.icon {
	color: rgba(255, 255, 255, 0.75);
	flex-shrink: 0;
}

.value {
	font-size: 14px;
	font-weight: 600;
	line-height: 1;
	letter-spacing: 0.01em;
}

.platform-label {
	font-size: 11px;
	font-weight: 500;
	line-height: 1;
}
</style>
