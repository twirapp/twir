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

const waveformBars = [
	{ key: 'signal-01', height: '42%', delay: '-1.1s', duration: '1.38s' },
	{ key: 'signal-02', height: '72%', delay: '-0.4s', duration: '1.12s' },
	{ key: 'signal-03', height: '54%', delay: '-0.8s', duration: '1.56s' },
	{ key: 'signal-04', height: '88%', delay: '-1.3s', duration: '1.28s' },
	{ key: 'signal-05', height: '64%', delay: '-0.2s', duration: '1.44s' },
	{ key: 'signal-06', height: '36%', delay: '-0.9s', duration: '1.18s' },
	{ key: 'signal-07', height: '78%', delay: '-1.5s', duration: '1.62s' },
	{ key: 'signal-08', height: '48%', delay: '-0.6s', duration: '1.32s' },
	{ key: 'signal-09', height: '92%', delay: '-1.2s', duration: '1.48s' },
	{ key: 'signal-10', height: '58%', delay: '-0.3s', duration: '1.24s' },
	{ key: 'signal-11', height: '70%', delay: '-1s', duration: '1.54s' },
	{ key: 'signal-12', height: '40%', delay: '-0.5s', duration: '1.16s' },
	{ key: 'signal-13', height: '82%', delay: '-1.4s', duration: '1.36s' },
	{ key: 'signal-14', height: '50%', delay: '-0.7s', duration: '1.58s' },
	{ key: 'signal-15', height: '68%', delay: '-0.1s', duration: '1.22s' },
	{ key: 'signal-16', height: '94%', delay: '-1.25s', duration: '1.46s' },
	{ key: 'signal-17', height: '46%', delay: '-0.45s', duration: '1.14s' },
	{ key: 'signal-18', height: '76%', delay: '-1.05s', duration: '1.52s' },
] as const

const { timing, elapsedLabel, durationLabel, progressStyle } = useTrackProgress(track)
const { paletteStyle } = useArtworkPalette(imageUrl, backgroundColor)
</script>

<template>
	<div
		v-if="track"
		class="signal-deck-shell"
		:data-mode="timing.mode"
		:style="paletteStyle"
	>
		<article
			class="signal-deck"
			:class="{ 'signal-deck--no-cover': !settings.showImage }"
		>
			<img
				v-if="settings.showImage"
				class="signal-deck-cover image"
				:src="imageUrl ?? '/overlays/images/play.png'"
				alt=""
				aria-hidden="true"
			>

			<div class="signal-deck-metadata np-palette-surface">
				<h2 class="name" :title="track.title">{{ track.title }}</h2>
				<div class="artist" :title="track.artist">{{ track.artist }}</div>
			</div>

			<div class="signal-deck-signal np-palette-surface">
				<span v-if="timing.mode === 'timed'" class="artist signal-deck-time">
					{{ elapsedLabel }}
				</span>
				<div class="signal-deck-waveform" aria-hidden="true">
					<span
						v-for="bar in waveformBars"
						:key="bar.key"
						class="signal-deck-bar"
						:style="{
							'--signal-height': bar.height,
							'--signal-delay': bar.delay,
							'--signal-duration': bar.duration,
						}"
					/>
				</div>
				<span v-if="timing.mode === 'timed'" class="artist signal-deck-time">
					{{ durationLabel }}
				</span>
			</div>

			<div
				class="signal-deck-progress np-progress np-palette-surface"
				:data-mode="timing.mode"
				:style="progressStyle"
				role="progressbar"
				aria-label="Playback progress"
				aria-valuemin="0"
				aria-valuemax="100"
				:aria-valuenow="timing.mode === 'timed' ? Math.round(timing.percent) : undefined"
			/>
		</article>
	</div>
</template>

<style scoped>
.signal-deck-shell {
	container-type: inline-size;
	container-name: signal-deck;
	width: min(100%, 46rem);
	background: transparent;
}

.signal-deck {
	display: grid;
	grid-template-columns: clamp(7rem, 24cqi, 9rem) minmax(0, 1fr);
	grid-template-rows: auto auto auto;
	gap: clamp(0.5rem, 1.8cqi, 0.75rem);
	padding: clamp(0.625rem, 2cqi, 0.875rem);
	border: 1px solid color-mix(in srgb, var(--np-accent) 20%, transparent);
	border-radius: clamp(0.875rem, 2.8cqi, 1.25rem);
	background-color: var(--np-surface);
	box-shadow: 0 0.875rem 2.25rem color-mix(in srgb, var(--np-glow) 32%, transparent);
	color: var(--np-text);
}

.signal-deck--no-cover {
	grid-template-columns: minmax(0, 1fr);
}

.signal-deck-cover {
	grid-column: 1;
	grid-row: 1 / 3;
	align-self: center;
	display: block;
	width: 100%;
	height: auto;
	aspect-ratio: 1;
	object-fit: cover;
	border-radius: clamp(0.625rem, 2cqi, 0.875rem);
	background-color: var(--np-surface-alt);
}

.signal-deck-metadata,
.signal-deck-signal {
	grid-column: 2;
	min-width: 0;
	border: 1px solid color-mix(in srgb, var(--np-muted) 13%, transparent);
	border-radius: clamp(0.625rem, 2cqi, 0.875rem);
	background-color: var(--np-surface-alt);
}

.signal-deck-metadata {
	display: grid;
	grid-row: 1;
	align-content: center;
	gap: clamp(0.1875rem, 1cqi, 0.375rem);
	padding: clamp(0.625rem, 2.2cqi, 0.875rem);
}

.signal-deck-signal {
	display: grid;
	grid-row: 2;
	grid-template-columns: auto minmax(0, 1fr) auto;
	align-items: center;
	gap: clamp(0.5rem, 1.8cqi, 0.75rem);
	padding: clamp(0.5rem, 1.8cqi, 0.75rem);
}

.signal-deck--no-cover .signal-deck-metadata,
.signal-deck--no-cover .signal-deck-signal {
	grid-column: 1;
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
	font-size: clamp(1rem, 3.4cqi, 1.375rem);
	line-height: 1.16;
	color: var(--np-text);
}

.artist {
	font-size: clamp(0.75rem, 2.5cqi, 0.9375rem);
	line-height: 1.25;
	color: var(--np-muted);
}

.signal-deck-time {
	font-size: clamp(0.6875rem, 2cqi, 0.8125rem);
	font-variant-numeric: tabular-nums;
	line-height: 1;
	color: var(--np-muted);
}

.signal-deck-waveform {
	display: flex;
	height: clamp(1.25rem, 4.5cqi, 2rem);
	min-width: 0;
	align-items: center;
	justify-content: center;
	gap: clamp(0.125rem, 0.65cqi, 0.25rem);
	overflow: hidden;
}

.signal-deck-bar {
	width: clamp(0.125rem, 0.55cqi, 0.25rem);
	height: var(--signal-height);
	flex: 0 0 auto;
	border-radius: 999px;
	background-color: var(--np-accent);
	box-shadow: 0 0 0.5rem color-mix(in srgb, var(--np-glow) 38%, transparent);
	transform-origin: center;
	animation: signal-deck-pulse var(--signal-duration) ease-in-out var(--signal-delay) infinite alternate;
}

[data-mode='ambient'] .signal-deck-signal {
	grid-template-columns: minmax(0, 1fr);
}

.signal-deck-progress {
	grid-column: 1 / -1;
	grid-row: 3;
	min-width: 0;
	margin-block-start: clamp(0.125rem, 0.6cqi, 0.25rem);
	background-color: var(--np-surface-alt);
}

@keyframes signal-deck-pulse {
	from {
		transform: scaleY(0.55);
	}
	to {
		transform: scaleY(1);
	}
}

@container signal-deck (max-width: 480px) {
	.signal-deck {
		grid-template-columns: 5.75rem minmax(0, 1fr);
	}

	.signal-deck--no-cover {
		grid-template-columns: minmax(0, 1fr);
	}

	.signal-deck-bar:nth-child(n + 13) {
		display: none;
	}
}

@container signal-deck (max-width: 180px) {
	.signal-deck {
		grid-template-columns: minmax(0, 1fr);
	}

	.signal-deck-cover {
		grid-column: 1;
		grid-row: 1;
		justify-self: center;
		width: min(100%, 5.75rem);
	}

	.signal-deck-metadata {
		grid-column: 1;
		grid-row: 2;
	}

	.signal-deck-signal {
		grid-column: 1;
		grid-row: 3;
	}

	.signal-deck-progress {
		grid-row: 4;
	}

	.signal-deck--no-cover .signal-deck-metadata {
		grid-row: 1;
	}

	.signal-deck--no-cover .signal-deck-signal {
		grid-row: 2;
	}

	.signal-deck--no-cover .signal-deck-progress {
		grid-row: 3;
	}
}

@media (prefers-reduced-motion: reduce) {
	.signal-deck-bar {
		animation: none;
		transform: scaleY(0.82);
	}
}
</style>
