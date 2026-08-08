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
			return 13
		case StreamStatsOverlayVariant.Large:
			return 18
		default:
			return 15
	}
})

const platformIconSize = computed(() => iconSize.value - 2)
</script>

<template>
	<div class="glass" :class="rootClass">
		<template v-for="(item, index) in items" :key="item.id">
			<div v-if="index > 0" class="divider" />
			<div class="item">
				<component
					:is="counterIcons[item.key]"
					:size="iconSize"
					:stroke-width="2.25"
					class="icon"
					:style="{ color: item.color }"
				/>
				<RollingNumber :value="item.value" class="value" />
				<template v-if="item.platformColor">
					<PlatformIcon
						v-if="platformIcons && item.platform"
						:platform="item.platform"
						:size="platformIconSize"
						:style="{ color: item.platformColor }"
					/>
					<span v-else class="platform-label" :style="{ color: item.platformColor }">
						{{ item.label }}
					</span>
				</template>
			</div>
		</template>
	</div>
</template>

<style scoped>
.glass {
	display: inline-flex;
	align-items: center;
	gap: 14px;
	padding: 9px 18px;
	border: 1px solid rgba(255, 255, 255, 0.1);
	border-radius: 9999px;
	background-color: rgba(10, 10, 14, 0.5);
	backdrop-filter: blur(10px);
	box-shadow:
		inset 0 1px 0 rgba(255, 255, 255, 0.08),
		0 4px 16px rgba(0, 0, 0, 0.35);
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

.root--compact {
	gap: 10px;
	padding: 7px 14px;
}

.root--compact .divider {
	height: 13px;
}

.root--compact .item {
	gap: 6px;
}

.root--compact .value {
	font-size: 12px;
}

.root--compact .platform-label {
	font-size: 10px;
}

.root--large {
	gap: 18px;
	padding: 12px 24px;
}

.root--large .divider {
	height: 20px;
}

.root--large .item {
	gap: 9px;
}

.root--large .value {
	font-size: 19px;
}

.root--large .platform-label {
	font-size: 14px;
}

.root--vertical {
	flex-direction: column;
	align-items: stretch;
	gap: 10px;
	border-radius: 16px;
}

.root--vertical .divider {
	width: auto;
	height: 1px;
}

.root--vertical .item {
	justify-content: flex-start;
}

.root--vertical .value {
	margin-right: auto;
}

.root--vertical.root--compact {
	gap: 7px;
}
</style>
