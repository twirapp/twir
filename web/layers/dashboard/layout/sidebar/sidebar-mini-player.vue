<script setup lang="ts">
import { useProfile } from '~~/layers/dashboard/api/auth'
import { useSongRequestGql } from '~~/layers/dashboard/composables/useSongRequestGql.js'
import { useGlobalYoutubePlayer } from '~~/layers/dashboard/composables/useGlobalYoutubePlayer.js'
import { Button } from '@/components/ui/button'
import { Slider } from '@/components/ui/slider'
import { convertMillisToTime } from '~~/layers/dashboard/helpers/convertMillisToTime.js'

const { data: profile } = useProfile()
const channelId = computed(() => profile.value?.selectedDashboardId ?? '')

const {
	playbackState,
	play,
	pause,
	skip,
	setVolume: setVolumeGql,
} = useSongRequestGql(channelId)

const {
	isPlaying,
	sliderTime,
	duration,
	sliderVolume,
	isMuted,
	playVideo,
	pauseVideo,
	seek,
	setVolume,
	toggleMute,
} = useGlobalYoutubePlayer()

const router = useRouter()
const localePath = useLocalePath()

const shouldShowMiniPlayer = computed(() => {
	return (
		playbackState.value != null &&
		router.currentRoute.value.path !== localePath('/dashboard/song-requests')
	)
})

const formattedTime = computed(() => {
	return `${convertMillisToTime(sliderTime.value * 1000)} / ${convertMillisToTime(duration.value * 1000)}`
})

function goToSongRequests() {
	router.push(localePath('/dashboard/song-requests'))
}

function handlePlayPause() {
	if (isPlaying.value) {
		pause()
		pauseVideo()
	} else {
		play(playbackState.value!.videoId)
		playVideo()
	}
}

function handleSkip() {
	skip()
}

function handleSeek(value: number[] | undefined) {
	if (!value) return
	seek(value[0])
}

function handleVolumeChange(value: number[] | undefined) {
	if (!value) return
	setVolume(value[0])
	setVolumeGql(value[0])
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
			<Icon name="lucide:list-music" class="size-4 shrink-0 text-muted-foreground" />
			<div class="flex flex-col min-w-0 flex-1">
				<span class="text-xs font-medium truncate">{{ playbackState?.title }}</span>
			</div>
		</button>

		<!-- Progress slider -->
		<div class="flex flex-col gap-1 px-1">
			<Slider
				:model-value="[sliderTime]"
				:step="1"
				:max="duration || 1"
				:disabled="!playbackState"
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
				:disabled="!playbackState"
				@click="handlePlayPause"
			>
				<Icon name="lucide:play"
					v-if="!isPlaying"
					class="size-3.5" />
				<Icon name="lucide:pause"
					v-else
					class="size-3.5" />
			</Button>

			<Button
				size="icon"
				variant="ghost"
				class="size-7"
				:disabled="!playbackState"
				@click="handleSkip"
			>
				<Icon name="lucide:skip-forward" class="size-3.5" />
			</Button>

			<div class="flex-1" />

			<Button
				size="icon"
				variant="ghost"
				class="size-7"
				@click="toggleMute"
			>
				<Icon name="lucide:volume2"
					v-if="!isMuted"
					class="size-3.5" />
				<Icon name="lucide:volume-x"
					v-else
					class="size-3.5" />
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
