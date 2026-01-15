<script lang="ts" setup>
import {
	Ban,
	Eye,
	EyeOff,
	Link,
	ListMusic,
	Loader2,
	Pause,
	Play,
	Settings,
	SkipForward,
	User,
	Volume2,
	VolumeX,
} from 'lucide-vue-next'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { useProfile } from '@/api/auth'
import { useYoutubeSocket } from '@/components/songRequests/hook.js'
import { useGlobalYoutubePlayer } from '@/composables/useGlobalYoutubePlayer.js'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardTitle } from '@/components/ui/card'
import { Slider } from '@/components/ui/slider'
import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
	AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { convertMillisToTime } from '@/helpers/convertMillisToTime.js'

const props = defineProps<{
	noCookie: boolean
	openSettingsModal: () => void
}>()

const { currentVideo, banSong, banUser } = useYoutubeSocket()

const {
	playerReady,
	isPlaying,
	sliderTime,
	duration,
	playerVisible,
	isMuted,
	sliderVolume,
	noCookie,
	initPlayer,
	togglePlay,
	seek,
	setVolume,
	toggleMute,
	playNext,
} = useGlobalYoutubePlayer()

const playerContainerRef = ref<HTMLElement | null>(null)

// Update noCookie setting
watch(() => props.noCookie, (value) => {
	noCookie.value = value
}, { immediate: true })

function updatePlayerPosition() {
	if (!playerContainerRef.value) return

	const rect = playerContainerRef.value.getBoundingClientRect()
	const globalContainer = document.getElementById('global-yt-player-container')

	if (globalContainer) {
		globalContainer.style.left = `${rect.left}px`
		globalContainer.style.top = `${rect.top}px`
		globalContainer.style.width = `${rect.width}px`
		globalContainer.style.height = `${rect.height}px`
		console.log('[Player] Updated position:', rect)
	}
}

onMounted(async () => {
	console.log('[Player] Mounted, initializing player if needed')

	// Initialize player if not already initialized
	if (!playerReady.value) {
		console.log('[Player] Initializing player...')
		await initPlayer()
	}

	// Wait a bit for DOM to settle
	await new Promise(resolve => setTimeout(resolve, 100))

	// Position the global player over our container
	updatePlayerPosition()

	// Update position on window resize
	window.addEventListener('resize', updatePlayerPosition)
	window.addEventListener('scroll', updatePlayerPosition, true)
})

onUnmounted(() => {
	console.log('[Player] Unmounted, moving iframe off-screen')
	window.removeEventListener('resize', updatePlayerPosition)
	window.removeEventListener('scroll', updatePlayerPosition, true)

	// Reset iframe position to off-screen
	const globalContainer = document.getElementById('global-yt-player-container')
	if (globalContainer) {
		globalContainer.style.left = '-9999px'
		globalContainer.style.top = '-9999px'
		globalContainer.style.width = '640px'
		globalContainer.style.height = '360px'
		console.log('[Player] Iframe moved off-screen')
	}
})

const formattedTime = computed(() => {
	return `${convertMillisToTime(sliderTime.value * 1000)} / ${convertMillisToTime(duration.value * 1000)}`
})

function handleSeek(value: number[] | undefined) {
	if (!value) return
	seek(value[0])
}

function handleVolumeChange(value: number[] | undefined) {
	if (!value) return
	setVolume(value[0])
}

const { data: profile } = useProfile()
const { t } = useI18n()
</script>

<template>
	<Card class="p-0">
		<CardContent class="p-0">
			<div class="flex flex-row justify-between items-center px-2 py-2 border-b">
				<CardTitle class="text-base">{{ t('songRequests.player.title') }}</CardTitle>
				<div class="flex gap-1">
					<Button
						variant="outline"
						size="icon"
						class="size-8"
						@click="playerVisible = !playerVisible"
					>
						<EyeOff v-if="playerVisible" class="size-4" />
						<Eye v-else class="size-4" />
					</Button>
					<Button variant="outline" size="icon" class="size-8" @click="openSettingsModal">
						<Settings class="size-4" />
					</Button>
				</div>
			</div>
			<div v-if="profile?.id !== profile?.selectedDashboardId" class="p-6 text-center">
				<p class="text-muted-foreground">{{ t('songRequests.player.noAccess') }}</p>
			</div>

			<div v-else>
				<div
					v-show="playerVisible"
					ref="playerContainerRef"
					class="h-[300px] overflow-hidden relative bg-black"
				>
					<!-- Global YouTube iframe will be positioned over this container -->
				</div>

				<div class="flex flex-col gap-4 py-5 px-6">
					<div class="flex gap-2 items-center">
						<Button
							size="icon"
							class="flex size-8 min-w-8"
							variant="secondary"
							:disabled="currentVideo == null"
							@click="togglePlay"
						>
							<Play v-if="!isPlaying" class="size-4" />
							<Pause v-else class="size-4" />
						</Button>

						<Button
							size="icon"
							class="flex size-8 min-w-8"
							variant="secondary"
							:disabled="currentVideo == null"
							@click="playNext"
						>
							<SkipForward class="size-4" />
						</Button>

						<Slider
							:model-value="[sliderTime]"
							:step="1"
							:max="duration || 1"
							:disabled="!currentVideo"
							@update:model-value="handleSeek"
						/>
						<span class="text-xs text-muted-foreground whitespace-nowrap">{{ formattedTime }}</span>
					</div>

					<div class="flex items-center gap-2">
						<Button
							size="icon"
							variant="secondary"
							class="size-8 min-w-8"
							@click="toggleMute"
						>
							<Volume2 v-if="!isMuted" class="size-4" />
							<VolumeX v-else class="size-4" />
						</Button>
						<Slider
							:model-value="[sliderVolume]"
							:step="1"
							:max="100"
							@update:model-value="handleVolumeChange"
						/>
					</div>
				</div>
			</div>
		</CardContent>

		<CardFooter class="flex-col items-start gap-4 border-t pt-2">
			<template v-if="currentVideo">
				<div class="flex flex-col gap-2 w-full">
					<div class="flex items-center gap-2">
						<ListMusic class="size-4 shrink-0" />
						<span class="truncate">{{ currentVideo?.title }}</span>
					</div>

					<div class="flex items-center gap-2">
						<User class="size-4 shrink-0" />
						<span>{{ currentVideo?.orderedByDisplayName || currentVideo?.orderedByName }}</span>
					</div>

					<div class="flex items-center gap-2">
						<Link class="size-4 shrink-0" />
						<a
							class="underline text-sm truncate"
							:href="currentVideo.songLink ?? `https://youtu.be/${currentVideo?.videoId}`"
							target="_blank"
						>
							{{ currentVideo.songLink || `youtu.be/${currentVideo?.videoId}` }}
						</a>
					</div>
				</div>

				<div class="flex gap-2 mb-2 justify-end w-full">
					<AlertDialog>
						<AlertDialogTrigger as-child>
							<Button size="sm" variant="destructive">
								<Ban class="size-4 mr-1" />
								{{ t('songRequests.ban.song') }}
							</Button>
						</AlertDialogTrigger>
						<AlertDialogContent>
							<AlertDialogHeader>
								<AlertDialogTitle>{{ t('songRequests.ban.songConfirm') }}</AlertDialogTitle>
							</AlertDialogHeader>
							<AlertDialogFooter>
								<AlertDialogCancel>{{ t('deleteConfirmation.cancel') }}</AlertDialogCancel>
								<AlertDialogAction @click="banSong(currentVideo.videoId)">
									{{ t('deleteConfirmation.confirm') }}
								</AlertDialogAction>
							</AlertDialogFooter>
						</AlertDialogContent>
					</AlertDialog>

					<AlertDialog>
						<AlertDialogTrigger as-child>
							<Button size="sm" variant="destructive">
								<Ban class="size-4 mr-1" />
								{{ t('songRequests.ban.user') }}
							</Button>
						</AlertDialogTrigger>
						<AlertDialogContent>
							<AlertDialogHeader>
								<AlertDialogTitle>{{ t('songRequests.ban.userConfirm') }}</AlertDialogTitle>
							</AlertDialogHeader>
							<AlertDialogFooter>
								<AlertDialogCancel>{{ t('deleteConfirmation.cancel') }}</AlertDialogCancel>
								<AlertDialogAction @click="banUser(currentVideo.orderedById)">
									{{ t('deleteConfirmation.confirm') }}
								</AlertDialogAction>
							</AlertDialogFooter>
						</AlertDialogContent>
					</AlertDialog>
				</div>
			</template>

			<div v-else class="flex flex-col items-center justify-center w-full py-4 gap-2">
				<Loader2 class="size-6 animate-spin text-muted-foreground" />
				<p class="text-muted-foreground text-sm">{{ t('songRequests.waiting') }}</p>
			</div>
		</CardFooter>
	</Card>
</template>
