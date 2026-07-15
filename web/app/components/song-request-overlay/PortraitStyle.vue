<script setup lang="ts">
import {
	type SongRequestOverlayVisualProps,
	getSongRequestPlaybackMetrics,
	getYouTubeThumbnailUrl,
} from './types'

const props = defineProps<SongRequestOverlayVisualProps>()

const metrics = computed(() => getSongRequestPlaybackMetrics(props.position, props.duration))
const thumbnailUrl = computed(() => getYouTubeThumbnailUrl(props.videoId))
const cardStyle = computed(() => ({
	'--song-accent': props.accentColor,
}))
</script>

<template>
	<div class="absolute inset-x-0 bottom-0 p-2 sm:p-4">
		<div
			class="song-request-portrait-card w-full max-w-[360px] overflow-hidden rounded-2xl border border-white/10 bg-black/85 text-white backdrop-blur-sm"
			:style="cardStyle"
		>
			<div class="bg-muted aspect-video w-full overflow-hidden">
				<img
					v-if="thumbnailUrl"
					:src="thumbnailUrl"
					:alt="title ? `${title} thumbnail` : 'Song thumbnail'"
					class="size-full object-cover"
				/>
			</div>

			<div class="p-4 sm:p-5">
				<p class="line-clamp-2 text-base leading-tight font-semibold sm:text-lg">
					{{ title }}
				</p>
				<p
					v-if="requester"
					class="mt-1 truncate text-xs text-white/60 sm:text-sm"
				>
					Requested by {{ requester }}
				</p>

				<div
					class="mt-4 flex items-center justify-between font-mono text-xs text-white/70 tabular-nums"
				>
					<span>{{ metrics.formattedPosition }}</span>
					<span>{{ metrics.formattedDuration }}</span>
				</div>
				<div class="mt-2 h-1.5 overflow-hidden rounded-full bg-white/15">
					<div
						class="h-full rounded-full transition-[width] duration-200 motion-reduce:transition-none"
						:style="{ width: `${metrics.progress}%`, backgroundColor: accentColor }"
					/>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
.song-request-portrait-card {
	box-shadow: 0 16px 48px color-mix(in srgb, var(--song-accent) 20%, transparent);
}
</style>
