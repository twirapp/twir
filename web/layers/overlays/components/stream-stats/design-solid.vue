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
			return 11
		case StreamStatsOverlayVariant.Large:
			return 16
		default:
			return 13
	}
})
</script>

<template>
	<div class="solid" :class="rootClass">
		<div v-for="item in items" :key="item.id" class="chip">
			<span class="icon-badge">
				<Icon
					:name="counterIcons[item.key]"
					class="icon"
					:style="{ width: iconSize + 'px', height: iconSize + 'px' }"
				/>
			</span>
			<RollingNumber :value="item.value" class="value" />
			<span v-if="item.platformColor" class="platform-label">
				<span class="platform-dot" :style="{ backgroundColor: item.platformColor }" />
				<PlatformIcon
					v-if="platformIcons && item.platform"
					:platform="item.platform"
					:size="iconSize - 1"
					class="platform-icon"
				/>
				<template v-else>{{ item.label }}</template>
			</span>
		</div>
	</div>
</template>

<style scoped>
.solid {
	display: inline-flex;
	align-items: center;
	gap: 8px;
	font-family: system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif;
}

.chip {
	display: flex;
	align-items: center;
	gap: 7px;
	padding: 5px 14px 5px 5px;
	border-radius: 9999px;
	background: linear-gradient(135deg, #7c3aed 0%, #db2777 100%);
	box-shadow: 0 4px 14px rgba(124, 58, 237, 0.4);
	color: #fff;
}

.icon-badge {
	display: flex;
	align-items: center;
	justify-content: center;
	width: 22px;
	height: 22px;
	border-radius: 9999px;
	background-color: rgba(255, 255, 255, 0.22);
	flex-shrink: 0;
}

.icon {
	color: #fff;
}

.value {
	font-size: 14px;
	font-weight: 700;
	line-height: 1;
}

.platform-label {
	display: flex;
	align-items: center;
	gap: 5px;
	font-size: 11px;
	font-weight: 600;
	line-height: 1;
	color: rgba(255, 255, 255, 0.8);
}

.platform-icon {
	color: rgba(255, 255, 255, 0.9);
}

.platform-dot {
	width: 6px;
	height: 6px;
	border-radius: 9999px;
	flex-shrink: 0;
}

.root--compact {
	gap: 6px;
}

.root--compact .chip {
	gap: 6px;
	padding: 4px 11px 4px 4px;
}

.root--compact .icon-badge {
	width: 18px;
	height: 18px;
}

.root--compact .value {
	font-size: 12px;
}

.root--compact .platform-label {
	font-size: 10px;
}

.root--compact .platform-dot {
	width: 5px;
	height: 5px;
}

.root--large {
	gap: 10px;
}

.root--large .chip {
	gap: 9px;
	padding: 7px 18px 7px 7px;
}

.root--large .icon-badge {
	width: 28px;
	height: 28px;
}

.root--large .value {
	font-size: 18px;
}

.root--large .platform-label {
	font-size: 14px;
}

.root--large .platform-dot {
	width: 7px;
	height: 7px;
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
