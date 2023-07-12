<script setup lang='ts'>
import { NTable } from 'naive-ui';

import Donationalerts from '@/components/integrations/donationalerts.vue';
import Faceit from '@/components/integrations/faceit.vue';
import Lastfm from '@/components/integrations/lastfm.vue';
import Spotify from '@/components/integrations/spotify.vue';
import Streamlabs from '@/components/integrations/streamlabs.vue';
import Vk from '@/components/integrations/vk.vue';
</script>

<template>
  <div class="integrations">
    <n-table :bordered="false" :single-line="false" size="large">
      <thead style="display: none">
        <tr>
          <th></th>
          <th></th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <Spotify />
        <Lastfm />
        <Vk />
        <Streamlabs />
        <Donationalerts />
        <Faceit />
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
	text-align: center;
}

.integrations table tr td:last-of-type {
	width: 70px;
}

.integrations table tr td {
	white-space: nowrap;   /* Prevent line breaks within cells */
	overflow: hidden;      /* Hide overflowing content */
	text-overflow: ellipsis;  /* Add an ellipsis (...) for overflow */
	border-right: 0;
}
</style>
