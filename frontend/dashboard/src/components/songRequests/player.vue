<script setup lang='ts'>
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
import { ref, computed, onMounted, onUnmounted, Ref, watch } from 'vue';

defineProps<{
	showSettingsModal: () => void
}>();

const player = ref<HTMLVideoElement | null>(null);

const plyr = computed(() => {
	if (!player.value) return null;

	return new Plyr(player.value);
});

onMounted(() => {
	if (!plyr.value) {
		throw new Error('Plyr is not initialized');
	}
});

onUnmounted(() => {
	if (!plyr.value) return;

	plyr.value!.destroy();
});

const isPlayerHidden = useLocalStorage('twirPlayerIsHidden', false);

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
            <n-button size="tiny" text round>
              <IconPlayerPlayFilled />
            </n-button>
            <n-button size="tiny" text round>
              <IconPlayerPlayFilled />
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
