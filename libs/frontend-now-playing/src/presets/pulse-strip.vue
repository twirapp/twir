<script setup lang="ts">
import { useElementSize } from '@vueuse/core'
import { computed, ref, toRef } from 'vue'

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
const metadataViewport = ref<HTMLElement | null>(null)
const metadata = ref<HTMLElement | null>(null)

const { timing, progressStyle } = useTrackProgress(track)
const { paletteStyle } = useArtworkPalette(imageUrl, backgroundColor)
const { width: viewportWidth } = useElementSize(metadataViewport)
const { width: metadataWidth } = useElementSize(metadata)
const metadataKey = computed(() => JSON.stringify([
	track.value?.title ?? '',
	track.value?.artist ?? '',
]))
const marqueeEnabled = computed(() => {
	const intrinsicWidth = Math.max(
		metadataWidth.value,
		metadata.value?.scrollWidth ?? 0,
	)
	return Math.ceil(intrinsicWidth) > Math.floor(viewportWidth.value)
})
</script>

<template>
	<div
		v-if="track"
		class="pulse-shell"
		:data-mode="timing.mode"
		:style="paletteStyle"
	>
		<div
			class="pulse-card np-palette-surface"
			:class="{ 'pulse-card--without-image': !settings.showImage }"
		>
			<img
				v-if="settings.showImage"
				class="pulse-cover image"
				:src="imageUrl ?? '/overlays/images/play.png'"
				alt=""
				aria-hidden="true"
			>

			<div ref="metadataViewport" class="pulse-viewport">
				<div
					:key="metadataKey"
					class="pulse-marquee-track"
					:class="{ 'pulse-marquee-track--active': marqueeEnabled }"
				>
					<div
						ref="metadata"
						class="pulse-metadata pulse-metadata--primary"
						:title="`${track.title} • ${track.artist}`"
					>
						<h2 class="name">{{ track.title }}</h2>
						<span class="pulse-separator" aria-hidden="true">•</span>
						<span class="artist">{{ track.artist }}</span>
					</div>

					<div
						v-if="marqueeEnabled"
						class="pulse-metadata pulse-metadata--duplicate"
						aria-hidden="true"
					>
						<h2 class="name">{{ track.title }}</h2>
						<span class="pulse-separator" aria-hidden="true">•</span>
						<span class="artist">{{ track.artist }}</span>
					</div>
				</div>
			</div>

			<div
				class="pulse-progress np-progress np-palette-surface"
				:data-mode="timing.mode"
				:style="progressStyle"
				role="progressbar"
				aria-label="Playback progress"
				aria-valuemin="0"
				aria-valuemax="100"
				:aria-valuenow="timing.mode === 'timed' ? Math.round(timing.percent) : undefined"
			/>
		</div>
	</div>
</template>

<style scoped>
.pulse-shell {
	container-type: inline-size;
	container-name: pulse-strip;
	width: 100%;
	max-width: 46rem;
	background: transparent;
}

.pulse-card {
	display: grid;
	grid-template-columns: auto minmax(0, 1fr) minmax(5.5rem, 0.34fr);
	align-items: center;
	gap: clamp(0.625rem, 2.5cqi, 1rem);
	padding: clamp(0.5rem, 2cqi, 0.75rem);
	border: 1px solid color-mix(in srgb, var(--np-accent) 18%, transparent);
	border-radius: clamp(0.75rem, 3cqi, 1rem);
	background-color: color-mix(in srgb, var(--np-surface) 94%, transparent);
	box-shadow: 0 0.625rem 1.75rem color-mix(in srgb, var(--np-glow) 42%, transparent);
	color: var(--np-text);
}

.pulse-card--without-image {
	grid-template-columns: minmax(0, 1fr) minmax(5.5rem, 0.34fr);
}

.pulse-cover {
	width: clamp(2.25rem, 8cqi, 2.75rem);
	aspect-ratio: 1;
	object-fit: cover;
	border-radius: clamp(0.5rem, 2cqi, 0.75rem);
	background-color: var(--np-surface-alt);
}

.pulse-viewport {
	min-width: 0;
	overflow: hidden;
}

.pulse-marquee-track {
	--pulse-marquee-gap: clamp(1.5rem, 6cqi, 3rem);
	--pulse-marquee-half-gap: clamp(0.75rem, 3cqi, 1.5rem);

	display: flex;
	width: max-content;
	align-items: baseline;
	gap: var(--pulse-marquee-gap);
}

.pulse-marquee-track--active {
	animation: pulse-marquee 12s linear infinite;
	will-change: transform;
}

.pulse-metadata {
	flex: 0 0 auto;
	width: max-content;
	max-width: none;
	white-space: nowrap;
}

.name {
	display: inline;
	margin: 0;
	font-size: clamp(0.875rem, 3.2cqi, 1.125rem);
	line-height: 1.2;
	color: var(--np-text);
}

.artist {
	font-size: clamp(0.75rem, 2.7cqi, 0.9375rem);
	line-height: 1.2;
	color: var(--np-muted);
}

.pulse-separator {
	margin-inline: clamp(0.375rem, 1.5cqi, 0.625rem);
	color: var(--np-accent);
}

.pulse-progress {
	min-width: 0;
}

@keyframes pulse-marquee {
	from {
		transform: translateX(0);
	}
	to {
		transform: translateX(calc(-50% - var(--pulse-marquee-half-gap)));
	}
}

@media (prefers-reduced-motion: reduce) {
	.pulse-marquee-track--active {
		width: 100%;
		animation: none;
		transform: none;
		will-change: auto;
	}

	.pulse-marquee-track--active .pulse-metadata--primary {
		width: 100%;
		min-width: 0;
		max-width: 100%;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.pulse-metadata--duplicate {
		display: none;
	}
}

@container pulse-strip (max-width: 420px) {
	.pulse-card,
	.pulse-card--without-image {
		grid-template-columns: auto minmax(0, 1fr);
	}

	.pulse-card--without-image .pulse-viewport {
		grid-column: 1 / -1;
	}

	.pulse-progress {
		grid-column: 1 / -1;
	}
}
</style>
