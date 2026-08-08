<script setup lang="ts">
import type {
	StreamStatsCounterItem,
	StreamStatsCounterKey,
} from '../../composables/stream-stats/use-stream-stats-counters.js'
import { StreamStatsOverlayVariant } from '~/gql/graphql.js'

import { variantRootClass } from './counter-meta.js'
import RollingNumber from './rolling-number.vue'

const props = withDefaults(
	defineProps<{
		items: StreamStatsCounterItem[]
		variant?: StreamStatsOverlayVariant
	}>(),
	{
		variant: StreamStatsOverlayVariant.Horizontal,
	},
)

const rootClass = computed(() => variantRootClass(props.variant))

const terminalLabels: Record<StreamStatsCounterKey, string> = {
	viewers: 'viewers',
	messages: 'msgs',
	uptime: 'uptime',
	subscribers: 'subs',
	followers: 'followers',
}

function itemLabel(item: StreamStatsCounterItem): string {
	return item.platformColor ? item.label.toLowerCase() : terminalLabels[item.key]
}
</script>

<template>
	<div class="terminal" :class="rootClass">
		<span v-for="item in items" :key="item.id" class="line">
			<span class="label" :style="item.platformColor ? { color: item.platformColor } : {}">
				{{ itemLabel(item) }}
			</span>
			<span class="bracket">[</span>
			<RollingNumber :value="item.value" class="value" />
			<span class="bracket">]</span>
		</span>
		<span class="cursor">_</span>
	</div>
</template>

<style scoped>
.terminal {
	display: inline-flex;
	align-items: center;
	gap: 2ch;
	padding: 8px 14px;
	border: 1px solid rgba(34, 197, 94, 0.25);
	border-radius: 6px;
	background-color: rgba(3, 10, 3, 0.75);
	box-shadow: 0 4px 16px rgba(0, 0, 0, 0.35);
	color: #22c55e;
	font-family: 'JetBrains Mono', 'Fira Code', ui-monospace, Consolas, monospace;
	font-size: 13px;
	line-height: 1.5;
}

.line {
	display: inline-flex;
	align-items: baseline;
	gap: 1ch;
	white-space: nowrap;
}

.label {
	color: rgba(34, 197, 94, 0.55);
}

.bracket {
	color: rgba(34, 197, 94, 0.55);
}

.value {
	font-weight: 700;
	color: #22c55e;
}

.cursor {
	color: #22c55e;
	font-weight: 700;
	animation: terminal-blink 1s step-end infinite;
}

@keyframes terminal-blink {
	0%,
	49% {
		opacity: 1;
	}
	50%,
	100% {
		opacity: 0;
	}
}

.root--compact {
	padding: 6px 10px;
	font-size: 11px;
}

.root--large {
	padding: 11px 18px;
	font-size: 17px;
}

.root--vertical {
	flex-direction: column;
	align-items: flex-start;
	gap: 0;
}

@media (prefers-reduced-motion: reduce) {
	.cursor {
		animation: none;
	}
}
</style>
