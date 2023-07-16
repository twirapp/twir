<script setup lang='ts'>
import {
	NGrid,
	NGridItem,
} from 'naive-ui';

import Player from '@/components/songRequests/player.vue';
</script>

<template>
  <n-grid :cols="12" :x-gap="15">
    <n-grid-item :span="6">
      <player />
    </n-grid-item>

    <n-grid-item :span="6">
      queue
    </n-grid-item>
  </n-grid>
</template>

