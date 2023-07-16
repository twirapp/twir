<script lang='ts' setup>
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
} from 'naive-ui';
import Plyr from 'plyr';
import { ref, onMounted, watch, onUnmounted, computed } from 'vue';

import { convertMillisToTime } from '@/components/songRequests/helpers.js';
import { Video } from '@/components/songRequests/hook.js';

const props = defineProps<{
	currentVideo: Video | null
	nextVideo: boolean
}>();

const emits = defineEmits<{
	next: []
}>();

const player = ref<HTMLVideoElement | null>(null);
const plyr = ref<Plyr | null>(null);

const playNext = () => {
	emits('next');

	plyr.value!.once('ready', () => {
		plyr.value!.play();
	});
};

const isPlaying = ref(false);
const sliderTime = ref(0);

onMounted(() => {
	if (!player.value) return;

	plyr.value = new Plyr(player.value, {
		controls: ['fullscreen', 'settings'],
		settings: ['quality', 'speed'],
		hideControls: true,
		clickToPlay: false,
	});

	plyr.value.on('play', () => {
		isPlaying.value = true;
	});

	plyr.value.on('pause', () => {
		isPlaying.value = false;
	});

	plyr.value.on('timeupdate', () => {
		sliderTime.value = plyr.value!.currentTime;
	});

	plyr.value.on('ended', () => {
		if (!props.nextVideo) return;
		playNext();
	});
});

watch(() => props.currentVideo, (video) => {
	if (!plyr.value) return;
	if (!video) {
		plyr.value.source = {
			type: 'video',
			sources: [],
		};
		plyr.value.stop();
		return;
	}

	// plyr.value!.once('ready', () => {
	// 	plyr.value!.play();
	// });

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

});

onUnmounted(() => {
	if (!plyr.value) return;
	plyr.value.destroy();
});

const playerDisplay = useLocalStorage<string>('twirPlayerIsHidden', 'block');
const isMuted = useLocalStorage('twirPlayerIsMuted', false);
const volume = useLocalStorage('twirPlayerVolume', 10);
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
</script>

<template>
  <n-card
    title="Current Song"
    content-style="padding: 0;"
    header-style="padding: 10px;"
    segmented
  >
    <template #header-extra>
      <n-space :wrap="false" :wrap-item="false">
        <n-button tertiary size="small" @click="playerDisplay = playerDisplay === 'block' ? 'none' : 'block'">
          <IconEyeOff v-if="playerDisplay === 'block'" />
          <IconEye v-else />
        </n-button>
        <n-button tertiary size="small" @click="() => {}">
          <IconSettings />
        </n-button>
      </n-space>
    </template>
    <video
      ref="player"
      :style="{
        height: '300px',
      }"
      class="plyr"
    />

    <n-space vertical class="card-content">
      <n-grid :cols="24" :x-gap="10" style="align-items: center; margin-top: 10px; margin-bottom: 10px" responsive="screen">
        <n-grid-item :span="4">
          <n-space align="center">
            <n-button
              size="tiny"
              text
              round
              :disabled="currentVideo == null"
              style="display: flex"
              @click="isPlaying ? plyr?.pause() : plyr?.play()"
            >
              <IconPlayerPlayFilled v-if="!isPlaying" />
              <IconPlayerPauseFilled v-else />
            </n-button>
            <n-button
              style="display: flex"
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

        <n-grid-item :span="14">
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

        <n-grid-item :span="6">
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

    <template #footer>
      <n-list v-if="currentVideo" :show-divider="false">
        <n-list-item>
          <template #prefix>
            <IconPlaylist class="card-icon" />
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

          <a :href="`https://youtu.be/${currentVideo?.videoId}`" class="card-song-link" target="_blank">youtu.be/{{ currentVideo?.videoId }}</a>
        </n-list-item>
      </n-list>
      <n-empty v-else description="Waiting for songs">
        <template #icon>
          <n-spin size="small" stroke="#959596" />
        </template>
      </n-empty>
    </template>
  </n-card>
</template>

<style scoped>
.card-content {
	padding-left: 15px;
	padding-right: 15px
}

.card-icon {
	display: flex
}

.card-song-link {
	color: #63e2b7;
	text-decoration: none
}
</style>

<style>
.plyr {
	display: v-bind(playerDisplay);
}
</style>
