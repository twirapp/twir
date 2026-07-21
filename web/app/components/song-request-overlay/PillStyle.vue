<script setup lang="ts">
import {
	type SongRequestOverlayVisualProps,
	getSongRequestPlaybackMetrics,
	getYouTubeThumbnailUrl,
} from './types'

const props = defineProps<SongRequestOverlayVisualProps>()

const metrics = computed(() => getSongRequestPlaybackMetrics(props.position, props.duration))
const thumbnailUrl = computed(() => getYouTubeThumbnailUrl(props.videoId))
</script>

<template>
	<div class="absolute inset-x-0 bottom-0 p-2 sm:p-4">
		<div
			class="flex min-h-14 w-full items-center gap-3 overflow-hidden rounded-full border border-white/10 bg-black/85 px-2 py-2 text-white shadow-xl backdrop-blur-sm sm:gap-4 sm:px-3"
		>
			<div class="bg-muted h-10 w-16 shrink-0 overflow-hidden rounded-full">
				<img
					v-if="thumbnailUrl"
					:src="thumbnailUrl"
					:alt="title ? `${title} thumbnail` : 'Song thumbnail'"
					class="size-full object-cover"
				/>
			</div>

			<p class="min-w-0 flex-1 truncate text-sm font-medium sm:text-base">
				{{ title }}
				<template v-if="requester">
					<span :style="{ color: accentColor }"> • </span>
					<span class="text-white/65">@{{ requester }}</span>
				</template>
			</p>

			<div class="w-20 shrink-0 sm:w-40">
				<div
					class="flex justify-between font-mono text-[10px] text-white/65 tabular-nums sm:text-xs"
				>
					<span>{{ metrics.formattedPosition }}</span>
					<span>{{ metrics.formattedDuration }}</span>
				</div>
				<div class="mt-1 h-1 overflow-hidden rounded-full bg-white/15">
					<div
						class="h-full rounded-full transition-[width] duration-200 motion-reduce:transition-none"
						:style="{ width: `${metrics.progress}%`, backgroundColor: accentColor }"
					/>
				</div>
			</div>
		</div>
	</div>
</template>
