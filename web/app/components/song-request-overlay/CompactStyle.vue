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
	<div class="absolute inset-x-0 bottom-0 p-2 sm:p-4">
		<div class="w-full max-w-[480px] overflow-hidden rounded-xl text-white shadow-2xl">
			<div
				class="aspect-video w-full overflow-hidden"
				:class="{ 'bg-muted': !$slots.media }"
			>
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
				</slot>
			</div>

			<div class="h-19 bg-black/85 p-3 backdrop-blur-sm sm:h-23 sm:p-4">
				<div class="flex min-w-0 items-end justify-between gap-3">
					<div class="min-w-0">
						<p class="truncate text-sm font-semibold sm:text-base">
							{{ title }}
						</p>
						<p
							v-if="requester"
							class="truncate text-xs text-white/60 sm:text-sm"
						>
							Requested by {{ requester }}
						</p>
					</div>
					<p class="shrink-0 font-mono text-xs text-white/70 tabular-nums">
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
	</div>
</template>
