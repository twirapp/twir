<script setup lang="ts">
import type { StreamStatsCounterItem } from '../../composables/stream-stats/use-stream-stats-counters.js'
import { StreamStatsOverlayVariant } from '~/gql/graphql.js'

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
			return 14
		case StreamStatsOverlayVariant.Large:
			return 20
		default:
			return 16
	}
})
</script>

<template>
	<div class="cards" :class="rootClass">
		<div v-for="item in items" :key="item.id" class="card">
			<div class="accent" :style="{ backgroundColor: item.color }" />
			<Icon
				:name="counterIcons[item.key]"
				class="icon"
				:style="{ width: iconSize + 'px', height: iconSize + 'px', color: item.color }"
			/>
			<RollingNumber :value="item.value" class="value" />
			<PlatformIcon
				v-if="platformIcons && item.platform && item.platformColor"
				:platform="item.platform"
				:size="iconSize - 3"
				:style="{ color: item.platformColor }"
			/>
			<span
				v-else
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
	position: relative;
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 4px;
	min-width: 76px;
	padding: 12px 16px 10px;
	border: 1px solid rgba(255, 255, 255, 0.08);
	border-radius: 12px;
	background-color: rgba(10, 10, 14, 0.6);
	backdrop-filter: blur(10px);
	box-shadow: 0 4px 16px rgba(0, 0, 0, 0.35);
	color: #fff;
	overflow: hidden;
}

.accent {
	position: absolute;
	top: 0;
	left: 0;
	right: 0;
	height: 2px;
}

.icon {
	flex-shrink: 0;
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

.root--compact {
	gap: 8px;
}

.root--compact .card {
	min-width: 64px;
	padding: 10px 12px 8px;
	gap: 3px;
}

.root--compact .value {
	font-size: 16px;
}

.root--compact .label {
	font-size: 8px;
}

.root--compact .label-platform {
	font-size: 9px;
}

.root--large {
	gap: 12px;
}

.root--large .card {
	min-width: 96px;
	padding: 14px 20px 12px;
	gap: 5px;
}

.root--large .value {
	font-size: 22px;
}

.root--large .label {
	font-size: 11px;
}

.root--large .label-platform {
	font-size: 12px;
}

.root--vertical {
	flex-direction: column;
	align-items: stretch;
}

.root--vertical .card {
	flex-direction: row;
	align-items: center;
	justify-content: flex-start;
	gap: 8px;
	min-width: 0;
	padding: 8px 14px;
}

.root--vertical .card .accent {
	top: 0;
	bottom: 0;
	right: auto;
	width: 2px;
	height: auto;
}

.root--vertical .value {
	margin-right: auto;
}
</style>
