<script setup lang='ts'>
import {
  NGrid,
  NGridItem,
  NModal,
} from 'naive-ui';
import { ref } from 'vue';

import { useYoutubeSocket } from '@/components/songRequests/hook.js';
import Player from '@/components/songRequests/player.vue';
import VideosQueue from '@/components/songRequests/queue.vue';
import SettingsModal from '@/components/songRequests/settings.vue';

const isSettingsModalOpened = ref(false);
const openSettingsModal = () => isSettingsModalOpened.value = true;

const {
	videos,
	currentVideo,
	nextVideo,
	deleteVideo,
	deleteAllVideos,
	moveVideo,
} = useYoutubeSocket();
</script>

<template>
  <n-grid cols="1 s:1 m:1 l:3" responsive="screen" :y-gap="15" :x-gap="15">
    <n-grid-item :span="1">
      <player
        :current-video="currentVideo"
        :next-video="videos.length > 1"
        :open-settings-modal="openSettingsModal"
        @next="nextVideo"
      />
    </n-grid-item>

    <n-grid-item :span="2">
      <videos-queue
        :queue="videos"
        @delete-video="(id) => deleteVideo(id)"
        @delete-all-videos="deleteAllVideos"
        @move-video="(id, index) => moveVideo(id, index)"
      />
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

