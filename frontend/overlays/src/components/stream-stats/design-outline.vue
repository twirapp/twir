<script setup lang="ts">
import { computed } from 'vue'

import type { StreamStatsCounterItem } from '@/composables/stream-stats/use-stream-stats-counters.js'
import { StreamStatsOverlayVariant } from '@/gql/graphql'

import { counterIcons, variantRootClass } from './counter-meta.js'
import PlatformIcon from './platform-icon.vue'
import RollingNumber from './rolling-number.vue'

const props = withDefaults(
	defineProps<{
		items: StreamStatsCounterItem[]
		variant?: StreamStatsOverlayVariant
		platformIcons?: boolean
	}>(),
	{
		variant: StreamStatsOverlayVariant.Horizontal,
		platformIcons: false,
	},
)

const rootClass = computed(() => variantRootClass(props.variant))

const iconSize = computed(() => {
	switch (props.variant) {
		case StreamStatsOverlayVariant.HorizontalCompact:
		case StreamStatsOverlayVariant.VerticalCompact:
			return 12
		case StreamStatsOverlayVariant.Large:
			return 18
		default:
			return 14
	}
})

function chipStyle(color: string) {
	return {
		borderColor: color,
		backgroundColor: `color-mix(in srgb, ${color} 5%, transparent)`,
	}
}

function faded(color: string): string {
	return `color-mix(in srgb, ${color} 60%, transparent)`
}
</script>

<template>
	<div class="outline" :class="rootClass">
		<div v-for="item in items" :key="item.id" class="chip" :style="chipStyle(item.color)">
			<component
				:is="counterIcons[item.key]"
				:size="iconSize"
				:stroke-width="2.25"
				class="icon"
				:style="{ color: item.color }"
			/>
			<RollingNumber :value="item.value" class="value" :style="{ color: item.color }" />
			<PlatformIcon
				v-if="platformIcons && item.platform && item.platformColor"
				:platform="item.platform"
				:size="iconSize - 2"
				:style="{ color: item.platformColor }"
			/>
			<span
				v-else
				class="label"
				:style="{ color: item.platformColor ? item.platformColor : faded(item.color) }"
			>
				{{ item.label }}
			</span>
		</div>
	</div>
</template>

<style scoped>
.outline {
	display: inline-flex;
	align-items: center;
	gap: 8px;
	font-family: system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif;
}

.chip {
	display: flex;
	align-items: center;
	gap: 7px;
	padding: 5px 13px;
	border: 1.5px solid transparent;
	border-radius: 9999px;
	box-shadow: 0 2px 10px rgba(0, 0, 0, 0.25);
}

.icon {
	flex-shrink: 0;
}

.value {
	font-size: 14px;
	font-weight: 700;
	line-height: 1;
}

.label {
	font-size: 9px;
	font-weight: 600;
	line-height: 1;
	letter-spacing: 0.12em;
	text-transform: uppercase;
}

.root--compact {
	gap: 6px;
}

.root--compact .chip {
	gap: 5px;
	padding: 4px 10px;
}

.root--compact .value {
	font-size: 12px;
}

.root--compact .label {
	font-size: 8px;
}

.root--large {
	gap: 10px;
}

.root--large .chip {
	gap: 9px;
	padding: 7px 17px;
}

.root--large .value {
	font-size: 19px;
}

.root--large .label {
	font-size: 11px;
}

.root--vertical {
	flex-direction: column;
	align-items: stretch;
}

.root--vertical .chip {
	justify-content: flex-start;
}

.root--vertical .value {
	margin-right: auto;
}
</style>
