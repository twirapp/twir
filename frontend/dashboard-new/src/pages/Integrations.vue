<script setup lang='ts'>
import { NTable, NSpace } from 'naive-ui';

import Spotify from '@/components/integrations/spotify.vue';
</script>

<template>
  <n-space justify="center">
    <n-table :bordered="false" :single-line="false" :striped="false" size="large" style="width: 50%">
      <thead style="display: none">
        <tr>
          <th></th>
          <th></th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <Spotify />
      </tbody>
    </n-table>
  </n-space>
</template>
