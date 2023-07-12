<script setup lang='ts'>
import { NBadge } from 'naive-ui';

import { useSpotifyIntegration } from '@/api/index.js';
import SpotifySVG from '@/assets/icons/integrations/spotify.svg?component';

const manager = useSpotifyIntegration();
const { data } = manager.useData();
</script>

<template>
  <tr>
    <td><SpotifySVG style="width: 30px" /></td>
    <td><n-badge :value="'Not logged in'" type="info" /></td>
    <td>asdasd</td>
  </tr>
</template>

