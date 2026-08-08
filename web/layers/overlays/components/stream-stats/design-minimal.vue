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
			return 13
		case StreamStatsOverlayVariant.Large:
			return 19
		default:
			return 15
	}
})
</script>

<template>
	<div class="minimal" :class="rootClass">
		<div v-for="item in items" :key="item.id" class="item">
			<Icon
				:name="counterIcons[item.key]"
				class="icon"
				:style="{ width: iconSize + 'px', height: iconSize + 'px', color: item.color }"
			/>
			<RollingNumber :value="item.value" class="value" />
			<template v-if="item.platformColor">
				<PlatformIcon
					v-if="platformIcons && item.platform"
					:platform="item.platform"
					:size="iconSize - 2"
					:style="{ color: item.platformColor }"
				/>
				<span v-else class="platform-label" :style="{ color: item.platformColor }">
					{{ item.label }}
				</span>
			</template>
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

.root--compact {
	gap: 13px;
}

.root--compact .item {
	gap: 5px;
}

.root--compact .value {
	font-size: 12px;
}

.root--compact .platform-label {
	font-size: 10px;
}

.root--large {
	gap: 24px;
}

.root--large .item {
	gap: 8px;
}

.root--large .value {
	font-size: 19px;
}

.root--large .platform-label {
	font-size: 14px;
}

.root--vertical {
	flex-direction: column;
	align-items: flex-start;
	gap: 8px;
}

.root--vertical.root--compact {
	gap: 6px;
}
</style>
