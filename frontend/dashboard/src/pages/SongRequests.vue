<script setup lang='ts'>
import {
	NGrid,
	NGridItem,
	NModal,
	NSpace,
	NButton,
} from 'naive-ui';
import { ref } from 'vue';

import { useYoutubeSocket } from '@/components/songRequests/hook.js';
import Player from '@/components/songRequests/player.vue';
import VideosQueue from '@/components/songRequests/queue.vue';
import SettingsModal from '@/components/songRequests/settings.vue';

const isSettingsModalOpened = ref(false);
const showSettingsModal = () => isSettingsModalOpened.value = true;

const { videos, currentVideo, nextVideo, deleteVideo } = useYoutubeSocket();
</script>

<template>
  <n-grid cols="1 s:1 m:1 l:12 xl:12" :x-gap="15" responsive="screen">
    <n-grid-item :span="4">
      <player :current-video="currentVideo" :next-video="videos.length > 1" @next="nextVideo" />
    </n-grid-item>

    <n-grid-item :span="8">
      <videos-queue :queue="videos" @delete-video="(id) => deleteVideo(id)" />
    </n-grid-item>
  </n-grid>

  <n-modal
    v-model:show="isSettingsModalOpened"
    :span="10"
    :mask-closable="false"
    :segmented="true"
    preset="card"
    title="Settings"
    :style="{ width: '600px',top: '50px' }"
  >
    <settings-modal />
  </n-modal>
</template>

