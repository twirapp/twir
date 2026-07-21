<script setup lang="ts">
import {
	type SongRequestOverlayVisualProps,
	getSongRequestPlaybackMetrics,
	getYouTubeThumbnailUrl,
} from './types'

const props = defineProps<SongRequestOverlayVisualProps>()

const metrics = computed(() => getSongRequestPlaybackMetrics(props.position, props.duration))
const thumbnailUrl = computed(() => getYouTubeThumbnailUrl(props.videoId))
const waveformHeights = ['h-2', 'h-4', 'h-3', 'h-5', 'h-2', 'h-4', 'h-3', 'h-5', 'h-3']
</script>

<template>
	<div class="absolute inset-x-0 bottom-0 p-2 sm:p-4">
		<div
			class="flex w-full max-w-[680px] overflow-hidden rounded-xl border border-white/10 bg-black/85 text-white shadow-2xl backdrop-blur-sm"
		>
			<div class="bg-muted size-28 shrink-0 overflow-hidden sm:size-36">
				<img
					v-if="thumbnailUrl"
					:src="thumbnailUrl"
					:alt="title ? `${title} thumbnail` : 'Song thumbnail'"
					class="size-full object-cover"
				/>
			</div>

			<div class="flex min-w-0 flex-1 flex-col justify-between p-3 sm:p-4">
				<div class="min-w-0">
					<p class="truncate text-sm font-semibold sm:text-lg">
						{{ title }}
					</p>
					<p
						v-if="requester"
						class="truncate text-xs text-white/60 sm:text-sm"
					>
						Requested by {{ requester }}
					</p>
				</div>

				<div>
					<div
						class="flex items-center gap-2 font-mono text-[10px] text-white/70 tabular-nums sm:gap-3 sm:text-xs"
					>
						<span>{{ metrics.formattedPosition }}</span>
						<div
							aria-hidden="true"
							class="flex h-6 min-w-0 flex-1 items-center justify-center gap-1 overflow-hidden"
						>
							<span
								v-for="(height, index) in waveformHeights"
								:key="index"
								class="w-1 rounded-full opacity-80 motion-reduce:animate-none"
								:class="[height, { 'animate-pulse': isPlaying }]"
								:style="{ backgroundColor: accentColor, animationDelay: `${index * 90}ms` }"
							/>
						</div>
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
	</div>
</template>
