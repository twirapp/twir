<script setup lang='ts'>
import {
  NGrid,
  NGridItem,
} from 'naive-ui';

import OBS from '@/components/overlays/obs.vue';
import TTS from '@/components/overlays/tts.vue';
</script>

<template>
  <div style="display: flex; align-items: center; justify-content: center;">
    <n-grid cols="1 s:1 m:1 l:2" :x-gap="24" :y-gap="10" responsive="screen">
      <n-grid-item :span="1">
        <TTS />
      </n-grid-item>
      <n-grid-item :span="1">
        <OBS />
      </n-grid-item>
    </n-grid>
  </div>
</template>

