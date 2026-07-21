<script setup lang="ts">
import { computed, toRef } from 'vue'

import { useArtworkPalette } from '../composables/use-artwork-palette.js'
import { useTrackProgress } from '../composables/use-track-progress.js'
import type { Settings, Track } from '../types.js'

interface Props {
	track?: Track | null
	settings: Settings
}

const props = defineProps<Props>()
const track = toRef(props, 'track')
const imageUrl = computed(() => track.value?.imageUrl)
const backgroundColor = computed(() => props.settings.backgroundColor)

const { timing, elapsedLabel, durationLabel, progressStyle } = useTrackProgress(track)
const { paletteStyle } = useArtworkPalette(imageUrl, backgroundColor)
</script>

<template>
	<div
		v-if="track"
		class="aura-shell"
		:data-mode="timing.mode"
		:style="paletteStyle"
	>
		<article class="aura-card np-palette-surface">
			<img
				v-if="settings.showImage"
				class="aura-cover image"
				:src="imageUrl ?? '/overlays/images/play.png'"
				alt=""
				aria-hidden="true"
			>

			<div class="aura-metadata np-palette-surface">
				<h2 class="name">{{ track.title }}</h2>
				<div class="artist">{{ track.artist }}</div>
			</div>

			<div class="aura-timing np-palette-surface">
				<span v-if="timing.mode === 'timed'" class="artist aura-time">{{ elapsedLabel }}</span>
				<div
					class="np-progress np-palette-surface"
					:data-mode="timing.mode"
					:style="progressStyle"
					role="progressbar"
					aria-label="Playback progress"
					aria-valuemin="0"
					aria-valuemax="100"
					:aria-valuenow="timing.mode === 'timed' ? Math.round(timing.percent) : undefined"
				/>
				<span v-if="timing.mode === 'timed'" class="artist aura-time">{{ durationLabel }}</span>
			</div>
		</article>
	</div>
</template>

<style scoped>
.aura-shell {
	container-type: inline-size;
	container-name: aura-stack;
	width: min(100%, 22.5rem);
	background: transparent;
}

.aura-card {
	display: grid;
	gap: clamp(0.625rem, 3cqi, 0.875rem);
	padding: clamp(0.625rem, 3cqi, 0.875rem);
	border: 1px solid color-mix(in srgb, var(--np-accent) 20%, transparent);
	border-radius: clamp(1rem, 5cqi, 1.5rem);
	background-color: color-mix(in srgb, var(--np-surface) 94%, transparent);
	box-shadow: 0 1rem 2.5rem color-mix(in srgb, var(--np-glow) 38%, transparent);
	color: var(--np-text);
}

.aura-cover {
	width: 100%;
	aspect-ratio: 1;
	object-fit: cover;
	border-radius: clamp(0.75rem, 4cqi, 1.125rem);
	background-color: var(--np-surface-alt);
}

.aura-metadata,
.aura-timing {
	min-width: 0;
	border-radius: clamp(0.75rem, 4cqi, 1rem);
	background-color: var(--np-surface-alt);
}

.aura-metadata {
	display: grid;
	gap: clamp(0.25rem, 1.5cqi, 0.4375rem);
	padding: clamp(0.75rem, 4cqi, 1rem);
}

.name,
.artist {
	min-width: 0;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

.name {
	margin: 0;
	font-size: clamp(1rem, 5cqi, 1.375rem);
	line-height: 1.18;
	color: var(--np-text);
}

.artist {
	font-size: clamp(0.8125rem, 3.6cqi, 1rem);
	line-height: 1.3;
	color: var(--np-muted);
}

.aura-timing {
	display: grid;
	grid-template-columns: auto minmax(0, 1fr) auto;
	align-items: center;
	gap: clamp(0.5rem, 2.5cqi, 0.75rem);
	padding: clamp(0.625rem, 3cqi, 0.875rem);
}

.aura-time {
	font-size: clamp(0.6875rem, 3cqi, 0.8125rem);
	font-variant-numeric: tabular-nums;
	line-height: 1;
	color: var(--np-muted);
}

[data-mode='ambient'] .aura-timing {
	grid-template-columns: minmax(0, 1fr);
}

@container aura-stack (max-width: 260px) {
	.aura-card {
		gap: 0.5rem;
		padding: 0.5rem;
	}

	.aura-metadata,
	.aura-timing {
		padding: 0.625rem;
	}
}
</style>
