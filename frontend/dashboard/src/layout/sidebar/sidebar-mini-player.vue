<script setup lang="ts">
import { ListMusic, Pause, Play, SkipForward, Volume2, VolumeX } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRouter } from 'vue-router'

import { useYoutubeSocket } from '@/components/songRequests/hook.js'
import { useGlobalYoutubePlayer } from '@/composables/useGlobalYoutubePlayer.js'
import { Button } from '@/components/ui/button'
import { Slider } from '@/components/ui/slider'
import { convertMillisToTime } from '@/helpers/convertMillisToTime.js'
const { currentVideo } = useYoutubeSocket()
const {
	isPlaying,
	hasEverPlayedSong,
	sliderTime,
	duration,
	sliderVolume,
	isMuted,
	togglePlay,
	playNext,
	seek,
	setVolume,
	toggleMute,
} = useGlobalYoutubePlayer()

const router = useRouter()

// Show mini-player only if:
// 1. User has ever played a song (to know they use the feature)
// 2. There's a current video
// 3. The video is actually playing OR has been played (has progress)
// 4. Not on the song-requests page
const shouldShowMiniPlayer = computed(() => {
	return (
		hasEverPlayedSong.value &&
		currentVideo.value != null &&
		currentVideo.value !== undefined &&
		(isPlaying.value || sliderTime.value > 0) &&
		router.currentRoute.value.path !== '/dashboard/song-requests'
	)
})

const formattedTime = computed(() => {
	return `${convertMillisToTime(sliderTime.value * 1000)} / ${convertMillisToTime(duration.value * 1000)}`
})

function goToSongRequests() {
	router.push('/dashboard/song-requests')
}

function handleSeek(value: number[] | undefined) {
	if (!value) return
	seek(value[0])
}

function handleVolumeChange(value: number[] | undefined) {
	if (!value) return
	setVolume(value[0])
}
</script>

<template>
	<div
		v-if="shouldShowMiniPlayer"
		class="flex flex-col gap-2 px-2 py-2 border-t"
	>
		<!-- Song info -->
		<button
			class="flex items-center gap-2 min-w-0 hover:bg-accent/50 rounded p-1.5 transition-colors text-left"
			@click="goToSongRequests"
		>
			<ListMusic class="size-4 shrink-0 text-muted-foreground" />
			<div class="flex flex-col min-w-0 flex-1">
				<span class="text-xs font-medium truncate">{{ currentVideo?.title }}</span>
				<span class="text-[10px] text-muted-foreground truncate">
					{{ currentVideo?.orderedByDisplayName || currentVideo?.orderedByName }}
				</span>
			</div>
		</button>

		<!-- Progress slider -->
		<div class="flex flex-col gap-1 px-1">
			<Slider
				:model-value="[sliderTime]"
				:step="1"
				:max="duration || 1"
				:disabled="!currentVideo"
				class="w-full"
				@update:model-value="handleSeek"
			/>
			<span class="text-[10px] text-muted-foreground text-center">{{ formattedTime }}</span>
		</div>

		<!-- Controls -->
		<div class="flex items-center gap-1 px-1">
			<Button
				size="icon"
				variant="ghost"
				class="size-7"
				:disabled="!currentVideo"
				@click="togglePlay"
			>
				<Play
					v-if="!isPlaying"
					class="size-3.5"
				/>
				<Pause
					v-else
					class="size-3.5"
				/>
			</Button>

			<Button
				size="icon"
				variant="ghost"
				class="size-7"
				:disabled="!currentVideo"
				@click="playNext"
			>
				<SkipForward class="size-3.5" />
			</Button>

			<div class="flex-1" />

			<Button
				size="icon"
				variant="ghost"
				class="size-7"
				@click="toggleMute"
			>
				<Volume2
					v-if="!isMuted"
					class="size-3.5"
				/>
				<VolumeX
					v-else
					class="size-3.5"
				/>
			</Button>
		</div>

		<!-- Volume slider -->
		<div class="px-1">
			<Slider
				:model-value="[sliderVolume]"
				:step="1"
				:max="100"
				class="w-full"
				@update:model-value="handleVolumeChange"
			/>
		</div>
	</div>
</template>
