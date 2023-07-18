<script setup lang='ts'>
import { IconBroadcast } from '@tabler/icons-vue';
import { NSkeleton, NCard, NSpace, NText, NModal, NButton, useMessage } from 'naive-ui';
import { ref, computed } from 'vue';

import Settings from './obs/settings.vue';

import { useObsOverlayManager, useProfile } from '@/api/index.js';

const messages = useMessage();
const obsSettingsManager = useObsOverlayManager();
const obsSettings = obsSettingsManager.getSettings();

const isModalOpened = ref(false);

const userProfile = useProfile();
const overlayLink = computed(() => {
	return `${window.location.origin}/overlays/${userProfile.data?.value?.apiKey}/obs`;
});

const copyOverlayLink = () => {
	navigator.clipboard.writeText(overlayLink.value);
	messages.success('Copied link url, paste it in obs as browser source');
	return overlayLink;
};
</script>

<template>
  <n-card
    class="overlay-item"
    content-style="padding: 0px" @click="isModalOpened = true"
  >
    <n-skeleton v-if="obsSettings.isLoading.value" size="large" :repeat="4" />
    <n-space v-else vertical align="center">
      <IconBroadcast style="width: 112px; height: 112px" />
      <n-text strong style="font-size: 50px">
        OBS
      </n-text>
    </n-space>
  </n-card>

  <n-modal
    v-model:show="isModalOpened"
    :mask-closable="false"
    :segmented="true"
    preset="card"
    title="OBS"
    content-style="padding: 10px; width: 100%"
    style="width: 500px; max-width: calc(100vw - 40px);"
  >
    <template #header-extra>
      <n-button secondary type="success" @click="copyOverlayLink">
        Copy link url
      </n-button>
    </template>

    <Settings />
  </n-modal>
</template>

