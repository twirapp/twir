<script lang="ts" setup>
import 'plyr/dist/plyr.css'
import {
	IconBan,
	IconEye,
	IconEyeOff,
	IconLink,
	IconPlayerPauseFilled,
	IconPlayerPlayFilled,
	IconPlayerSkipForwardFilled,
	IconPlaylist,
	IconSettings,
	IconUser,
	IconVolume,
	IconVolume3,
} from '@tabler/icons-vue'
import { useLocalStorage } from '@vueuse/core'
import {
	NCard,
	NEmpty,
	NList,
	NListItem,
	NPopconfirm,
	NResult,
	NSlider,
	NSpin,
} from 'naive-ui'
import Plyr from 'plyr'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { useProfile } from '@/api/index.js'
import { useYoutubeSocket } from '@/components/songRequests/hook.js'
import { Button } from '@/components/ui/button'
import { convertMillisToTime } from '@/helpers/convertMillisToTime.js'

const props = defineProps<{
	noCookie: boolean
	openSettingsModal: () => void
}>()

const {
	currentVideo,
	nextVideo,
	sendPlaying,
	banSong,
	banUser,
} = useYoutubeSocket()

const player = ref<HTMLVideoElement>()
const plyr = ref<Plyr>()

function playNext() {
	nextVideo()

	plyr.value!.once('ready', () => {
		plyr.value!.play()
	})
}

const isPlaying = ref(false)
const sliderTime = ref(0)

const volume = useLocalStorage('twirPlayerVolume', 10)
const playerDisplay = useLocalStorage('twirPlayerIsHidden', 'block')
const isMuted = useLocalStorage('twirPlayerIsMuted', false)

watch(isMuted, (v) => {
	if (!plyr.value) return
	plyr.value.muted = v
})

onMounted(() => {
	if (!player.value) return

	plyr.value = new Plyr(player.value, {
		controls: ['fullscreen', 'settings'],
		settings: ['quality', 'speed'],
		hideControls: true,
		clickToPlay: false,
		youtube: { noCookie: props.noCookie },
	})

	plyr.value.on('play', () => {
		isPlaying.value = true
		plyr.value!.volume = volume.value / 100
		plyr.value!.muted = isMuted.value
		sendPlaying()
	})

	plyr.value.on('pause', () => {
		isPlaying.value = false
	})

	plyr.value.on('timeupdate', () => {
		sliderTime.value = plyr.value!.currentTime
	})

	plyr.value.on('ended', () => {
		// if (!props.nextVideo) return;
		playNext()
	})
})

const isFirstLoad = ref(true)
watch(currentVideo, (video) => {
	if (!plyr.value) return
	if (!video) {
		plyr.value.source = {
			type: 'video',
			sources: [],
		}
		plyr.value.stop()
		return
	}

	if (!isFirstLoad.value) {
		plyr.value!.once('ready', () => {
			plyr.value!.play()
		})
	}

	plyr.value.source = {
		type: 'video',
		sources: [
			{
				src: `https://www.youtube.com/watch?v=${video.videoId}`,
				provider: 'youtube',
			},
		],
		title: '',
	}
	isFirstLoad.value = false
})

onUnmounted(() => {
	if (!plyr.value) return
	plyr.value.destroy()
})

watch(volume, (value) => {
	if (!plyr.value) return

	if (value === 0) {
		isMuted.value = true
		plyr.value.muted = true
	} else {
		isMuted.value = false
		plyr.value!.volume = value / 100
	}
})

const sliderVolume = computed(() => {
	if (isMuted.value) return 0
	return volume.value
})

function formatLabelTime(v: number) {
	return `${convertMillisToTime(v * 1000)}/${convertMillisToTime((plyr.value?.duration ?? 0) * 1000)}`
}

const { data: profile } = useProfile()
const { t } = useI18n()
</script>

<template>
	<NCard
		:title="t('songRequests.player.title')"
		content-style="padding: 0;"
		header-style="padding: 10px;"
		segmented
	>
		<div v-if="profile?.id !== profile?.selectedDashboardId" class="p-2.5">
			<NResult
				status="404"
				:title="t('songRequests.player.noAccess')"
				size="small"
			>
			</NResult>
		</div>

		<div v-else>
			<video
				ref="player"
				:style="{
					height: '300px',
				}"
				class="plyr"
			></video>

			<div class="flex flex-col gap-4 py-5 px-6">
				<div class="flex gap-2 items-center">
					<Button
						size="icon"
						class="flex size-8 min-w-8"
						variant="secondary"
						:disabled="currentVideo == null"
						@click="isPlaying ? plyr?.pause() : plyr?.play()"
					>
						<IconPlayerPlayFilled v-if="!isPlaying" class="size-4" />
						<IconPlayerPauseFilled v-else class="size-4" />
					</Button>

					<Button
						size="icon"
						class="flex size-8 min-w-8"
						variant="secondary"
						:disabled="currentVideo == null"
						@click="playNext"
					>
						<IconPlayerSkipForwardFilled class="size-4" />
					</Button>

					<NSlider
						v-model:value="sliderTime"
						:format-tooltip="formatLabelTime"
						:step="1"
						:max="plyr?.duration ?? 0"
						placement="bottom"
						:disabled="!currentVideo"
						@update-value="(v) => {
							plyr!.currentTime = v
						}"
					/>
				</div>

				<div class="flex items-center gap-2">
					<Button
						size="icon"
						variant="secondary"
						class="size-8 min-w-8"
						@click="isMuted = !isMuted"
					>
						<IconVolume v-if="!isMuted" class="size-4" />
						<IconVolume3 v-else class="size-4" />
					</Button>
					<NSlider :value="sliderVolume" :step="1" @update-value="(v) => volume = v" />
				</div>
			</div>
		</div>

		<template #header-extra>
			<div class="flex gap-2">
				<Button
					variant="secondary"
					size="icon"
					class="size-8"
					@click="playerDisplay = playerDisplay === 'block' ? 'none' : 'block'"
				>
					<IconEyeOff v-if="playerDisplay === 'block'" class="size-4" />
					<IconEye v-else class="size-4" />
				</Button>
				<Button
					variant="secondary"
					size="icon"
					class="size-8"
					@click="openSettingsModal"
				>
					<IconSettings class="size-4" />
				</Button>
			</div>
		</template>
		<template #footer>
			<template v-if="currentVideo">
				<NList :show-divider="false">
					<NListItem>
						<template #prefix>
							<IconPlaylist class="flex" />
						</template>

						{{ currentVideo?.title }}
					</NListItem>

					<NListItem>
						<template #prefix>
							<IconUser class="card-icon" />
						</template>

						{{ currentVideo?.orderedByDisplayName || currentVideo?.orderedByName }}
					</NListItem>

					<NListItem>
						<template #prefix>
							<IconLink class="card-icon" />
						</template>

						<a
							class="underline"
							:href="currentVideo.songLink ?? `https://youtu.be/${currentVideo?.videoId}`"
							target="_blank"
						>
							{{ currentVideo.songLink || `youtu.be/${currentVideo?.videoId}` }}
						</a>
					</NListItem>
				</NList>
				<div class="flex gap-2 justify-end">
					<NPopconfirm
						:positive-text="t('deleteConfirmation.confirm')"
						:negative-text="t('deleteConfirmation.cancel')"
						@positive-click="() => banSong(currentVideo.videoId)"
					>
						<template #trigger>
							<Button size="sm" variant="destructive">
								<div class="flex gap-1 items-center">
									<IconBan />
									{{ t('songRequests.ban.song') }}
								</div>
							</Button>
						</template>

						{{ t('songRequests.ban.songConfirm') }}
					</NPopconfirm>
					<NPopconfirm
						:positive-text="t('deleteConfirmation.confirm')"
						:negative-text="t('deleteConfirmation.cancel')"
						@positive-click="() => banUser(currentVideo.orderedById)
						"
					>
						<template #trigger>
							<Button size="sm" variant="destructive">
								<div class="flex gap-1 items-center">
									<IconBan />
									{{ t('songRequests.ban.user') }}
								</div>
							</Button>
						</template>

						{{ t('songRequests.ban.userConfirm') }}
					</NPopconfirm>
				</div>
			</template>
			<NEmpty v-else :description="t('songRequests.waiting')">
				<template #icon>
					<NSpin size="small" stroke="#959596" />
				</template>
			</NEmpty>
		</template>
	</NCard>
</template>

<style>
.plyr {
	display: v-bind(playerDisplay);
}
</style>
