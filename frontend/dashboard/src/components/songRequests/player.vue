<script setup lang='ts'>
import 'plyr/dist/plyr.css';
import { useVideoQueue } from '@mellkam/vue-plyr-queue';
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
} from 'naive-ui';
import Plyr from 'plyr';
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';

import { convertMillisToTime } from '@/components/songRequests/helpers.js';
import { useQueue } from '@/components/songRequests/hook.js';

defineProps<{
	showSettingsModal: () => void
}>();

const player = ref<HTMLVideoElement | null>(null);
const playerDisplay = useLocalStorage<string>('twirPlayerIsHidden', 'block');

const isPlaying = ref(false);

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

const sliderTime = ref(0);

const marks = ref({});

const setMarks = (duration: number) => {
	marks.value = {
		'0': '0:00',
		[duration.toString()]: convertMillisToTime(duration * 1000),
	};
};

const formatLabelTime = (v: number) => {
	return `${convertMillisToTime(v * 1000)}/${convertMillisToTime((plyr.value?.duration ?? 0) * 1000)}`;
};
const duration = ref(0);

const plyr = computed(() => {
	if (!player.value) return null;

	const p = new Plyr(player.value, {
		controls: ['fullscreen', 'settings'],
		settings: ['quality', 'speed'],
		hideControls: true,
		clickToPlay: false,
	});

	p.on('play', () => {
		isPlaying.value = true;
	});

	p.on('pause', () => {
		isPlaying.value = false;
	});

	p.on('ready', () => {
		duration.value = p.duration ?? 0;
		setMarks(duration.value);
	});

	p.on('timeupdate', () => {
		sliderTime.value = p.currentTime;
	});

	return p;
});

onMounted(() => {
	if (!plyr.value) {
		throw new Error('Plyr is not initialized');
	}
});
onUnmounted(() => plyr.value?.destroy());

const s = useQueue(plyr);

const { currentVideo, queue, nextVideo } = useVideoQueue(
	{
		plyr,
		initialQueue: [
			{ id: '1', src: 'https://www.youtube.com/watch?v=2-1ymGpV_1A&list=LL&index=4' },
			{ id: '1', src: 'https://www.youtube.com/watch?v=P4ALDytLAXQ' },
		],
		defaultProvider: 'youtube',
	},
);

const nextVideoAndAutoplay = () => {
	nextVideo();
	plyr.value?.once('ready', () => {
		plyr.value?.play();
	});
};

const canSkip = computed(() => {
	return currentVideo.value != null || queue.value.length >= 1;
});
</script>

<template>
  <n-card
    title="Card Slots Demo"
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
        <n-button tertiary size="small" @click="showSettingsModal()">
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
        <n-grid-item :span="3">
          <n-space :wrap-item="false" :wrap="false" align="center">
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
              :disabled="!canSkip" @click="nextVideoAndAutoplay()"
            >
              <IconPlayerSkipForwardFilled />
            </n-button>
          </n-space>
        </n-grid-item>

        <n-grid-item :span="15">
          <n-slider
            v-model:value="sliderTime"
            :format-tooltip="formatLabelTime"
            :step="1"
            :max="duration"
            placement="bottom"
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
      <n-list :show-divider="false">
        <n-list-item>
          <template #prefix>
            <IconPlaylist class="card-icon" />
          </template>

          Anna Yvette - Shooting Star [Forza Horizon 4 Pulse] - Synthwave, Retrowave, Synthpop
        </n-list-item>

        <n-list-item>
          <template #prefix>
            <IconUser class="card-icon" />
          </template>

          Satont
        </n-list-item>

        <n-list-item>
          <template #prefix>
            <IconLink class="card-icon" />
          </template>

          <a href="https://youtu.be/ZXgHcdLM7lM" class="card-song-link" target="_blank">youtu.be/ZXgHcdLM7lM</a>
        </n-list-item>
      </n-list>
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
