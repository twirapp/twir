<script setup lang='ts'>
import {
  NGrid,
  NGridItem,
	useThemeVars,
} from 'naive-ui';
import { computed } from 'vue';

import Tts from '@/components/overlays/tts.vue';

const themeVars = useThemeVars();
const primaryColor = computed(() => themeVars.value.hoverColor);
</script>

<template>
  <n-grid cols="1 s:2 m:6 l:15 xl:21" :x-gap="10" :y-gap="10" responsive="screen">
    <n-grid-item :span="3">
      <tts />
    </n-grid-item>
    <n-grid-item :span="3">
      <tts />
    </n-grid-item>
  </n-grid>
</template>

<style>
.overlay-item:hover {
	cursor: pointer;
	background: v-bind(primaryColor);
}
</style>
