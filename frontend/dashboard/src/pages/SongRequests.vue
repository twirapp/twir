<script setup lang='ts'>
import {
	NGrid,
	NGridItem,
	NModal,
	NSpace,
	NButton,
} from 'naive-ui';
import { ref } from 'vue';

import Player from '@/components/songRequests/player.vue';
import SettingsModal from '@/components/songRequests/settings.vue';

const isSettingsModalOpened = ref(false);
const showSettingsModal = () => isSettingsModalOpened.value = true;
</script>

<template>
  <n-grid cols="1 s:1 m:1 l:1 xl:12" :x-gap="15" responsive="screen">
    <n-grid-item :span="4">
      <player :showSettingsModal="showSettingsModal" />
    </n-grid-item>

    <n-grid-item>
      queue
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

