<script setup lang='ts'>
import { NTable } from 'naive-ui';

import Spotify from '@/components/integrations/spotify.vue';
</script>

<template>
  <div class="integrations">
    <n-table :bordered="false" :single-line="false" :striped="true" size="large">
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
  </div>
</template>

<style>
.integrations {
	display: flex;
	justify-content: center;
	align-items: center;
	height: 100%;
}

.integrations table {
	width: 50%;
}

.integrations table tr td:first-of-type {
	width: 70px;
}

.integrations table tr td {
	white-space: nowrap;   /* Prevent line breaks within cells */
	overflow: hidden;      /* Hide overflowing content */
	text-overflow: ellipsis;  /* Add an ellipsis (...) for overflow */
}
</style>
