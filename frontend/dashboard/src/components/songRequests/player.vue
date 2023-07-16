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

defineProps<{
	showSettingsModal: () => void
}>();

const player = ref<HTMLVideoElement | null>(null);
const isPlayerHidden = useLocalStorage('twirPlayerIsHidden', false);

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

const isFirstVideo = ref(true);
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
		isFirstVideo.value = false;
	});

	p.on('pause', () => {
		isPlaying.value = false;
	});

	return p;
});
onMounted(() => {
	if (!plyr.value) {
		throw new Error('Plyr is not initialized');
	}
});

onUnmounted(() => plyr.value?.destroy());

const { currentVideo, queue, nextVideo, removeVideo, addVideo } = useVideoQueue(
	{
		plyr,
		initialQueue: [
			{ id: '1', src: 'https://www.youtube.com/watch?v=2-1ymGpV_1A&list=LL&index=4' },
			{ id: '1', src: 'https://www.youtube.com/watch?v=BFtseA7Hs4g' },
		],
		defaultProvider: 'youtube',
		onNextVideo: () => {
			if (isFirstVideo.value) return;
			plyr.value?.once('ready', () => {
				plyr.value?.play();
			});
		},
		onRemoveVideo: (video) => {
			console.log(video);
		},
	},
);

const canSkip = computed(() => {
	return currentVideo.value != null || queue.value.length >= 1;
});

const length = ref(0);
</script>

<template>
  <n-card
    title="Card Slots Demo"
    content-style="padding: 0;"
    header-style="padding: 10px;"
  >
    <template #header-extra>
      <n-space>
        <n-button tertiary size="small" @click="isPlayerHidden = !isPlayerHidden">
          <IconEyeOff v-if="!isPlayerHidden" />
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
        display: isPlayerHidden ? 'none' : 'block',
      }"
    />

    <n-space vertical class="card-content">
      <n-grid :cols="24" :x-gap="10" style="align-items: center">
        <n-grid-item :span="3">
          <n-space>
            <n-button
              size="tiny"
              text
              round
              :disabled="currentVideo == null"
              @click="isPlaying ? plyr?.pause() : plyr?.play()"
            >
              <IconPlayerPlayFilled v-if="!isPlaying" />
              <IconPlayerPauseFilled v-else />
            </n-button>
            <n-button size="tiny" text round :disabled="!canSkip" @click="nextVideo">
              <IconPlayerSkipForwardFilled />
            </n-button>
          </n-space>
        </n-grid-item>

        <n-grid-item :span="15">
          <n-slider v-model:value="length" :step="1" :marks="{ 100: '100'}" />
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
    </n-space>
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
