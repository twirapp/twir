<script lang="ts" setup>
import { useProfile } from '~~/layers/dashboard/api/auth'
import { useSongRequestsApi } from '~~/layers/dashboard/api/song-requests'
import { useDebouncedSeek } from '~~/layers/dashboard/composables/useDebouncedSeek.js'
import { useDebouncedVolume } from '~~/layers/dashboard/composables/useDebouncedVolume.js'
import { useSongRequestGql } from '~~/layers/dashboard/composables/useSongRequestGql.js'
import { convertMillisToTime } from '~~/layers/dashboard/helpers/convertMillisToTime.js'

import type { SongRequestsSettingsOpts } from '~/gql/graphql.js'

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
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardTitle } from '@/components/ui/card'
import { Slider } from '@/components/ui/slider'

const props = defineProps<{
	openSettingsModal: () => void
}>()

const { data: profile } = useProfile()
const channelId = computed(() => profile.value?.selectedDashboardId ?? '')

const {
	playbackState,
	queue,
	play,
	pause,
	skip,
	setVolume,
	deleteFromQueue,
	updatePosition,
} = useSongRequestGql(channelId)

const currentVideo = computed(() => {
	return playbackState.value?.videoId ? playbackState.value : null
})

const currentQueueItem = computed(() => {
	if (!currentVideo.value) return null
	return queue.value.find((item) => item.id === currentVideo.value?.videoId) ?? null
})

const hasPlayableVideo = computed(() => currentVideo.value != null || queue.value.length > 0)

const duration = computed(() => currentQueueItem.value?.durationSeconds ?? 0)

const {
	localPosition,
	handleSeekInput,
	syncFromServer: syncSeekFromServer,
} = useDebouncedSeek(
	() => channelId.value,
	({ position }) => updatePosition(position),
)

watch(
	() => currentVideo.value?.position,
	(position) => {
		if (position !== undefined) {
			syncSeekFromServer(position)
		}
	},
	{ immediate: true },
)

const {
	localVolume,
	handleVolumeInput,
	syncFromServer: syncVolumeFromServer,
} = useDebouncedVolume(
	() => channelId.value,
	({ volume }) => setVolume(volume),
)

watch(
	() => currentVideo.value?.volume,
	(volume) => {
		if (volume !== undefined) {
			syncVolumeFromServer(volume)
		}
	},
	{ immediate: true },
)

const progressPercent = computed(() => {
	const duration = currentQueueItem.value?.durationSeconds ?? 0
	if (!currentVideo.value || duration <= 0) return 0
	return Math.min((localPosition.value / duration) * 100, 100)
})

const formattedTime = computed(() => {
	const duration = currentQueueItem.value?.durationSeconds ?? 0
	return `${convertMillisToTime(localPosition.value * 1000)} / ${convertMillisToTime(duration * 1000)}`
})

function getFirstQueuedVideoId(): string {
	return queue.value[0]?.id ?? ''
}

function handlePlayPause() {
	if (currentVideo.value?.isPlaying) {
		void pause()
		return
	}

	const videoId = currentVideo.value?.videoId || getFirstQueuedVideoId()
	if (!videoId) return

	void play(videoId)
}

function handleSeek(value: number[] | undefined) {
	const position = value?.[0]
	if (position === undefined || !currentVideo.value) return

	handleSeekInput(position)
}

function handleVolumeChange(value: number[] | undefined) {
	const volume = value?.[0]
	if (volume === undefined) return

	handleVolumeInput(volume)
}

const youtubeModuleManager = useSongRequestsApi()
const { data: youtubeSettings } = youtubeModuleManager.useSongRequestQuery()
const youtubeModuleUpdater = youtubeModuleManager.useSongRequestMutation()

watch(
	() => youtubeSettings.value?.songRequests?.volume,
	(volume) => {
		if (volume !== undefined && !currentVideo.value) {
			syncVolumeFromServer(volume)
		}
	},
	{ immediate: true },
)

const { t } = useI18n()

function toSettingsOpts(settings: Record<string, unknown>): SongRequestsSettingsOpts {
	const { channelApiKey, __typename, ...rest } = settings
	return rest as unknown as SongRequestsSettingsOpts
}

async function banUser(userId: string) {
	if (!youtubeSettings.value?.songRequests) return
	const settings = toSettingsOpts(youtubeSettings.value.songRequests as Record<string, unknown>)

	await youtubeModuleUpdater.executeMutation({
		opts: {
			...settings,
			denyList: {
				...settings.denyList,
				users: [...settings.denyList!.users, userId],
			},
		},
	})

	const userVideos = queue.value.filter((video) => video.orderedByName === userId)
	for (const video of userVideos) {
		deleteFromQueue(video.id)
	}
}

async function banSong(videoId: string) {
	if (!youtubeSettings.value?.songRequests) return
	const settings = toSettingsOpts(youtubeSettings.value.songRequests as Record<string, unknown>)

	await youtubeModuleUpdater.executeMutation({
		opts: {
			...settings,
			denyList: {
				...settings.denyList,
				songs: [...settings.denyList!.songs, videoId],
			},
		},
	})

	deleteFromQueue(videoId)
}
</script>

<template>
	<Card class="p-0">
		<CardContent class="p-0">
			<div class="flex flex-row items-center justify-between border-b px-2 py-2">
				<CardTitle class="text-base">{{ t('songRequests.player.title') }}</CardTitle>
				<Button
					variant="outline"
					size="icon"
					class="size-8"
					@click="props.openSettingsModal"
				>
					<Icon
						name="lucide:settings"
						class="size-4"
					/>
				</Button>
			</div>

			<div class="flex flex-col gap-4 px-6 py-5">
				<div class="flex items-center gap-2">
					<Button
						size="icon"
						class="flex size-8 min-w-8"
						variant="secondary"
						:disabled="!hasPlayableVideo"
						@click="handlePlayPause"
					>
						<Icon
							:name="currentVideo?.isPlaying ? 'lucide:pause' : 'lucide:play'"
							class="size-4"
						/>
					</Button>

					<Button
						size="icon"
						class="flex size-8 min-w-8"
						variant="secondary"
						:disabled="!currentVideo"
						@click="skip"
					>
						<Icon
							name="lucide:skip-forward"
							class="size-4"
						/>
					</Button>

					<Slider
						:model-value="[localPosition]"
						:step="1"
						:max="duration || 1"
						:disabled="!currentVideo"
						@update:model-value="handleSeek"
					/>
					<span class="text-muted-foreground text-xs whitespace-nowrap">{{ formattedTime }}</span>
				</div>

				<div class="flex items-center gap-2">
					<Icon
						name="lucide:volume-2"
						class="size-4 shrink-0 text-muted-foreground"
					/>
					<Slider
						:model-value="[localVolume]"
						:step="1"
						:max="100"
						@update:model-value="handleVolumeChange"
					/>
				</div>

				<template v-if="currentVideo">
					<div class="flex items-center gap-2 text-sm">
						<Icon
							name="lucide:radio"
							class="size-4 shrink-0 text-primary"
						/>
						<span class="truncate font-medium">{{ currentVideo.title }}</span>
					</div>

					<div class="flex flex-col gap-1">
						<div class="h-2 overflow-hidden rounded-full bg-muted">
							<div
								class="h-full rounded-full bg-primary transition-[width] duration-500"
								:style="{ width: `${progressPercent}%` }"
							/>
						</div>
						<span class="text-muted-foreground text-xs">{{ formattedTime }}</span>
					</div>

					<div class="text-muted-foreground text-sm">Playback is handled by the overlay.</div>
				</template>

				<div
					v-else
					class="flex w-full flex-col items-center justify-center gap-2 py-4"
				>
					<Icon
						name="lucide:music"
						class="text-muted-foreground size-6"
					/>
					<p class="text-muted-foreground text-sm">{{ t('songRequests.waiting') }}</p>
				</div>
			</div>
		</CardContent>

		<CardFooter class="flex-col items-start gap-4 border-t pt-2">
			<template v-if="currentVideo">
				<div class="flex w-full flex-col gap-2">
					<div
						v-if="currentQueueItem"
						class="flex items-center gap-2"
					>
						<Icon
							name="lucide:user"
							class="size-4 shrink-0"
						/>
						<span>{{ currentQueueItem.orderedByDisplayName || currentQueueItem.orderedByName }}</span>
					</div>

					<div
						v-if="currentQueueItem"
						class="flex items-center gap-2"
					>
						<Icon
							name="lucide:link"
							class="size-4 shrink-0"
						/>
						<a
							class="truncate text-sm underline"
							:href="currentQueueItem.songLink ?? `https://youtu.be/${currentVideo.videoId}`"
							target="_blank"
						>
							{{ currentQueueItem.songLink || `youtu.be/${currentVideo.videoId}` }}
						</a>
					</div>
				</div>

				<div class="mb-2 flex w-full justify-end gap-2">
					<Button
						size="sm"
						variant="secondary"
						@click="skip"
					>
						<Icon
							name="lucide:skip-forward"
							class="mr-1 size-4"
						/>
						Skip
					</Button>

					<AlertDialog>
						<AlertDialogTrigger as-child>
							<Button
								size="sm"
								variant="destructive"
							>
								<Icon
									name="lucide:ban"
									class="mr-1 size-4"
								/>
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

					<AlertDialog v-if="currentQueueItem">
						<AlertDialogTrigger as-child>
							<Button
								size="sm"
								variant="destructive"
							>
								<Icon
									name="lucide:ban"
									class="mr-1 size-4"
								/>
								{{ t('songRequests.ban.user') }}
							</Button>
						</AlertDialogTrigger>
						<AlertDialogContent>
							<AlertDialogHeader>
								<AlertDialogTitle>{{ t('songRequests.ban.userConfirm') }}</AlertDialogTitle>
							</AlertDialogHeader>
							<AlertDialogFooter>
								<AlertDialogCancel>{{ t('deleteConfirmation.cancel') }}</AlertDialogCancel>
								<AlertDialogAction @click="banUser(currentQueueItem.orderedByName)">
									{{ t('deleteConfirmation.confirm') }}
								</AlertDialogAction>
							</AlertDialogFooter>
						</AlertDialogContent>
					</AlertDialog>
				</div>
			</template>
		</CardFooter>
	</Card>
</template>
