<script lang="ts" setup>
import { useTitle } from '@vueuse/core';

import Lastfm from '@/components/integrations/Lastfm.vue';
import Spotify from '@/components/integrations/Spotify.vue';
import Vk from '@/components/integrations/Vk.vue';

const title = useTitle();
title.value = 'Tsuwari - Integrations';
</script>

<template>
  <div class="grid xl:grid-cols-4 lg:grid-cols-4 md:grid-cols-2 grid-cols-1 gap-2 m-3">
    <Spotify />

    <div class="grid grid-cols-1">
      <div><Vk /></div>
      <div class="mt-1">
        <Lastfm />
      </div>
    </div>
  </div>
</template>

<style scoped>
@media (max-width: 426px) {
  .small-screen {
    width: 100%;
  }
}
</style>
