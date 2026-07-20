<script setup lang="ts">
import {
	type SongRequestOverlayVisualProps,
	getSongRequestPlaybackMetrics,
	getYouTubeThumbnailUrl,
} from './types'

const props = defineProps<SongRequestOverlayVisualProps>()

defineSlots<{
	media?: (props: { thumbnailUrl: string }) => unknown
}>()

const metrics = computed(() => getSongRequestPlaybackMetrics(props.position, props.duration))
const thumbnailUrl = computed(() => getYouTubeThumbnailUrl(props.videoId))
</script>

<template>
	<div class="absolute inset-0 flex min-h-0 flex-col overflow-hidden">
		<div class="relative min-h-0 flex-1 overflow-hidden">
			<slot
				name="media"
				:thumbnail-url="thumbnailUrl"
			>
				<img
					v-if="thumbnailUrl"
					:src="thumbnailUrl"
					:alt="title ? `${title} thumbnail` : 'Song thumbnail'"
					class="size-full object-cover"
				/>
				<div
					v-else
					class="bg-muted size-full"
				/>
			</slot>
		</div>

		<div class="shrink-0 bg-black px-4 py-3 text-white sm:px-6 sm:py-4">
			<div class="flex min-w-0 items-end justify-between gap-4">
				<div class="min-w-0">
					<p class="truncate text-base font-semibold sm:text-xl">
						{{ title }}
					</p>
					<p
						v-if="requester"
						class="truncate text-xs text-white/65 sm:text-sm"
					>
						Requested by {{ requester }}
					</p>
				</div>
				<p class="shrink-0 font-mono text-xs text-white/75 tabular-nums sm:text-sm">
					{{ metrics.formattedPosition }} / {{ metrics.formattedDuration }}
				</p>
			</div>
			<div class="mt-3 h-1 overflow-hidden rounded-full bg-white/15">
				<div
					class="h-full rounded-full transition-[width] duration-200 motion-reduce:transition-none"
					:style="{ width: `${metrics.progress}%`, backgroundColor: accentColor }"
				/>
			</div>
		</div>
	</div>
</template>
