<script setup lang='ts'>
import {
  NGrid,
  NGridItem,
	useThemeVars,
} from 'naive-ui';
import { computed } from 'vue';

import OBS from '@/components/overlays/obs.vue';
import TTS from '@/components/overlays/tts.vue';

const themeVars = useThemeVars();
const primaryColor = computed(() => themeVars.value.hoverColor);
</script>

<template>
  <n-grid cols="1 s:2 m:3 l:5" :x-gap="10" :y-gap="10" responsive="screen">
    <n-grid-item :span="1">
      <TTS />
    </n-grid-item>
    <n-grid-item :span="1">
      <OBS />
    </n-grid-item>
  </n-grid>
</template>

<style>
.overlay-item:hover {
	cursor: pointer;
	background: v-bind(primaryColor);
}
</style>
