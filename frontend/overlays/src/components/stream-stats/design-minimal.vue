<script setup lang="ts">
import type { StreamStatsCounterItem } from '@/composables/stream-stats/use-stream-stats-counters.js'

import CounterValue from './counter-value.vue'
import { counterIcons } from './counter-meta.js'

defineProps<{
	items: StreamStatsCounterItem[]
}>()
</script>

<template>
	<div class="minimal">
		<div v-for="item in items" :key="item.id" class="item">
			<component :is="counterIcons[item.key]" :size="15" :stroke-width="2.5" class="icon" />
			<CounterValue :value="item.value" :animated="item.key !== 'uptime'" class="value" />
			<span
				v-if="item.platformColor"
				class="platform-label"
				:style="{ color: item.platformColor }"
			>
				{{ item.label }}
			</span>
		</div>
	</div>
</template>

<style scoped>
.minimal {
	display: inline-flex;
	align-items: center;
	gap: 18px;
	color: #fff;
	font-family: system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif;
}

.item {
	display: flex;
	align-items: center;
	gap: 6px;
	text-shadow:
		0 1px 3px rgba(0, 0, 0, 0.9),
		0 0 2px rgba(0, 0, 0, 0.9);
}

.icon {
	flex-shrink: 0;
	filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.9));
}

.value {
	font-size: 14px;
	font-weight: 600;
	line-height: 1;
}

.platform-label {
	font-size: 11px;
	font-weight: 600;
	line-height: 1;
}
</style>
