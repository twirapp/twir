<script lang="ts" setup>
import 'plyr/dist/plyr.css';
import {
	IconEyeOff,
	IconEye,
	IconSettings,
	IconPlaylist,
	IconUser,
	IconLink,
	IconVolume,
	IconVolume3,
	IconPlayerPlayFilled,
	IconPlayerSkipForwardFilled,
	IconPlayerPauseFilled,
	IconBan,
} from '@tabler/icons-vue';
import { useLocalStorage } from '@vueuse/core';
import {
	NCard,
	NButton,
	NSpace,
	NList,
	NListItem,
	NSlider,
	NGrid,
	NGridItem,
	NEmpty,
	NSpin,
	NResult,
	NPopconfirm,
} from 'naive-ui';
import { storeToRefs } from 'pinia';
import Plyr from 'plyr';
import { ref, onMounted, watch, onUnmounted, computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { useProfile } from '@/api/index.js';
import { useYoutubeSocket } from '@/components/songRequests/hook.js';
import { convertMillisToTime } from '@/helpers/convertMillisToTime.js';

const props = defineProps<{
	noCookie: boolean
	openSettingsModal: () => void
}>();

const socket = useYoutubeSocket();
const { currentVideo } = storeToRefs(socket);

const player = ref<HTMLVideoElement>();
const plyr = ref<Plyr>();

const playNext = () => {
	socket.nextVideo();

	plyr.value!.once('ready', () => {
		plyr.value!.play();
	});
};

const isPlaying = ref(false);
const sliderTime = ref(0);
const volume = useLocalStorage('twirPlayerVolume', 10);

onMounted(() => {
	if (!player.value) return;

	plyr.value = new Plyr(player.value, {
		controls: ['fullscreen', 'settings'],
		settings: ['quality', 'speed'],
		hideControls: true,
		clickToPlay: false,
		youtube: { noCookie: props.noCookie },
	});

	plyr.value.on('play', () => {
		isPlaying.value = true;
		plyr.value!.volume = volume.value / 100;
		socket.sendPlaying();
	});

	plyr.value.on('pause', () => {
		isPlaying.value = false;
	});

	plyr.value.on('timeupdate', () => {
		sliderTime.value = plyr.value!.currentTime;
	});

	plyr.value.on('ended', () => {
		// if (!props.nextVideo) return;
		playNext();
	});
});

const isFirstLoad = ref(true);
watch(currentVideo, (video) => {
	if (!plyr.value) return;
	if (!video) {
		plyr.value.source = {
			type: 'video',
			sources: [],
		};
		plyr.value.stop();
		return;
	}

	if (!isFirstLoad.value) {
		plyr.value!.once('ready', () => {
			plyr.value!.play();
		});
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
	};
	isFirstLoad.value = false;
});

onUnmounted(() => {
	if (!plyr.value) return;
	plyr.value.destroy();
});

const playerDisplay = useLocalStorage<string>('twirPlayerIsHidden', 'block');
const isMuted = useLocalStorage('twirPlayerIsMuted', false);

watch(volume, (value) => {
	if (!plyr.value) return;

	if (value === 0) {
		isMuted.value = true;
		plyr.value.muted = true;
	} else {
		isMuted.value = false;
		plyr.value!.volume = value / 100;
	}
});

watch(isMuted, (v) => {
	if (!plyr.value) return;
	plyr.value.muted = v;
});

const sliderVolume = computed(() => {
	if (isMuted.value) return 0;
	return volume.value;
});
const formatLabelTime = (v: number) => {
	return `${convertMillisToTime(v * 1000)}/${convertMillisToTime((plyr.value?.duration ?? 0) * 1000)}`;
};

const { data: profile } = useProfile();

const { t } = useI18n();
</script>

<template>
	<n-card
		:title="t('songRequests.player.title')"
		content-style="padding: 0;"
		header-style="padding: 10px;"
		segmented
	>
		<div v-if="profile?.id != profile?.selectedDashboardId" class="p-2.5">
			<n-result
				status="404"
				:title="t('songRequests.player.noAccess')"
				size="small"
			>
			</n-result>
		</div>

		<div v-else>
			<video
				ref="player"
				:style="{
					height: '300px',
				}"
				class="plyr"
			></video>

			<n-space vertical class="px-3.5">
				<n-grid
					cols="1 s:1 m:1 l:1 xl:5"
					responsive="screen"
					:x-gap="10"
					:y-gap="10"
					class="items-center my-2.5"
				>
					<n-grid-item :span="1" class="w-full">
						<n-space align="center" justify="center">
							<n-button
								size="tiny"
								text
								round
								:disabled="currentVideo == null"
								class="flex"
								@click="isPlaying ? plyr?.pause() : plyr?.play()"
							>
								<IconPlayerPlayFilled v-if="!isPlaying" />
								<IconPlayerPauseFilled v-else />
							</n-button>
							<n-button
								class="flex"
								size="tiny"
								text
								round
								:disabled="currentVideo == null"
								@click="playNext"
							>
								<IconPlayerSkipForwardFilled />
							</n-button>
						</n-space>
					</n-grid-item>

					<n-grid-item :span="2">
						<n-slider
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
					</n-grid-item>

					<n-grid-item :span="2">
						<n-space :wrap-item="false" :wrap="false" align="center">
							<n-button size="tiny" text round>
								<IconVolume v-if="!isMuted" @click="isMuted = true" />
								<IconVolume3 v-else @click="isMuted = false" />
							</n-button>
							<n-slider :value="sliderVolume" :step="1" @update-value="(v) => volume = v" />
						</n-space>
					</n-grid-item>
				</n-grid>
			</n-space>
		</div>

		<template #header-extra>
			<n-space :wrap="false" :wrap-item="false">
				<n-button
					tertiary size="small"
					@click="playerDisplay = playerDisplay === 'block' ? 'none' : 'block'"
				>
					<IconEyeOff v-if="playerDisplay === 'block'" />
					<IconEye v-else />
				</n-button>
				<n-button tertiary size="small" @click="openSettingsModal">
					<IconSettings />
				</n-button>
			</n-space>
		</template>
		<template #footer>
			<template v-if="currentVideo">
				<n-list :show-divider="false">
					<n-list-item>
						<template #prefix>
							<IconPlaylist class="flex" />
						</template>

						{{ currentVideo?.title }}
					</n-list-item>

					<n-list-item>
						<template #prefix>
							<IconUser class="card-icon" />
						</template>

						{{ currentVideo?.orderedByDisplayName || currentVideo?.orderedByName }}
					</n-list-item>

					<n-list-item>
						<template #prefix>
							<IconLink class="card-icon" />
						</template>

						<n-button
							tag="a"
							type="primary"
							text
							:href="currentVideo.songLink ?? `https://youtu.be/${currentVideo?.videoId}`"
							target="_blank"
						>
							{{ currentVideo.songLink || `youtu.be/${currentVideo?.videoId}` }}
						</n-button>
					</n-list-item>
				</n-list>
				<n-space justify="end">
					<n-popconfirm
						:positive-text="t('deleteConfirmation.confirm')"
						:negative-text="t('deleteConfirmation.cancel')"
						@positive-click="() => socket.banSong(currentVideo.videoId)"
					>
						<template #trigger>
							<n-button
								secondary
								type="warning"
							>
								<div class="flex gap-1 items-center">
									<IconBan />
									{{ t('songRequests.ban.song') }}
								</div>
							</n-button>
						</template>

						{{ t('songRequests.ban.songConfirm') }}
					</n-popconfirm>
					<n-popconfirm
						:positive-text="t('deleteConfirmation.confirm')"
						:negative-text="t('deleteConfirmation.cancel')"
						@positive-click="() => socket.banUser(currentVideo.orderedById)
						"
					>
						<template #trigger>
							<n-button
								secondary
								type="error"
							>
								<div class="flex gap-1 items-center">
									<IconBan />
									{{ t('songRequests.ban.user') }}
								</div>
							</n-button>
						</template>

						{{ t('songRequests.ban.userConfirm') }}
					</n-popconfirm>
				</n-space>
			</template>
			<n-empty v-else :description="t('songRequests.waiting')">
				<template #icon>
					<n-spin size="small" stroke="#959596" />
				</template>
			</n-empty>
		</template>
	</n-card>
</template>

<style>
.plyr {
	display: v-bind(playerDisplay);
}
</style>
